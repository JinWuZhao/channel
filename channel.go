package channel

import (
	"fmt"
	"sync"
	"time"
)

type Channel struct {
	ch              chan Cargo
	isClosed        bool
	name            string
	sendDeadLine    *time.Time
	receiveDeadLine *time.Time
	checkInterval   time.Duration
	m               sync.RWMutex
}

func (c *Channel) Name() string {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.name
}

func (c *Channel) Capacity() int {
	c.m.RLock()
	defer c.m.RUnlock()
	return len(c.ch)
}

func (c *Channel) SetSendDeadLine(t *time.Time) {
	c.m.Lock()
	defer c.m.Unlock()
	c.sendDeadLine = t
}

func (c *Channel) SetReceiveDeadLine(t *time.Time) {
	c.m.Lock()
	defer c.m.Unlock()
	c.receiveDeadLine = t
}

func (c *Channel) SetCheckInterval(duration time.Duration) {
	c.m.Lock()
	defer c.m.Unlock()
	c.checkInterval = duration
}

func (c *Channel) IsClosed() bool {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.isClosed
}

func (c *Channel) String() string {
	c.m.RLock()
	defer c.m.RUnlock()
	return fmt.Sprintf("(*Channel Cargo){ Name: %s, Capacity: %d, IsClosed: %t }",
		c.name, len(c.ch), c.isClosed)
}

func (c *Channel) setName(name string) {
	c.m.Lock()
	defer c.m.Unlock()
	c.name = name
}

func (c *Channel) close() {
	c.m.Lock()
	defer c.m.Unlock()
	c.isClosed = true
}

func (c *Channel) send(cargo Cargo) error {
	return nil
}

func (c *Channel) recv() (Cargo, error) {
	return nil, nil
}

func (c *Channel) trySend(cargo Cargo) bool {
	c.m.RLock()
	defer c.m.RUnlock()
	if c.isClosed {
		return false
	}

	var ok bool
	select {
	case c.ch <- cargo:
		ok = true
	default:
	}
	return ok
}

func (c *Channel) tryRecv() (Cargo, bool) {
	c.m.RLock()
	defer c.m.RUnlock()
	if c.isClosed {
		return nil, false
	}

	var ok bool
	var cargo Cargo
	select {
	case cargo = <-c.ch:
		ok = true
	default:
	}
	return cargo, ok
}

func (c *Channel) cancel() {
	return
}
