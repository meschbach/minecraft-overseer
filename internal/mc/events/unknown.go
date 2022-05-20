package events

type UnknownLogEntry struct {
	Line string
}

func (u *UnknownLogEntry) String() string {
	return u.Line
}

func (u *UnknownLogEntry) IsOperations() bool {
	return true
}
