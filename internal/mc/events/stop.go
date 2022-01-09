package events

import "fmt"

type StoppingEntry struct {
}

func (s *StoppingEntry) AsString() string {
	return fmt.Sprintf("Stopping server")
}
