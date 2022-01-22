package config

import (
	"context"
	"errors"
	"github.com/meschbach/minecraft-overseer/internal/discord"
	"github.com/meschbach/minecraft-overseer/internal/junk"
	"github.com/meschbach/minecraft-overseer/internal/mc"
)

//DiscordSpecV1 is V1 of the discord configuration specification
//TODO: move this into Discord
//really requires machinery to build a command interpreter system
type DiscordSpecV1 struct {
	AuthSpecFile string `json:"auth-file,omitempty"`
	Guild        string `json:"guild"`
	Channel      string `json:"channel"`
}

func (d *DiscordSpecV1) interpret(config *RuntimeConfig) error {
	if d.Guild == "" {
		return errors.New("guild is empty")
	}
	if d.Channel == "" {
		return errors.New("channel is empty")
	}

	var manifest DiscordAuthSpecV1
	if err := junk.ParseJSONFile(d.AuthSpecFile, &manifest); err != nil {
		return err
	}
	config.subsystems = append(config.subsystems, &discordLogger{
		token:   manifest.Token,
		guild:   d.Guild,
		channel: d.Channel,
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
	token   string
	guild   string
	channel string
}

//TODO 1: this is not needed, driving point for probably wrong abstraction
func (d *discordLogger) Start(systemContext context.Context, instance *mc.Instance) error {
	return nil
}

func (d *discordLogger) OnGameStart(systemContext context.Context, game *mc.RunningGame) error {
	logger, err := discord.NewLogger(d.token, d.guild, d.channel)
	if err != nil {
		return err
	}
	go logger.Ingest(game.Reactor.Logs)
	return nil
}
