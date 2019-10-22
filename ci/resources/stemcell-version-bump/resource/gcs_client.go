package resource

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCSClient struct {
	client *storage.Client
	ctx    context.Context
}

func NewGCSClient(jsonKey string) (GCSClient, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(jsonKey)))
	if err != nil {
		return GCSClient{}, fmt.Errorf("failed to create GCS storage client: %w", err)
	}

	return GCSClient{
		client: client,
		ctx:    ctx,
	}, nil
}

func (gcsclient GCSClient) Get(bucketName string, objectPath string) ([]byte, error) {
	rc, err := gcsclient.client.Bucket(bucketName).Object(objectPath).NewReader(gcsclient.ctx)
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

	wc := gcsclient.client.Bucket(bucketName).Object(objectPath).NewWriter(gcsclient.ctx)
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
