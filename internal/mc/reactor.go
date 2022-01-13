package mc

import (
	"fmt"
	"github.com/meschbach/minecraft-overseer/internal/mc/events"
)

type ConsoleOperationFunc = func(logEntry <-chan events.LogEntry, stdin chan<- string)

type ConsoleOperation interface {
	Apply(logEntry <-chan events.LogEntry, stdin chan<- string)
}

type ConsoleReactor struct {
	Logs              *events.LogDispatcher
	PendingOperations chan ConsoleOperation
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

	//Controlling actor
	go func() {
		consumer := make(chan events.LogEntry)
		defer close(consumer)

		done := dispatcher.Add(consumer)
		defer done()

		for op := range operations {
			fmt.Printf("Next operation %#v\n", op)
			op.Apply(consumer, stdin)
		}
	}()

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

func (w *WaitForStart) Apply(logEntry <-chan events.LogEntry, stdin chan<- string) {
	for e := range logEntry {
		switch e.(type) {
		case *events.StartedEntry:
			fmt.Println("[Event] Game started")
			return
		}
	}
}

type EnsureUserOperators struct {
	Users []string
}

func (w *EnsureUserOperators) Apply(logEntry <-chan events.LogEntry, stdin chan<- string) {
	fmt.Printf("Ensuring users are oeprators")
	for _, operator := range w.Users {
		stdin <- fmt.Sprintf("whitelist add %s", operator)
		stdin <- fmt.Sprintf("op %s", operator)
	}
}

type EnsureWhitelistAdd struct {
	Users []string
}

func (w *EnsureWhitelistAdd) Apply(logEntry <-chan events.LogEntry, stdin chan<- string) {
	fmt.Printf("Ensuring users are on the whitelist")
	for _, operator := range w.Users {
		stdin <- fmt.Sprintf("whitelist add %s", operator)
	}
}
