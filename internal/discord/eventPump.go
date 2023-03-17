package discord

import (
	"context"
	"github.com/meschbach/minecraft-overseer/internal/mc/events"
)

type eventPump struct {
	dispatcher *events.LogDispatcher
	filter     func(entry events.LogEntry) bool
	sink       chan<- string
}

func (e *eventPump) Serve(ctx context.Context) error {
	gameEvents := make(chan events.LogEntry, 128)
	done := e.dispatcher.Add("DiscordEventLogger", gameEvents)
	defer done()

	for {
		select {
		case gameEvent := <-gameEvents:
			if e.filter(gameEvent) {
				switch gameEvent.(type) {
				case *events.UnknownLogEntry:
				default:
					select {
					case e.sink <- gameEvent.String():
					case <-ctx.Done():
						return ctx.Err()
					}
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
