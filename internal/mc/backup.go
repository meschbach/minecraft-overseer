package mc

import (
	"context"
	"github.com/meschbach/go-junk-bucket/sub"
	"path/filepath"
)

type BackupTarget interface {
	StoreArtifact(ctx context.Context, sourceFileName string) error
}

func (i *Instance) BackupInstance(ctx context.Context, target BackupTarget) error {
	//Resolve the parent directory
	parent := filepath.Dir(i.GameDirectory)
	base := filepath.Base(i.GameDirectory)
	//Generate the minio
	fileName := "minio.tgz"
	backupFileName := filepath.Join(parent, fileName)
	tar := sub.NewSubcommand("tar", []string{"czf", backupFileName, base})
	tar.WithOption(&sub.WorkingDir{Where: parent})
	if err := tar.PumpToStandard("minio@" + backupFileName); err != nil {
		return err
	}
	//Actually perform minio
	err := target.StoreArtifact(ctx, backupFileName)
	if err != nil {
		return err
	}
	return nil
}
