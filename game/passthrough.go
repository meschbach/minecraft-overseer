package game

type passthroughTranslator struct {

}

func (*passthroughTranslator) translate(input string) LogEntry {
	return &UnknownLogEntry{Line: input}
}
