package mc

import (
	"fmt"
	"github.com/meschbach/go-junk-bucket/sub"
	"github.com/meschbach/minecraft-overseer/internal/junk"
	"os"
	"path/filepath"
)

type Instance struct {
	GameDirectory string
}

func NewInstance(gameDirectory string) (*Instance, error) {
	absoluteGameDirectory, err := filepath.Abs(gameDirectory)
	if err != nil {
		return nil, err
	}

	return &Instance{
		GameDirectory: absoluteGameDirectory,
	}, nil
}

func (i *Instance) PrepareRunning() (*RunningGame, error) {
	stdoutChannel := make(chan string, 16)
	echoStdoutChannel := make(chan string, 16)
	gameEventsInput := make(chan string, 16)
	stdout := &junk.StringBroadcast{
		Input: stdoutChannel,
		Out:   []chan<- string{echoStdoutChannel, gameEventsInput},
	}
	stderr := make(chan string, 16)
	stdin := make(chan string, 16)
	go func() {
		fmt.Println("<<stderr initialized>>")
		for msg := range stderr {
			fmt.Fprintf(os.Stderr, "<<stderr>> %s\n", msg)
		}
	}()

	//standard output
	go stdout.RunLoop()

	go func() {
		fmt.Println("<<stdout initialized>>")
		for msg := range echoStdoutChannel {
			fmt.Printf("<<stdout>> %s\n", msg)
		}
	}()

	reactor := StartReactor(gameEventsInput, stdin)

	cmd := sub.NewSubcommand("java", []string{
		"-Dlog4j2.formatMsgNoLookups=true", "-Dlog4j.configurationFile=/log4j.xml",
		"-jar", "minecraft_server.jar",
		"--nogui"})
	cmd.WithOption(&sub.WorkingDir{Where: i.GameDirectory})

	return &RunningGame{
		start: func() error {
			return cmd.Interact(stdin, stdoutChannel, stdin)
		},
		Stdout:  stdout,
		Reactor: reactor,
	}, nil
}

type RunningGame struct {
	start   func() error
	Stdout  *junk.StringBroadcast
	Reactor *ConsoleReactor
}

func (r *RunningGame) Run() error {
	return r.start()
}
