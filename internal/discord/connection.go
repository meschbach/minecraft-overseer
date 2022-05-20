package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/meschbach/minecraft-overseer/internal/mc/events"
)

type connection struct {
	guildName string
	//userChannel receives messages of interest to general users
	userChannel string
	//opChannel receives operational messages like startup, connections, etc
	opChannel string
	subsystem *EventLogger
}

func (c *connection) onReadyEvent(s *discordgo.Session, event *discordgo.Ready) {
	found := 0
	fmt.Println("Discord client ready")
	//TODO: Viewing Guild names require additional permissions
	if len(s.State.Guilds) != 1 {
		fmt.Printf("Warning: Has %d guilds.  Assuming each guild in list is %q\n", len(s.State.Guilds), c.guildName)
	}

	var opsChannelID string
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

			if channel.Name == c.userChannel {
				found++
				c.onUserChannelFound(s, channel.ID)
			}
			if channel.Name == c.opChannel {
				found++
				opsChannelID = channel.ID
				c.onOpsChannelFound(s, opsChannelID)
			}
		}
	}
	if found < 2 {
		fmt.Printf("Warning: Could not find guild %q with channel %q or %q\n", c.guildName, c.userChannel, c.opChannel)
	} else {
		if _, err := s.ChannelMessageSend(opsChannelID, "Overseer <-> Discord connection established."); err != nil {
			fmt.Printf("Warning: Unable to send initial message because %s\n", err.Error())
		}
	}
}

func (c *connection) onUserChannelFound(s *discordgo.Session, channelID string) {
	go c.subsystem.pumpMessagesOut(s, channelID, func(entry events.LogEntry) bool {
		return !entry.IsOperations()
	})
}

func (c *connection) onOpsChannelFound(s *discordgo.Session, channelID string) {
	go c.subsystem.pumpMessagesOut(s, channelID, func(entry events.LogEntry) bool {
		return entry.IsOperations()
	})
}
