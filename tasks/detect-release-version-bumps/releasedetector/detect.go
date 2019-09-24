package releasedetector

import (
	"io/ioutil"

	"github.com/cloudfoundry/runtime-ci/task-libs/bosh"
	"gopkg.in/yaml.v2"
)

type ReleaseDetector struct {
	prevDeploymentManifestPath string
	nextDeploymentManifestPath string

	prevReleases    []bosh.Release
	nextReleases    []bosh.Release
	updatedReleases []bosh.Release
}

func NewReleaseDetector() *ReleaseDetector {
	return &ReleaseDetector{}
}

func (r *ReleaseDetector) Load(prevDeploymentManifestPath, nextDeploymentManifestPath string) error {
	var err error
	r.prevReleases, err = readManifest(prevDeploymentManifestPath)
	if err != nil {
		return err
	}
	r.nextReleases, err = readManifest(nextDeploymentManifestPath)
	if err != nil {
		return err
	}

	return nil
}

func readManifest(manifestPath string) ([]bosh.Release, error) {
	content, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}

	manifest, err := bosh.NewManifestFromFile(content)
	if err != nil {
		return nil, err
	}

	return manifest.Releases, nil
}

func (r *ReleaseDetector) DetectUpdatedReleases() {
	prevReleases := map[string]string{}

	for _, prevRelease := range r.prevReleases {
		prevReleases[prevRelease.Name] = prevRelease.Version
	}

	for _, nextRelease := range r.nextReleases {
		if prevVersion, ok := prevReleases[nextRelease.Name]; !ok {
			r.updatedReleases = append(r.updatedReleases, nextRelease)
		} else if prevVersion != nextRelease.Version {
			r.updatedReleases = append(r.updatedReleases, nextRelease)
		}
	}
}

func (r ReleaseDetector) Write() ([]byte, error) {
	return yaml.Marshal(bosh.Manifest{
		Releases: r.updatedReleases,
	})
}
