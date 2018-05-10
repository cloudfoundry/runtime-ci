package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("main", func() {
	var (
		pathToBinary      string
		inputPath         string
		outputPath        string
		CommitMessagePath string
		buildDir          string
		emptyDir          string

		err error
	)

	BeforeEach(func() {
		pathToBinary, err = gexec.Build("github.com/cloudfoundry/runtime-ci/util/update-manifest-releases")
		Expect(err).NotTo(HaveOccurred())

		buildDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		emptyDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		CommitMessagePath = os.Getenv("COMMIT_MESSAGE_PATH")
		os.Setenv("COMMIT_MESSAGE_PATH", "commit-message.txt")
	})

	AfterEach(func() {
		os.Setenv("COMMIT_MESSAGE_PATH", CommitMessagePath)
	})

	Context("opsfile", func() {
		const (
			originalOpsFile string = `
- type: replace
  path: /releases/-
  value:
    name: release1
    url: original-release1-url
    version: original-release1-version
    sha1: original-release1-sha
- type: replace
  path: /releases/-
  value:
    name: release4
    url: original-release4-url
    version: original-release4-version
    sha1: original-release4-sha
`
			expectedOpsFile string = `
- type: replace
  path: /releases/-
  value:
    name: release1
    url: original-release1-url
    version: original-release1-version
    sha1: original-release1-sha
- type: replace
  path: /releases/-
  value:
    name: release4
    url: new-release4-url
    version: new-release4-version
    sha1: new-release4-sha
`
		)

		BeforeEach(func() {
			for _, dir := range []string{
				"commit-message",
				"original-ops-file",
				"updated-ops-file",
			} {
				err = os.Mkdir(filepath.Join(buildDir, dir), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())
			}

			err = ioutil.WriteFile(filepath.Join(buildDir, "original-ops-file", "original_ops_file.yml"), []byte(originalOpsFile), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			inputPath = os.Getenv("ORIGINAL_OPS_FILE_PATH")
			outputPath = os.Getenv("UPDATED_OPS_FILE_PATH")

			os.Setenv("ORIGINAL_OPS_FILE_PATH", "original_ops_file.yml")
			os.Setenv("UPDATED_OPS_FILE_PATH", "updated_ops_file.yml")

			for _, release := range []map[string]string{
				{"name": "release4", "version": "new-release4-version", "url": "new-release4-url", "sha1": "new-release4-sha"},
			} {
				releaseDir := filepath.Join(buildDir, fmt.Sprintf("%s-release", release["name"]))
				err = os.Mkdir(releaseDir, os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				for _, value := range []string{"version", "url", "sha1"} {
					err = ioutil.WriteFile(filepath.Join(releaseDir, value), []byte(release[value]), os.ModePerm)
					Expect(err).NotTo(HaveOccurred())
				}
			}
		})

		AfterEach(func() {
			os.Setenv("ORIGINAL_OPS_FILE_PATH", inputPath)
			os.Setenv("UPDATED_OPS_FILE_PATH", outputPath)
		})

		It("updates the original ops file with new releases", func() {
			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--target", "opsfile"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			updatedOpsFile, err := ioutil.ReadFile(filepath.Join(buildDir, "updated-ops-file", "updated_ops_file.yml"))
			Expect(err).NotTo(HaveOccurred())

			Expect(updatedOpsFile).To(MatchYAML(expectedOpsFile))
		})

		It("writes the commit message to COMMIT_MESSAGE_PATH", func() {
			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--target", "opsfile"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			commitMessage, err := ioutil.ReadFile(filepath.Join(buildDir, "commit-message", "commit-message.txt"))
			Expect(err).NotTo(HaveOccurred())

			Expect(string(commitMessage)).To(Equal("Updated opsfile with release4-release new-release4-version"))
		})

		Context("failure cases", func() {
			It("errors when the build dir does not exist", func() {
				fakeDirName := fmt.Sprintf("fake-dir-%v", time.Now().Unix())
				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", fakeDirName, "--target", "opsfile"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring(fmt.Sprintf("%v: no such file or directory", fakeDirName)))
			})

			It("errors when the original ops file does not exist", func() {
				os.Setenv("ORIGINAL_OPS_FILE_PATH", emptyDir)

				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", emptyDir, "--target", "opsfile"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring("no such file or directory"))
			})

			It("errors when the directory to write out the commit message does not exist", func() {
				os.Setenv("COMMIT_MESSAGE_PATH", filepath.Join(emptyDir, "doesnt-exist"))

				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--target", "opsfile"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring("doesnt-exist: no such file or directory"))
			})

			It("errors when the directory to write out the updated ops file path does not exist", func() {
				os.Setenv("UPDATED_OPS_FILE_PATH", filepath.Join(emptyDir, "doesnt-exist"))

				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--target", "opsfile"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring("doesnt-exist: no such file or directory"))
			})
		})
	})

	Context("manifest", func() {
		const (
			expectedReleasesAndStemcells string = `
name: cf-deployment
releases:
- name: release1
  url: original-release1-url
  version: original-release1-version
  sha1: original-release1-sha
- name: release2
  url: original-release2-url
  version: original-release2-version
  sha1: original-release2-sha
- name: release3
  url: new-release3-url
  version: new-release3-version
  sha1: new-release3-sha
- name: release4
  url: new-release4-url
  version: new-release4-version
  sha1: new-release4-sha
stemcells:
- alias: default
  os: ubuntu-trusty
  version: updated-stemcell-version
`

			expectedSingleReleaseAndStemcells string = `
name: cf-deployment
releases:
- name: release1
  url: original-release1-url
  version: original-release1-version
  sha1: original-release1-sha
- name: release2
  url: original-release2-url
  version: original-release2-version
  sha1: original-release2-sha
- name: release3
  url: new-release3-url
  version: new-release3-version
  sha1: new-release3-sha
stemcells:
- alias: default
  os: ubuntu-trusty
  version: updated-stemcell-version
`

			originalManifest string = `
name: cf-deployment
releases:
- name: release1
  url: original-release1-url
  version: original-release1-version
  sha1: original-release1-sha
- name: release2
  url: original-release2-url
  version: original-release2-version
  sha1: original-release2-sha
stemcells:
- alias: default
  os: ubuntu-trusty
  version: original-stemcell-version
`
		)

		BeforeEach(func() {
			for _, dir := range []string{
				"commit-message",
				"deployment-configuration",
				"updated-deployment-manifest",
				"stemcell",
			} {
				err = os.Mkdir(filepath.Join(buildDir, dir), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())
			}

			err = ioutil.WriteFile(filepath.Join(buildDir, "deployment-configuration", "original-manifest.yml"), []byte(originalManifest), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			for _, release := range []map[string]string{
				{"name": "release1", "version": "original-release1-version", "url": "original-release1-url", "sha1": "original-release1-sha"},
				{"name": "release2", "version": "original-release2-version", "url": "original-release2-url", "sha1": "original-release2-sha"},
				{"name": "release3", "version": "new-release3-version", "url": "new-release3-url", "sha1": "new-release3-sha"},
				{"name": "release4", "version": "new-release4-version", "url": "new-release4-url", "sha1": "new-release4-sha"},
			} {
				releaseDir := filepath.Join(buildDir, fmt.Sprintf("%s-release", release["name"]))
				err = os.Mkdir(releaseDir, os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				for _, value := range []string{"version", "url", "sha1"} {
					err = ioutil.WriteFile(filepath.Join(releaseDir, value), []byte(release[value]), os.ModePerm)
					Expect(err).NotTo(HaveOccurred())
				}
			}

			err = ioutil.WriteFile(filepath.Join(buildDir, "stemcell", "version"), []byte("updated-stemcell-version"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			inputPath = os.Getenv("ORIGINAL_DEPLOYMENT_MANIFEST_PATH")
			outputPath = os.Getenv("UPDATED_DEPLOYMENT_MANIFEST_PATH")

			os.Setenv("ORIGINAL_DEPLOYMENT_MANIFEST_PATH", "original-manifest.yml")
			os.Setenv("UPDATED_DEPLOYMENT_MANIFEST_PATH", "updated-manifest.yml")
		})

		AfterEach(func() {
			os.Setenv("ORIGINAL_DEPLOYMENT_MANIFEST_PATH", inputPath)
			os.Setenv("UPDATED_DEPLOYMENT_MANIFEST_PATH", outputPath)
		})

		It("only updates the manifest with the release passed in", func() {
			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--release", "release3"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			updatedManifest, err := ioutil.ReadFile(filepath.Join(buildDir, "updated-deployment-manifest", "updated-manifest.yml"))
			Expect(err).NotTo(HaveOccurred())

			Expect(updatedManifest).To(MatchYAML(expectedSingleReleaseAndStemcells))
		})

		It("updates the given manifest with new releases and stemcells", func() {
			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			updatedManifest, err := ioutil.ReadFile(filepath.Join(buildDir, "updated-deployment-manifest", "updated-manifest.yml"))
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

			Expect(string(commitMessage)).To(Equal("Updated manifest with release3-release new-release3-version, release4-release new-release4-version, ubuntu-trusty stemcell updated-stemcell-version"))
		})

		Context("failure cases", func() {
			It("errors when the build dir does not exist", func() {
				fakeDirName := fmt.Sprintf("fake-dir-%v", time.Now().Unix())
				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", fakeDirName}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring(fmt.Sprintf("%v: no such file or directory", fakeDirName)))
			})

			It("errors when the deployment manifest does not exist", func() {
				os.Setenv("ORIGINAL_DEPLOYMENT_MANIFEST_PATH", emptyDir)

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
				os.Setenv("UPDATED_DEPLOYMENT_MANIFEST_PATH", filepath.Join(emptyDir, "doesnt-exist"))

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
})
