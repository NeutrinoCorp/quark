package pkg

// Handler handles any Event coming from an specific topic(s)
//
// Returned boolean works as Acknowledgement mechanism (Ack or NAck).
//
// - If true, Quark will send success to the actual message broker/queue system.
//
// - If false, Quark will not send any success response to the actual message broker/queue system.
//
//	Ack mechanism is available in specific providers
type Handler interface {
	// ServeEvent handles any Event coming from an specific topic(s)
	//
	//
	// Returned boolean works as Acknowledgement mechanism (Ack or NAck).
	//
	// - If true, Quark will send success to the actual message broker/queue system.
	//
	// - If false, Quark will not send any success response to the actual message broker/queue system.
	//
	//	Ack mechanism is available in specific providers
	ServeEvent(EventWriter, *Event) bool
}

// HandlerFunc handles any Event coming from an specific topic(s).
//
// Returned boolean works as Acknowledgement mechanism (Ack or NAck).
//
// - If true, Quark will send success to the actual message broker/queue system.
//
// - If false, Quark will not send any success response to the actual message broker/queue system.
//
//	Ack mechanism is available in specific providers
type HandlerFunc func(EventWriter, *Event) bool
