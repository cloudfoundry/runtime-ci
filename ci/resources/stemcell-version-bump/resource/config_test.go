package resource_test

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	"stemcell-version-bump/resource"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCheckInRequest(t *testing.T) {
	type checkNewResourceFunc func(*testing.T, resource.CheckInRequest, error)

	checks := func(cs ...checkNewResourceFunc) []checkNewResourceFunc { return cs }

	expectResource := func(expectedResource resource.CheckInRequest) checkNewResourceFunc {
		return func(t *testing.T, actualResource resource.CheckInRequest, _ error) {
			t.Helper()

			assert.Equal(t, expectedResource, actualResource)
		}
	}

	expectNoError := func(t *testing.T, _ resource.CheckInRequest, actualErr error) {
		t.Helper()

		require.NoError(t, actualErr)
	}

	expectError := func(expectedErr string) checkNewResourceFunc {
		return func(t *testing.T, _ resource.CheckInRequest, actualErr error) {
			t.Helper()

			assert.EqualError(t, actualErr, expectedErr)
		}
	}

	expectWrappedError := func(expectedOuter string, expectedInner error) checkNewResourceFunc {
		return func(t *testing.T, _ resource.CheckInRequest, actualErr error) {
			t.Helper()

			require.Error(t, actualErr)

			assert.Contains(t, actualErr.Error(), expectedOuter)

			actualInner := errors.Unwrap(actualErr)
			require.Error(t, actualInner)

			assert.IsType(t, expectedInner, actualInner)
		}
	}

	tests := []struct {
		name   string
		inArg  io.Reader
		checks []checkNewResourceFunc
	}{
		{
			"simple happy case",
			strings.NewReader(`{
				"source": {
					"json_key": "some-json-key",
					"bucket_name": "some-bucket-name",
          "file_name": "some-file-name"
				},
				"version": {
					"type": "some-current-type",
					"version": "some-current-version"
				}
			}`),
			checks(
				expectNoError,
				expectResource(
					resource.CheckInRequest{
						Source: resource.Source{
							JSONKey:    "some-json-key",
							BucketName: "some-bucket-name",
							FileName:   "some-file-name",
						},
						Version: resource.Version{
							Type:    "some-current-type",
							Version: "some-current-version",
						},
					},
				),
			),
		},

		{
			"invalid json provided",
			strings.NewReader(`%%%`),
			checks(
				expectWrappedError("decoding json", new(json.SyntaxError)),
			),
		},

		{
			"missing single required field",
			strings.NewReader(`{
				"source": {
					"bucket_name": "some-bucket-name",
          "file_name": "some-file-name"
				},
				"version": {
					"type": "some-current-type",
					"version": "some-current-version"
				}
			}`),
			checks(
				expectError("missing required fields: 'json_key'"),
			),
		},

		{
			"missing multiple required fields",
			strings.NewReader(`{}`),
			checks(
				expectError("missing required fields: 'json_key', 'bucket_name', 'file_name'"),
			),
		},
	}

	for _, test := range tests {
		arg, checks := test.inArg, test.checks
		t.Run(test.name, func(t *testing.T) {
			actualOutput, actualErr := resource.NewCheckInRequest(arg)

			for _, check := range checks {
				check(t, actualOutput, actualErr)
			}
		})
	}
}

func TestNewOutputCheckInRequest(t *testing.T) {
	type checkNewResourceFunc func(*testing.T, resource.OutRequest, error)

	checks := func(cs ...checkNewResourceFunc) []checkNewResourceFunc { return cs }

	expectResource := func(expectedResource resource.OutRequest) checkNewResourceFunc {
		return func(t *testing.T, actualResource resource.OutRequest, _ error) {
			t.Helper()

			assert.Equal(t, expectedResource, actualResource)
		}
	}
	expectNoError := func(t *testing.T, _ resource.OutRequest, actualErr error) {
		t.Helper()

		require.NoError(t, actualErr)
	}
	expectError := func(expectedErr string) checkNewResourceFunc {
		return func(t *testing.T, _ resource.OutRequest, actualErr error) {
			t.Helper()

			assert.EqualError(t, actualErr, expectedErr)
		}
	}
	expectWrappedError := func(expectedOuter string, expectedInner error) checkNewResourceFunc {
		return func(t *testing.T, _ resource.OutRequest, actualErr error) {
			t.Helper()

			require.Error(t, actualErr)

			assert.Contains(t, actualErr.Error(), expectedOuter)

			actualInner := errors.Unwrap(actualErr)
			require.Error(t, actualInner)

			assert.IsType(t, expectedInner, actualInner)
		}
	}

	tests := []struct {
		name   string
		inArg  io.Reader
		checks []checkNewResourceFunc
	}{
		{
			"simple happy case",
			strings.NewReader(`{
				"source": {
					"json_key": "some-json-key",
					"bucket_name": "some-bucket-name",
          "file_name": "some-file-name"
				},
				"params": {
					"type_file": "some-type-file-path",
					"version_file": "some-version-file-path"
				}
			}`),
			checks(
				expectNoError,
				expectResource(
					resource.OutRequest{
						Source: resource.Source{
							JSONKey:    "some-json-key",
							BucketName: "some-bucket-name",
							FileName:   "some-file-name",
						},
						Params: resource.OutParams{
							VersionFile: "some-version-file-path",
							TypeFile:    "some-type-file-path",
						},
					},
				),
			),
		},

		{
			"invalid json provided",
			strings.NewReader(`%%%`),
			checks(
				expectWrappedError("decoding json", new(json.SyntaxError)),
			),
		},

		{
			"missing single required field",
			strings.NewReader(`{
				"source": {
					"bucket_name": "some-bucket-name",
          "file_name": "some-file-name",
					"json_key": "some-json-key"
				},
				"params": {
					"version_file": "some-version-file-path"
				}
			}`),
			checks(
				expectError("missing required fields: 'params.type_file'"),
			),
		},

		{
			"missing multiple required fields",
			strings.NewReader(`{}`),
			checks(
				expectError("missing required fields: 'json_key', 'bucket_name', 'file_name', 'params.version_file', 'params.type_file'"),
			),
		},
	}

	for _, test := range tests {
		arg, checks := test.inArg, test.checks
		t.Run(test.name, func(t *testing.T) {
			actualOutput, actualErr := resource.NewOutRequest(arg)

			for _, check := range checks {
				check(t, actualOutput, actualErr)
			}
		})
	}
}
