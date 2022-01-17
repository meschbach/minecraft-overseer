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

func NewLogger(token string, targetChannel string) (*EventLogger, error) {
	client, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	subsystem := &EventLogger{
		client:     client,
		eventQueue: make(chan events.LogEntry, 16),
	}
	client.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		fmt.Println("Discord client ready")
		for _, guild := range s.State.Guilds {
			fmt.Printf("\tGuild %q\n", guild.Name)
			channels, _ := s.GuildChannels(guild.ID)
			for _, c := range channels {
				// Check if channel is a guild text channel and not a voice or DM channel
				if c.Type != discordgo.ChannelTypeGuildText {
					continue
				}
				fmt.Printf("\t\tChannel %q\n", c.Name)

				if c.Name == targetChannel {
					s.ChannelMessageSend(
						c.ID,
						fmt.Sprintf("Overseer connected."),
					)
					go func() {
						for event := range subsystem.eventQueue {
							switch event.(type) {
							case *events.UnknownLogEntry:
							default:
								s.ChannelMessageSend(c.ID, event.String())
							}
						}
					}()
				}
			}
		}
	})
	if err := client.Open(); err != nil {
		return nil, err
	}
	return subsystem, nil
}

func (e *EventLogger) Ingest(dispatcher *events.LogDispatcher) {
	dispatcher.Add(e.eventQueue)
}
