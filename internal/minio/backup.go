package minio

import (
	"context"
	"fmt"
	"github.com/meschbach/minecraft-overseer/internal/mc"
	"github.com/minio/minio-go/v7"
	"os"
)

type poisonPill struct {
}

type takeBackup struct {
	name string
}

type BackupAgent struct {
	client    *minio.Client
	bucket    string
	keyPrefix string
	mailbox   chan interface{}
	game      *mc.Instance
}

type AccessError struct {
	Bucket     string
	KeyPrefix  string
	Underlying error
}

func (a *AccessError) Error() string {
	return fmt.Sprintf("failed to access bucket %q with key prefix %q because %s", a.Bucket, a.KeyPrefix, a.Underlying.Error())
}

func (b *BackupAgent) verify(ctx context.Context) error {
	//TODO: faster & more reliable way to ensure we have reasonable connectivity?
	listResult := b.client.ListObjects(ctx, b.bucket, minio.ListObjectsOptions{
		Prefix:    b.keyPrefix,
		Recursive: true,
	})
	for obj := range listResult {
		if obj.Err != nil {
			return &AccessError{
				Bucket:     b.bucket,
				KeyPrefix:  b.keyPrefix,
				Underlying: obj.Err,
			}
		}
	}
	return nil
}

func (b *BackupAgent) runActor() {
	fmt.Println("Running backup agent.")
	for msg := range b.mailbox {
		switch t := msg.(type) {
		case *poisonPill:
			break
		case *takeBackup:
			if err := b.doBackup(t.name); err != nil {
				panic(err)
			}
		default:
			panic(fmt.Errorf("unknown message type %#v", msg))
		}
	}
}

func (b *BackupAgent) doBackup(name string) error {
	ctx, done := context.WithCancel(context.Background())
	defer done()

	fmt.Printf("Performing Minio backup to name %q\n", name)
	err := b.game.BackupInstance(ctx, &BackupMinioTarget{
		client: b.client,
		Bucket: b.bucket,
		Key:    b.keyPrefix + name,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Backup to %q completed\n", name)
	return nil
}

type BackupMinioTarget struct {
	client *minio.Client
	Bucket string
	Key    string
}

func (b *BackupMinioTarget) StoreArtifact(ctx context.Context, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return err
	}

	_, err = b.client.PutObject(ctx, b.Bucket, b.Key, file, fileStat.Size(), minio.PutObjectOptions{ContentType: "application/tar+gzip"})
	return err
}
