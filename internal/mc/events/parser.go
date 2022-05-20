package events

import "strings"

func ParseLogEntry(rawEntry string) LogEntry {
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
	if strings.HasPrefix(message, "Starting minecraft server version ") {
		return &StartingEntry{
			message[len("Starting minecraft server version "):],
		}
	}
	if strings.HasPrefix(message, "Done (") {
		prefix := "Done ("
		afterPrefix := message[len(prefix):]
		timingEndIndex := strings.Index(afterPrefix, ")!")
		time := afterPrefix[:timingEndIndex]
		return &StartedEntry{
			TimeTaken: time,
		}
	}
	if message == "Stopping the server" {
		return &StoppingEntry{}
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
	if strings.HasPrefix(message, "<") {
		endUserName := strings.Index(message, ">")
		name := message[1:endUserName]
		remainder := message[endUserName+2:]
		return &UserSaidEvent{
			Speaker: name,
			Message: remainder,
		}
	}

	if strings.HasSuffix(message, "was killed by Witch using magic") {
		return &GenericDeathMessage{Message: message}
	}
	if strings.HasSuffix(message, "was slain by Zombie") {
		return &GenericDeathMessage{Message: message}
	}
	return &UnknownLogEntry{Line: message}
}
