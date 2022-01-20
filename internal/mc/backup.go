package mc

import (
	"fmt"
	"github.com/meschbach/go-junk-bucket/sub"
)

func BackupInstance(gameDirectory string, bucket string, key string) error {
	tar := sub.NewSubcommand("tar", []string{"czf", "backup.tgz", gameDirectory})
	if err := tar.PumpToStandard("backup@" + key); err != nil {
		return err
	}
	return fmt.Errorf("todo")
}
