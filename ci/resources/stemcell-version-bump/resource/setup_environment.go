package resource

import (
	"fmt"
	"io/ioutil"
	"os"
)

func SetupEnvironment(jsonKey string) error {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return fmt.Errorf("failed to create temp file for GCS credentials: %w", err)
	}
	credPath := f.Name()
	defer f.Close()

	_, err = f.Write([]byte(jsonKey))
	if err != nil {
		return fmt.Errorf("failed to write JSON key to temp file: %w", err)
	}

	err = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credPath)
	if err != nil {
		return fmt.Errorf("failed to set GOOGLE_APPLICATION_CREDENTIALS env var: %w", err)
	}

	return nil
}
