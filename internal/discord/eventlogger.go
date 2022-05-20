package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/meschbach/minecraft-overseer/internal/mc/events"
)

type EventLogger struct {
	//IPC layer to Discord
	client     *discordgo.Session
	eventQueue chan events.LogEntry
}

func (e *EventLogger) pumpMessagesOut(s *discordgo.Session, channelID string, filter func(entry events.LogEntry) bool) {
	fmt.Printf("Discord Event Pump activated channel %s\n", channelID)
	for event := range e.eventQueue {
		if !filter(event) {
			continue
		}
		switch event.(type) {
		case *events.UnknownLogEntry:
		default:
			_, err := s.ChannelMessageSend(channelID, event.String())
			if err != nil {
				fmt.Printf("WWW Failed to send Discord message because %q", err.Error())
			}
		}
	}
	fmt.Println("Stopped discord event pump")
}

func (e *EventLogger) Ingest(dispatcher *events.LogDispatcher) {
	dispatcher.Add("DiscordEventLogger", e.eventQueue)
}
