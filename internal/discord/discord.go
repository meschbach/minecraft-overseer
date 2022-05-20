package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/meschbach/minecraft-overseer/internal/mc/events"
)

func NewLogger(token string, guildName string, targetChannel string) (*EventLogger, error) {
	client, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	subsystem := &EventLogger{
		client:     client,
		eventQueue: make(chan events.LogEntry, 16),
	}
	connectionHandler := &connection{
		guildName:     guildName,
		targetChannel: targetChannel,
		subsystem:     subsystem,
	}
	client.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		connectionHandler.onReadyEvent(s, event)
	})
	if err := client.Open(); err != nil {
		return nil, err
	}
	return subsystem, nil
}
