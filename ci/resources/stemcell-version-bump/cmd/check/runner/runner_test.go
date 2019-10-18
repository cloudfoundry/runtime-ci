package runner_test

import (
	"errors"
	"testing"

	"stemcell-version-bump/cmd/check/runner"
	"stemcell-version-bump/cmd/check/runner/runnerfakes"
	"stemcell-version-bump/resource"

	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	var returner = func(output string, err error) runner.Getter {
		getter := new(runnerfakes.FakeGetter)

		getter.GetReturns([]byte(output), err)

		return getter
	}

	type checkCheckFunc func(*testing.T, string, error)
	checks := func(cs ...checkCheckFunc) []checkCheckFunc { return cs }

	var expectNoError = func(t *testing.T, _ string, actualErr error) {
		if !assert.NoError(t, actualErr) {
			t.FailNow()
		}
	}

	var expectError = func(expectedErr string) checkCheckFunc {
		return func(t *testing.T, _ string, actualErr error) {
			assert.EqualError(t, actualErr, expectedErr)
		}
	}

	var expectVersionList = func(expectedVersions string) checkCheckFunc {
		return func(t *testing.T, actualOutput string, actualErr error) {
			assert.JSONEq(t, expectedVersions, actualOutput)
		}
	}

	type in struct {
		config resource.Config
		getter runner.Getter
	}

	type testcase struct {
		name   string
		inArg  in
		checks []checkCheckFunc
	}
	tests := []testcase{
		testcase{
			"initial check, no previous version",
			in{
				config: resource.Config{
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

		testcase{
			"check with version bump that matches the type filter",
			in{
				config: resource.Config{
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

		testcase{
			"check with version bump that does not match the type filter",
			in{
				config: resource.Config{
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

		testcase{
			"check with no version bump",
			in{
				config: resource.Config{
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

		testcase{
			"fail to get version info",
			in{
				config: resource.Config{
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
				config: resource.Config{
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
			actualOutput, actualErr := runner.Check(arg.config, arg.getter)

			for _, check := range checks {
				check(t, actualOutput, actualErr)
			}
		})
	}
}
