package config

import (
	"encoding/json"
	"os"
)

type ManifestV2 struct {
	Type       string
	Version    string
	ServerURL  string   `json:"server-url"`
	DefaultOps []string `json:"default-operators"`
}


type ManifestV1 struct {
	Repository string
	Plugins    []string
	Forge      string
}

type Manifest struct {
	V1 *ManifestV1 `json:"v1,omitempty"`
	V2 *ManifestV2 `json:"v2,omitempty"`
}

func ParseManifest(manifest *Manifest, fileName string) error {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, manifest); err != nil {
		return err
	}
	return nil
}
