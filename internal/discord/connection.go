package discord

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/meschbach/minecraft-overseer/internal/mc/events"
	"github.com/thejerf/suture/v4"
)

type connection struct {
	guildName string
	//userChannel receives messages of interest to general users
	userChannel string
	//opChannel receives operational messages like startup, connections, etc
	opChannel string
	subsystem *EventLogger
	//userCommands will interpret a message from a client and respond as appropriate
	userCommands     chan<- discordgo.Message
	userReplies      <-chan string
	ParentSupervisor *suture.Supervisor
}

func (c *connection) onReadyEvent(s *discordgo.Session, event *discordgo.Ready) {
	found := 0
	fmt.Printf("Discord client ready, searching for {ops: %q, user: %q}\n", c.opChannel, c.userChannel)
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
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if m.ChannelID == channelID {
			c.userCommands <- *m.Message
		}
	})

	outgoingSupervisor := suture.NewSimple("outgoing")
	outgoingSupervisor.Add(&outgoingDispatcher{s, channelID, c.userReplies})
	c.ParentSupervisor.Add(outgoingSupervisor)
}

func (c *connection) onOpsChannelFound(s *discordgo.Session, channelID string) {
	go c.subsystem.pumpMessagesOut(s, channelID, func(entry events.LogEntry) bool {
		return entry.IsOperations()
	})
}

type outgoingDispatcher struct {
	s         *discordgo.Session
	channelID string
	queue     <-chan string
}

func (o *outgoingDispatcher) Serve(ctx context.Context) error {
	for {
		select {
		case msg := <-o.queue:
			_, err := o.s.ChannelMessageSend(o.channelID, msg, discordgo.WithContext(ctx))
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
