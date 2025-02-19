package runner_test

import (
	"errors"
	"testing"

	"stemcell-version-bump/cmd/check/runner"
	"stemcell-version-bump/cmd/check/runner/runnerfakes"
	"stemcell-version-bump/resource"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheck(t *testing.T) {
	var returner = func(output string, err error) runner.Getter {
		getter := new(runnerfakes.FakeGetter)

		getter.GetReturns([]byte(output), err)

		return getter
	}

	type checkCheckFunc func(*testing.T, string, error)

	checks := func(cs ...checkCheckFunc) []checkCheckFunc { return cs }
	expectNoError := func(t *testing.T, _ string, actualErr error) {
		t.Helper()

		require.NoError(t, actualErr)
	}
	expectError := func(expectedErr string) checkCheckFunc {
		return func(t *testing.T, _ string, actualErr error) {
			t.Helper()

			assert.EqualError(t, actualErr, expectedErr)
		}
	}
	expectVersionList := func(expectedVersions string) checkCheckFunc {
		return func(t *testing.T, actualOutput string, _ error) {
			t.Helper()

			assert.JSONEq(t, expectedVersions, actualOutput)
		}
	}

	type in struct {
		request resource.CheckInRequest
		getter  runner.Getter
	}

	type testcase struct {
		name   string
		inArg  in
		checks []checkCheckFunc
	}

	tests := []testcase{
		{
			"initial check, no previous version",
			in{
				request: resource.CheckInRequest{
					Source: resource.Source{
						TypeFilter: "minor",
					},
					Version: resource.Version{
						Type:    "",
						Version: "",
					},
				},
				getter: returner(`{"type": "minor", "version": "1.2"}`, nil),
			},
			checks(
				expectNoError,
				expectVersionList(`[{"type": "minor", "version": "1.2"}]`),
			),
		},

		{
			"check with version bump that matches the type filter",
			in{
				request: resource.CheckInRequest{
					Source: resource.Source{
						TypeFilter: "minor",
					},
					Version: resource.Version{
						Type:    "minor",
						Version: "1.1",
					},
				},
				getter: returner(`{"type": "minor", "version": "1.2"}`, nil),
			},
			checks(
				expectNoError,
				expectVersionList(`[{"type": "minor", "version": "1.2"}]`),
			),
		},

		{
			"check with version bump that does not match the type filter",
			in{
				request: resource.CheckInRequest{
					Source: resource.Source{
						TypeFilter: "minor",
					},
					Version: resource.Version{
						Type:    "minor",
						Version: "1.1",
					},
				},
				getter: returner(`{"type": "major", "version": "2.1"}`, nil),
			},
			checks(
				expectNoError,
				expectVersionList(`[]`),
			),
		},

		{
			"check with no version bump",
			in{
				request: resource.CheckInRequest{
					Version: resource.Version{
						Type:    "minor",
						Version: "1.1",
					},
				},
				getter: returner(`{"type": "minor", "version": "1.1"}`, nil),
			},
			checks(
				expectNoError,
				expectVersionList(`[]`),
			),
		},

		{
			"fail to get version info",
			in{
				request: resource.CheckInRequest{
					Source: resource.Source{
						BucketName: "some-bucket",
						FileName:   "path/to/file",
					},
					Version: resource.Version{
						Type:    "minor",
						Version: "1.1",
					},
				},
				getter: returner(``, errors.New("failed-to-get")),
			},
			checks(
				expectError("failed to fetch version info from bucket/file (some-bucket, path/to/file): failed-to-get"),
			),
		},

		{
			"bad version info",
			in{
				request: resource.CheckInRequest{
					Source: resource.Source{
						BucketName: "some-bucket",
						FileName:   "path/to/file",
					},
					Version: resource.Version{
						Type:    "minor",
						Version: "1.1",
					},
				},
				getter: returner(`%%%`, nil),
			},
			checks(
				expectError("failed to unmarshal version info file: invalid character '%' looking for beginning of value"),
			),
		},
	}

	for _, test := range tests {
		arg, checks := test.inArg, test.checks
		t.Run(test.name, func(t *testing.T) {
			actualOutput, actualErr := runner.Check(arg.request, arg.getter)

			for _, check := range checks {
				check(t, actualOutput, actualErr)
			}
		})
	}
}
