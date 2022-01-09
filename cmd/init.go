package cmd

import (
	"context"
	"fmt"
	"github.com/meschbach/minecraft-overseer/internal/config"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
)

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

			manifest := &config.Manifest{}
			if err := config.ParseManifest(manifest, manifestFile); err != nil {
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
