package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/meschbach/minecraft-overseer/internal/discord/commands"
	"github.com/meschbach/minecraft-overseer/internal/mc/events"
	"github.com/thejerf/suture/v4"
)

type SystemConfig struct {
	// Token is the discord token used to authenticate the bot.  Token should not include `Bot ` prefix.
	Token string
	// InstanceName is the game identifier to use when communicating about the instance
	InstanceName string
	PortSpec     string
	Guild        GuildConfig `json:"guild,omitempty"`
}

type GuildConfig struct {
	GuildName   string
	UserChannel string
	OpsChannel  string
}

type Runtime struct {
	Parent     *suture.Supervisor
	Dispatcher *events.LogDispatcher
}

func NewSystem(rt *Runtime, spec SystemConfig) {
	commandsInput := make(chan discordgo.Message, 64)
	userOutput := make(chan string, 128)
	userOutput <- fmt.Sprintf("Starting!")
	connections := suture.NewSimple("connections")

	connectorSupervisor := suture.NewSimple("connector")
	connectorSupervisor.Add(&connector{
		token:           spec.Token,
		connectionsTree: connections,
		connectionFactory: func(parent *suture.Supervisor, s *discordgo.Session) suture.Service {
			return &onConnection{
				session:      s,
				sessionTree:  parent,
				guildName:    spec.Guild.GuildName,
				userChannel:  spec.Guild.UserChannel,
				userMessages: commandsInput,
				userReplies:  userOutput,
				opChannel:    spec.Guild.OpsChannel,
				instanceName: spec.InstanceName,
			}
		},
	})

	commandSystem := commands.NewCommandSystem(commands.Config{
		PortSpec: spec.PortSpec,
	}, commandsInput, userOutput)

	pumps := suture.NewSimple("pumps")
	pumps.Add(&eventPump{
		dispatcher: rt.Dispatcher,
		filter:     nonOperationEvents,
		sink:       userOutput,
	})
	pumps.Add(&eventPump{
		dispatcher: rt.Dispatcher,
		filter:     operationEvents,
		sink:       userOutput,
	})

	root := suture.NewSimple("discord")
	root.Add(connectorSupervisor)
	root.Add(connections)
	root.Add(pumps)
	root.Add(commandSystem)
	rt.Parent.Add(root)
}

func nonOperationEvents(entry events.LogEntry) bool {
	return !entry.IsOperations()
}

func operationEvents(entry events.LogEntry) bool {
	return entry.IsOperations()
}
