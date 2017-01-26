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
		CommitMessagePath           string
		buildDir                    string
		emptyDir                    string
	)

	BeforeEach(func() {
		var err error

		pathToBinary, err = gexec.Build("github.com/cloudfoundry/runtime-ci/scripts/ci/create-binaries-manifest-section")
		Expect(err).NotTo(HaveOccurred())

		buildDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		emptyDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		for _, dir := range []string{
			"commit-message",
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

			for file, value := range map[string]string{"version": "original-release-version", "url": "original-release-url", "sha1": "original-release-sha"} {
				err = ioutil.WriteFile(filepath.Join(releaseDir, file), []byte(value), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())
			}
		}

		err = ioutil.WriteFile(filepath.Join(buildDir, "stemcell", "version"), []byte("updated-stemcell-version"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		DeploymentConfigurationPath = os.Getenv("DEPLOYMENT_CONFIGURATION_PATH")
		DeploymentManifestPath = os.Getenv("DEPLOYMENT_MANIFEST_PATH")
		CommitMessagePath = os.Getenv("COMMIT_MESSAGE_PATH")

		os.Setenv("DEPLOYMENT_CONFIGURATION_PATH", "original-manifest.yml")
		os.Setenv("DEPLOYMENT_MANIFEST_PATH", "updated-manifest.yml")
		os.Setenv("COMMIT_MESSAGE_PATH", "commit-message.txt")
	})

	AfterEach(func() {
		os.Setenv("DEPLOYMENT_CONFIGURATION_PATH", DeploymentConfigurationPath)
		os.Setenv("DEPLOYMENT_MANIFEST_PATH", DeploymentManifestPath)
		os.Setenv("COMMIT_MESSAGE_PATH", CommitMessagePath)
	})

	It("updates the given manifest with new releases and stemcells", func() {
		session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir}...), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(gexec.Exit())
		Expect(session.ExitCode()).To(Equal(0))

		updatedManifest, err := ioutil.ReadFile(filepath.Join(buildDir, "deployment-manifest", "updated-manifest.yml"))
		Expect(err).NotTo(HaveOccurred())

		Expect(updatedManifest).To(MatchYAML(expectedReleasesAndStemcells))
	})

	It("writes the commit message to COMMIT_MESSAGE_PATH", func() {
		session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir}...), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(gexec.Exit())
		Expect(session.ExitCode()).To(Equal(0))

		commitMessage, err := ioutil.ReadFile(filepath.Join(buildDir, "commit-message", "commit-message.txt"))
		Expect(err).NotTo(HaveOccurred())

		Expect(string(commitMessage)).To(Equal("Updated ubuntu-trusty stemcell"))
	})

	Context("failure cases", func() {
		It("errors when the deployment manifest does not exist", func() {
			os.Setenv("DEPLOYMENT_CONFIGURATION_PATH", emptyDir)

			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", emptyDir}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(1))

			Expect(string(session.Err.Contents())).To(ContainSubstring("no such file or directory"))
		})

		It("errors when the directory to write out the commit message does not exist", func() {
			os.Setenv("COMMIT_MESSAGE_PATH", filepath.Join(emptyDir, "doesnt-exist"))

			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(1))

			Expect(string(session.Err.Contents())).To(ContainSubstring("doesnt-exist: no such file or directory"))
		})

		It("errors when the directory to write out the updated manifest does not exist", func() {
			os.Setenv("DEPLOYMENT_MANIFEST_PATH", filepath.Join(emptyDir, "doesnt-exist"))

			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(1))

			Expect(string(session.Err.Contents())).To(ContainSubstring("doesnt-exist: no such file or directory"))
		})

		It("errors when a required stemcell file is missing", func() {
			err := os.Remove(filepath.Join(buildDir, "stemcell", "version"))
			Expect(err).NotTo(HaveOccurred())

			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(1))

			Expect(string(session.Err.Contents())).To(ContainSubstring("stemcell/version: no such file or directory"))
		})
	})
})
