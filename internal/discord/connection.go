package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type connection struct {
	guildName     string
	targetChannel string
	subsystem     *EventLogger
}

func (c *connection) onReadyEvent(s *discordgo.Session, event *discordgo.Ready) {
	found := false
	fmt.Println("Discord client ready")
	for _, guild := range s.State.Guilds {
		//fmt.Printf("\tGuild %q\n", guild.Name)
		//TODO: this should be configurable based on grants
		//if guild.Name != guildName {
		//	continue
		//}

		channels, _ := s.GuildChannels(guild.ID)
		for _, channel := range channels {
			// Check if channel is a guild text channel and not a voice or DM channel
			if channel.Type != discordgo.ChannelTypeGuildText {
				continue
			}
			fmt.Printf("\t\tChannel %q\n", channel.Name)

			if channel.Name == c.targetChannel {
				found = true
				s.ChannelMessageSend(
					channel.ID,
					fmt.Sprintf("Overseer connected."),
				)
				go c.subsystem.pumpMessagesOut(s, channel.ID)
			}
		}
	}
	if !found {
		fmt.Printf("Warning: Could not find guild %q with channel %q\n", c.guildName, c.targetChannel)
	}
}
