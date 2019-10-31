package runner_test

import (
	"errors"
	"testing"

	"stemcell-version-bump/cmd/in/runner"
	"stemcell-version-bump/cmd/in/runner/runnerfakes"
	"stemcell-version-bump/resource"

	"github.com/stretchr/testify/assert"
)

func TestIn(t *testing.T) {
	var returner = func(output string, err error) runner.Getter {
		getter := new(runnerfakes.FakeGetter)

		getter.GetReturns([]byte(output), err)

		return getter
	}

	type checkInFunc func(*testing.T, string, error)
	checks := func(cs ...checkInFunc) []checkInFunc { return cs }

	var expectNoError = func(t *testing.T, _ string, actualErr error) {
		if !assert.NoError(t, actualErr) {
			t.FailNow()
		}
	}

	var expectError = func(expectedErr string) checkInFunc {
		return func(t *testing.T, _ string, actualErr error) {
			assert.EqualError(t, actualErr, expectedErr)
		}
	}

	var expectResource = func(expectedResourceJSON string) checkInFunc {
		return func(t *testing.T, actualOutput string, actualErr error) {
			assert.JSONEq(t, expectedResourceJSON, actualOutput)
		}
	}

	type in struct {
		request resource.CheckInRequest
		getter  runner.Getter
	}

	type testcase struct {
		name   string
		inArg  in
		checks []checkInFunc
	}
	tests := []testcase{
		testcase{
			"happy path, fetch succeeds",
			in{
				request: resource.CheckInRequest{
					Source: resource.Source{
						TypeFilter: "minor",
					},
					Version: resource.Version{
						Type:    "minor",
						Version: "1.2",
					},
				},
				getter: returner(`{"type": "minor", "version": "1.2"}`, nil),
			},
			checks(
				expectNoError,
				expectResource(`{"version": {"type": "minor", "version": "1.2"}}`),
			),
		},

		testcase{
			"fail to fetch resource",
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

		testcase{
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

		testcase{
			"old version request",
			in{
				request: resource.CheckInRequest{
					Source: resource.Source{
						BucketName: "some-bucket",
						FileName:   "path/to/file",
					},
					Version: resource.Version{
						Type:    "major",
						Version: "1.1",
					},
				},
				getter: returner(`{"type": "minor", "version": "1.2"}`, nil),
			},
			checks(
				expectError("failed to retrieve specified version: requested {major 1.1}, found {minor 1.2}"),
			),
		},
	}

	for _, test := range tests {
		arg, checks := test.inArg, test.checks
		t.Run(test.name, func(t *testing.T) {
			actualOutput, actualErr := runner.In(arg.request, arg.getter)

			for _, check := range checks {
				check(t, actualOutput, actualErr)
			}
		})
	}
}
