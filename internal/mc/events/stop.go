package events

import "fmt"

type StoppingEntry struct {
}

func (s *StoppingEntry) String() string {
	return fmt.Sprintf("Stopping server")
}
