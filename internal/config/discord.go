package config

import (
	"context"
	"errors"
	"fmt"
	junk "github.com/meschbach/go-junk-bucket/pkg/files"
	"github.com/meschbach/minecraft-overseer/internal/discord"
	"github.com/meschbach/minecraft-overseer/internal/mc"
	"github.com/thejerf/suture/v4"
	"os"
)

// DiscordSpecV1 is V1 of the discord configuration specification
// TODO: move this into Discord
// really requires machinery to build a command interpreter system
type DiscordSpecV1 struct {
	AuthSpecFile string `json:"auth-file,omitempty"`
	Guild        string `json:"guild"`
	Channel      string `json:"channel"`
	OpsChannel   string `json:"ops-channel"`
}

func (d *DiscordSpecV1) interpret(config *RuntimeConfig) error {
	if d.Guild == "" {
		return errors.New("guild is empty")
	}
	if d.Channel == "" {
		return errors.New("channel is empty")
	}
	if len(d.OpsChannel) == 0 {
		d.OpsChannel = d.Channel
	}

	instanceName := "Minecraft "
	if name, has := os.LookupEnv("INSTANCE_NAME"); has {
		instanceName = name + " "
	}

	portSpec := ""
	if value, has := os.LookupEnv("PORT_SPEC"); has {
		portSpec = value
	}

	var manifest DiscordAuthSpecV1
	if err := junk.ParseJSONFile(d.AuthSpecFile, &manifest); err != nil {
		return err
	}
	config.subsystems = append(config.subsystems, &discordLogger{
		token:        manifest.Token,
		guild:        d.Guild,
		channel:      d.Channel,
		opsChannel:   d.OpsChannel,
		instanceName: instanceName,
		portSpec:     portSpec,
	})
	return nil
}

func (d *DiscordSpecV1) ParseAuthFile(spec *DiscordAuthSpecV1) error {
	return junk.ParseJSONFile(d.AuthSpecFile, spec)
}

type DiscordAuthSpecV1 struct {
	Token string
}

type discordLogger struct {
	token        string
	guild        string
	channel      string
	opsChannel   string
	instanceName string
	portSpec     string
}

// TODO 1: this is not needed, driving point for probably wrong abstraction
func (d *discordLogger) Start(systemContext context.Context, instance *mc.Instance) error {
	return nil
}

func (d *discordLogger) OnGameStart(systemContext context.Context, game *mc.RunningGame) error {
	actors := suture.NewSimple("discord")
	actors.Add(&startedNotice{})
	discord.NewSystem(&discord.Runtime{
		Parent:     actors,
		Dispatcher: game.Reactor.Logs,
	}, discord.SystemConfig{
		InstanceName: d.instanceName,
		Token:        d.token,
		PortSpec:     d.portSpec,
		Guild: discord.GuildConfig{
			GuildName:   d.guild,
			UserChannel: d.channel,
			OpsChannel:  d.opsChannel,
		},
	})
	go func() {
		<-systemContext.Done()
		fmt.Println("[discord] Game context closing.")
	}()
	go func() {
		fmt.Printf("[discord] Starting.\n")
		err := actors.Serve(systemContext)
		if err != nil {
			fmt.Printf("Discord errored: %s\n", err.Error())
		}
	}()
	return nil
}

type startedNotice struct {
}

func (s *startedNotice) Serve(ctx context.Context) error {
	fmt.Printf("[discord] Actor system started.\n")
	return suture.ErrDoNotRestart
}
