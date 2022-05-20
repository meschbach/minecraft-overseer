package events

type LogEntry interface {
	//IsOperations identifies if the log entry is related to the operations of the Minecraft instance.
	IsOperations() bool
	String() string
}
