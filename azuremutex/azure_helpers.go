package azuremutex

import (
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"net/url"
)

func (m AzureMutex) createContainerURL() (*azblob.ContainerURL, error) {
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", m.accountName, m.containerName))

	credential, err := azblob.NewSharedKeyCredential(m.accountName, m.accountKey)
	if err != nil {
		return nil, err
	}
	containerURL := azblob.NewContainerURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{}))
	return &containerURL, nil
}

func (m AzureMutex) createContainerIfNotExists() error {
	if m.containerReference == nil {
		return fmt.Errorf("containerURL not initialized")
	}

	_, err := m.containerReference.Create(m.ctx, nil, azblob.PublicAccessNone)

	// TODO: Use azblob.StorageErrorCodeContainerAlreadyExists here
	if stgErr, ok := err.(azblob.StorageError); ok && stgErr.ServiceCode() == "ContainerAlreadyExists" {
		return nil
	}

	return err
}
