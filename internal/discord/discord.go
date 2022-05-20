package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/meschbach/minecraft-overseer/internal/mc/events"
)

type Config struct {
	Token         string
	GuildName     string
	TargetChannel string
}

func NewLogger(config Config) (*EventLogger, error) {
	client, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, err
	}

	subsystem := &EventLogger{
		client:     client,
		eventQueue: make(chan events.LogEntry, 16),
	}
	connectionHandler := &connection{
		guildName:     config.GuildName,
		targetChannel: config.TargetChannel,
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
