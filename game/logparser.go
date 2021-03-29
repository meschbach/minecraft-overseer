package game

import (
	"fmt"
	"strings"
)

type LogEntry interface {
	AsString() string
}

type UnknownLogEntry struct {
	Line string
}

func (u *UnknownLogEntry) AsString() string {
	return u.Line
}

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

type StoppingEntry struct {
}

func (s *StoppingEntry) AsString() string {
	return fmt.Sprintf("Stopping server")
}

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

func parseLogEntry(rawEntry string) LogEntry {
	var entry string
	if rawEntry[len(rawEntry)-1:] == "\n" {
		entry = rawEntry[:len(rawEntry)-1]
	} else {
		entry = rawEntry
	}

	firstSpace := strings.Index(entry, " ")
	if firstSpace == -1 {
		return &UnknownLogEntry{Line: entry}
	}
	afterDate := entry[firstSpace+1:]
	colonIndex := strings.Index(afterDate, ": ")
	if colonIndex == -1 {
		return &UnknownLogEntry{Line: entry}
	}

	message := afterDate[colonIndex+2:]
	if strings.HasPrefix(message,"Starting minecraft server version ") {
		return &StartingEntry{
			message[len("Starting minecraft server version "):],
		}
	}
	if strings.HasPrefix(message,"Done (") && strings.HasSuffix(message,"! For help, type \"help\" or \"?\"") {
		prefix := "Done ("
		afterPrefix := message[len(prefix):]
		timingEndIndex := strings.Index(afterPrefix,")!")
		time := afterPrefix[:timingEndIndex]
		return &StartedEntry {
			TimeTaken: time,
		}
	}
	if message == "Stopping the server" {
		return &StoppingEntry {
		}
	}
	if strings.HasSuffix(message, " joined the game") {
		endUserName := strings.Index(message, " ")
		name := message[:endUserName]
		return &UserJoinedEntry{
			User: name,
		}
	}
	if strings.HasSuffix(message, "left the game") {
		endUserName := strings.Index(message, " ")
		name := message[:endUserName]
		return &UserLeftEvent{
			User: name,
		}
	}
	return &UnknownLogEntry{Line: message}
}
