package game

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

type startCommand struct {
}

func (*startCommand) run(state *internalState, game *Game) error {
	if state.state != idleState {
		return nil
	}
	state.state = startingState

	out, err := state.serviceProcess.StdoutPipe()
	if err != nil {
		return err
	}
	serverOutput := bufio.NewReader(out)
	go game.pumpStream(serverOutput)

	stderrRaw, err := state.serviceProcess.StderrPipe()
	if err != nil {
		return err
	}
	stderrReader := bufio.NewReader(stderrRaw)
	go game.pumpStream(stderrReader)

	stdinRaw, err := state.serviceProcess.StdinPipe()
	if err != nil { return err }
	stdin := bufio.NewWriter(stdinRaw)
	go state.pumpCommands(stdin, game)

	if err := state.serviceProcess.Start(); err != nil {
		return err
	}
	state.state = runningState
	go game.postCleanup(state.serviceProcess)
	return nil
}

func (i *Game) pumpStream(serverOutput *bufio.Reader)  {
	for {
		line, err := serverOutput.ReadString('\n')
		if err == io.EOF {
			//TODO: Determine how to properly handle closed service
			return
			//close(i.ServiceMessage)
		} else if err != nil {
			//TODO: Shut things down gracefully
			panic(err)
		} else {
			i.ServiceMessage <- line
		}
	}
}

func (i *Game) postCleanup(proc *exec.Cmd) {
	err := proc.Wait()
	if err != nil {
		if errorCode, ok := err.(*exec.ExitError); ok {
			serviceMessage := fmt.Sprintf("Exit code of %d", errorCode.ExitCode())
			i.ServiceMessage <- serviceMessage
		} else {
			i.ServiceMessage <- err.Error()
		}
	}
}

func (i *internalState) pumpCommands(out *bufio.Writer, game *Game) {
	for {
		cmd := <- i.commands
		game.ServiceMessage <- fmt.Sprintf("[command] '%s'", cmd)
		output := fmt.Sprintf("%s\n", cmd)
		count, err := out.WriteString(output)
		if err != nil {
			panic(err)
		}
		if count < len(cmd) {
			panic("write shorter than command")
		}
		if err := out.Flush(); err != nil {
			panic(err)
		}
	}
}