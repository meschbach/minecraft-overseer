package ws

type OutputMessage interface {
	asWireMessage() []byte
}

type StringOutput struct {
	message string
}

func (s StringOutput) asWireMessage() []byte {
	return []byte(s.message)
}

type Message interface {
	apply(hub *Hub) OutputMessage
}

type StopMessage struct{}

func (m StopMessage) apply(hub *Hub) OutputMessage {
	hub.overseer.fsm <- Stop
	return &StringOutput{message: "ok"}
}

type StartMessage struct{}

func (m StartMessage) apply(hub *Hub) OutputMessage {
	hub.overseer.fsm <- Start
	return &StringOutput{message: "ok"}
}
