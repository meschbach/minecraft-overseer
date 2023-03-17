package discord

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type outgoingDispatcher struct {
	s         *discordgo.Session
	channelID string
	queue     <-chan string
	prefix    string
}

func (o *outgoingDispatcher) Serve(ctx context.Context) error {
	fmt.Printf("[discord] Starting outgoing dispatcher to %q\n", o.channelID)
	for {
		select {
		case msg := <-o.queue:
			_, err := o.s.ChannelMessageSend(o.channelID, o.prefix+msg, discordgo.WithContext(ctx))
			if err != nil {
				fmt.Printf("[discord] failed to send because %s\n", err.Error())
				return err
			}
		case <-ctx.Done():
			fmt.Printf("[discord] dispatch to %q is done\n", o.channelID)
			return ctx.Err()
		}
	}
}
