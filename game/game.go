package game

import (
	"github.com/meschbach/minecraft-overseer/internal/mc/events"
	"os/exec"
)

const (
	idleState = iota
	startingState
	runningState
	stoppingState
	errorState
)

type gameState int

// Game manages a single instance of Minecraft.
//
// This is originally a derative work of https://levelup.gitconnected.com/lets-build-a-minecraft-server-wrapper-in-go-122c087e0023
// however this did not perform a number of functions to match what is needed here.
type Game struct {
	ServiceMessage chan events.LogEntry
	reactor        chan gameCommand
}

type internalState struct {
	state          gameState
	commands       chan string
	serviceProcess *exec.Cmd
}

type gameCommand interface {
	run(state *internalState, game *Game) error
}

func NewInstance(baseDirectory string) *Game {
	command := exec.Command("java", "-jar", "minecraft_server.jar", "nogui")
	command.Dir = baseDirectory
	i := &Game{
		ServiceMessage: make(chan events.LogEntry),
		reactor:        make(chan gameCommand),
	}
	go i.runState(command)
	return i
}

func (g *Game) runState(serviceProcess *exec.Cmd) {
	state := &internalState{
		state:          idleState,
		commands:       make(chan string),
		serviceProcess: serviceProcess,
	}
	for {
		msg := <-g.reactor
		err := msg.run(state, g)
		if err != nil {
			state.state = errorState
			//TODO: Cleanup
		}
	}
}

func (g *Game) Start() {
	g.reactor <- &startCommand{}
}

func (g *Game) Stop() {
	g.reactor <- &stopCommand{}
}
