package runner_test

import (
	"errors"
	"testing"

	"stemcell-version-bump/cmd/out/runner"
	"stemcell-version-bump/cmd/out/runner/runnerfakes"
	"stemcell-version-bump/resource"

	"github.com/stretchr/testify/assert"
)

func TestOut(t *testing.T) {
	var returner = func(err error) runner.Putter {
		putter := new(runnerfakes.FakePutter)

		putter.PutReturns(err)

		return putter
	}

	type checkOutFunc func(*testing.T, runner.Putter, error)
	checks := func(cs ...checkOutFunc) []checkOutFunc { return cs }

	var expectNoError = func(t *testing.T, _ runner.Putter, actualErr error) {
		if !assert.NoError(t, actualErr) {
			t.FailNow()
		}
	}

	var expectError = func(expectedErr string) checkOutFunc {
		return func(t *testing.T, _ runner.Putter, actualErr error) {
			assert.EqualError(t, actualErr, expectedErr)
		}
	}

	var expectPutArgs = func(bucketName, fileName, output string) checkOutFunc {
		return func(t *testing.T, fakePutter runner.Putter, actualErr error) {
			fake, ok := fakePutter.(*runnerfakes.FakePutter)
			if !assert.Truef(t, ok, "expected %T to be of type '*runnerfakes.FakePutter'", fakePutter) {
				t.FailNow()
			}

			actualBucket, actualFileName, actualOutput := fake.PutArgsForCall(0)

			assert.Equal(t, bucketName, actualBucket)
			assert.Equal(t, fileName, actualFileName)
			assert.Equal(t, output, string(actualOutput))
		}
	}

	type in struct {
		config  resource.OutRequest
		putter  runner.Putter
		content []byte
	}

	type testcase struct {
		name   string
		inArg  in
		checks []checkOutFunc
	}
	tests := []testcase{
		testcase{
			"happy path, post succeeds",
			in{
				config: resource.OutRequest{
					Source: resource.Source{
						BucketName: "some-bucket",
						FileName:   "path/to/file",
					},
				},
				putter:  returner(nil),
				content: []byte("my-version-file"),
			},
			checks(
				expectNoError,
				expectPutArgs("some-bucket", "path/to/file", "my-version-file"),
			),
		},

		testcase{
			"fail to fetch resource",
			in{
				config: resource.OutRequest{
					Source: resource.Source{
						BucketName: "some-bucket",
						FileName:   "path/to/file",
					},
				},
				putter:  returner(errors.New("failed-to-put")),
				content: []byte("my-version-file"),
			},
			checks(
				expectError("updating version info in bucket/file (some-bucket, path/to/file): failed-to-put"),
			),
		},
	}

	for _, test := range tests {
		arg, checks := test.inArg, test.checks
		t.Run(test.name, func(t *testing.T) {
			_ = arg
			actualErr := runner.Out(arg.config, arg.putter, arg.content)

			for _, check := range checks {
				check(t, arg.putter, actualErr)
			}
		})
	}
}
