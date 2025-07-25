package internal

type Observer interface {
	Update(Message)
	GetID() string
}
