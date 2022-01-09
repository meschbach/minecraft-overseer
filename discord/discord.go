package discord

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/meschbach/minecraft-overseer/internal/mc/events"
)

type DiscordOutput struct {
	client *discordgo.Session
	//internal reactor stuff
	ready bool
}

func NewDiscordClient(ctx context.Context, token string) (*DiscordOutput, error) {
	client, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	subsystem := &DiscordOutput{
		client: client,
		ready:  false,
	}
	client.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		fmt.Println("Discord client ready")
		subsystem.ready = true
		for _, guild := range s.State.Guilds {
			fmt.Printf("\tGuild %q\n", guild.Name)
			channels, _ := s.GuildChannels(guild.ID)
			for _, c := range channels {
				// Check if channel is a guild text channel and not a voice or DM channel
				if c.Type != discordgo.ChannelTypeGuildText {
					continue
				}
				fmt.Printf("\t\tChannel %q\n", c.Name)

				//if c.Name == "minecraft-talk-and-chat" {
				//	s.ChannelMessageSend(
				//		c.ID,
				//		fmt.Sprintf("This is a test.  This is only a test.  Otherwise this would be giving info about the minecrat server."),
				//	)
				//}
			}
		}
	})
	if err := client.Open(); err != nil {
		return nil, err
	}
	return subsystem, nil
}

func (d *DiscordOutput) OnEvent(e events.LogEntry) {

}
