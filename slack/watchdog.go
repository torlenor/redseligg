package slack

import (
	"sync"
	"time"

	"github.com/torlenor/abylebotter/logging"
)

type failCallback func()

type watchdog struct {
	failCallback failCallback

	timer *time.Timer

	food chan bool
	done chan bool

	startStopMutex sync.Mutex
	isRunning      bool
}

// SetFailCallback function to be called when watchdog times out
func (w *watchdog) SetFailCallback(c failCallback) *watchdog {
	w.failCallback = c
	return w
}

// Start the watchdog
func (w *watchdog) Start(interval time.Duration) {
	w.startStopMutex.Lock()
	defer w.startStopMutex.Unlock()
	w.isRunning = true

	w.food = make(chan bool)
	w.done = make(chan bool)

	w.timer = time.NewTimer(interval)
	go func(interval time.Duration) {
		for {
			select {
			case <-w.done:
				if !w.timer.Stop() {
					<-w.timer.C
				}
				return
			case <-w.food:
				if !w.timer.Stop() {
					<-w.timer.C
				}
				w.timer.Reset(interval)
			case <-w.timer.C:
				w.startStopMutex.Lock()
				w.isRunning = false
				w.startStopMutex.Unlock()

				if w.failCallback != nil {
					w.failCallback()
				} else {
					logging.Get("Watchdog").Infof("Watchdog not fed in time. No callback set")
				}
				return
			}
		}
	}(interval)
}

// Stop the watchdog
func (w *watchdog) Stop() {
	w.startStopMutex.Lock()
	defer w.startStopMutex.Unlock()

	if w.isRunning {
		w.done <- true
		w.isRunning = false
	}
}

// Feed the dog
func (w *watchdog) Feed() {
	w.startStopMutex.Lock()
	defer w.startStopMutex.Unlock()

	if w.isRunning {
		w.food <- true
	}
}
