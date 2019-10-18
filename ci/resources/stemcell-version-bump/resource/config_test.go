package resource_test

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	"stemcell-version-bump/resource"

	"github.com/stretchr/testify/assert"
)

func TestNewResource(t *testing.T) {
	type checkNewResourceFunc func(*testing.T, resource.Config, error)
	checks := func(cs ...checkNewResourceFunc) []checkNewResourceFunc { return cs }

	var expectResource = func(expectedResource resource.Config) checkNewResourceFunc {
		return func(t *testing.T, actualResource resource.Config, _ error) {
			assert.Equal(t, expectedResource, actualResource)
		}
	}

	var expectNoError = func(t *testing.T, _ resource.Config, actualErr error) {
		if !assert.NoError(t, actualErr) {
			t.FailNow()
		}
	}

	var expectError = func(expectedErr string) checkNewResourceFunc {
		return func(t *testing.T, _ resource.Config, actualErr error) {
			assert.EqualError(t, actualErr, expectedErr)
		}
	}

	var expectWrappedError = func(expectedOuter string, expectedInner error) checkNewResourceFunc {
		return func(t *testing.T, _ resource.Config, actualErr error) {
			if !assert.Error(t, actualErr) {
				t.FailNow()
			}

			assert.Contains(t, actualErr.Error(), expectedOuter)

			actualInner := errors.Unwrap(actualErr)
			if !assert.Error(t, actualInner) {
				t.FailNow()
			}

			assert.IsType(t, expectedInner, actualInner)
		}
	}

	type testcase struct {
		name   string
		inArg  io.Reader
		checks []checkNewResourceFunc
	}
	tests := []testcase{
		testcase{
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
					resource.Config{
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

		testcase{
			"invalid json provided",
			strings.NewReader(`%%%`),
			checks(
				expectWrappedError("decoding json", new(json.SyntaxError)),
			),
		},

		testcase{
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

		testcase{
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
			actualOutput, actualErr := resource.NewConfig(arg)

			for _, check := range checks {
				check(t, actualOutput, actualErr)
			}
		})
	}
}
