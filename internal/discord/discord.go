package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/meschbach/minecraft-overseer/internal/mc/events"
	"github.com/thejerf/suture/v4"
)

type Config struct {
	// Token is the discord token used to authenticate the bot.  Token should not include `Bot ` prefix.
	Token string
	// GuildName is ideally the discord server to connect to.  This is a TODO as I have not figured out how to request
	// via the API
	GuildName string
	// UserChannel is the human name of the user channel to accept user input and relay messages too.
	UserChannel string
	// OpsChannel is the human name place to push operations related messages.
	OpChannel string
	// UserInterpreter is the target interpreter for user commands
	UserInterpreter chan<- discordgo.Message
	// Outgoing are user messages intended to be sent to the Discord system
	Outgoing         <-chan string
	ParentSupervisor *suture.Supervisor
}

func NewLogger(config Config) (*EventLogger, error) {
	client, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return nil, err
	}

	subsystem := &EventLogger{
		client:     client,
		eventQueue: make(chan events.LogEntry, 128),
	}
	connectionHandler := &connection{
		guildName:        config.GuildName,
		userChannel:      config.UserChannel,
		opChannel:        config.OpChannel,
		subsystem:        subsystem,
		userCommands:     config.UserInterpreter,
		ParentSupervisor: config.ParentSupervisor,
		userReplies:      config.Outgoing,
	}
	client.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		connectionHandler.onReadyEvent(s, event)
	})
	if err := client.Open(); err != nil {
		return nil, err
	}
	return subsystem, nil
}
