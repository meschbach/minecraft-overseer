package game

type passthroughTranslator struct {
	prefix string
}

func (p *passthroughTranslator) translate(input string) LogEntry {
	return &UnknownLogEntry{Line: p.prefix + input}
}
