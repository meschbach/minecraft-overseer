package game

import "github.com/meschbach/minecraft-overseer/internal/mc/events"

type passthroughTranslator struct {
	prefix string
}

func (p *passthroughTranslator) translate(input string) events.LogEntry {
	return &events.UnknownLogEntry{Line: p.prefix + input}
}
