package events

import "fmt"

type UserJoinedEntry struct {
	User string
}

func (s *UserJoinedEntry) AsString() string {
	return fmt.Sprintf("User joined %s", s.User)
}

type UserLeftEvent struct {
	User string
}

func (s *UserLeftEvent) AsString() string {
	return fmt.Sprintf("User left %s", s.User)
}
