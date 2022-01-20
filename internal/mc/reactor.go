package mc

import (
	"fmt"
	"github.com/meschbach/minecraft-overseer/internal/mc/events"
)

type ConsoleOperationFunc = func(buffer LogEntryBuffer, stdin chan<- string)

type LogEntryBuffer interface {
	WaitFor(func(entry events.LogEntry) bool)
}

type ConsoleOperation interface {
	Apply(buffer LogEntryBuffer, stdin chan<- string)
}

type ConsoleReactor struct {
	Logs              *events.LogDispatcher
	PendingOperations chan ConsoleOperation
}

func (c *ConsoleReactor) consumeOperations(stdin chan<- string) {
	for op := range c.PendingOperations {
		c.consumeOperation(op, stdin)
	}
}

func (c *ConsoleReactor) consumeOperation(op ConsoleOperation, stdin chan<- string) {
	consumer := make(chan events.LogEntry, 32)
	defer close(consumer)

	done := c.Logs.Add("ConsoleReactorOperation", consumer)
	defer done()

	buffer := consoleEventBuffer{input: consumer}
	op.Apply(&buffer, stdin)
	buffer.drain()
}

type consoleEventBuffer struct {
	input chan events.LogEntry
}

func (c *consoleEventBuffer) WaitFor(filter func(entry events.LogEntry) bool) {
	fmt.Printf("Waiting on messages\n")
	for e := range c.input {
		if filter(e) {
			fmt.Printf("Found\n")
			return
		}
	}
}

func (c *consoleEventBuffer) drain() {
	for {
		select {
		case <-c.input:
		default:
			return
		}
	}
}

func StartReactor(stdout <-chan string, stdin chan<- string) *ConsoleReactor {
	gameChannel := make(chan events.LogEntry, 16)
	operations := make(chan ConsoleOperation, 16)
	dispatcher := events.NewLogDispatcher()
	go dispatcher.Consume(gameChannel)
	reactor := &ConsoleReactor{
		Logs:              dispatcher,
		PendingOperations: operations,
	}
	go reactor.consumeOperations(stdin)

	//Parser
	go func() {
		for line := range stdout {
			entry := events.ParseLogEntry(line)
			gameChannel <- entry
		}
	}()

	return reactor
}

type WaitForStart struct{}

func (w *WaitForStart) Apply(buffer LogEntryBuffer, stdin chan<- string) {
	buffer.WaitFor(func(entry events.LogEntry) bool {
		switch entry.(type) {
		case *events.StartedEntry:
			return true
		default:
			return false
		}
	})
	fmt.Printf("Game started.\n")
}

type EnsureUserOperators struct {
	Users []string
}

func (w *EnsureUserOperators) Apply(buffer LogEntryBuffer, stdin chan<- string) {
	for _, operator := range w.Users {
		stdin <- fmt.Sprintf("whitelist add %s", operator)
		stdin <- fmt.Sprintf("op %s", operator)
	}
}

type EnsureWhitelistAdd struct {
	Users []string
}

func (w *EnsureWhitelistAdd) Apply(buffer LogEntryBuffer, stdin chan<- string) {
	for _, operator := range w.Users {
		stdin <- fmt.Sprintf("whitelist add %s", operator)
	}
}
