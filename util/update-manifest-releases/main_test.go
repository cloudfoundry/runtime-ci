package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("main", func() {
	var (
		inputPath         string
		outputPath        string
		CommitMessagePath string
		buildDir          string
		emptyDir          string
	)

	BeforeEach(func() {
		var err error
		buildDir, err = os.MkdirTemp("", "")
		Expect(err).NotTo(HaveOccurred())

		emptyDir, err = os.MkdirTemp("", "")
		Expect(err).NotTo(HaveOccurred())

		CommitMessagePath = os.Getenv("COMMIT_MESSAGE_PATH")

		err = os.Setenv("COMMIT_MESSAGE_PATH", "commit-message.txt")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := os.Setenv("COMMIT_MESSAGE_PATH", CommitMessagePath)
		Expect(err).NotTo(HaveOccurred(), "The original COMMIT_MESSAGE_PATH env var should be set back")
	})

	Context("opsfile", func() {
		const (
			originalOpsFile = `
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
			expectedOpsFile = `
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
			anotherOriginalOpsFileWithRelease4 = `
- type: replace
  path: /releases/-
  value:
    name: release2
    url: original-release2-url
    version: original-release2-version
    sha1: original-release2-sha
- type: replace
  path: /releases/-
  value:
    name: release4
    url: original-release4-url
    version: original-release4-version
    sha1: original-release4-sha
`
			anotherExpectedOpsFileWithRelease4 = `
- type: replace
  path: /releases/-
  value:
    name: release2
    url: original-release2-url
    version: original-release2-version
    sha1: original-release2-sha
- type: replace
  path: /releases/-
  value:
    name: release4
    url: new-release4-url
    version: new-release4-version
    sha1: new-release4-sha
`
			opsFileWithoutRelease4 = `
- type: replace
  path: /releases/-
  value:
    name: release2
    url: original-release2-url
    version: original-release2-version
    sha1: original-release2-sha
- type: replace
  path: /releases/-
  value:
    name: release5
    url: original-release5-url
    version: original-release5-version
    sha1: original-release5-sha
`
		)

		BeforeEach(func() {
			for _, dir := range []string{
				"original-ops-file",
				"updated-ops-file",
			} {
				err := os.Mkdir(filepath.Join(buildDir, dir), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())
			}

			err := os.WriteFile(filepath.Join(buildDir, "original-ops-file", "original_ops_file.yml"), []byte(originalOpsFile), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			inputPath = os.Getenv("ORIGINAL_OPS_FILE_PATH")
			err = os.Setenv("ORIGINAL_OPS_FILE_PATH", "original_ops_file.yml")
			Expect(err).NotTo(HaveOccurred())

			outputPath = os.Getenv("UPDATED_OPS_FILE_PATH")
			err = os.Setenv("UPDATED_OPS_FILE_PATH", "updated_ops_file.yml")
			Expect(err).NotTo(HaveOccurred())

			for _, release := range []map[string]string{
				{"name": "release1", "version": "original-release1-version", "url": "original-release1-url", "sha1": "original-release1-sha"},
				{"name": "release4", "version": "new-release4-version", "url": "new-release4-url", "sha1": "new-release4-sha"},
			} {
				releaseDir := filepath.Join(buildDir, fmt.Sprintf("%s-release", release["name"]))
				err = os.Mkdir(releaseDir, os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				for _, value := range []string{"version", "url", "sha1"} {
					err = os.WriteFile(filepath.Join(releaseDir, value), []byte(release[value]), os.ModePerm)
					Expect(err).NotTo(HaveOccurred())
				}
			}
		})

		AfterEach(func() {
			err := os.Setenv("ORIGINAL_OPS_FILE_PATH", inputPath)
			Expect(err).NotTo(HaveOccurred(), "The original ORIGINAL_OPS_FILE_PATH env var should be set back")
			err = os.Setenv("UPDATED_OPS_FILE_PATH", outputPath)
			Expect(err).NotTo(HaveOccurred(), "The original UPDATED_OPS_FILE_PATH env var should be set back")
		})

		Context("when there is more than one ops file containing desired release", func() {
			BeforeEach(func() {
				err := os.Unsetenv("ORIGINAL_OPS_FILE_PATH")
				Expect(err).NotTo(HaveOccurred())
				err = os.Unsetenv("UPDATED_OPS_FILE_PATH")
				Expect(err).NotTo(HaveOccurred())

				anotherOpsFileDir := filepath.Join(buildDir, "original-ops-file", "nested-dir")

				err = os.MkdirAll(anotherOpsFileDir, os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(filepath.Join(anotherOpsFileDir, "another_original_ops_file.yml"), []byte(anotherOriginalOpsFileWithRelease4), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(filepath.Join(buildDir, "original-ops-file", "ops_file_that_should_stay_the_same.yml"), []byte(opsFileWithoutRelease4), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				for _, release := range []map[string]string{
					{"name": "release2", "version": "original-release2-version", "url": "original-release2-url", "sha1": "original-release2-sha"},
					{"name": "release5", "version": "original-release5-version", "url": "original-release5-url", "sha1": "original-release5-sha"},
				} {
					releaseDir := filepath.Join(buildDir, fmt.Sprintf("%s-release", release["name"]))
					err = os.Mkdir(releaseDir, os.ModePerm)
					Expect(err).NotTo(HaveOccurred())

					for _, value := range []string{"version", "url", "sha1"} {
						err = os.WriteFile(filepath.Join(releaseDir, value), []byte(release[value]), os.ModePerm)
						Expect(err).NotTo(HaveOccurred())
					}
				}
			})

			It("updates both ops files with the release and doesn't include non-updated ops files in the output directory", func() {
				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "original-ops-file", "--output-dir", "updated-ops-file", "--target", "opsfile", "--release", "release4"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 5*time.Second).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(0))

				updatedOpsFile1, err := os.ReadFile(filepath.Join(buildDir, "updated-ops-file", "original_ops_file.yml"))
				Expect(err).NotTo(HaveOccurred())

				updatedOpsFile2, err := os.ReadFile(filepath.Join(buildDir, "updated-ops-file", "nested-dir", "another_original_ops_file.yml"))
				Expect(err).NotTo(HaveOccurred())

				_, err = os.ReadFile(filepath.Join(buildDir, "updated-ops-file", "ops_file_that_should_stay_the_same.yml"))
				Expect(err).To(HaveOccurred())

				Expect(updatedOpsFile1).To(MatchYAML(expectedOpsFile))
				Expect(updatedOpsFile2).To(MatchYAML(anotherExpectedOpsFileWithRelease4))
			})

			It("writes the commit message to COMMIT_MESSAGE_PATH", func() {
				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "original-ops-file", "--output-dir", "updated-ops-file", "--target", "opsfile", "--release", "release4"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 5*time.Second).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(0))

				commitMessage, err := os.ReadFile(filepath.Join(buildDir, "commit-message.txt"))
				Expect(err).NotTo(HaveOccurred())

				Expect(string(commitMessage)).To(Equal("Updated ops file(s) with release4-release new-release4-version"))
			})

			It("ignores cf-deployment.yml", func() {
				manifest := `
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
				err := os.WriteFile(filepath.Join(buildDir, "original-ops-file", "cf-deployment.yml"), []byte(manifest), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "original-ops-file", "--output-dir", "updated-ops-file", "--target", "opsfile"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 5*time.Second).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(0))

				updatedOpsFile1, err := os.ReadFile(filepath.Join(buildDir, "updated-ops-file", "original_ops_file.yml"))
				Expect(err).NotTo(HaveOccurred())

				updatedOpsFile2, err := os.ReadFile(filepath.Join(buildDir, "updated-ops-file", "nested-dir", "another_original_ops_file.yml"))
				Expect(err).NotTo(HaveOccurred())

				_, err = os.ReadFile(filepath.Join(buildDir, "updated-ops-file", "ops_file_that_should_stay_the_same.yml"))
				Expect(err).To(HaveOccurred())

				Expect(updatedOpsFile1).To(MatchYAML(expectedOpsFile))
				Expect(updatedOpsFile2).To(MatchYAML(anotherExpectedOpsFileWithRelease4))
			})

			It("doesn't error if the scripts directory contains yml files", func() {
				vars := `
some_client: some_value
`
				err := os.Mkdir(filepath.Join(buildDir, "original-ops-file", "scripts"), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(filepath.Join(buildDir, "original-ops-file", "scripts", "vars-store.yml"), []byte(vars), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "original-ops-file", "--output-dir", "updated-ops-file", "--target", "opsfile"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 5*time.Second).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(0))
			})

			It("doesn't error if the units directory contains yml files", func() {
				vars := `
some_client: some_value
`
				err := os.Mkdir(filepath.Join(buildDir, "original-ops-file", "units"), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				err = os.WriteFile(filepath.Join(buildDir, "original-ops-file", "units", "vars-store.yml"), []byte(vars), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "original-ops-file", "--output-dir", "updated-ops-file", "--target", "opsfile"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 5*time.Second).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(0))
			})
		})

		It("updates the original ops file with new releases", func() {
			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "original-ops-file", "--output-dir", "updated-ops-file", "--target", "opsfile"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 5*time.Second).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			updatedOpsFile, err := os.ReadFile(filepath.Join(buildDir, "updated-ops-file", "updated_ops_file.yml"))
			Expect(err).NotTo(HaveOccurred())

			Expect(updatedOpsFile).To(MatchYAML(expectedOpsFile))
		})

		It("writes the commit message to COMMIT_MESSAGE_PATH", func() {
			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "original-ops-file", "--output-dir", "updated-ops-file", "--target", "opsfile"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 5*time.Second).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			commitMessage, err := os.ReadFile(filepath.Join(buildDir, "commit-message.txt"))
			Expect(err).NotTo(HaveOccurred())

			Expect(string(commitMessage)).To(Equal("Updated ops file(s) with release4-release new-release4-version"))
		})

		It("creates nested directories when the directory to write out the updated ops file path does not exist", func() {
			updatedOpsFilePath := filepath.Join("doesnt-exist", "updated_ops_file.yml")
			err := os.Setenv("UPDATED_OPS_FILE_PATH", updatedOpsFilePath)
			Expect(err).NotTo(HaveOccurred())

			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "original-ops-file", "--output-dir", "updated-ops-file", "--target", "opsfile"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 5*time.Second).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			updatedOpsFile, err := os.ReadFile(filepath.Join(buildDir, "updated-ops-file", updatedOpsFilePath))
			Expect(err).NotTo(HaveOccurred())

			Expect(updatedOpsFile).To(MatchYAML(expectedOpsFile))
		})

		Context("failure cases", func() {
			It("errors when the build dir does not exist", func() {
				fakeDirName := fmt.Sprintf("fake-dir-%v", time.Now().Unix())
				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", fakeDirName, "--input-dir", "original-ops-file", "--output-dir", "updated-ops-file", "--target", "opsfile"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 5*time.Second).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring(fmt.Sprintf("%v: no such file or directory", fakeDirName)))
			})

			It("errors when the original ops file does not exist", func() {
				err := os.Setenv("ORIGINAL_OPS_FILE_PATH", emptyDir)
				Expect(err).NotTo(HaveOccurred())

				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", emptyDir, "--input-dir", "original-ops-file", "--output-dir", "updated-ops-file", "--target", "opsfile"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 5*time.Second).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring("no such file or directory"))
			})

			It("errors when the directory to write out the commit message does not exist", func() {
				err := os.Setenv("COMMIT_MESSAGE_PATH", filepath.Join(emptyDir, "doesnt-exist"))
				Expect(err).NotTo(HaveOccurred())

				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "original-ops-file", "--output-dir", "updated-ops-file", "--target", "opsfile"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 5*time.Second).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring("doesnt-exist: no such file or directory"))
			})
		})
	})

	Context("manifest", func() {
		const (
			expectedReleases = `
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
  version: original-stemcell-version
`

			expectedSingleRelease = `
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
  version: original-stemcell-version
`

			originalManifest = `
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
				"cf-deployment",
				"updated-cf-deployment",
				"stemcell",
			} {
				err := os.Mkdir(filepath.Join(buildDir, dir), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())
			}

			err := os.WriteFile(filepath.Join(buildDir, "cf-deployment", "original-manifest.yml"), []byte(originalManifest), os.ModePerm)
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
					err = os.WriteFile(filepath.Join(releaseDir, value), []byte(release[value]), os.ModePerm)
					Expect(err).NotTo(HaveOccurred())
				}
			}

			err = os.WriteFile(filepath.Join(buildDir, "stemcell", "version"), []byte("updated-stemcell-version"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			err = os.WriteFile(filepath.Join(buildDir, "stemcell", "url"), []byte("https://foo.com/bosh-stemcell-ubuntu-trusty.tgz"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			inputPath = os.Getenv("ORIGINAL_DEPLOYMENT_MANIFEST_PATH")
			err = os.Setenv("ORIGINAL_DEPLOYMENT_MANIFEST_PATH", "original-manifest.yml")
			Expect(err).NotTo(HaveOccurred())

			outputPath = os.Getenv("UPDATED_DEPLOYMENT_MANIFEST_PATH")
			err = os.Setenv("UPDATED_DEPLOYMENT_MANIFEST_PATH", "updated-manifest.yml")
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			err := os.Setenv("ORIGINAL_DEPLOYMENT_MANIFEST_PATH", inputPath)
			Expect(err).NotTo(HaveOccurred(), "The original ORIGINAL_DEPLOYMENT_MANIFEST_PATH env var should be set back")
			err = os.Setenv("UPDATED_DEPLOYMENT_MANIFEST_PATH", outputPath)
			Expect(err).NotTo(HaveOccurred(), "The original UPDATED_DEPLOYMENT_MANIFEST_PATH env var should be set back")
		})

		It("only updates the manifest with the release passed in", func() {
			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment", "--release", "release3"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 5*time.Second).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			updatedManifest, err := os.ReadFile(filepath.Join(buildDir, "updated-cf-deployment", "updated-manifest.yml"))
			Expect(err).NotTo(HaveOccurred())

			Expect(updatedManifest).To(MatchYAML(expectedSingleRelease))
		})

		It("updates the given manifest with new releases and stemcells", func() {
			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 5*time.Second).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			updatedManifest, err := os.ReadFile(filepath.Join(buildDir, "updated-cf-deployment", "updated-manifest.yml"))
			Expect(err).NotTo(HaveOccurred())

			Expect(updatedManifest).To(MatchYAML(expectedReleases))
		})

		It("writes the commit message to COMMIT_MESSAGE_PATH", func() {
			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 5*time.Second).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			commitMessage, err := os.ReadFile(filepath.Join(buildDir, "commit-message.txt"))
			Expect(err).NotTo(HaveOccurred())

			Expect(string(commitMessage)).To(Equal("Updated manifest with release3-release new-release3-version, release4-release new-release4-version"))
		})

		It("creates nested directory when the directory to write out the updated manifest does not exist", func() {
			err := os.Setenv("UPDATED_DEPLOYMENT_MANIFEST_PATH", filepath.Join("doesnt-exist", "updated-manifest.yml"))
			Expect(err).NotTo(HaveOccurred())

			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 5*time.Second).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			updatedManifest, err := os.ReadFile(filepath.Join(buildDir, "updated-cf-deployment", "doesnt-exist", "updated-manifest.yml"))
			Expect(err).NotTo(HaveOccurred())

			Expect(updatedManifest).To(MatchYAML(expectedReleases))
		})

		Context("failure cases", func() {
			It("errors when the build dir does not exist", func() {
				fakeDirName := fmt.Sprintf("fake-dir-%v", time.Now().Unix())
				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", fakeDirName, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 5*time.Second).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring(fmt.Sprintf("%v: no such file or directory", fakeDirName)))
			})

			It("errors when the deployment manifest does not exist", func() {
				err := os.Setenv("ORIGINAL_DEPLOYMENT_MANIFEST_PATH", emptyDir)
				Expect(err).NotTo(HaveOccurred())

				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", emptyDir, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 5*time.Second).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring("no such file or directory"))
			})

			It("errors when the directory to write out the commit message does not exist", func() {
				err := os.Setenv("COMMIT_MESSAGE_PATH", filepath.Join(emptyDir, "doesnt-exist"))
				Expect(err).NotTo(HaveOccurred())

				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 5*time.Second).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring("doesnt-exist: no such file or directory"))
			})
		})
	})

	Context("stemcell", func() {
		const (
			expectedStemcells = `
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
  version: updated-stemcell-version
`

			originalManifest = `
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
				"cf-deployment",
				"updated-cf-deployment",
				"stemcell",
			} {
				err := os.Mkdir(filepath.Join(buildDir, dir), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())
			}

			err := os.WriteFile(filepath.Join(buildDir, "cf-deployment", "original-manifest.yml"), []byte(originalManifest), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			err = os.WriteFile(filepath.Join(buildDir, "stemcell", "version"), []byte("updated-stemcell-version"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			err = os.WriteFile(filepath.Join(buildDir, "stemcell", "url"), []byte("https://foo.com/bosh-stemcell-ubuntu-trusty.tgz"), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			inputPath = os.Getenv("ORIGINAL_DEPLOYMENT_MANIFEST_PATH")
			err = os.Setenv("ORIGINAL_DEPLOYMENT_MANIFEST_PATH", "original-manifest.yml")
			Expect(err).NotTo(HaveOccurred())

			outputPath = os.Getenv("UPDATED_DEPLOYMENT_MANIFEST_PATH")
			err = os.Setenv("UPDATED_DEPLOYMENT_MANIFEST_PATH", "updated-manifest.yml")
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			err := os.Setenv("ORIGINAL_DEPLOYMENT_MANIFEST_PATH", inputPath)
			Expect(err).NotTo(HaveOccurred(), "The original ORIGINAL_DEPLOYMENT_MANIFEST_PATH env var should be set back")
			err = os.Setenv("UPDATED_DEPLOYMENT_MANIFEST_PATH", outputPath)
			Expect(err).NotTo(HaveOccurred(), "The original UPDATED_DEPLOYMENT_MANIFEST_PATH env var should be set back")
		})

		It("only updates the manifest with the new stemcell", func() {
			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment", "--target", "stemcell"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 5*time.Second).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			updatedManifest, err := os.ReadFile(filepath.Join(buildDir, "updated-cf-deployment", "updated-manifest.yml"))
			Expect(err).NotTo(HaveOccurred())

			Expect(updatedManifest).To(MatchYAML(expectedStemcells))
		})

		It("writes the commit message to COMMIT_MESSAGE_PATH", func() {
			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment", "--target", "stemcell"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 5*time.Second).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			commitMessage, err := os.ReadFile(filepath.Join(buildDir, "commit-message.txt"))
			Expect(err).NotTo(HaveOccurred())

			Expect(string(commitMessage)).To(Equal("Updated manifest with ubuntu-trusty stemcell updated-stemcell-version"))
		})

		It("creates nested directory when the directory to write out the updated manifest does not exist", func() {
			err := os.Setenv("UPDATED_DEPLOYMENT_MANIFEST_PATH", filepath.Join("doesnt-exist", "updated-manifest.yml"))
			Expect(err).NotTo(HaveOccurred())

			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment", "--target", "stemcell"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 5*time.Second).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			updatedManifest, err := os.ReadFile(filepath.Join(buildDir, "updated-cf-deployment", "doesnt-exist", "updated-manifest.yml"))
			Expect(err).NotTo(HaveOccurred())

			Expect(updatedManifest).To(MatchYAML(expectedStemcells))
		})

		Context("failure cases", func() {
			It("errors when the build dir does not exist", func() {
				fakeDirName := fmt.Sprintf("fake-dir-%v", time.Now().Unix())
				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", fakeDirName, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment", "--target", "stemcell"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 5*time.Second).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring(fmt.Sprintf("%v: no such file or directory", fakeDirName)))
			})

			It("errors when the deployment manifest does not exist", func() {
				err := os.Setenv("ORIGINAL_DEPLOYMENT_MANIFEST_PATH", emptyDir)
				Expect(err).NotTo(HaveOccurred())

				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", emptyDir, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment", "--target", "stemcell"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 5*time.Second).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring("no such file or directory"))
			})

			It("errors when the directory to write out the commit message does not exist", func() {
				err := os.Setenv("COMMIT_MESSAGE_PATH", filepath.Join(emptyDir, "doesnt-exist"))
				Expect(err).NotTo(HaveOccurred())

				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment", "--target", "stemcell"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 5*time.Second).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring("doesnt-exist: no such file or directory"))
			})

			It("errors when a required stemcell file is missing", func() {
				err := os.Remove(filepath.Join(buildDir, "stemcell", "version"))
				Expect(err).NotTo(HaveOccurred())

				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment", "--target", "stemcell"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 5*time.Second).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring("stemcell/version: no such file or directory"))
			})
		})
	})

	Context("compiled releases opsfile", func() {
		const (
			originalOpsFile string = `
---
- path: /releases/name=release1
  type: replace
  value:
    name: release1
    url: https://storage.googleapis.com/cf-deployment-compiled-releases/release1-0.0.0-stemcell1-0.0-20180808-202210-307673159.tgz
    version: 0.0.0
    sha1: 4ee0dfe1f1b9acd14c18863061268f4156c291a4
    stemcell:
      os: stemcell1
      version: "0.0"
- path: /releases/name=release2
  type: replace
  value:
    name: release2
    url: https://storage.googleapis.com/cf-deployment-compiled-releases/release2-0.0.1-stemcell1-0.0-20180808-202210-307673159.tgz
    version: 0.0.1
    sha1: 5ee0dfe1f1b9acd14c18863061268f4156c291a4
    stemcell:
      os: stemcell1
      version: "0.0"
`
			expectedOpsFile string = `
---
- path: /releases/name=release1
  type: replace
  value:
    name: release1
    url: https://storage.googleapis.com/cf-deployment-compiled-releases/release1-0.2.0-stemcell2-2.0-20180808-195254-497840039.tgz
    version: 0.2.0
    sha1: 8867c88b56e0bfb82cffaf15a66bc8d107d6754a
    stemcell:
      os: stemcell2
      version: "2.0"
- path: /releases/name=release2
  type: replace
  value:
    name: release2
    url: https://storage.googleapis.com/cf-deployment-compiled-releases/release2-0.0.1-stemcell1-0.0-20180808-202210-307673159.tgz
    version: 0.0.1
    sha1: 5ee0dfe1f1b9acd14c18863061268f4156c291a4
    stemcell:
      os: stemcell1
      version: "0.0"
`
		)

		BeforeEach(func() {
			for _, dir := range []string{
				"original-compiled-releases-ops-file",
				"updated-compiled-releases-ops-file",
			} {
				err := os.Mkdir(filepath.Join(buildDir, dir), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())
			}

			err := os.WriteFile(filepath.Join(buildDir, "original-compiled-releases-ops-file", "original_ops_file.yml"), []byte(originalOpsFile), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			inputPath = os.Getenv("ORIGINAL_OPS_FILE_PATH")
			err = os.Setenv("ORIGINAL_OPS_FILE_PATH", "original_ops_file.yml")
			Expect(err).NotTo(HaveOccurred())

			outputPath = os.Getenv("UPDATED_OPS_FILE_PATH")
			err = os.Setenv("UPDATED_OPS_FILE_PATH", "updated_ops_file.yml")
			Expect(err).NotTo(HaveOccurred())

			for _, release := range []map[string]string{
				{"name": "release1", "version": "0.2.0", "url": "https://storage.googleapis.com/cf-deployment-compiled-releases/release1-0.0.0-stemcell1-0.0-20180808-202210-307673159.tgz", "sha1": "4ee0dfe1f1b9acd14c18863061268f4156c291a4"},
			} {
				releaseDir := filepath.Join(buildDir, fmt.Sprintf("%s-release", release["name"]))
				err = os.Mkdir(releaseDir, os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				for _, value := range []string{"version", "url", "sha1"} {
					err = os.WriteFile(filepath.Join(releaseDir, value), []byte(release[value]), os.ModePerm)
					Expect(err).NotTo(HaveOccurred())
				}

				compiledReleaseDir := filepath.Join(buildDir, fmt.Sprintf("%s-compiled-release-tarball", release["name"]))
				err = os.Mkdir(compiledReleaseDir, os.ModePerm)
				Expect(err).NotTo(HaveOccurred())

				compiledReleaseTarballName := "release1-0.2.0-stemcell2-2.0-20180808-195254-497840039.tgz"
				err = os.WriteFile(filepath.Join(compiledReleaseDir, compiledReleaseTarballName), []byte("anything"), os.ModePerm)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		AfterEach(func() {
			err := os.Setenv("ORIGINAL_OPS_FILE_PATH", inputPath)
			Expect(err).NotTo(HaveOccurred(), "The original ORIGINAL_OPS_FILE_PATH env var should be set back")
			err = os.Setenv("UPDATED_OPS_FILE_PATH", outputPath)
			Expect(err).NotTo(HaveOccurred(), "The original UPDATED_OPS_FILE_PATH env var should be set back")
		})

		It("updates the original ops file with new releases", func() {
			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "original-compiled-releases-ops-file", "--output-dir", "updated-compiled-releases-ops-file", "--target", "compiledReleasesOpsfile"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 5*time.Second).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			updatedOpsFile, err := os.ReadFile(filepath.Join(buildDir, "updated-compiled-releases-ops-file", "updated_ops_file.yml"))
			Expect(err).NotTo(HaveOccurred())

			Expect(updatedOpsFile).To(MatchYAML(expectedOpsFile))
		})

		It("does not overwrite the commit message if it says that there are changes", func() {
			err := os.WriteFile(filepath.Join(buildDir, os.Getenv("COMMIT_MESSAGE_PATH")), []byte("previous commit message with changes"), 0666)
			Expect(err).NotTo(HaveOccurred())

			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "original-compiled-releases-ops-file", "--output-dir", "updated-compiled-releases-ops-file", "--target", "compiledReleasesOpsfile"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 5*time.Second).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			commitMessage, err := os.ReadFile(filepath.Join(buildDir, "commit-message.txt"))
			Expect(err).NotTo(HaveOccurred())

			Expect(string(commitMessage)).To(Equal("previous commit message with changes"))
		})

		It("updates commit message if current commit message says that there are no changes", func() {
			err := os.WriteFile(filepath.Join(buildDir, os.Getenv("COMMIT_MESSAGE_PATH")), []byte("No manifest release or stemcell version updates"), 0666)
			Expect(err).NotTo(HaveOccurred())

			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "original-compiled-releases-ops-file", "--output-dir", "updated-compiled-releases-ops-file", "--target", "compiledReleasesOpsfile"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 5*time.Second).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			commitMessage, err := os.ReadFile(filepath.Join(buildDir, "commit-message.txt"))
			Expect(err).NotTo(HaveOccurred())

			Expect(string(commitMessage)).To(Equal("Updated compiled releases with release1 0.2.0"))
		})

		It("creates nested directory when the directory to write out the updated ops file path does not exist", func() {
			err := os.Setenv("UPDATED_OPS_FILE_PATH", filepath.Join("doesnt-exist", "updated_ops_file.yml"))
			Expect(err).NotTo(HaveOccurred())

			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "original-compiled-releases-ops-file", "--output-dir", "updated-compiled-releases-ops-file", "--target", "compiledReleasesOpsfile"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session, 5*time.Second).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			updatedOpsFile, err := os.ReadFile(filepath.Join(buildDir, "updated-compiled-releases-ops-file", "doesnt-exist", "updated_ops_file.yml"))
			Expect(err).NotTo(HaveOccurred())

			Expect(updatedOpsFile).To(MatchYAML(expectedOpsFile))
		})

		Context("failure cases", func() {
			It("errors when the build dir does not exist", func() {
				fakeDirName := fmt.Sprintf("fake-dir-%v", time.Now().Unix())
				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", fakeDirName, "--input-dir", "original-compiled-releases-ops-file", "--output-dir", "updated-compiled-releases-ops-file", "--target", "compiledReleasesOpsfile"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 5*time.Second).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring(fmt.Sprintf("%v: no such file or directory", fakeDirName)))
			})

			It("errors when the original ops file does not exist", func() {
				err := os.Setenv("ORIGINAL_OPS_FILE_PATH", emptyDir)
				Expect(err).NotTo(HaveOccurred())

				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", emptyDir, "--input-dir", "original-compiled-releases-ops-file", "--output-dir", "updated-compiled-releases-ops-file", "--target", "compiledReleasesOpsfile"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session, 5*time.Second).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring("no such file or directory"))
			})
		})
	})
})
