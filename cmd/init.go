package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
)

type ManifestV1 struct {
	Repository string
	Plugins    []string
	Forge      string
}

func parseManifest(ctx context.Context, manifest *Manifest, fileName string) error {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, manifest); err != nil {
		return err
	}
	return nil
}

//TODO: Use context aware gets
func downloadFile(ctx context.Context, fileURL string, to string) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	sink, err := os.Create(to)
	if err != nil {
		return err
	}
	defer sink.Close()

	_, err = io.Copy(sink, resp.Body)
	return err
}

func newInitCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init <manifest>",
		Short:   "Sets up Minecraft within a container, caching files for deployment",
		PreRunE: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manifestFile := args[0]
			ctx := cmd.Context()

			manifest := &Manifest{}
			if err := parseManifest(ctx, manifest, manifestFile); err != nil {
				return err
			}

			repository := manifest.V1.Repository
			fmt.Printf("Using %q as base for all files", repository)

			err := downloadFile(ctx, repository+"/"+manifest.V1.Forge, manifest.V1.Forge)
			if err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
