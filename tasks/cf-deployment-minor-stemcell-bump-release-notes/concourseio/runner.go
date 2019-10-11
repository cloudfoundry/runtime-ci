package concourseio

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry/runtime-ci/task-libs/bosh"
)

type Runner struct {
	In  Inputs
	Out Outputs
}

type Inputs struct {
	CFDeploymentDir   string
	ReleaseVersionDir string
	StemcellDir       string
}

type Outputs struct {
	ReleaseNotesDir string
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

func (r Runner) ReadStemcellInfoFromManifest(stemcellAlias string) (bosh.Stemcell, error) {
	content, err := ioutil.ReadFile(filepath.Join(r.In.CFDeploymentDir, "cf-deployment.yml"))
	if err != nil {
		return bosh.Stemcell{}, fmt.Errorf("failed to read cf-deployment.yml: %w", err)
	}

	manifest, err := bosh.NewManifestFromFile(content)
	if err != nil {
		return bosh.Stemcell{}, fmt.Errorf("failed to unmarshal cf-deployment.yml: %w", err)
	}

	for _, stemcell := range manifest.Stemcells {
		if stemcell.Alias == stemcellAlias {
			return stemcell, nil
		}
	}

	return bosh.Stemcell{}, fmt.Errorf("failed to find stemcell version for alias %q", stemcellAlias)
}

// Not tested due to simplicty
func (r Runner) ReadStemcellFromResource() (bosh.Stemcell, error) {
	return bosh.NewStemcellFromInput(r.In.StemcellDir)
}

func (r Runner) ValidateStemcellBump(oldStemcell, newStemcell bosh.Stemcell) error {
	if oldStemcell.OS != newStemcell.OS {
		return fmt.Errorf("stemcell os mismatch: found %q and %q", oldStemcell.OS, newStemcell.OS)
	}
	return nil
}

func (r Runner) GenerateReleaseNotes(oldStemcell, newStemcell bosh.Stemcell) error {
	template := `## Stemcell Updates
| Release | Old Version | New Version |
| - | - | - |
| %s | %s | %s |
`
	content := fmt.Sprintf(template, strings.ReplaceAll(oldStemcell.OS, "-", " "), oldStemcell.Version, newStemcell.Version)
	err := ioutil.WriteFile(filepath.Join(r.Out.ReleaseNotesDir, "body.txt"), []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write release notes file: %w", err)
	}

	return nil
}

func (r Runner) GenerateReleaseName() error {
	releaseVersion, err := ioutil.ReadFile(filepath.Join(r.In.ReleaseVersionDir, "version"))
	if err != nil {
		return fmt.Errorf("failed to read release version: %w", err)
	}

	err = ioutil.WriteFile(filepath.Join(r.Out.ReleaseNotesDir, "name.txt"), []byte(fmt.Sprintf("v%s", releaseVersion)), 0644)
	if err != nil {
		return fmt.Errorf("failed to write release name file: %w", err)
	}

	return nil
}

func setupInputs(buildDir string) (Inputs, error) {
	cfDeploymentDir, err := buildSubDir(buildDir, "cf-deployment-master")
	if err != nil {
		return Inputs{}, err
	}

	releaseVersionDir, err := buildSubDir(buildDir, "release-version")
	if err != nil {
		return Inputs{}, err
	}

	stemcellDir, err := buildSubDir(buildDir, "stemcell")
	if err != nil {
		return Inputs{}, err
	}

	return Inputs{cfDeploymentDir, releaseVersionDir, stemcellDir}, nil
}

func setupOutputs(buildDir string) (Outputs, error) {
	ReleaseNotesDir, err := buildSubDir(buildDir, "cf-deployment-minor-stemcell-bump-release-notes")
	if err != nil {
		return Outputs{}, err
	}

	return Outputs{ReleaseNotesDir}, nil
}

func buildSubDir(buildDir, subDir string) (string, error) {
	dir := filepath.Join(buildDir, subDir)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("missing sub directory '%s' in build directory '%s'", subDir, buildDir)
	}

	return dir, nil
}
