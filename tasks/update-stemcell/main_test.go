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

	const (
		expectedManifest = `
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

		originalOpsFile string = `
---
- path: /releases/name=release1/url
  type: replace
  value: https://storage.googleapis.com/cf-deployment-compiled-releases/release1-0.0.0-stemcell1-0.0-20180808-202210-307673159.tgz
- path: /releases/name=release1/version
  type: replace
  value: 0.0.0
- path: /releases/name=release1/sha1
  type: replace
  value: 4ee0dfe1f1b9acd14c18863061268f4156c291a4
- path: /releases/name=release1/stemcell?
  type: replace
  value:
    os: stemcell1
    version: "0.0"
- path: /releases/name=release2/url
  type: replace
  value: https://storage.googleapis.com/cf-deployment-compiled-releases/release2-0.0.1-stemcell1-0.0-20180808-202210-307673159.tgz
- path: /releases/name=release2/version
  type: replace
  value: 0.0.1
- path: /releases/name=release2/sha1
  type: replace
  value: 5ee0dfe1f1b9acd14c18863061268f4156c291a4
- path: /releases/name=release2/stemcell?
  type: replace
  value:
    os: stemcell1
    version: "0.0"
`
		expectedOpsFile string = `
---
- path: /releases/name=release1/url
  type: replace
  value: https://storage.googleapis.com/cf-deployment-compiled-releases/release1-0.0.0-stemcell2-2.0-20180808-195254-497840039.tgz
- path: /releases/name=release1/version
  type: replace
  value: 0.0.0
- path: /releases/name=release1/sha1
  type: replace
  value: 2062e90b3ea10a86ff666a76c41aa0d9e9d88f4e
- path: /releases/name=release1/stemcell?
  type: replace
  value:
    os: stemcell2
    version: "2.0"
- path: /releases/name=release2/url
  type: replace
  value: https://storage.googleapis.com/cf-deployment-compiled-releases/release2-0.0.1-stemcell2-2.0-20180808-195254-497840039.tgz
- path: /releases/name=release2/version
  type: replace
  value: 0.0.1
- path: /releases/name=release2/sha1
  type: replace
  value: 5dfc64fff3b07c7c01ebd39706ec3cf3e6c37464
- path: /releases/name=release2/stemcell?
  type: replace
  value:
    os: stemcell2
    version: "2.0"
`
	)

	BeforeEach(func() {
		pathToBinary, err = gexec.Build("github.com/cloudfoundry/runtime-ci/tasks/update-stemcell")
		Expect(err).NotTo(HaveOccurred())

		buildDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		emptyDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		CommitMessagePath = os.Getenv("COMMIT_MESSAGE_PATH")
		os.Setenv("COMMIT_MESSAGE_PATH", "commit-message.txt")

		for _, dir := range []string{
			"cf-deployment",
			"updated-cf-deployment",
			"stemcell",
		} {
			err = os.Mkdir(filepath.Join(buildDir, dir), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
		}

		err = ioutil.WriteFile(filepath.Join(buildDir, "cf-deployment", "original-manifest.yml"), []byte(originalManifest), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		err = ioutil.WriteFile(filepath.Join(buildDir, "stemcell", "version"), []byte("updated-stemcell-version"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		err = ioutil.WriteFile(filepath.Join(buildDir, "stemcell", "url"), []byte("https://foo.com/bosh-stemcell-ubuntu-trusty.tgz"), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		inputPath = os.Getenv("ORIGINAL_DEPLOYMENT_MANIFEST_PATH")
		outputPath = os.Getenv("UPDATED_DEPLOYMENT_MANIFEST_PATH")

		os.Setenv("ORIGINAL_DEPLOYMENT_MANIFEST_PATH", "original-manifest.yml")
		os.Setenv("UPDATED_DEPLOYMENT_MANIFEST_PATH", "updated-manifest.yml")

		err = ioutil.WriteFile(filepath.Join(buildDir, "cf-deployment", "original_ops_file.yml"), []byte(originalOpsFile), os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		inputPath = os.Getenv("ORIGINAL_OPS_FILE_PATH")
		outputPath = os.Getenv("UPDATED_OPS_FILE_PATH")

		os.Setenv("ORIGINAL_OPS_FILE_PATH", "original_ops_file.yml")
		os.Setenv("UPDATED_OPS_FILE_PATH", "updated_ops_file.yml")

		for _, release := range []map[string]string{
			{"name": "release1", "version": "0.0.0"},
			{"name": "release2", "version": "0.0.1"},
		} {
			compiledReleaseDir := filepath.Join(buildDir, fmt.Sprintf("%s-compiled-release-tarball", release["name"]))
			err = os.Mkdir(compiledReleaseDir, os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			compiledReleaseTarballName := fmt.Sprintf("%s-%s-stemcell2-2.0-20180808-195254-497840039.tgz", release["name"], release["version"])
			err = ioutil.WriteFile(filepath.Join(compiledReleaseDir, compiledReleaseTarballName), []byte(release["name"]), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
		}
	})

	AfterEach(func() {
		os.Setenv("COMMIT_MESSAGE_PATH", CommitMessagePath)

		os.Setenv("ORIGINAL_DEPLOYMENT_MANIFEST_PATH", inputPath)
		os.Setenv("UPDATED_DEPLOYMENT_MANIFEST_PATH", outputPath)

		os.Setenv("ORIGINAL_OPS_FILE_PATH", inputPath)
		os.Setenv("UPDATED_OPS_FILE_PATH", outputPath)
	})

	Context("required parameters", func() {
		Context("when the build-dir flag is missing", func() {
			It("returns an error", func() {
				session, err := gexec.Start(exec.Command(pathToBinary), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1), "expected command to fail")

				Expect(string(session.Err.Contents())).To(ContainSubstring("missing required flag: build-dir"))
			})
		})

		Context("when the input-dir flag is missing", func() {
			It("returns an error", func() {
				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1), "expected command to fail")

				Expect(string(session.Err.Contents())).To(ContainSubstring("missing required flag: input-dir"))
			})
		})

		Context("when the output-dir flag is missing", func() {
			It("returns an error", func() {
				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "cf-deployment"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1), "expected command to fail")

				Expect(string(session.Err.Contents())).To(ContainSubstring("missing required flag: output-dir"))
			})
		})

		Context("when the path to the input deployment manifest is missing", func() {
			BeforeEach(func() {
				os.Setenv("ORIGINAL_DEPLOYMENT_MANIFEST_PATH", "")
			})

			It("returns an error", func() {
				fakeDirName := fmt.Sprintf("fake-dir-%v", time.Now().Unix())
				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", fakeDirName, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring("missing path to input deployment manifest"))

			})
		})

		Context("when the path to the output deployment manifest is missing", func() {
			BeforeEach(func() {
				os.Setenv("UPDATED_DEPLOYMENT_MANIFEST_PATH", "")
			})

			It("returns an error", func() {
				fakeDirName := fmt.Sprintf("fake-dir-%v", time.Now().Unix())
				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", fakeDirName, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring("missing path to output deployment manifest"))

			})
		})

		Context("when the path to the input compiled release ops-file is missing", func() {
			BeforeEach(func() {
				os.Setenv("ORIGINAL_OPS_FILE_PATH", "")
			})

			It("returns an error", func() {
				fakeDirName := fmt.Sprintf("fake-dir-%v", time.Now().Unix())
				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", fakeDirName, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring("missing path to input compiled release ops-file"))

			})
		})

		Context("when the path to the output compiled release ops-file is missing", func() {
			BeforeEach(func() {
				os.Setenv("UPDATED_OPS_FILE_PATH", "")
			})

			It("returns an error", func() {
				fakeDirName := fmt.Sprintf("fake-dir-%v", time.Now().Unix())
				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", fakeDirName, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring("missing path to output compiled release ops-file"))

			})
		})
	})

	Context("stemcell", func() {
		It("only updates the manifest with the new stemcell", func() {
			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			updatedManifest, err := ioutil.ReadFile(filepath.Join(buildDir, "updated-cf-deployment", "updated-manifest.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedManifest).To(MatchYAML(expectedManifest))
		})

		It("writes the commit message to COMMIT_MESSAGE_PATH", func() {
			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			commitMessage, err := ioutil.ReadFile(filepath.Join(buildDir, "commit-message.txt"))
			Expect(err).NotTo(HaveOccurred())

			Expect(string(commitMessage)).To(Equal("Updated manifest with ubuntu-trusty stemcell updated-stemcell-version"))
		})

		It("creates nested directory when the directory to write out the updated manifest does not exist", func() {
			os.Setenv("UPDATED_DEPLOYMENT_MANIFEST_PATH", filepath.Join("doesnt-exist", "updated-manifest.yml"))

			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			updatedManifest, err := ioutil.ReadFile(filepath.Join(buildDir, "updated-cf-deployment", "doesnt-exist", "updated-manifest.yml"))
			Expect(err).NotTo(HaveOccurred())

			Expect(updatedManifest).To(MatchYAML(expectedManifest))
		})

		Context("failure cases", func() {
			It("errors when the build dir does not exist", func() {
				fakeDirName := fmt.Sprintf("fake-dir-%v", time.Now().Unix())
				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", fakeDirName, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring(fmt.Sprintf("%v/cf-deployment/original-manifest.yml: no such file or directory", fakeDirName)))
			})

			It("errors when the deployment manifest does not exist", func() {
				os.Setenv("ORIGINAL_DEPLOYMENT_MANIFEST_PATH", emptyDir)

				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", emptyDir, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring("no such file or directory"))
			})

			It("errors when the directory to write out the commit message does not exist", func() {
				os.Setenv("COMMIT_MESSAGE_PATH", filepath.Join(emptyDir, "doesnt-exist"))

				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring("doesnt-exist: no such file or directory"))
			})

			It("errors when a required stemcell file is missing", func() {
				err := os.Remove(filepath.Join(buildDir, "stemcell", "version"))
				Expect(err).NotTo(HaveOccurred())

				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring("stemcell/version: no such file or directory"))
			})
		})
	})

	Context("compiled releases opsfile", func() {
		It("updates the original ops file with new stemcell", func() {
			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment", "--target", "compiledReleasesOpsfile"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			updatedOpsFile, err := ioutil.ReadFile(filepath.Join(buildDir, "updated-cf-deployment", "updated_ops_file.yml"))
			Expect(err).NotTo(HaveOccurred())

			Expect(updatedOpsFile).To(MatchYAML(expectedOpsFile))
		})

		It("creates nested directory when the directory to write out the updated ops file path does not exist", func() {
			os.Setenv("UPDATED_OPS_FILE_PATH", filepath.Join("doesnt-exist", "updated_ops_file.yml"))

			session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", buildDir, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment", "--target", "compiledReleasesOpsfile"}...), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit())
			Expect(session.ExitCode()).To(Equal(0))

			updatedOpsFile, err := ioutil.ReadFile(filepath.Join(buildDir, "updated-cf-deployment", "doesnt-exist", "updated_ops_file.yml"))
			Expect(err).NotTo(HaveOccurred())

			Expect(updatedOpsFile).To(MatchYAML(expectedOpsFile))
		})

		Context("failure cases", func() {
			It("errors when the original ops file does not exist", func() {
				os.Setenv("ORIGINAL_OPS_FILE_PATH", emptyDir)

				session, err := gexec.Start(exec.Command(pathToBinary, []string{"--build-dir", emptyDir, "--input-dir", "cf-deployment", "--output-dir", "updated-cf-deployment", "--target", "compiledReleasesOpsfile"}...), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit())
				Expect(session.ExitCode()).To(Equal(1))

				Expect(string(session.Err.Contents())).To(ContainSubstring("no such file or directory"))
			})
		})
	})
})
