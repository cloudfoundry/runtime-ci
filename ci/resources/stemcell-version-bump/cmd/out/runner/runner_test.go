package runner_test

import (
	"errors"
	"io/ioutil"
	"os"
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

func TestReadVersionBump(t *testing.T) {
	type checkReadVersionBumpFunc func(*testing.T, []byte, error)
	checks := func(cs ...checkReadVersionBumpFunc) []checkReadVersionBumpFunc { return cs }

	var expectNoError = func(t *testing.T, _ []byte, actualErr error) {
		if !assert.NoError(t, actualErr) {
			t.FailNow()
		}
	}

	var expectError = func(expectedErr string) checkReadVersionBumpFunc {
		return func(t *testing.T, _ []byte, actualErr error) {
			assert.EqualError(t, actualErr, expectedErr)
		}
	}

	_ = expectError

	var expectWrappedError = func(expectedOuter string, expectedInner error) checkReadVersionBumpFunc {
		return func(t *testing.T, _ []byte, actualErr error) {
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

	var expectStemcellBumpTypeContent = func(expectedContent []byte) checkReadVersionBumpFunc {
		return func(t *testing.T, actualContent []byte, _ error) {
			assert.JSONEq(t, string(expectedContent), string(actualContent))
		}
	}

	type setup struct {
		versionPath     string
		versionContent  []byte
		bumpTypePath    string
		bumpTypeContent []byte
	}

	type in struct {
		config resource.OutRequest
	}

	type testcase struct {
		name   string
		setup  setup
		inArg  in
		checks []checkReadVersionBumpFunc
	}
	tests := []testcase{
		testcase{
			"happy path, read succeeds",
			setup{
				versionPath:     "version-file",
				versionContent:  []byte(`some-version`),
				bumpTypePath:    "bump-type-file",
				bumpTypeContent: []byte(`minor`),
			},
			in{
				config: resource.OutRequest{
					Source: resource.Source{
						BucketName: "some-bucket",
						FileName:   "path/to/file",
					},
					Params: resource.OutParams{
						VersionFile: "version-file",
						TypeFile:    "bump-type-file",
					},
				},
			},
			checks(
				expectNoError,
				expectStemcellBumpTypeContent([]byte(`{"version": "some-version", "type": "minor"}`)),
			),
		},

		testcase{
			"fail to find required version file",
			setup{},
			in{
				config: resource.OutRequest{
					Source: resource.Source{
						BucketName: "some-bucket",
						FileName:   "path/to/file",
					},
					Params: resource.OutParams{
						VersionFile: "missing-version-file",
						TypeFile:    "bump-type-file",
					},
				},
			},
			checks(
				expectWrappedError("reading version file:", new(os.PathError)),
			),
		},

		testcase{
			"fail to find required bump type file",
			setup{
				versionPath: "version-file",
			},
			in{
				config: resource.OutRequest{
					Source: resource.Source{
						BucketName: "some-bucket",
						FileName:   "path/to/file",
					},
					Params: resource.OutParams{
						VersionFile: "version-file",
						TypeFile:    "missing-bump-type-file",
					},
				},
			},
			checks(
				expectWrappedError("reading bump type file:", new(os.PathError)),
			),
		},

		testcase{
			"fail on invalid bump type",
			setup{
				versionPath:     "version-file",
				versionContent:  []byte(`some-version`),
				bumpTypePath:    "some-bad-bump-type-file",
				bumpTypeContent: []byte(`some-bad-bump-type`),
			},
			in{
				config: resource.OutRequest{
					Source: resource.Source{
						BucketName: "some-bucket",
						FileName:   "path/to/file",
					},
					Params: resource.OutParams{
						VersionFile: "version-file",
						TypeFile:    "some-bad-bump-type-file",
					},
				},
			},
			checks(
				expectError(`invalid bump type: "some-bad-bump-type"`),
			),
		},
	}

	for _, test := range tests {
		setup, arg, checks := test.setup, test.inArg, test.checks
		t.Run(test.name, func(t *testing.T) {
			tmpDir, err := ioutil.TempDir("", test.name+"-")
			if err != nil {
				t.Fatalf("failed setup: %v", err)
			}
			err = os.Chdir(tmpDir)
			if err != nil {
				t.Fatalf("failed setup: %v", err)
			}

			if setup.versionPath != "" {
				err := ioutil.WriteFile(setup.versionPath, setup.versionContent, 0755)
				if err != nil {
					t.Fatalf("failed setup: %v", err)
				}
			}

			if setup.bumpTypePath != "" {
				err := ioutil.WriteFile(setup.bumpTypePath, setup.bumpTypeContent, 0755)
				if err != nil {
					t.Fatalf("failed setup: %v", err)
				}
			}

			actualContent, actualErr := runner.ReadVersionBump(arg.config)

			for _, check := range checks {
				check(t, actualContent, actualErr)
			}

			err = os.RemoveAll(tmpDir)
			if err != nil {
				t.Fatalf("failed cleanup: %v", err)
			}
		})
	}
}
