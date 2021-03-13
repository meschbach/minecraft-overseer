package main

import "os/exec"

const (
	Stop = 0
	Start = 1
)

type StateMachine struct {
	running bool
	output chan string
	fsm chan int
	quit chan int
	instance RunningInstance
}

type RunningInstance struct {

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

		mc := exec.Command("java","-jar", "papermc.jar")
		if err := mc.Start(); err != nil {
			s.output <- "[overseer] Failed to start process"
			s.output <- err.Error()
		}

		go func() {
			err := mc.Wait()
			s.output <- "[overseer] Minecraft exited"
			if err != nil {
				s.output <- "[overseer] Error reported while exiting -- " + err.Error()
			}
		}()
	}
}


func (s *StateMachine) stopMinecraft() {
	if s.running {
		s.output <- "[overseer] Stopping"
		s.running = false
	} else {
		s.output <- "[overseer] Not running"
	}
}
