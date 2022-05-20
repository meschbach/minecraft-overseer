package events

import "fmt"

type GenericDeathMessage struct {
	Message string
}

func (g *GenericDeathMessage) String() string {
	return fmt.Sprintf("%s", g.Message)
}

func (g *GenericDeathMessage) IsOperations() bool {
	return false
}
