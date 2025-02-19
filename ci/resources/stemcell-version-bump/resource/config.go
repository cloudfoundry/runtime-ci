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

type CheckInRequest struct {
	Source  Source
	Version Version
}

type OutRequest struct {
	Source Source
	Params OutParams
}

type OutParams struct {
	VersionFile string `json:"version_file"`
	TypeFile    string `json:"type_file"`
}

func NewCheckInRequest(in io.Reader) (CheckInRequest, error) {
	var resource CheckInRequest

	err := json.NewDecoder(in).Decode(&resource)
	if err != nil {
		return CheckInRequest{}, fmt.Errorf("decoding json: %w", err)
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
		return CheckInRequest{}, fmt.Errorf("missing required fields: '%s'", strings.Join(missingFields, "', '"))
	}

	return resource, nil
}

func NewOutRequest(in io.Reader) (OutRequest, error) {
	var resource OutRequest

	err := json.NewDecoder(in).Decode(&resource)
	if err != nil {
		return OutRequest{}, fmt.Errorf("decoding json: %w", err)
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

	if resource.Params.VersionFile == "" {
		missingFields = append(missingFields, "params.version_file")
	}

	if resource.Params.TypeFile == "" {
		missingFields = append(missingFields, "params.type_file")
	}

	if len(missingFields) > 0 {
		return OutRequest{}, fmt.Errorf("missing required fields: '%s'", strings.Join(missingFields, "', '"))
	}

	return resource, nil
}
