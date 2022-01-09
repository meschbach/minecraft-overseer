package game

import "github.com/meschbach/minecraft-overseer/internal/mc/events"

type parsedTranslator struct {
}

func (*parsedTranslator) translate(input string) events.LogEntry {
	return events.ParseLogEntry(input)
}
