package discord

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/thejerf/suture/v4"
)

type connectorState struct {
	config       *connector
	ops          chan connectorStateOp
	hasBeenReady bool
	serviceToken *suture.ServiceToken
}

func (c *connectorState) init(config *connector) {
	c.config = config
	c.ops = make(chan connectorStateOp, 16)
	c.hasBeenReady = false
	c.serviceToken = nil
}

func (c *connectorState) onConnected(s *discordgo.Session) {
	c.perform(func(ctx context.Context) error {
		fmt.Printf("[discord] Connected %q (ready? %t)\n", s.State.SessionID, c.hasBeenReady)
		if c.hasBeenReady {
			if err := c.startConnectionHandlers(s); err != nil {
				return err
			}
		}
		return nil
	})
}

func (c *connectorState) onDisconnected() {
	c.perform(func(ctx context.Context) error {
		fmt.Println("[discord] Disconnected")
		return c.shutdownConnection()
	})
}

func (c *connectorState) onReady(s *discordgo.Session) {
	c.perform(func(ctx context.Context) error {
		if !c.hasBeenReady {
			c.hasBeenReady = true
			if err := c.startConnectionHandlers(s); err != nil {
				return err
			}
		}
		return nil
	})
}

func (c *connectorState) startConnectionHandlers(s *discordgo.Session) error {
	if err := c.shutdownConnection(); err != nil {
		return err
	}
	sessionTree := suture.NewSimple("connection-" + s.State.SessionID)
	factory := suture.NewSimple("factory")
	factory.Add(c.config.connectionFactory(sessionTree, s))
	sessionTree.Add(factory)

	newToken := c.config.connectionsTree.Add(sessionTree)
	c.serviceToken = &newToken
	return nil
}

func (c *connectorState) shutdownConnection() error {
	if c.serviceToken != nil {
		if err := c.config.connectionsTree.Remove(*c.serviceToken); err != nil {
			return err
		}
		c.serviceToken = nil
	}
	return nil
}

func (c *connectorState) perform(op func(ctx context.Context) error) {
	c.ops <- &connectorFuncOp{op: op}
}

type connectorStateOp interface {
	perform(ctx context.Context) error
}

type connectorFuncOp struct {
	op func(ctx context.Context) error
}

func (c *connectorFuncOp) perform(ctx context.Context) error {
	return c.op(ctx)
}
