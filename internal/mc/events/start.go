package events

import "fmt"

type StartingEntry struct {
	Version string
}

func (s *StartingEntry) String() string {
	return fmt.Sprintf("Starting %s", s.Version)
}

func (s *StartingEntry) IsOperations() bool {
	return true
}

type StartedEntry struct {
	TimeTaken string
}

func (s *StartedEntry) String() string {
	return fmt.Sprintf("Started")
}

func (s *StartedEntry) IsOperations() bool {
	return true
}
