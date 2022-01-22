package minio

import (
	"context"
	"fmt"
	"github.com/meschbach/minecraft-overseer/internal/junk"
	"github.com/meschbach/minecraft-overseer/internal/mc"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"time"
)

type BackupSpecV1 struct {
	Frequency       string `json:"frequency"`
	CredentialsFile string `json:"credentials-file"`
	Bucket          string
	KeyPrefix       string `json:"key-prefix"`
}

func (m *BackupSpecV1) Start(overseer context.Context, game *mc.Instance) error {
	fmt.Println("Starting backup agent")
	//Parse configuration file
	credentials := BackupCredentialsFileV1{}
	if err := junk.ParseJSONFile(m.CredentialsFile, &credentials); err != nil {
		return err
	}

	//setup and verify connectivity
	client, err := credentials.buildClient()
	if err != nil {
		return &BackupConfigurationError{err}
	}

	mailbox := make(chan interface{})
	agent := &BackupAgent{
		client:    client,
		bucket:    m.Bucket,
		keyPrefix: m.KeyPrefix,
		mailbox:   mailbox,
		game:      game,
	}
	fmt.Println("Verifying backup agent connectivity")
	if err := agent.verify(overseer); err != nil {
		return &BackupConfigurationError{err}
	}

	//spawn loop
	go func() {
		<-overseer.Done()
		fmt.Println("Underlying application context is done, exiting backup agent")
		mailbox <- &poisonPill{}
	}()

	//Backup frequency
	if len(m.Frequency) == 0 {
		m.Frequency = "1h0m0s"
	}
	backupFrequency, err := time.ParseDuration(m.Frequency)
	if err != nil {
		return err
	}
	fmt.Printf("Backing up game instance very %s\n", backupFrequency)

	//TODO: should be extracted into state machine actor to manage
	ticker := time.Tick(backupFrequency)
	go func() {
		for t := range ticker {
			mailbox <- &takeBackup{
				name: t.Format(time.RFC3339) + ".tgz",
			}
		}
	}()
	go agent.runActor()

	return nil
}

//TODO 1: another example where this design may be wrong
func (m *BackupSpecV1) OnGameStart(gameContext context.Context, game *mc.RunningGame) error {
	return nil
}

type BackupCredentialsFileV1 struct {
	UseTLS    bool   `json:"use-tls"`
	KeyID     string `json:"access-key"`
	KeySecret string `json:"key-secret"`
	Endpoint  string `json:"endpoint"`
}

func (b *BackupCredentialsFileV1) buildClient() (*minio.Client, error) {
	fmt.Printf("Creating new Minio client with key %q to %q (tls: %t)\n", b.KeyID, b.Endpoint, b.UseTLS)
	return minio.New(b.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(b.KeyID, b.KeySecret, ""),
		Secure: b.UseTLS,
	})
}

type BackupConfigurationError struct {
	Underlying error
}

func (b *BackupConfigurationError) Error() string {
	return fmt.Sprintf("failed to configure because %s", b.Underlying.Error())
}
