package event

type Observer interface {
	Notify(Event)
}
