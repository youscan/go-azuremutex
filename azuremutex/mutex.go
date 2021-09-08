package azuremutex

import (
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	log "github.com/sirupsen/logrus"
	"net/url"
)

type AzureMutex interface {
	Acquire(key string) error
	Release(key string)
}

type mutexState struct {
	accountName string
	accountKey  string
	container   string
	ctx         context.Context
}

func NewMutex(accountName string, accountKey string, container string) AzureMutex {
	return mutexState{
		accountName: accountName,
		accountKey:  accountKey,
		container:   container,
		ctx:         context.Background(),
	}
}

func (m mutexState) Acquire(key string) error {

	_, err := m.createContainerIfNotExists()
	if err != nil {
		return err
	}

	return nil
}

func (m mutexState) Release(key string) {
	log.Println("Release not implemented")
}

func (m mutexState) createContainerIfNotExists() (*azblob.ContainerURL, error) {
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", m.accountName, m.container))

	credential, err := azblob.NewSharedKeyCredential(m.accountName, m.accountKey)
	if err != nil {
		return nil, err
	}
	containerURL := azblob.NewContainerURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{}))

	_, err = containerURL.Create(m.ctx, nil, azblob.PublicAccessNone)

	if err == nil {
		return &containerURL, nil
	}

	if stgErr, ok := err.(azblob.StorageError); ok && stgErr.ServiceCode() == "ContainerAlreadyExists" {
		return &containerURL, nil
	}

	return nil, err
}
