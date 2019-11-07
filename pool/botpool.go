package pool

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/torlenor/abylebotter/api"
	"github.com/torlenor/abylebotter/logging"
	"github.com/torlenor/abylebotter/platform"
)

// BotPool holds a set of bots uniquely identified by an ID
type BotPool struct {
	log *logrus.Entry

	bots  map[string]platform.Bot
	mutex sync.Mutex

	controlAPI *api.API
}

// NewBotPool creates a BotPool with default values.
func NewBotPool(controlAPI *api.API) *BotPool {
	b := &BotPool{
		bots:       make(map[string]platform.Bot),
		log:        logging.Get("BotPool"),
		controlAPI: controlAPI,
	}

	controlAPI.AttachModuleGet("/bots", b.getBotsEndpoint)

	return b
}

// AddViaID will use a BotConfigProvider to get all the necessary information
// to create the bot for the given Bot ID.
// This is an idempotent function
func (b *BotPool) AddViaID(id string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// TODO

	return nil
}

// RemoveViaID removes a bot via its ID.
// This is an idempotent function.
func (b *BotPool) RemoveViaID(id string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	delete(b.bots, id)
}

// Add a bot with the provided id. If a bot with that ID already exists and error is returned.
func (b *BotPool) Add(id string, bot platform.Bot) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if _, ok := b.bots[id]; ok {
		return fmt.Errorf("Bot with ID %s already exists", id)
	}

	return nil
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

// Run the bots in the pool. If a new bot is added to the pool, it will be automatically started.
func (b *BotPool) Run(ctx context.Context) (err error) {
	b.log.Info("BotPool started")
	// TODO

	<-ctx.Done()

	b.log.Info("BotPool shutdown")
	b.log.Info("BotPool exited properly")

	return err
}
