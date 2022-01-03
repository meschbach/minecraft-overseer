package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/magiconair/properties"
	"log"
	"os"
	"path"
	"sync"

	"github.com/meschbach/go-junk-bucket/sub"
	"github.com/spf13/cobra"
)

func internalError(ws *websocket.Conn, msg string, err error) {
	log.Println(msg, err)
	ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}

type ManifestV2 struct {
	Type      string
	Version   string
	ServerURL string `json:"server-url"`
}

type Manifest struct {
	V1 *ManifestV1 `json:"v1,omitempty"`
	V2 *ManifestV2 `json:"v2,omitempty"`
}

func LoadJSONFile(fileName string, out interface{}) error {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, out)
}

func withoutFile(baseDir string, file string, perform func(fileName string) error) error {
	serverFile := path.Join(baseDir, file)
	_, err := os.Stat(serverFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return perform(serverFile)
		} else {
			return err
		}
	}
	return err
}

func pipeCommandOutput(prefix string, program string, args ...string) error {
	var pumpsDone sync.WaitGroup

	toClose := make([]chan string, 0)
	defer func() {
		for _, pipe := range toClose {
			close(pipe)
		}
		pumpsDone.Wait()
	}()
	buildPump := func(name string) chan string {
		pipe := make(chan string, 8)
		toClose = append(toClose, pipe)
		pumpsDone.Add(1)
		go sub.PumpPrefixedChannel(name, pipe, &pumpsDone)
		return pipe
	}

	stdoutPrefix := fmt.Sprintf("<<%s.stdout>>", prefix)
	stderrPrefix := fmt.Sprintf("<<%s.stderr>>", prefix)
	stdout := buildPump(stdoutPrefix)
	stderr := buildPump(stderrPrefix)
	initCmd := sub.NewSubcommand(program, args)
	return initCmd.Run(stdout, stderr)
}

func initV2(ctx context.Context, configFile string, gameDir string) error {
	fmt.Println("II    Ensuring initialized.")
	//Load && Parse Manifest file
	var manifest Manifest
	if err := LoadJSONFile(configFile, &manifest); err != nil {
		return err
	}
	if manifest.V1 != nil {
		return errors.New("v1 is no longer supported")
	}
	if manifest.V2 == nil {
		return errors.New("v2 must be provided")
	}

	switch manifest.V2.Type {
	case "vanilla":
	default:
		return errors.New("only vanilla supported right now")
	}

	// change to the game directory
	if err := os.Chdir(gameDir); err != nil {
		return err
	}

	//Ensure game file is downloaded
	if err := withoutFile(gameDir, "minecraft_server.jar", func(fileName string) error {
		fmt.Println("Server JAR does not exist.  Downloading.")
		return downloadFile(ctx, manifest.V2.ServerURL, fileName)
	}); err != nil {
		return err
	}

	//Has the configuration been seeded?
	if err := withoutFile(gameDir, "server.properties", func(fileName string) error {
		return pipeCommandOutput("config-default", "java", "-Dlog4j.configurationFile=/log4j.xml", "-jar", "minecraft_server.jar", "--initSettings", "--nogui")
	}); err != nil {
		return err
	}

	//Override the eula
	eulaProps, err := properties.LoadFile("eula.txt", properties.UTF8)
	if err != nil {
		return err
	}
	eulaValue, hasEula := eulaProps.Get("eula")
	if !hasEula || eulaValue != "true" {
		_, ok, err := eulaProps.Set("eula", "true")
		if err != nil {
			return err
		}
		if !ok {
			panic("unable to set key")
		}

		renderedEula := eulaProps.String()
		if err := os.WriteFile("eula.txt", []byte(renderedEula), 0777); err != nil {
			return err
		}
	}

	return nil
}

func RunServer(initCtx context.Context, opts *serverOpts) error {
	if err := initV2(initCtx, opts.fs.configFile, opts.fs.gameDir); err != nil {
		return err
	}
	err := pipeCommandOutput("game", "java", "-Dlog4j.configurationFile=/log4j.xml", "-jar", "minecraft_server.jar", "--nogui")
	if err != nil {
		return err
	}
	return nil
}

type serverOpts struct {
	httpAddress  string
	discordToken string
	fs           struct {
		configFile string
		gameDir    string
	}
}

func newServerCommands() *cobra.Command {
	opts := &serverOpts{}
	run := &cobra.Command{
		Use:   "server",
		Short: "Begins the Overseer service",
		//PreRunE: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			startup := cmd.Context()
			return RunServer(startup, opts)
		},
	}
	run.PersistentFlags().StringVar(&opts.httpAddress, "http-bind", "127.0.0.1:8080", "Port to bind webhost too")
	run.PersistentFlags().StringVar(&opts.discordToken, "discord-token", "", "Enables connecting to Discord")
	run.PersistentFlags().StringVarP(&opts.fs.gameDir, "game-dir", "d", "/mc/instance", "Game directory")
	run.PersistentFlags().StringVarP(&opts.fs.configFile, "config-file", "c", "/mc/config/manifest.json", "Configuration manifest for game")
	return run
}
