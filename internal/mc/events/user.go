package events

import "fmt"

type UserJoinedEntry struct {
	User string
}

func (s *UserJoinedEntry) String() string {
	return fmt.Sprintf("User joined %s", s.User)
}

func (s *UserJoinedEntry) IsOperations() bool {
	return false
}

type UserLeftEvent struct {
	User string
}

func (s *UserLeftEvent) String() string {
	return fmt.Sprintf("User left %s", s.User)
}

func (s *UserLeftEvent) IsOperations() bool {
	return false
}

type UserSaidEvent struct {
	Speaker string
	Message string
}

func (s *UserSaidEvent) String() string {
	return fmt.Sprintf("<%s> %s", s.Speaker, s.Message)
}

func (s *UserSaidEvent) IsOperations() bool {
	return false
}
