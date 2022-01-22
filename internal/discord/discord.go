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

func NewLogger(token string, guildName string, targetChannel string) (*EventLogger, error) {
	client, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	subsystem := &EventLogger{
		client:     client,
		eventQueue: make(chan events.LogEntry, 16),
	}
	client.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		found := false
		fmt.Println("Discord client ready")
		for _, guild := range s.State.Guilds {
			//fmt.Printf("\tGuild %q\n", guild.Name)
			//TODO: this should be configurable based on grants
			//if guild.Name != guildName {
			//	continue
			//}

			channels, _ := s.GuildChannels(guild.ID)
			for _, c := range channels {
				// Check if channel is a guild text channel and not a voice or DM channel
				if c.Type != discordgo.ChannelTypeGuildText {
					continue
				}
				fmt.Printf("\t\tChannel %q\n", c.Name)

				if c.Name == targetChannel {
					found = true
					s.ChannelMessageSend(
						c.ID,
						fmt.Sprintf("Overseer connected."),
					)
					go subsystem.pumpMessagesOut(s, c.ID)
				}
			}
		}
		if !found {
			fmt.Printf("Warning: Could not find guild %q with channel %q\n", guildName, targetChannel)
		}
	})
	if err := client.Open(); err != nil {
		return nil, err
	}
	return subsystem, nil
}

func (e *EventLogger) pumpMessagesOut(s *discordgo.Session, channelID string) {
	fmt.Printf("Discord Event Pump activated channel %s\n", channelID)
	for event := range e.eventQueue {
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
