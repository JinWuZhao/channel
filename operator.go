package channel

func Make(capacity int) *Channel {
	return &Channel{
		ch: make(chan Cargo, capacity),
		checkInterval: defaultCheckInterval,
	}
}

func Close(ch *Channel) {

}

func Send(ch *Channel, cargo Cargo) error {
	return nil
}

func Recv(ch *Channel) (cargo Cargo, err error) {
	return
}

func MultiSend(ch *Channel, cargos []Cargo) error {
	return nil
}

func MultiRecv(ch *Channel) (cargos []Cargo, err error) {
	return
}

func TrySend(ch *Channel, cargo Cargo) bool {
	return false
}

func TryRecv(ch *Channel) (cargo Cargo, ok bool) {
	return
}

func Expand(ch *Channel, capacity int) *Channel {
	return nil
}

func Merge(chs ...*Channel) *Channel {
	return nil
}

func ForRecv(ch *Channel, f func(cargo Cargo) error) error {
	return nil
}

func Map(ch *Channel, f func(cargo Cargo) Cargo) *Channel {
	return nil
}

func Pipe(src *Channel, dst *Channel) *Channel {
	return nil
}
