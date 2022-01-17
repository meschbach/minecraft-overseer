package config

import "github.com/meschbach/minecraft-overseer/internal/junk"

type DiscordManifestSpec struct {
	AuthSpecFile string `json:"auth-file,omitempty"`
}

func (d *DiscordManifestSpec) ParseAuthFile(spec *DiscordAuthSpec) error {
	return junk.ParseJSONFile(d.AuthSpecFile, spec)
}

type DiscordAuthSpec struct {
	Token string
}
