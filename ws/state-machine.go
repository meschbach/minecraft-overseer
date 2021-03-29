package ws

import (
	"fmt"
	"github.com/meschbach/minecraft-overseer/game"
	"os"
	"path"
)

const (
	Stop = 0
	Start = 1
)

type StateMachine struct {
	running bool
	output chan string
	fsm chan int
	quit chan int
	game *game.Game
}

func newStateMachine( output chan string ) *StateMachine {
	return &StateMachine{
		running: false,
		output: output,
		fsm: make(chan int),
		quit: make(chan int),
	}
}

func (s *StateMachine) run() {
	println("[overseer] Starting")
	for {
		select {
		case instruction := <- s.fsm:
			switch instruction {
			case Start: s.startMinecraft()
			case Stop: s.stopMinecraft()
			default:
				println("Message", instruction)
			}
		case <- s.quit:
			println("quit")
			return
		}
	}
}

func (s *StateMachine) startMinecraft() {
	if s.running {
		s.output <- "[overseer] Already running"
	} else {
		s.output <- "[overseer] Starting"
		s.running = true

		pwd, err := os.Getwd()
		if err != nil {
			s.running = false
			s.output <- "[overseer] Failed to get current working directory"
			s.output <- err.Error()
			return
		}

		s.game = game.NewInstance(path.Join(pwd,"w"))
		s.game.Start()
		go s.pumpMessages()
		s.output <- "[overseer] Waiting for instance to start..."
	}
}

func (s *StateMachine) pumpMessages() {
	for {
		msg, ok := <- s.game.ServiceMessage
		if !ok {
			break
		}
		formatted := fmt.Sprintf("[instance] %s",msg)
		s.output <- formatted
	}
	s.output <- "[instance] channel consumed, exiting."
}

func (s *StateMachine) stopMinecraft() {
	if s.running {
		s.output <- "[overseer] Stopping"
		s.game.Stop()
		s.running = false
	} else {
		s.output <- "[overseer] Not running"
	}
}
