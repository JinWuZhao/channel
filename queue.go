package channel

import (
	"container/list"
	"sync"
	"time"
)

type exit struct {
	l       sync.RWMutex
	opening bool
	cond    *sync.Cond
}

func (u *exit) Open() {
	u.l.Lock()
	defer u.l.Unlock()
	if u.cond != nil && !u.opening {
		u.opening = true
		u.cond.Broadcast()
	}
}

func (u *exit) Opening() bool {
	u.l.RLock()
	defer u.l.RUnlock()
	return u.opening
}

func (u *exit) SetCond(cond *sync.Cond) {
	u.l.Lock()
	defer u.l.Unlock()
	u.cond = cond
}

type blockingQueue struct {
	closed   bool
	capacity int
	list     *list.List
	single   Cargo
	l        *sync.RWMutex
	seal     *sync.Cond
}

func newBlockingQueue(capacity int) *blockingQueue {
	if capacity < 0 {
		panic("capacity can not be negative")
	}
	bq := &blockingQueue{
		capacity: capacity,
		list:     list.New(),
		l:        new(sync.RWMutex),
	}
	bq.seal = sync.NewCond(bq.l)
	return bq
}

func (bq *blockingQueue) checkLen(pushing bool) bool {
	if pushing {
		if bq.capacity == 0 {
			if bq.single != nil {
				return false
			} else {
				return true
			}
		}
		if bq.list.Len() >= bq.capacity {
			return false
		} else {
			return true
		}
	} else {
		if bq.capacity == 0 {
			if bq.single == nil {
				return false
			} else {
				return true
			}
		}
		if bq.list.Len() <= 0 {
			return false
		} else {
			return true
		}
	}
}

func (bq *blockingQueue) check(pushing bool, exit *exit) bool {
	if pushing && bq.closed {
		return false
	}

check:
	if !bq.checkLen(pushing) {
		if !pushing && bq.closed {
			return false
		}
		if exit != nil {
			exit.SetCond(bq.seal)
		}
		bq.seal.Wait()
		if exit != nil {
			exit.SetCond(nil)
		}
		if pushing && bq.closed {
			return false
		}
		if exit != nil && exit.Opening() {
			return false
		}
		goto check
	}
	if exit != nil && exit.Opening() {
		return false
	}
	return true
}

func (bq *blockingQueue) unblock() {
	bq.seal.Broadcast()
}

func (bq *blockingQueue) Push(cargo Cargo, exit *exit) bool {
	bq.l.Lock()
	defer bq.l.Unlock()

	if bq.capacity == 0 {
		if !bq.check(true, exit) {
			return false
		}
		if bq.capacity == 0 {
			bq.single = cargo
		}
	}

	bq.unblock()

	if !bq.check(true, exit) {
		return false
	}

	if bq.capacity > 0 {
		bq.list.PushBack(cargo)
	}

	bq.unblock()
	return true
}

func (bq *blockingQueue) Pop(exit *exit) (Cargo, bool) {
	bq.l.Lock()
	defer bq.l.Unlock()

	bq.unblock()

	if !bq.check(false, exit) {
		return nil, false
	}

	var cargo Cargo

	if bq.capacity == 0 {
		cargo = bq.single
		bq.single = nil
	} else {
		front := bq.list.Front()
		cargo = front.Value.(Cargo)
		bq.list.Remove(front)
	}

	bq.unblock()
	return cargo, true
}

func (bq *blockingQueue) Expand(size int) bool {
	bq.l.Lock()
	defer bq.l.Unlock()

	if bq.closed {
		return false
	}

	if bq.capacity+size < 0 {
		return false
	}

	if bq.list.Len() > bq.capacity+size {
		return false
	}

	bq.capacity += size

	bq.unblock()
	return true
}

func (bq *blockingQueue) TryPush(cargo Cargo) bool {
	bq.l.Lock()
	defer bq.l.Unlock()

	if bq.closed {
		return false
	}

	if bq.capacity == 0 {
		if bq.single != nil {
			return false
		}
		bq.single = cargo
	}

	bq.l.Unlock()

	bq.unblock()
	time.Sleep(time.Millisecond)

	bq.l.Lock()

	if bq.closed {
		return false
	}

	if !bq.checkLen(true) {
		bq.single = nil
		return false
	}

	if bq.capacity > 0 {
		bq.list.PushBack(cargo)
	}

	bq.unblock()
	return true
}

func (bq *blockingQueue) TryPop() (Cargo, bool) {
	bq.l.Lock()
	defer bq.l.Unlock()

	bq.l.Unlock()

	bq.unblock()
	time.Sleep(time.Millisecond)

	bq.l.Lock()

	if !bq.checkLen(false) {
		return nil, false
	}

	var cargo Cargo

	if bq.capacity == 0 {
		cargo = bq.single
		bq.single = nil
	} else {
		front := bq.list.Front()
		cargo = front.Value.(Cargo)
		bq.list.Remove(front)
	}

	bq.unblock()
	return cargo, true
}

func (bq *blockingQueue) Close() {
	bq.l.Lock()
	defer bq.l.Unlock()

	bq.closed = true

	bq.unblock()
}

func (bq *blockingQueue) Closed() bool {
	bq.l.RLock()
	defer bq.l.RUnlock()

	return bq.closed
}

func (bq *blockingQueue) Cap() int {
	bq.l.RLock()
	defer bq.l.RUnlock()

	return bq.capacity
}

func (bq *blockingQueue) Len() int {
	bq.l.RLock()
	defer bq.l.RUnlock()

	return bq.list.Len()
}
