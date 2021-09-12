package azuremutex

import (
	"context"
	"fmt"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

type MutexOptions struct {
	AccountName        string
	AccountKey         string
	ContainerName      string
	UseStorageEmulator bool
}

type AzureMutex struct {
	ctx                context.Context
	options            MutexOptions
	leaseId            string
	containerReference *azblob.ContainerURL
}

func NewMutex(options MutexOptions) *AzureMutex {
	return NewMutexWithContext(options, context.Background())
}

func NewMutexWithContext(options MutexOptions, ctx context.Context) *AzureMutex {
	return &AzureMutex{
		options: options,
		ctx:     ctx,
	}
}

func (m *AzureMutex) Acquire(key string, leaseDuration int32) error {

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
	if stgErr, ok := err.(azblob.StorageError); ok && stgErr.ServiceCode() == "LeaseIdMissing" {
		return NewLeaseAlreadyPresentError(err)
	}
	if err != nil {
		return err
	}

	response, err := blob.AcquireLease(m.ctx, "", leaseDuration, azblob.ModifiedAccessConditions{})
	if stgErr, ok := err.(azblob.StorageError); ok && stgErr.ServiceCode() == "LeaseAlreadyPresent" {
		return NewLeaseAlreadyPresentError(err)
	}
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
