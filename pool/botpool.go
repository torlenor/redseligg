package pool

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/torlenor/abylebotter/api"
	"github.com/torlenor/abylebotter/logging"
	"github.com/torlenor/abylebotter/platform"
	"github.com/torlenor/abylebotter/providers"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"
)

// BotPool holds a set of bots uniquely identified by an ID
type BotPool struct {
	log *logrus.Entry

	context       context.Context
	childRoutines *errgroup.Group

	bots        map[string]platform.Bot
	botContexts map[string]context.Context
	shutdownFns map[string]context.CancelFunc

	wg sync.WaitGroup

	mutex sync.Mutex

	controlAPI *api.API

	botProvider *providers.BotProvider

	checkerStop chan bool

	isRunning bool
}

// NewBotPool creates a BotPool with default values.
func NewBotPool(controlAPI *api.API, botProvider *providers.BotProvider) (*BotPool, error) {
	b := &BotPool{
		bots:        make(map[string]platform.Bot),
		botContexts: make(map[string]context.Context),
		shutdownFns: make(map[string]context.CancelFunc),

		log: logging.Get("BotPool"),

		controlAPI: controlAPI,

		botProvider: botProvider,
	}

	if controlAPI != nil {
		controlAPI.AttachModuleGet("/bots", b.getBotsEndpoint)
		controlAPI.AttachModulePost("/bots", b.postBotsEndpoint)

		controlAPI.AttachModuleGet("/bots/{botId}", b.getBotEndPoint)
		controlAPI.AttachModuleDelete("/bots/{botId}", b.deleteBotEndpoint)
	}

	return b, nil
}

// AddViaID will use a BotConfigProvider to get all the necessary information
// to create the bot for the given Bot ID.
// This is an idempotent function
func (b *BotPool) AddViaID(id string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if _, ok := b.bots[id]; ok {
		return fmt.Errorf("Bot with ID %s already exists", id)
	}

	bot, err := b.botProvider.GetBot(id)
	if err != nil {
		return fmt.Errorf("Not possible to add bot with id %s: %s", id, err)
	}

	b.bots[id] = bot

	b.startSingle(id)

	return nil
}

// RemoveViaID removes a bot via its ID.
// This is an idempotent function.
func (b *BotPool) RemoveViaID(id string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.log.Debugf("Removing bot %s", id)

	if _, ok := b.bots[id]; !ok {
		b.log.Debugf("Does not exist, ignoring %s", id)
		return
	}

	b.stopSingle(id)

	delete(b.bots, id)
}

// GetBotIDs returns all known BotIDs in no particular order
func (b *BotPool) GetBotIDs() []string {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.bots == nil {
		return []string{}
	}

	ids := []string{}

	for key := range b.bots {
		ids = append(ids, key)
	}

	return ids
}

func (b *BotPool) startSingle(id string) {
	if !b.isRunning {
		b.log.Infof("Not starting Bot with id %s because we are not running", id)
		return
	}

	var bot platform.Bot
	var ok bool
	if bot, ok = b.bots[id]; !ok {
		b.log.Warnf("Cannot start Bot with ID %s, does not exist", id)
	}

	ctx, botShutdownFn := context.WithCancel(b.context)
	b.botContexts[id] = ctx
	b.shutdownFns[id] = botShutdownFn

	b.childRoutines.Go(func() error {
		if !b.isRunning {
			return nil
		}

		if err := bot.Run(ctx); err != nil {
			if err != context.Canceled {
				b.log.Errorln("Stopped bot", "reason", err)
			} else {
				b.log.Infoln("Stopped bot", "reason", err)
			}
		}

		return nil
	})
}

func (b *BotPool) stopSingle(id string) {
	if _, ok := b.bots[id]; !ok {
		b.log.Warnf("Cannot stop, Bot with ID %s does not exist", id)
	}

	b.shutdownFns[id]()

	delete(b.shutdownFns, id)
	delete(b.botContexts, id)
}

func (b *BotPool) restartSingle(id string) {
	b.stopSingle(id)
	b.startSingle(id)
}

// checkBots continuously monitors the bots and tries to restart them if they are not healthy
func (b *BotPool) checkBots(interval time.Duration, stop chan bool) {
	b.log.Debugf("Bot monitoring started")
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-stop:
			ticker.Stop()
			b.log.Debugf("Bot monitoring stopped")
			return
		case <-ticker.C:
			b.log.Tracef("Running bots check")
			b.mutex.Lock()
			for id, bot := range b.bots {
				b.log.Tracef("Checking bot with ID %s", id)
				if !bot.GetInfo().Healthy {
					b.log.Warnf("Bot %s unhealthy. Restarting it", id)
					b.restartSingle(id)
				}
			}
			b.mutex.Unlock()
		}
	}
}

// Run the BotPool. If a new bot is added to the pool, it will be automatically started.
func (b *BotPool) Run(ctx context.Context) (err error) {
	b.log.Info("BotPool started")

	childRoutines, childCtx := errgroup.WithContext(ctx)
	b.childRoutines = childRoutines
	b.context = childCtx

	b.checkerStop = make(chan bool)
	go func() {
		b.wg.Add(1)
		b.checkBots(5*time.Second, b.checkerStop)
		defer b.wg.Done()
	}()

	b.isRunning = true

	<-ctx.Done()

	b.isRunning = false

	b.checkerStop <- true
	b.wg.Wait()

	if err := b.childRoutines.Wait(); err != nil && !xerrors.Is(err, context.Canceled) {
		b.log.Errorln("Failed waiting for services to shutdown", "err", err)
	}

	b.log.Info("BotPool shutdown")
	b.log.Info("BotPool exited properly")

	return err
}
