package config

import (
	"fmt"
	"github.com/meschbach/minecraft-overseer/internal/junk"
)

type ManifestV2 struct {
	Type        string
	Version     string
	ServerURL   string                `json:"server-url"`
	DefaultOps  []string              `json:"default-operators"`
	Allowed     []string              `json:"allowed-users"`
	DiscordList []DiscordManifestSpec `json:"discord"`
	BackupSpec  *BackupSpecV1         `json:"backup,omitempty"`
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
