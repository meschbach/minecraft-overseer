package config

import (
	"context"
	"github.com/meschbach/minecraft-overseer/internal/mc"
)

type Subsystem interface {
	Start(systemContext context.Context, instance *mc.Instance) error
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
