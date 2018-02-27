package channel

import (
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func TestBlockingQueue_PushCap0Exit(t *testing.T) {
	bq := newBlockingQueue(0)
	exit := new(exit)
	go func() {
		time.Sleep(time.Second * 1)
		exit.Open()
	}()
	ok := bq.Push("cargo", exit)
	if ok {
		t.Error("unexpected Push result")
	}
}

func TestBlockingQueue_PushCap0Succ(t *testing.T) {
	bq := newBlockingQueue(0)
	srcCargo := "cargo"
	go func() {
		time.Sleep(time.Second * 1)
		dstCargo, ok := bq.Pop(nil)
		if !ok {
			t.Error("Pop failed")
		}
		if srcCargo != dstCargo {
			t.Errorf("unexpected cargo: %s", dstCargo)
		}
	}()
	ok := bq.Push(srcCargo, nil)
	if !ok {
		t.Error("Push failed")
	}
}

func TestBlockingQueue_PushSucc(t *testing.T) {
	bq := newBlockingQueue(1)
	ok := bq.Push("cargo", nil)
	if !ok {
		t.Error("Push failed")
	}
}

func TestBlockingQueue_PushExit(t *testing.T) {
	bq := newBlockingQueue(1)
	ok := bq.Push("cargo", nil)
	if !ok {
		t.Error("Push failed")
	}
	exit := new(exit)
	go func() {
		time.Sleep(time.Second * 1)
		exit.Open()
	}()
	ok = bq.Push("cargo2", exit)
	if ok {
		t.Error("unexpeceted Push result")
	}
}

func TestBlockingQueue_PopCap0Exit(t *testing.T) {
	bq := newBlockingQueue(0)
	exit := new(exit)
	go func() {
		time.Sleep(time.Second * 1)
		exit.Open()
	}()
	_, ok := bq.Pop(exit)
	if ok {
		t.Error("unexpected Push result")
	}
}

func TestBlockingQueue_PopCap0Succ(t *testing.T) {
	bq := newBlockingQueue(0)
	srcCargo := "cargo"
	go func() {
		time.Sleep(time.Second * 1)
		ok := bq.Push(srcCargo, nil)
		if !ok {
			t.Error("Push failed")
		}
	}()
	dstCargo, ok := bq.Pop(nil)
	if !ok {
		t.Error("Pop failed")
	}
	if srcCargo != dstCargo {
		t.Errorf("unexpected cargo: %s", dstCargo)
	}
}

func TestBlockingQueue_PopSucc(t *testing.T) {
	bq := newBlockingQueue(1)
	srcCargo := "cargo"
	ok := bq.Push(srcCargo, nil)
	if !ok {
		t.Error("Push failed")
	}
	dstCargo, ok := bq.Pop(nil)
	if !ok {
		t.Error("Pop failed")
	}
	if srcCargo != dstCargo {
		t.Errorf("unexpected cargo: %s", dstCargo)
	}
}

func TestBlockingQueue_PopExit(t *testing.T) {
	bq := newBlockingQueue(1)
	exit := new(exit)
	go func() {
		time.Sleep(time.Second * 1)
		exit.Open()
	}()
	_, ok := bq.Pop(exit)
	if ok {
		t.Error("unexpeceted Push result")
	}
}

func TestBlockingQueue_Expand(t *testing.T) {
	bq := newBlockingQueue(0)
	go func() {
		time.Sleep(time.Second * 1)
		if !bq.Expand(4) {
			t.Error("Expand failed")
		}
	}()
	ok := bq.Push("cargo", nil)
	if !ok {
		t.Error("Push failed")
	}
}

func TestBlockingQueue_TryPushCap0Failed(t *testing.T) {
	bq := newBlockingQueue(0)
	if bq.TryPush("cargo") {
		t.Error("unexpected TryPush result")
	}
}

func TestBlockingQueue_TryPushCap0Succ(t *testing.T) {
	bq := newBlockingQueue(0)
	go func() {
		if !bq.TryPush("cargo") {
			t.Error("unexpected TryPush result")
			return
		}
	}()
	if _, ok := bq.Pop(nil); !ok {
		t.Error("Pop failed")
	}
	time.Sleep(time.Second * 1)
}

func TestBlockingQueue_TryPushSucc(t *testing.T) {
	bq := newBlockingQueue(1)
	if !bq.TryPush("cargo") {
		t.Error("unexpected TryPush result")
	}
}

func TestBlockingQueue_TryPopFailed(t *testing.T) {
	bq := newBlockingQueue(0)
	if _, ok := bq.TryPop(); ok {
		t.Error("unexpected TryPop result")
	}
}

func TestBlockingQueue_TryPopSucc(t *testing.T) {
	bq := newBlockingQueue(1)
	if !bq.TryPush("cargo") {
		t.Error("unexpected TryPush result")
		return
	}
	if _, ok := bq.TryPop(); !ok {
		t.Error("unexpected TryPop result")
	}
}

func TestBlockingQueue_Close(t *testing.T) {
	bq := newBlockingQueue(0)
	bq.Close()
	if bq.Push("", nil) {
		t.Error("unexpected Push result")
		return
	}
	if _, ok := bq.Pop(nil); ok {
		t.Error("unexpected Pop result")
	}
}

func TestBlockingQueue(t *testing.T) {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(4))
	N := 200
	if testing.Short() {
		N = 20
	}
	for chanCap := 0; chanCap < N; chanCap++ {
		{
			// Ensure that receive from empty chan blocks.
			q := newBlockingQueue(chanCap)
			recv := false
			go func() {
				q.Pop(nil)
				recv = true
			}()
			time.Sleep(time.Millisecond)
			if recv {
				t.Fatalf("chan[%d]: receive from empty chan", chanCap)
			}
			// Ensure that non-blocking receive does not block.
			if _, ok := q.TryPop(); ok {
				t.Fatalf("chan[%d]: receive from empty chan", chanCap)
			}
			q.Push(0, nil)
		}

		{
			// Ensure that send to full chan blocks.
			q := newBlockingQueue(chanCap)
			for i := 0; i < chanCap; i++ {
				q.Push(i, nil)
			}
			sent := uint32(0)
			go func() {
				q.Push(0, nil)
				atomic.StoreUint32(&sent, 1)
			}()
			time.Sleep(time.Millisecond)
			if atomic.LoadUint32(&sent) != 0 {
				t.Fatalf("chan[%d]: send to full chan", chanCap)
			}
			// Ensure that non-blocking send does not block.
			if q.TryPush(0) {
				t.Fatalf("chan[%d]: send to full chan", chanCap)
			}
			q.Pop(nil)
		}

		{
			// Ensure that we receive 0 from closed chan.
			q := newBlockingQueue(chanCap)
			for i := 0; i < chanCap; i++ {
				q.Push(i, nil)
			}
			q.Close()
			for i := 0; i < chanCap; i++ {
				v, _ := q.Pop(nil)
				if v != i {
					t.Fatalf("chan[%d]: received %v, expected %v", chanCap, v, i)
				}
			}
			if v, _ := q.Pop(nil); v != nil {
				t.Fatalf("chan[%d]: received %v, expected %v", chanCap, v, 0)
			}
		}

		{
			// Ensure that close unblocks receive.
			q := newBlockingQueue(chanCap)
			done := make(chan bool)
			go func() {
				v, ok := q.Pop(nil)
				done <- v == nil && ok == false
			}()
			time.Sleep(time.Millisecond)
			q.Close()
			if !<-done {
				t.Fatalf("chan[%d]: received non zero from closed chan", chanCap)
			}
		}

		{
			// Send 100 integers,
			// ensure that we receive them non-corrupted in FIFO order.
			q := newBlockingQueue(chanCap)
			go func() {
				for i := 0; i < 100; i++ {
					q.Push(i, nil)
				}
			}()
			for i := 0; i < 100; i++ {
				v, _ := q.Pop(nil)
				if v != i {
					t.Fatalf("chan[%d]: received %v, expected %v", chanCap, v, i)
				}
			}

			// Send 1000 integers in 4 goroutines,
			// ensure that we receive what we send.
			const P = 4
			const L = 1000
			for p := 0; p < P; p++ {
				go func() {
					for i := 0; i < L; i++ {
						q.Push(i, nil)
					}
				}()
			}
			done := make(chan map[int]int)
			for p := 0; p < P; p++ {
				go func() {
					recv := make(map[int]int)
					for i := 0; i < L; i++ {
						vv, _ := q.Pop(nil)
						v := vv.(int)
						recv[v] = recv[v] + 1
					}
					done <- recv
				}()
			}
			recv := make(map[int]int)
			for p := 0; p < P; p++ {
				for k, v := range <-done {
					recv[k] = recv[k] + v
				}
			}
			if len(recv) != L {
				t.Fatalf("chan[%d]: received %v values, expected %v", chanCap, len(recv), L)
			}
			for _, v := range recv {
				if v != P {
					t.Fatalf("chan[%d]: received %v values, expected %v", chanCap, v, P)
				}
			}
		}

		{
			// Test len/cap.
			q := newBlockingQueue(chanCap)
			if q.Len() != 0 || q.Cap() != chanCap {
				t.Fatalf("chan[%d]: bad len/cap, expect %v/%v, got %v/%v", chanCap, 0, chanCap, q.Len(), q.Cap())
			}
			for i := 0; i < chanCap; i++ {
				q.Push(i, nil)
			}
			if q.Len() != chanCap || q.Cap() != chanCap {
				t.Fatalf("chan[%d]: bad len/cap, expect %v/%v, got %v/%v", chanCap, chanCap, chanCap, q.Len(), q.Cap())
			}
		}

	}
}