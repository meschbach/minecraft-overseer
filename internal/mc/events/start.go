package events

import "fmt"

type StartingEntry struct {
	Version string
}

func (s *StartingEntry) AsString() string {
	return fmt.Sprintf("Starting %s", s.Version)
}

type StartedEntry struct {
	TimeTaken string
}

func (s *StartedEntry) AsString() string {
	return fmt.Sprintf("Started")
}
