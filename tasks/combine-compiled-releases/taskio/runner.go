package taskio

import "github.com/cloudfoundry/runtime-ci/task-libs/taskfs"

type Runner struct {
	In  Inputs
	Out Outputs
}

type Inputs struct {
	compiledReleasesPrevDir string
	compiledReleasesNextDir string
	cfDeploymentNextDir     string
}

type Outputs struct {
	compliledReleasesCombinedDir string
}

func NewRunner(buildDir string) (Runner, error) {
	bd := taskfs.BuildDir(buildDir)

	inputs, err := setupInputs(bd)
	if err != nil {
		return Runner{}, err
	}

	outputs, err := setupOutputs(bd)
	if err != nil {
		return Runner{}, err
	}

	return Runner{In: inputs, Out: outputs}, nil
}

func setupInputs(buildDir taskfs.BuildDir) (Inputs, error) {
	compiledReleasesPrevDir, err := buildDir.SubDir("compiled-releases-prev")
	if err != nil {
		return Inputs{}, err
	}

	compiledReleasesNextDir, err := buildDir.SubDir("compiled-releases-next")
	if err != nil {
		return Inputs{}, err
	}

	cfDeploymentNextDir, err := buildDir.SubDir("cf-deployment-next")
	if err != nil {
		return Inputs{}, err
	}

	return Inputs{
		compiledReleasesPrevDir: compiledReleasesPrevDir,
		compiledReleasesNextDir: compiledReleasesNextDir,
		cfDeploymentNextDir:     cfDeploymentNextDir,
	}, nil
}

func setupOutputs(buildDir taskfs.BuildDir) (Outputs, error) {
	compliledReleasesCombinedDir, err := buildDir.SubDir("compiled-releases-combined")
	if err != nil {
		return Outputs{}, err
	}

	return Outputs{compliledReleasesCombinedDir: compliledReleasesCombinedDir}, nil
}

func (r Runner) Run() error {
	return nil
}
