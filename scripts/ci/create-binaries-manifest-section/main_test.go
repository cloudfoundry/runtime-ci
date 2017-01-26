package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	expectedReleasesAndStemcells string = `
name: cf-deployment
releases:
- name: capi
  url: updated-release-url
  version: updated-release-version
  sha1: updated-release-sha
- name: consul
  url: updated-release-url
  version: updated-release-version
  sha1: updated-release-sha
- name: diego
  url: updated-release-url
  version: updated-release-version
  sha1: updated-release-sha
- name: etcd
  url: updated-release-url
  version: updated-release-version
  sha1: updated-release-sha
- name: loggregator
  url: updated-release-url
  version: updated-release-version
  sha1: updated-release-sha
- name: nats
  url: updated-release-url
  version: updated-release-version
  sha1: updated-release-sha
- name: cf-mysql
  url: updated-release-url
  version: updated-release-version
  sha1: updated-release-sha
- name: routing
  url: updated-release-url
  version: updated-release-version
  sha1: updated-release-sha
- name: uaa
  url: updated-release-url
  version: updated-release-version
  sha1: updated-release-sha
- name: garden-runc
  url: updated-release-url
  version: updated-release-version
  sha1: updated-release-sha
- name: cflinuxfs2-rootfs
  url: updated-release-url
  version: updated-release-version
  sha1: updated-release-sha
- name: binary-buildpack
  url: updated-release-url
  version: updated-release-version
  sha1: updated-release-sha
- name: dotnet-core-buildpack
  url: updated-release-url
  version: updated-release-version
  sha1: updated-release-sha
- name: go-buildpack
  url: updated-release-url
  version: updated-release-version
  sha1: updated-release-sha
- name: java-buildpack
  url: updated-release-url
  version: updated-release-version
  sha1: updated-release-sha
- name: nodejs-buildpack
  url: updated-release-url
  version: updated-release-version
  sha1: updated-release-sha
- name: php-buildpack
  url: updated-release-url
  version: updated-release-version
  sha1: updated-release-sha
- name: python-buildpack
  url: updated-release-url
  version: updated-release-version
  sha1: updated-release-sha
- name: ruby-buildpack
  url: updated-release-url
  version: updated-release-version
  sha1: updated-release-sha
- name: staticfile-buildpack
  url: updated-release-url
  version: updated-release-version
  sha1: updated-release-sha
stemcells:
- alias: default
  os: ubuntu-trusty
  version: updated-stemcell-version
`

	releaseAndStemcellStub string = `
name: cf-deployment
releases:
- name: capi
  url: original-release-url
  version: original-release-version
  sha1: original-release-sha
- name: consul
  url: original-release-url
  version: original-release-version
  sha1: original-release-sha
- name: diego
  url: original-release-url
  version: original-release-version
  sha1: original-release-sha
- name: etcd
  url: original-release-url
  version: original-release-version
  sha1: original-release-sha
- name: loggregator
  url: original-release-url
  version: original-release-version
  sha1: original-release-sha
- name: nats
  url: original-release-url
  version: original-release-version
  sha1: original-release-sha
- name: cf-mysql
  url: original-release-url
  version: original-release-version
  sha1: original-release-sha
- name: routing
  url: original-release-url
  version: original-release-version
  sha1: original-release-sha
- name: uaa
  url: original-release-url
  version: original-release-version
  sha1: original-release-sha
- name: garden-runc
  url: original-release-url
  version: original-release-version
  sha1: original-release-sha
- name: cflinuxfs2-rootfs
  url: original-release-url
  version: original-release-version
  sha1: original-release-sha
- name: binary-buildpack
  url: original-release-url
  version: original-release-version
  sha1: original-release-sha
- name: dotnet-core-buildpack
  url: original-release-url
  version: original-release-version
  sha1: original-release-sha
- name: go-buildpack
  url: original-release-url
  version: original-release-version
  sha1: original-release-sha
- name: java-buildpack
  url: original-release-url
  version: original-release-version
  sha1: original-release-sha
- name: nodejs-buildpack
  url: original-release-url
  version: original-release-version
  sha1: original-release-sha
- name: php-buildpack
  url: original-release-url
  version: original-release-version
  sha1: original-release-sha
- name: python-buildpack
  url: original-release-url
  version: original-release-version
  sha1: original-release-sha
- name: ruby-buildpack
  url: original-release-url
  version: original-release-version
  sha1: original-release-sha
- name: staticfile-buildpack
  url: original-release-url
  version: original-release-version
  sha1: original-release-sha
stemcells:
- alias: default
  os: ubuntu-trusty
  version: original-stemcell-version
`
)

var _ = Describe("main", func() {
	var (
		pathToBinary                string
		DeploymentConfigurationPath string
		DeploymentManifestPath      string
		buildDir                    string
	)

	BeforeEach(func() {
		var err error

		pathToBinary, err = gexec.Build("github.com/cloudfoundry/runtime-ci/scripts/ci/create-binaries-manifest-section")
		Expect(err).NotTo(HaveOccurred())

		buildDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		for _, dir := range []string{
			"deployment-configuration",
			"deployment-manifest",
			"stemcell",
		} {
			err = os.Mkdir(filepath.Join(buildDir, dir), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
		}

		err = ioutil.WriteFile(filepath.Join(buildDir, "deployment-configuration", "original-manifest.yml"), []byte(releaseAndStemcellStub), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		for _, release := range []string{
			"capi",
			"consul",
			"diego",
			"etcd",
			"loggregator",
			"nats",
			"cf-mysql",
			"routing",
			"uaa",
			"garden-runc",
			"cflinuxfs2-rootfs",
			"binary-buildpack",
			"dotnet-core-buildpack",
			"go-buildpack",
			"java-buildpack",
			"nodejs-buildpack",
			"php-buildpack",
			"python-buildpack",
			"ruby-buildpack",
			"staticfile-buildpack",
		} {
			releaseDir := filepath.Join(buildDir, fmt.Sprintf("%s-release", release))
			err = os.Mkdir(releaseDir, os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			for file, value := range map[string]string{"version": "updated-release-version", "url": "updated-release-url", "sha1": "updated-release-sha"} {
				err = ioutil.WriteFile(filepath.Join(releaseDir, file), []byte(value), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())
			}
		}

		err = ioutil.WriteFile(filepath.Join(buildDir, "stemcell", "version"), []byte("updated-stemcell-version"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		DeploymentConfigurationPath = os.Getenv("DEPLOYMENT_CONFIGURATION_PATH")
		DeploymentManifestPath = os.Getenv("DEPLOYMENT_MANIFEST_PATH")
	})

	AfterEach(func() {
		os.Setenv("DEPLOYMENT_CONFIGURATION_PATH", DeploymentConfigurationPath)
		os.Setenv("DEPLOYMENT_MANIFEST_PATH", DeploymentManifestPath)
	})

	It("updates the given manifest with new releases and stemcells", func() {
		os.Setenv("DEPLOYMENT_CONFIGURATION_PATH", "original-manifest.yml")
		os.Setenv("DEPLOYMENT_MANIFEST_PATH", "updated-manifest.yml")

		session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir}...), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(gexec.Exit())
		Expect(session.ExitCode()).To(Equal(0))

		updatedManifest, err := ioutil.ReadFile(filepath.Join(buildDir, "deployment-manifest", "updated-manifest.yml"))
		Expect(err).NotTo(HaveOccurred())

		Expect(updatedManifest).To(MatchYAML(expectedReleasesAndStemcells))
	})
})
