package concourseio

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Runner struct {
	In  Inputs
	Out Outputs

	UpdatedReleases []string
}

type Inputs struct {
	cfDeploymentPrevDir string
	cfDeploymentNextDir string
}

type Outputs struct {
	releaseListDir string
}

func NewRunner(buildDir string) (Runner, error) {
	inputs, err := setupInputs(buildDir)
	if err != nil {
		return Runner{}, err
	}

	outputs, err := setupOutputs(buildDir)
	if err != nil {
		return Runner{}, err
	}

	return Runner{In: inputs, Out: outputs}, nil
}

func setupInputs(buildDir string) (Inputs, error) {
	cfDeploymentPrevDir, err := buildSubDir(buildDir, "cf-deployment-prev")
	if err != nil {
		return Inputs{}, err
	}

	cfDeploymentNextDir, err := buildSubDir(buildDir, "cf-deployment-next")
	if err != nil {
		return Inputs{}, err
	}

	return Inputs{cfDeploymentPrevDir, cfDeploymentNextDir}, nil
}

func setupOutputs(buildDir string) (Outputs, error) {
	releaseListDir, err := buildSubDir(buildDir, "release-list")
	if err != nil {
		return Outputs{}, err
	}

	return Outputs{releaseListDir}, nil
}

func buildSubDir(buildDir, subDir string) (string, error) {
	dir := filepath.Join(buildDir, subDir)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("missing sub directory '%s' in build directory '%s'", subDir, buildDir)
	}

	return dir, nil
}

//go:generate counterfeiter . ReleaseDetector
type ReleaseDetector interface {
	Load(string, string) error
	DetectUpdatedReleases()
	Write() ([]byte, error)
}

func (r Runner) Run(detector ReleaseDetector) error {
	err := detector.Load(
		filepath.Join(r.In.cfDeploymentPrevDir, "cf-deployment.yml"),
		filepath.Join(r.In.cfDeploymentNextDir, "cf-deployment.yml"),
	)
	if err != nil {
		return err
	}

	detector.DetectUpdatedReleases()
	releaseList, err := detector.Write()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(r.Out.releaseListDir, "releases.yml"), releaseList, 0644)
	if err != nil {
		return err
	}

	return nil
}
