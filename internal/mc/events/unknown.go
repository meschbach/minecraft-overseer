package events

type UnknownLogEntry struct {
	Line string
}

func (u *UnknownLogEntry) AsString() string {
	return u.Line
}
