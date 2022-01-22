package config

import "github.com/meschbach/minecraft-overseer/internal/minio"

type BackupSpecV1 struct {
	Minio *minio.BackupSpecV1 `json:"minio"`
}
