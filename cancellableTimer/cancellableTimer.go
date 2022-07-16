package cancellableTimer

import (
	"sync"
	"time"
)

type CancellableTimer struct {
	kill           bool
	TimerTriggered bool
	active         bool
	cancel         chan bool
	wg             *sync.WaitGroup
}

func NewCancellableTimer() *CancellableTimer {
	return &CancellableTimer{
		cancel:         make(chan bool),
		active:         false,
		kill:           false,
		TimerTriggered: false,
	}
}

func (c *CancellableTimer) Start(d time.Duration, t Timeoutable) {
	if c.active {
		return
	}
	go func() {
		// If its killed before we set false again
		c.kill = false
		c.TimerTriggered = false
		c.active = true
		c.wg = &sync.WaitGroup{}
		timer := time.NewTimer(d)
		defer func() {
			c.active = false
			if !timer.Stop() && !c.TimerTriggered {
				<-timer.C
			}
			if !c.kill {
				t.WhenTimeout()
			}
		}()
		select {
		case <-timer.C:
			c.TimerTriggered = true
			// Waits if a move in process
			c.wg.Wait()
			return
		case <-c.cancel:
			return
		}
	}()
}

func (c *CancellableTimer) Cancel() {
	c.kill = true
	// If timer already triggered there is no need to cancel
	if c.TimerTriggered {
		return
	}
	c.cancel <- true
}

func (c *CancellableTimer) Lock() error {
	if !c.active {
		return &NotActiveError{}
	}
	c.wg.Add(1)
	return nil
}

func (c *CancellableTimer) Unlock(kill bool) {
	if !c.active {
		return
	}
	// Kill is gonna be false if move isn't playable
	if kill {
		c.Cancel()
	}
	c.wg.Done()
}
