package resource_test

import (
	"io/ioutil"
	"os"
	"stemcell-version-bump/resource"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupEnvironment(t *testing.T) {
	type checkSetupEnvironmentFunc func(*testing.T, error)
	checks := func(cs ...checkSetupEnvironmentFunc) []checkSetupEnvironmentFunc { return cs }

	var expectNoError = func(t *testing.T, actualErr error) {
		if !assert.NoError(t, actualErr) {
			t.FailNow()
		}
	}

	var expectJsonKeyWasConfigured = func(expectedJSONKey string) checkSetupEnvironmentFunc {
		return func(t *testing.T, actualErr error) {
			value := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
			content, err := ioutil.ReadFile(value)
			if !assert.NoError(t, err) {
				t.FailNow()
			}

			assert.Equal(t, string(content), expectedJSONKey)
		}
	}

	type in struct {
		jsonKey string
	}

	type testcase struct {
		name   string
		inArg  in
		checks []checkSetupEnvironmentFunc
	}
	tests := []testcase{
		testcase{
			"it sets the GCS key",
			in{
				jsonKey: "somekey",
			},
			checks(
				expectNoError,
				expectJsonKeyWasConfigured("somekey"),
			),
		},
	}

	for _, test := range tests {
		arg, checks := test.inArg, test.checks
		t.Run(test.name, func(t *testing.T) {
			actualErr := resource.SetupEnvironment(arg.jsonKey)

			for _, check := range checks {
				check(t, actualErr)
			}
		})
	}
}
