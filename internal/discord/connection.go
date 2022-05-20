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
	//TODO: Viewing Guild names require additional permissions
	if len(s.State.Guilds) != 1 {
		fmt.Printf("Warning: Has %d guilds.  Assuming each guild in list is %q\n", len(s.State.Guilds), c.guildName)
	}

	for _, guild := range s.State.Guilds {
		//TODO: Filters on guilds should probably be a tunable

		channels, err := s.GuildChannels(guild.ID)
		if err != nil {
			fmt.Printf("Error when attempting to list Guild %q channels: %s", guild.ID, err.Error())
			continue
		}
		for _, channel := range channels {
			// Check if channel is a guild text channel and not a voice or DM channel
			if channel.Type != discordgo.ChannelTypeGuildText {
				continue
			}

			if channel.Name == c.targetChannel {
				found = true
				c.onChannelFound(s, channel.ID)
			}
		}
	}
	if !found {
		fmt.Printf("Warning: Could not find guild %q with channel %q\n", c.guildName, c.targetChannel)
	}
}

func (c *connection) onChannelFound(s *discordgo.Session, channelID string) {
	go c.subsystem.pumpMessagesOut(s, channelID)
}
