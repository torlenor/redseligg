package utils

import (
	"sync"
	"time"

	"github.com/torlenor/redseligg/logging"
)

type failCallback func()

// Watchdog is a generic dog that watches and barks if something fails (by calling a fail callback)
type Watchdog struct {
	failCallback failCallback

	timer *time.Timer

	food chan bool
	done chan bool

	startStopMutex sync.Mutex
	isRunning      bool
}

// SetFailCallback function to be called when watchdog times out
func (w *Watchdog) SetFailCallback(c failCallback) *Watchdog {
	w.failCallback = c
	return w
}

// Start the watchdog; in case it is not fed, it will keep
// notifying until it is stopped.
func (w *Watchdog) Start(interval time.Duration) {
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
				if w.failCallback != nil {
					go w.failCallback()
				} else {
					logging.Get("Watchdog").Infof("Watchdog not fed in time, but no callback set")
				}
			}
		}
	}(interval)
}

// Stop the watchdog
func (w *Watchdog) Stop() {
	w.startStopMutex.Lock()
	defer w.startStopMutex.Unlock()

	if w.isRunning {
		w.done <- true
		w.isRunning = false
	}
}

// Feed the dog
func (w *Watchdog) Feed() {
	w.startStopMutex.Lock()
	defer w.startStopMutex.Unlock()

	if w.isRunning {
		w.food <- true
	}
}
