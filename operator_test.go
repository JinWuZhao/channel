package channel

import "testing"

func TestSendTimeout(t *testing.T) {
	ch := Make(0)
	err := Send(ch, "cargo")
	if err != ErrTimeout {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestSendClosed(t *testing.T) {
	ch := Make(0)
	Close(ch)
	err := Send(ch, "cargo")
	if err != ErrClosed {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestSendSucc(t *testing.T) {
	ch := Make(1)
	srcCargo := "cargo"
	err := Send(ch, srcCargo)
	if err != nil {
		t.Errorf("send failed: %s", err)
		return
	}
	dstCargo, err := Recv(ch)
	if err != nil {
		t.Errorf("recv failed: %s", err)
		return
	}
	if srcCargo != dstCargo {
		t.Errorf("unexpected cargo: %+v", dstCargo)
		return
	}
}

func TestRecvTimeout(t *testing.T) {
	ch := Make(0)
	_, err := Recv(ch)
	if err != ErrTimeout {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

func TestRecvClosed(t *testing.T) {
	ch := Make(0)
	Close(ch)
	_, err := Recv(ch)
	if err != ErrClosed {
		t.Errorf("unexpected error: %s", err)
		return
	}
}

func TestRecvSucc(t *testing.T) {
	ch := Make(1)
	srcCargo := "cargo"
	err := Send(ch, srcCargo)
	if err != nil {
		t.Errorf("send failed: %s", err)
		return
	}
	dstCargo, err := Recv(ch)
	if err != nil {
		t.Errorf("recv failed: %s", err)
		return
	}
	if srcCargo != dstCargo {
		t.Errorf("unexpected cargo: %+v", dstCargo)
		return
	}
}

func TestMultiSend(t *testing.T) {
	ch := Make(3)
	cargos := []Cargo{"cargo1", "cargo2", "cargo3"}
	err := MultiSend(ch, cargos)
	if err != nil {
		t.Errorf("MultiSend failed: %s", err)
	}
}

func TestMultiRecv(t *testing.T) {
	ch := Make(3)
	srcCargos := []Cargo{"cargo1", "cargo2", "cargo3"}
	err := MultiSend(ch, srcCargos)
	if err != nil {
		t.Errorf("MultiSend failed: %s", err)
		return
	}
	dstCargos, err := MultiRecv(ch)
	if err != nil {
		t.Errorf("MultiRecv failed: %s", err)
		return
	}
	for index, cargo := range dstCargos {
		if srcCargos[index] != cargo {
			t.Errorf("unexpected cargo: %+v", cargo)
			return
		}
	}
}

func TestTrySendFailed(t *testing.T) {
	ch := Make(0)
	ok := TrySend(ch, "cargo")
	if ok {
		t.Errorf("unexpected result: %+v", ok)
	}
}

func TestTrySendSucc(t *testing.T) {
	ch := Make(1)
	ok := TrySend(ch, "cargo")
	if !ok {
		t.Errorf("unexpected result: %+v", ok)
	}
}

func TestTryRecvFailed(t *testing.T) {
	ch := Make(0)
	_, ok := TryRecv(ch)
	if ok {
		t.Errorf("unexpected result: %+v", ok)
	}
}

func TestTryRecvSucc(t *testing.T) {
	ch := Make(1)
	srcCargo := "cargo"
	ok := TrySend(ch, srcCargo)
	if !ok {
		t.Error("unexpected TrySend failed")
		return
	}
	dstCargo, ok := TryRecv(ch)
	if !ok {
		t.Errorf("unexpected result: %+v", ok)
		return
	}
	if srcCargo != dstCargo {
		t.Errorf("unexpected cargo: %+v", dstCargo)
		return
	}
}

func TestExpand(t *testing.T) {
	ch := Make(0)
	ch = Expand(ch, 1)
	err := Send(ch, "cargo")
	if err != nil {
		t.Errorf("send failed: %s", err)
	}
}

func TestMerge(t *testing.T) {
	chs := [3]*Channel{Make(1), Make(1), Make(1)}
	mch := Merge(chs[:]...)
	srcCargos := [3]Cargo{"cargo1", "cargo2", "cargo3"}
	for index, cargo := range srcCargos {
		err := Send(chs[index], cargo)
		if err != nil {
			t.Errorf("send failed: %s", err)
			return
		}
	}
	for index := range [3]int{} {
		cargo, err := Recv(mch)
		if err != nil {
			t.Errorf("recv failed: %s", err)
			return
		}
		if cargo != srcCargos[index] {
			t.Errorf("unexpected cargo: %+v", cargo)
			return
		}
	}
}

func TestForRecv(t *testing.T) {
	ch := Make(3)
	srcCargos := []Cargo{"cargo1", "cargo2", "cargo3"}
	err := MultiSend(ch, srcCargos)
	if err != nil {
		t.Errorf("MultiSend failed: %s", err)
		return
	}
	index := 2
	dstCargos := make([]Cargo, 3)
	err = ForRecv(ch, func(cargo Cargo) error {
		dstCargos[index] = cargo
		index--
		return nil
	})
	for index, cargo := range dstCargos {
		if srcCargos[index] != cargo {
			t.Errorf("undexpected cargo: %+v", cargo)
			return
		}
	}
	if err != nil {
		t.Errorf("ForRecv failed: %s", err)
	}
}

func TestMap(t *testing.T) {
	ch := Make(3)
	srcCargos := []Cargo{"cargo1", "cargo2", "cargo3"}
	err := MultiSend(ch, srcCargos)
	if err != nil {
		t.Errorf("MultiSend failed: %s", err)
		return
	}
	mapF := func(cargo Cargo) Cargo {
		return cargo.(string) + "_mapped"
	}
	mapCh := Map(ch, mapF)
	dstCargos, err := MultiRecv(mapCh)
	if err != nil {
		t.Errorf("MultiRecv failed: %s", err)
		return
	}
	for index, cargo := range srcCargos {
		if dstCargos[index] != mapF(cargo) {
			t.Errorf("unexpected cargo: %+v", dstCargos[index])
			return
		}
	}
}

func TestPipe(t *testing.T) {
	srcCh := Make(1)
	dstCh := Make(1)
	pipedCh := Pipe(srcCh, dstCh)
	srcCargo := "cargo"
	err := Send(pipedCh, srcCargo)
	if err != nil {
		t.Errorf("Send failed: %s", err)
		return
	}
	dstCargo, err := Recv(pipedCh)
	if err != nil {
		t.Errorf("Recv failed: %s", err)
		return
	}
	if srcCargo != dstCargo {
		t.Errorf("undexpected cargo: %s", dstCargo)
		return
	}
}
