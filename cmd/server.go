package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/magiconair/properties"
	"github.com/meschbach/go-junk-bucket/sub"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
)

func internalError(ws *websocket.Conn, msg string, err error) {
	log.Println(msg, err)
	ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}

type ManifestV2 struct {
	Type       string
	Version    string
	ServerURL  string   `json:"server-url"`
	DefaultOps []string `json:"default-operators"`
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
	initCmd := sub.NewSubcommand(program, args)
	return initCmd.PumpToStandard(prefix)
}

type runtimeConfig struct {
	operators []string
}

func initV2(ctx context.Context, configFile string, gameDir string) (runtimeConfig, error) {
	fmt.Println("II    Ensuring initialized.")
	//Load && Parse Manifest file
	var manifest Manifest
	var config runtimeConfig
	if err := LoadJSONFile(configFile, &manifest); err != nil {
		return config, err
	}
	if manifest.V1 != nil {
		return config, errors.New("v1 is no longer supported")
	}
	if manifest.V2 == nil {
		return config, errors.New("v2 must be provided")
	}

	config.operators = manifest.V2.DefaultOps

	switch manifest.V2.Type {
	case "vanilla":
	default:
		return config, errors.New("only vanilla supported right now")
	}

	// change to the game directory
	if err := os.Chdir(gameDir); err != nil {
		return config, err
	}

	//Ensure game file is downloaded
	if err := withoutFile(gameDir, "minecraft_server.jar", func(fileName string) error {
		fmt.Println("Server JAR does not exist.  Downloading.")
		return downloadFile(ctx, manifest.V2.ServerURL, fileName)
	}); err != nil {
		return config, err
	}

	//Has the configuration been seeded?
	if err := withoutFile(gameDir, "server.properties", func(fileName string) error {
		err := pipeCommandOutput("config-default", "java", "-Dlog4j.configurationFile=/log4j.xml", "-jar", "minecraft_server.jar", "--initSettings", "--nogui")
		if err != nil {
			return err
		}

		//Set server config defaults
		serverProperties, err := properties.LoadFile("server.properties", properties.UTF8)
		if err != nil {
			return err
		}
		serverProperties.SetValue("sync-chunk-writes", "false")
		serverProperties.SetValue("motd", "Minecraft Overseer provisioned world")
		serverProperties.SetValue("white-list", "true")
		serverProperties.SetValue("spawn-protection", "1")
		if err := os.WriteFile("server.properties", []byte(serverProperties.String()), 0700); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return config, err
	}

	//Override the eula
	eulaProps, err := properties.LoadFile("eula.txt", properties.UTF8)
	if err != nil {
		return config, err
	}
	eulaValue, hasEula := eulaProps.Get("eula")
	if !hasEula || eulaValue != "true" {
		_, ok, err := eulaProps.Set("eula", "true")
		if err != nil {
			return config, err
		}
		if !ok {
			panic("unable to set key")
		}

		renderedEula := eulaProps.String()
		if err := os.WriteFile("eula.txt", []byte(renderedEula), 0700); err != nil {
			return config, err
		}
	}

	return config, nil
}

func RunProgram(initCtx context.Context, opts *serverOpts) error {
	runtimeConfig, err := initV2(initCtx, opts.fs.configFile, opts.fs.gameDir)
	if err != nil {
		return err
	}
	fmt.Printf("Passed configuration: %v\n", runtimeConfig)

	stdout := make(chan string, 16)
	stderr := make(chan string, 16)
	cmd := sub.NewSubcommand("java", []string{
		"-Dlog4j2.formatMsgNoLookups=true", "-Dlog4j.configurationFile=/log4j.xml",
		"-jar", "minecraft_server.jar",
		"--nogui"})
	go func() {
		fmt.Println("<<stderr initialized>>")
		for msg := range stderr {
			fmt.Fprintf(os.Stderr, "<<stderr>> %s\n", msg)
		}
	}()
	go func() {
		fmt.Println("<<stdout initialized>>")
		for msg := range stdout {
			fmt.Printf("<<stdout>> %s\n", msg)
		}
	}()

	stdin := make(chan string, 16)
	go func() {
		for _, operator := range runtimeConfig.operators {
			stdin <- fmt.Sprintf("whitelist add %s", operator)
			stdin <- fmt.Sprintf("op %s", operator)
		}
	}()

	err = cmd.Interact(stdin, stdout, stderr)
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
			return RunProgram(startup, opts)
		},
	}
	run.PersistentFlags().StringVar(&opts.httpAddress, "http-bind", "127.0.0.1:8080", "Port to bind webhost too")
	run.PersistentFlags().StringVar(&opts.discordToken, "discord-token", "", "Enables connecting to Discord")
	run.PersistentFlags().StringVarP(&opts.fs.gameDir, "game-dir", "d", "/mc/instance", "Game directory")
	run.PersistentFlags().StringVarP(&opts.fs.configFile, "config-file", "c", "/mc/config/manifest.json", "Configuration manifest for game")
	return run
}