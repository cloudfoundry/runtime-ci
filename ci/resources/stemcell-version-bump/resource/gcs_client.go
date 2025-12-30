package resource

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCSClient struct {
	client *storage.Client
}

func NewGCSClient(jsonKey string) (GCSClient, error) {
	client, err := storage.NewClient(context.TODO(), option.WithCredentialsJSON([]byte(jsonKey)))
	if err != nil {
		return GCSClient{}, fmt.Errorf("failed to create GCS storage client: %w", err)
	}

	return GCSClient{
		client: client,
	}, nil
}

func (gcsclient GCSClient) Get(bucketName string, objectPath string) ([]byte, error) {
	rc, err := gcsclient.client.Bucket(bucketName).Object(objectPath).NewReader(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to get version info from GCS (%s, %s): %w", bucketName, objectPath, err)
	}
	defer rc.Close() //nolint:errcheck

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("failed to read version info from reader: %w", err)
	}

	return data, nil
}

func (gcsclient GCSClient) Put(bucketName string, objectPath string, contents []byte) error {
	wc := gcsclient.client.Bucket(bucketName).Object(objectPath).NewWriter(context.TODO())

	_, err := wc.Write(contents)
	if err != nil {
		return fmt.Errorf("failed to write version info to writer: %w", err)
	}

	err = wc.Close()
	if err != nil {
		return fmt.Errorf("failed to close GCS writer: %w", err)
	}

	return nil
}
