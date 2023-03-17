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
	connectionFactory func(parent *suture.Supervisor, s *discordgo.Session) suture.Service
}

func (c *connector) Serve(ctx context.Context) error {
	fmt.Printf("[discord] Attempting to connect.\n")
	events := make(chan connectorOp, 16)
	client, err := discordgo.New("Bot " + c.token)
	if err != nil {
		return err
	}

	var token *suture.ServiceToken
	client.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		events <- &connectorFuncOp{op: func(ctx context.Context) error {
			fmt.Printf("[discord] Connected with session %q\n", s.State.SessionID)
			if token != nil {
				if err := c.connectionsTree.Remove(*token); err != nil {
					return err
				}
				token = nil
			}
			sessionTree := suture.NewSimple("connection-" + s.State.SessionID)
			factory := suture.NewSimple("factory")
			factory.Add(c.connectionFactory(sessionTree, s))
			sessionTree.Add(factory)

			newToken := c.connectionsTree.Add(sessionTree)
			token = &newToken
			return nil
		}}
	})
	client.AddHandler(func(s *discordgo.Session, event *discordgo.Disconnect) {
		events <- &connectorFuncOp{op: func(ctx context.Context) error {
			fmt.Println("[discord] Disconnected")
			if token != nil {
				if err := c.connectionsTree.Remove(*token); err != nil {
					return err
				}
				token = nil
			}
			return nil
		}}
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
		case op := <-events:
			if err := op.perform(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

type connectorOp interface {
	perform(ctx context.Context) error
}

type connectorFuncOp struct {
	op func(ctx context.Context) error
}

func (c *connectorFuncOp) perform(ctx context.Context) error {
	return c.op(ctx)
}
