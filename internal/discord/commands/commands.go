package commands

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/thejerf/suture/v4"
)

func NewCommandSystem(cfg Config, input chan discordgo.Message, output chan string) *suture.Supervisor {
	ipSystem := make(chan ipReq, 16)

	ips := suture.NewSimple("ip-address")
	ips.Add(&ipCommand{
		input:    ipSystem,
		portSpec: cfg.PortSpec,
	})

	dispatcher := suture.NewSimple("interpreter")
	dispatcher.Add(&interpreter{
		input:  input,
		output: output,
		ips:    ipSystem,
	})

	root := suture.NewSimple("commands")
	root.Add(dispatcher)
	root.Add(ips)
	return root
}

type interpreter struct {
	input  chan discordgo.Message
	output chan string
	ips    chan ipReq
}

func (i *interpreter) Serve(ctx context.Context) error {
	for {
		select {
		case in := <-i.input:
			if err := i.consumeMessage(ctx, in); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (i *interpreter) consumeMessage(ctx context.Context, msg discordgo.Message) error {
	content := msg.Content
	if content == "ip" {
		select {
		case i.ips <- ipReq{out: i.output}:
		case <-ctx.Done():
			return ctx.Err()
		default:
			select {
			case i.output <- "[ip] too many requests, dropping":
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	return nil
}
