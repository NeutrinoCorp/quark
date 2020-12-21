package pkg

// Handler handles any Event coming from an specific topic(s)
type Handler interface {
	ServeEvent(EventWriter, *Event)
}

// HandlerFunc handles any Event coming from an specific topic(s)
type HandlerFunc func(EventWriter, *Event)
