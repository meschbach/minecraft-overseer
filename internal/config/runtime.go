package config

import (
	"context"
	"github.com/meschbach/minecraft-overseer/internal/mc"
)

// Subsystem is a plugin point
// TODO 1: Re-evaluate if Start and OnGameStart make sense
type Subsystem interface {
	Start(systemContext context.Context, instance *mc.Instance) error
	OnGameStart(systemContext context.Context, game *mc.RunningGame) error
}

type RuntimeConfig struct {
	subsystems []Subsystem
}

func (r *RuntimeConfig) Start(systemContext context.Context, instance *mc.Instance) error {
	for _, subsystem := range r.subsystems {
		if err := subsystem.Start(systemContext, instance); err != nil {
			return err
		}
	}
	return nil
}

func (r *RuntimeConfig) OnGameStart(systemContext context.Context, game *mc.RunningGame) error {
	for _, subsystem := range r.subsystems {
		if err := subsystem.OnGameStart(systemContext, game); err != nil {
			return err
		}
	}
	return nil
}
