package runner_test

import (
	"errors"
	"os"
	"testing"

	"stemcell-version-bump/cmd/out/runner"
	"stemcell-version-bump/cmd/out/runner/runnerfakes"
	"stemcell-version-bump/resource"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadVersion(t *testing.T) {
	var returner = func(err error) runner.Putter {
		putter := new(runnerfakes.FakePutter)

		putter.PutReturns(err)

		return putter
	}

	type checkUploadVersionFunc func(*testing.T, runner.Putter, error)

	checks := func(cs ...checkUploadVersionFunc) []checkUploadVersionFunc { return cs }

	expectNoError := func(t *testing.T, _ runner.Putter, actualErr error) {
		t.Helper()

		require.NoError(t, actualErr)
	}
	expectError := func(expectedErr string) checkUploadVersionFunc {
		return func(t *testing.T, _ runner.Putter, actualErr error) {
			t.Helper()

			assert.EqualError(t, actualErr, expectedErr)
		}
	}
	expectPutArgs := func(bucketName, fileName, output string) checkUploadVersionFunc {
		return func(t *testing.T, fakePutter runner.Putter, _ error) {
			t.Helper()

			fake, ok := fakePutter.(*runnerfakes.FakePutter)
			require.Truef(t, ok, "expected %T to be of type '*runnerfakes.FakePutter'", fakePutter)

			actualBucket, actualFileName, actualUploadVersionput := fake.PutArgsForCall(0)

			assert.Equal(t, bucketName, actualBucket)
			assert.Equal(t, fileName, actualFileName)
			assert.JSONEq(t, output, string(actualUploadVersionput))
		}
	}

	type in struct {
		request     resource.OutRequest
		putter      runner.Putter
		versionBump resource.Version
	}

	type testcase struct {
		name   string
		inArg  in
		checks []checkUploadVersionFunc
	}

	tests := []testcase{
		{
			"happy path, post succeeds",
			in{
				request: resource.OutRequest{
					Source: resource.Source{
						BucketName: "some-bucket",
						FileName:   "path/to/file",
					},
				},
				putter:      returner(nil),
				versionBump: resource.Version{Version: "some-version", Type: "minor"},
			},
			checks(
				expectNoError,
				expectPutArgs("some-bucket", "path/to/file", `{"version":"some-version","type":"minor"}`),
			),
		},

		{
			"fail to put resource",
			in{
				request: resource.OutRequest{
					Source: resource.Source{
						BucketName: "some-bucket",
						FileName:   "path/to/file",
					},
				},
				putter:      returner(errors.New("failed-to-put")),
				versionBump: resource.Version{Version: "some-version", Type: "minor"},
			},
			checks(
				expectError("updating version info in bucket/file (some-bucket, path/to/file): failed-to-put"),
			),
		},
	}

	for _, test := range tests {
		arg, checks := test.inArg, test.checks
		t.Run(test.name, func(t *testing.T) {
			actualErr := runner.UploadVersion(arg.request, arg.putter, arg.versionBump)

			for _, check := range checks {
				check(t, arg.putter, actualErr)
			}
		})
	}
}

func TestGenerateResourceOutput(t *testing.T) {
	type checkGenerateResourceOutputFunc func(*testing.T, string, error)

	checks := func(cs ...checkGenerateResourceOutputFunc) []checkGenerateResourceOutputFunc { return cs }

	expectNoError := func(t *testing.T, _ string, actualErr error) {
		t.Helper()

		require.NoError(t, actualErr)
	}
	expectGenerateResourceOutputput := func(output string) checkGenerateResourceOutputFunc {
		return func(t *testing.T, actualGenerateResourceOutputput string, _ error) {
			t.Helper()

			assert.JSONEq(t, output, actualGenerateResourceOutputput)
		}
	}

	type in struct {
		version resource.Version
	}

	type testcase struct {
		name   string
		inArg  in
		checks []checkGenerateResourceOutputFunc
	}

	tests := []testcase{
		{
			"happy path, output generation succeeds",
			in{
				version: resource.Version{Version: "some-version", Type: "minor"},
			},
			checks(
				expectNoError,
				expectGenerateResourceOutputput(`{"version":{"version":"some-version","type":"minor"}}`),
			),
		},
	}

	for _, test := range tests {
		arg, checks := test.inArg, test.checks
		t.Run(test.name, func(t *testing.T) {
			actualOutput, actualErr := runner.GenerateResourceOutput(arg.version)

			for _, check := range checks {
				check(t, actualOutput, actualErr)
			}
		})
	}
}

func TestNewVersion(t *testing.T) {
	type checkNewVersionFunc func(*testing.T, resource.Version, error)

	checks := func(cs ...checkNewVersionFunc) []checkNewVersionFunc { return cs }
	expectNoError := func(t *testing.T, _ resource.Version, actualErr error) {
		t.Helper()

		require.NoError(t, actualErr)
	}
	expectError := func(expectedErr string) checkNewVersionFunc {
		return func(t *testing.T, _ resource.Version, actualErr error) {
			t.Helper()

			assert.EqualError(t, actualErr, expectedErr)
		}
	}
	expectWrappedError := func(expectedOuter string, expectedInner error) checkNewVersionFunc {
		return func(t *testing.T, _ resource.Version, actualErr error) {
			t.Helper()

			require.Error(t, actualErr)

			assert.Contains(t, actualErr.Error(), expectedOuter)

			actualInner := errors.Unwrap(actualErr)
			require.Error(t, actualInner)

			assert.IsType(t, expectedInner, actualInner)
		}
	}
	expectStemcellBumpTypeContent := func(expectedVersion resource.Version) checkNewVersionFunc {
		return func(t *testing.T, actualVersion resource.Version, _ error) {
			t.Helper()

			assert.Equal(t, expectedVersion, actualVersion)
		}
	}

	type setup struct {
		versionPath     string
		versionContent  []byte
		bumpTypePath    string
		bumpTypeContent []byte
	}

	type in struct {
		request resource.OutRequest
	}

	type testcase struct {
		name   string
		setup  setup
		inArg  in
		checks []checkNewVersionFunc
	}

	tests := []testcase{
		{
			"happy path, read succeeds",
			setup{
				versionPath:     "version-file",
				versionContent:  []byte(`some-version`),
				bumpTypePath:    "bump-type-file",
				bumpTypeContent: []byte(`minor`),
			},
			in{
				request: resource.OutRequest{
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
				expectStemcellBumpTypeContent(resource.Version{Version: "some-version", Type: "minor"}),
			),
		},

		{
			"fail to read required version file",
			setup{},
			in{
				request: resource.OutRequest{
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

		{
			"fail to read required bump type file",
			setup{
				versionPath: "version-file",
			},
			in{
				request: resource.OutRequest{
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

		{
			"fail on invalid bump type",
			setup{
				versionPath:     "version-file",
				versionContent:  []byte(`some-version`),
				bumpTypePath:    "some-bad-bump-type-file",
				bumpTypeContent: []byte(`some-bad-bump-type`),
			},
			in{
				request: resource.OutRequest{
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
			tmpDir := t.TempDir()

			err := os.Chdir(tmpDir)
			if err != nil {
				t.Fatalf("failed setup: %v", err)
			}

			if setup.versionPath != "" {
				err := os.WriteFile(setup.versionPath, setup.versionContent, 0755)
				if err != nil {
					t.Fatalf("failed setup: %v", err)
				}
			}

			if setup.bumpTypePath != "" {
				err := os.WriteFile(setup.bumpTypePath, setup.bumpTypeContent, 0755)
				if err != nil {
					t.Fatalf("failed setup: %v", err)
				}
			}

			actualVersion, actualErr := runner.NewVersion(arg.request)

			for _, check := range checks {
				check(t, actualVersion, actualErr)
			}

			err = os.RemoveAll(tmpDir)
			if err != nil {
				t.Fatalf("failed cleanup: %v", err)
			}
		})
	}
}
