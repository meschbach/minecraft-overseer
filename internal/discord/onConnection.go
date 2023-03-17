package discord

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/thejerf/suture/v4"
)

type onConnection struct {
	session     *discordgo.Session
	sessionTree *suture.Supervisor

	guildName    string
	userChannel  string
	userMessages chan<- discordgo.Message
	//userReplies is the channel to send messages on
	userReplies  <-chan string
	opChannel    string
	instanceName string
}

func (o *onConnection) Serve(ctx context.Context) error {
	s := o.session
	found := 0
	//todo: logging
	fmt.Printf("Discord client ready, searching for {ops: %q, user: %q}\n", o.opChannel, o.userChannel)
	//TODO: Viewing Guild names require additional permissions
	if len(s.State.Guilds) != 1 {
		//todo: logging
		fmt.Printf("Warning: Has %d guilds.  Assuming each guild in list is %q\n", len(s.State.Guilds), o.guildName)
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

			if channel.Name == o.userChannel {
				found++
				fmt.Printf("[discord] Found user channel at %q, spawning workers.\n", channel.ID)
				o.onUserChannelFound(channel.ID)
			}
			if channel.Name == o.opChannel {
				found++
				opsChannelID = channel.ID
				fmt.Printf("[discord] Found ops channel at %q, spawning workers.\n", channel.ID)
				o.onOpsChannelFound(opsChannelID)
			}
		}
	}
	if found < 2 {
		fmt.Printf("[discord] Warning: Could not find guild %q with channel %q or %q\n", o.guildName, o.userChannel, o.opChannel)
	} else {
		if _, err := s.ChannelMessageSend(opsChannelID, "Overseer <-> Discord connection established."); err != nil {
			fmt.Printf("[discord] Warning: Unable to send initial message because %s\n", err.Error())
		}
	}
	return suture.ErrDoNotRestart
}

func (o *onConnection) onUserChannelFound(channelID string) {
	incoming := suture.NewSimple("incoming")
	incoming.Add(&onIncomingMessage{
		session:       o.session,
		onlyChannelID: channelID,
		sink:          o.userMessages,
	})
	o.sessionTree.Add(incoming)

	outgoingSupervisor := suture.NewSimple("outgoing")
	outgoingSupervisor.Add(&outgoingDispatcher{o.session, channelID, o.userReplies, o.instanceName})
	o.sessionTree.Add(outgoingSupervisor)
}

func (o *onConnection) onOpsChannelFound(channelID string) {

}

type onIncomingMessage struct {
	session       *discordgo.Session
	onlyChannelID string
	sink          chan<- discordgo.Message
}

func (i *onIncomingMessage) Serve(ctx context.Context) error {
	done := i.session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if m.ChannelID == i.onlyChannelID {
			i.sink <- *m.Message
		}
	})
	defer done()

	<-ctx.Done()
	return ctx.Err()
}
