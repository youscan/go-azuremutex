package azuremutex

import (
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"net/url"
)

const (
	emulatorAccountName = "devstoreaccount1"
	emulatorAccountKey  = "Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw=="
)

func (m AzureMutex) createContainerURL() (*azblob.ContainerURL, error) {
	var (
		u          *url.URL
		credential *azblob.SharedKeyCredential
		err        error
	)
	if m.options.UseStorageEmulator {
		u, _ = url.Parse(fmt.Sprintf("http://127.0.0.1:10000/%s/%s", emulatorAccountName, m.options.ContainerName))
		credential, err = azblob.NewSharedKeyCredential(emulatorAccountName, emulatorAccountKey)
	} else {
		u, _ = url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", m.options.AccountName, m.options.ContainerName))
		credential, err = azblob.NewSharedKeyCredential(m.options.AccountName, m.options.AccountKey)
	}
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
