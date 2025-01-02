package azmutex

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/lease"
	"net/url"
)

const (
	emulatorAccountName = "devstoreaccount1"
	emulatorAccountKey  = "Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw=="
)

type LeaseAlreadyPresentError struct {
	Err error
}

func NewLeaseAlreadyPresentError(err error) *LeaseAlreadyPresentError {
	return &LeaseAlreadyPresentError{
		Err: err,
	}
}

func (e *LeaseAlreadyPresentError) Error() string {
	return "lease already present"
}

type MutexOptions struct {
	AccountName        string
	AccountKey         string
	ContainerName      string
	UseStorageEmulator bool
	LogFunc            func(message string)
}

type AzureMutex struct {
	ctx          context.Context
	options      MutexOptions
	leaseClients map[string]*lease.BlobClient
}

func NewMutex(options MutexOptions) *AzureMutex {
	return NewMutexWithContext(options, context.Background())
}

func NewMutexWithContext(options MutexOptions, ctx context.Context) *AzureMutex {
	return &AzureMutex{
		options:      options,
		ctx:          ctx,
		leaseClients: make(map[string]*lease.BlobClient),
	}
}

func (m *AzureMutex) createClient() (*azblob.Client, error) {
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

func (m *AzureMutex) createContainerIfNotExists(containerClient *container.Client) error {
	_, err := containerClient.Create(m.ctx, &azblob.CreateContainerOptions{
		Access: nil, // Public access disabled when not specified explicitly
	})
	var respErr *azcore.ResponseError
	if errors.As(err, &respErr) && respErr.ErrorCode == "ContainerAlreadyExists" {
		return nil
	}
	return err
}

func (m *AzureMutex) Acquire(key string, leaseDuration int32) error {
	var err error
	client, err := m.createClient()
	if err != nil {
		return err
	}
	containerClient := client.ServiceClient().NewContainerClient(m.options.ContainerName)
	err = m.createContainerIfNotExists(containerClient)
	if err != nil {
		return err
	}

	_, err = client.UploadBuffer(m.ctx, m.options.ContainerName, key, []byte{}, &azblob.UploadBufferOptions{})

	var stgErr *azcore.ResponseError
	if errors.As(err, &stgErr) && stgErr.ErrorCode == "LeaseIdMissing" {
		return NewLeaseAlreadyPresentError(err)
	}
	if err != nil {
		return err
	}

	blobClient := containerClient.NewBlobClient(key)
	leaseClient, err := lease.NewBlobClient(blobClient, &lease.BlobClientOptions{})
	if err != nil {
		return err
	}
	_, err = leaseClient.AcquireLease(m.ctx, leaseDuration, &lease.BlobAcquireOptions{})
	if errors.As(err, &stgErr) && stgErr.ErrorCode == "LeaseAlreadyPresent" {
		return NewLeaseAlreadyPresentError(err)
	}
	if err != nil {
		return err
	}

	m.leaseClients[key] = leaseClient
	return nil
}

func (m *AzureMutex) Renew(key string) error {
	leaseClient, exists := m.leaseClients[key]
	if !exists {
		return fmt.Errorf("lock not acquired for key: %s", key)
	}
	_, err := leaseClient.RenewLease(m.ctx, &lease.BlobRenewOptions{})

	return err
}

func (m *AzureMutex) Release(key string) error {
	leaseClient, exists := m.leaseClients[key]
	if !exists {
		return fmt.Errorf("lock not acquired for key: %s", key)
	}

	_, err := leaseClient.ReleaseLease(m.ctx, &lease.BlobReleaseOptions{})

	return err
}
