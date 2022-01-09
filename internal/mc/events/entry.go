package events

type LogEntry interface {
	//TODO: rename to String so it is comptable with Golang
	AsString() string
}
