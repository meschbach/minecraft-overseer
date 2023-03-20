package discord

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/thejerf/suture/v4"
)

type connector struct {
	token             string
	connectionsTree   *suture.Supervisor
	connectionFactory func(parent *suture.Supervisor, s *discordgo.Session, firstConnection bool) suture.Service
}

func (c *connector) Serve(ctx context.Context) error {
	fmt.Printf("[discord] Attempting to connect.\n")
	client, err := discordgo.New("Bot " + c.token)
	if err != nil {
		return err
	}

	state := &connectorState{}
	state.init(c)

	client.AddHandler(func(s *discordgo.Session, event *discordgo.Connect) {
		state.onConnected(s)
	})
	client.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		state.onReady(s)
	})
	client.AddHandler(func(s *discordgo.Session, event *discordgo.Disconnect) {
		state.onDisconnected()
	})
	//todo: figure out how to make this context aware
	if err := client.Open(); err != nil {
		return err
	}
	defer func() {
		//todo: log when there is an error
		client.Close()
	}()
	fmt.Printf("[discord] responding to events.\n")

	for {
		select {
		case op := <-state.ops:
			if err := op.perform(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
