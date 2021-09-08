package azuremutex

import (
	"context"
	"fmt"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

type AzureMutex struct {
	ctx                context.Context
	accountName        string
	accountKey         string
	containerName      string
	leaseId            string
	containerReference *azblob.ContainerURL
}

const (
	leaseDuration      = 30
)

func NewMutex(accountName string, accountKey string, container string) *AzureMutex {
	mutex := AzureMutex{
		accountName:   accountName,
		accountKey:    accountKey,
		containerName: container,
		ctx:           context.Background(),
	}
	return &mutex
}

func (m *AzureMutex) Acquire(key string) error {

	var err error
	m.containerReference, err = m.createContainerURL()
	if err != nil {
		return err
	}
	err = m.createContainerIfNotExists()
	if err != nil {
		return err
	}

	blob := m.containerReference.NewBlockBlobURL(key)
	_, err = azblob.UploadBufferToBlockBlob(m.ctx, []byte{}, blob, azblob.UploadToBlockBlobOptions{})
	if err != nil {
		return err
	}

	response, err := blob.AcquireLease(m.ctx, "", leaseDuration, azblob.ModifiedAccessConditions{})
	if err != nil {
		return err
	}

	m.leaseId = response.LeaseID()
	return nil
}

func (m *AzureMutex) Renew(key string) error {
	if m.containerReference == nil {
		return fmt.Errorf("lock not aquired")
	}

	blob := m.containerReference.NewBlockBlobURL(key)
	_, err := blob.RenewLease(m.ctx, m.leaseId, azblob.ModifiedAccessConditions{})

	return err
}

func (m *AzureMutex) Release(key string) error {
	if m.containerReference == nil {
		return fmt.Errorf("lock not aquired")
	}

	blob := m.containerReference.NewBlockBlobURL(key)
	_, err := blob.ReleaseLease(m.ctx, m.leaseId, azblob.ModifiedAccessConditions{})

	return err
}
