package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/magiconair/properties"
	"github.com/meschbach/minecraft-overseer/internal/config"
	"os"
)

type runtimeConfig struct {
	operators []string
}

func initV2(ctx context.Context, configFile string, gameDir string) (runtimeConfig, error) {
	fmt.Println("II    Ensuring initialized.")
	//Load && Parse Manifest file
	var manifest config.Manifest
	var configLater runtimeConfig
	if err := config.ParseManifest(&manifest, configFile); err != nil {
		return configLater, err
	}
	if manifest.V1 != nil {
		return configLater, errors.New("v1 is no longer supported")
	}
	if manifest.V2 == nil {
		return configLater, errors.New("v2 must be provided")
	}

	configLater.operators = manifest.V2.DefaultOps

	switch manifest.V2.Type {
	case "vanilla":
	default:
		return configLater, errors.New("only vanilla supported right now")
	}

	// change to the game directory
	if err := os.Chdir(gameDir); err != nil {
		return configLater, err
	}

	//Ensure game file is downloaded
	if err := withoutFile(gameDir, "minecraft_server.jar", func(fileName string) error {
		fmt.Println("Server JAR does not exist.  Downloading.")
		return downloadFile(ctx, manifest.V2.ServerURL, fileName)
	}); err != nil {
		return configLater, err
	}

	//Has the configuration been seeded?
	if err := withoutFile(gameDir, "server.properties", func(fileName string) error {
		err := pipeCommandOutput("configLater-default", "java", "-Dlog4j.configurationFile=/log4j.xml", "-jar", "minecraft_server.jar", "--initSettings", "--nogui")
		if err != nil {
			return err
		}

		//Set server configLater defaults
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
		return configLater, err
	}

	//Override the eula
	eulaProps, err := properties.LoadFile("eula.txt", properties.UTF8)
	if err != nil {
		return configLater, err
	}
	eulaValue, hasEula := eulaProps.Get("eula")
	if !hasEula || eulaValue != "true" {
		_, ok, err := eulaProps.Set("eula", "true")
		if err != nil {
			return configLater, err
		}
		if !ok {
			panic("unable to set key")
		}

		renderedEula := eulaProps.String()
		if err := os.WriteFile("eula.txt", []byte(renderedEula), 0700); err != nil {
			return configLater, err
		}
	}

	return configLater, nil
}
