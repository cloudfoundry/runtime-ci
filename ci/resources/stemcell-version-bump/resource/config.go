package resource

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Source struct {
	JSONKey    string `json:"json_key"`
	BucketName string `json:"bucket_name"`
	FileName   string `json:"file_name"`
	TypeFilter string `json:"type_filter"`
}

type Version struct {
	Type    string `json:"type"`
	Version string `json:"version"`
}

type Config struct {
	Source  Source
	Version Version
}

func NewConfig(in io.Reader) (Config, error) {
	var resource Config
	err := json.NewDecoder(in).Decode(&resource)
	if err != nil {
		return Config{}, fmt.Errorf("decoding json: %w", err)
	}

	var missingFields []string

	if resource.Source.JSONKey == "" {
		missingFields = append(missingFields, "json_key")
	}

	if resource.Source.BucketName == "" {
		missingFields = append(missingFields, "bucket_name")
	}

	if resource.Source.FileName == "" {
		missingFields = append(missingFields, "file_name")
	}

	if len(missingFields) > 0 {
		return Config{}, fmt.Errorf("missing required fields: '%s'", strings.Join(missingFields, "', '"))
	}

	return resource, nil
}
