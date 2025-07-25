package internal

type Subject interface {
	Register(observer Observer)
	Deregister(observer Observer)
	NotifyAll(message Message)
}

type Message struct {
	Topic string
	Data any
}