package azmutex

import (
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

const (
	emulatorAccountName = "devstoreaccount1"
	emulatorAccountKey  = "Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw=="
)

func (m AzureMutex) createClient() (*azblob.Client, error) {
	var (
		u          *url.URL
		credential *azblob.SharedKeyCredential
		err        error
	)
	if m.options.UseStorageEmulator {
		u, _ = url.Parse(fmt.Sprintf("http://127.0.0.1:10000/%s", emulatorAccountName))
		credential, err = azblob.NewSharedKeyCredential(emulatorAccountName, emulatorAccountKey)
	} else {
		u, _ = url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/", m.options.AccountName))
		credential, err = azblob.NewSharedKeyCredential(m.options.AccountName, m.options.AccountKey)
	}
	if err != nil {
		return nil, err
	}
	client, err := azblob.NewClientWithSharedKeyCredential(u.String(), credential, &azblob.ClientOptions{})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (m AzureMutex) createContainerIfNotExists() error {
	if m.client == nil {
		return fmt.Errorf("containerClient not initialized")
	}
	containerClient := m.client.ServiceClient().NewContainerClient(m.options.ContainerName)
	_, err := containerClient.Create(m.ctx, &azblob.CreateContainerOptions{
		Access: nil, // Public access disabled when not specified explicitly
	})
	var respErr *azcore.ResponseError
	if errors.As(err, &respErr) && respErr.ErrorCode == "ContainerAlreadyExists" {
		return nil
	}
	return err
}
