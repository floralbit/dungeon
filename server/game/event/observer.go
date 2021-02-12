package event

type Observer interface {
	Notify(Event)
}

// Observers ...
var Observers = []Observer{}

func NotifyObservers(e Event) {
	for _, o := range Observers {
		o.Notify(e)
	}
}
