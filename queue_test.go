package channel

import (
	"testing"
	"time"
)

func TestBlockingQueue_PushCap0Exit(t *testing.T) {
	bq := newBlockingQueue(0)
	exit := new(exit)
	go func() {
		time.Sleep(time.Second * 5)
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
		time.Sleep(time.Second * 5)
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
		time.Sleep(time.Second * 5)
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
		time.Sleep(time.Second * 5)
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
		time.Sleep(time.Second * 5)
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
		time.Sleep(time.Second * 5)
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
		time.Sleep(time.Second * 3)
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
		//time.Sleep(time.Second * 3)
		if !bq.TryPush("cargo") {
			t.Error("unexpected TryPush result")
			return
		}
		t.Log("TryPush Succ")
	}()
	if _, ok := bq.Pop(nil); !ok {
		t.Error("Pop failed")
	}
	time.Sleep(time.Second * 3)
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
