package config

import (
	"fmt"
	junk "github.com/meschbach/go-junk-bucket/pkg/files"
)

type ManifestV2 struct {
	Type         string
	Version      string
	ServerURL    string          `json:"server-url"`
	DefaultOps   []string        `json:"default-operators"`
	Allowed      []string        `json:"allowed-users"`
	DiscordList  []DiscordSpecV1 `json:"discord"`
	BackupSpec   *BackupSpecV1   `json:"backup,omitempty"`
	InstanceSpec *InstanceSpecV1 `json:"instance,omitempty"`
}

func (m *ManifestV2) interpret(config *RuntimeConfig) error {
	if m.BackupSpec == nil {
		fmt.Printf("II\tNo backup strategy configured.  Skipping.")
	} else {
		if m.BackupSpec.Minio == nil {
			return fmt.Errorf("unsupported minio config %#v", m.BackupSpec)
		}
		config.subsystems = append(config.subsystems, m.BackupSpec.Minio)
	}

	for _, discordConfig := range m.DiscordList {
		if err := discordConfig.interpret(config); err != nil {
			return err
		}
	}

	return nil
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

func (m *Manifest) Interpret(config *RuntimeConfig) error {
	if m.V2 == nil {
		return fmt.Errorf("only v2 supported now")
	}

	return m.V2.interpret(config)
}

func ParseManifest(manifest *Manifest, fileName string) error {
	return junk.ParseJSONFile(fileName, manifest)
}
