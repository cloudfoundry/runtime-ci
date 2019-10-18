package resource

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"cloud.google.com/go/storage"
)

type GCSClient struct {
	client *storage.Client
}

func NewGCSClient() (GCSClient, error) {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		return GCSClient{}, fmt.Errorf("failed to create GCS storage client: %w", err)
	}

	return GCSClient{
		client: client,
	}, nil
}

func (gcsclient GCSClient) Get(bucketName string, objectPath string) ([]byte, error) {
	rc, err := gcsclient.client.Bucket(bucketName).Object(objectPath).NewReader(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get version info from GCS (%s, %s): %w", bucketName, objectPath, err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("failed to read version info from reader: %w", err)
	}

	return data, nil
}

func (gcsclient GCSClient) Put(bucketName string, objectPath string, contents []byte) error {
	f := bytes.NewReader(contents)

	wc := gcsclient.client.Bucket(bucketName).Object(objectPath).NewWriter(context.Background())
	_, err := io.Copy(wc, f)
	if err != nil {
		return fmt.Errorf("failed to write version info to writer: %w", err)
	}

	err = wc.Close()
	if err != nil {
		return fmt.Errorf("failed to close GCS writer: %w", err)
	}

	return nil
}
