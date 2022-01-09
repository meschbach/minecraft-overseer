package main

import (
	"context"
	"errors"
	"github.com/meschbach/go-junk-bucket/sub"
	"io"
	"net/http"
	"os"
	"path"
)

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
