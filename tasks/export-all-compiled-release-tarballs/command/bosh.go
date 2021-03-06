package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os/exec"
	"strings"
)

type BoshCLI struct{}

func (cli BoshCLI) Cmd(name string, args ...string) (io.Reader, error) {
	boshArgs := append([]string{name}, args...)
	cmd := exec.Command("bosh", boshArgs...)

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf

	err := cmd.Run()
	if err != nil {
		return nil, parseErr(outBuf, err)
	}

	return outBuf, nil
}

func parseErr(r io.Reader, runErr error) error {
	var output struct {
		Blocks []string
		Lines  []string
	}

	err := json.NewDecoder(r).Decode(&output)
	if err != nil {
		return err
	}

	if len(output.Blocks) > 0 {
		for _, block := range output.Blocks {
			if strings.HasPrefix(block, "Error:") {
				return errors.New(block)
			}
		}
	} else {
		var errLines []string
		for _, line := range output.Lines {
			if strings.HasPrefix(line, "Using environment") {
				continue
			}
			if strings.HasPrefix(line, "Exit code") {
				continue
			}

			errLines = append(errLines, line)
		}
		return errors.New(strings.Join(errLines, "\n"))
	}

	return runErr
}
