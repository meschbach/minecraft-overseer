package game

type parsedTranslator struct {
}

func (*parsedTranslator) translate(input string) LogEntry {
	return ParseLogEntry(input)
}
