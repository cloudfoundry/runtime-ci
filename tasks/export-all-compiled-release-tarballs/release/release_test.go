package release_test

import (
	"bytes"
	"errors"
	"io"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/runtime-ci/tasks/export-all-compiled-release-tarballs/release"
	"github.com/cloudfoundry/runtime-ci/tasks/export-all-compiled-release-tarballs/release/releasefakes"
	"github.com/cloudfoundry/runtime-ci/tasks/export-all-compiled-release-tarballs/stemcell"
)

var _ = Describe("List", func() {
	var (
		fakeCLI *releasefakes.FakeBoshCLI

		returnedReader io.Reader
		returnedError  error

		actualReleases []release.Release
		actualErr      error
	)

	BeforeEach(func() {
		fakeCLI = new(releasefakes.FakeBoshCLI)

		returnedReader = new(bytes.Buffer)
		returnedError = nil
	})

	JustBeforeEach(func() {
		fakeCLI.CmdReturns(returnedReader, returnedError)

		actualReleases, actualErr = release.List(fakeCLI)
	})

	It("should call `bosh releases --json`", func() {
		Expect(fakeCLI.CmdCallCount()).To(Equal(1), "expected bosh cli calls")
		name, args := fakeCLI.CmdArgsForCall(0)

		Expect(name).To(Equal("releases"), "expected command name")
		Expect(args).To(ConsistOf("--json"), "expected command args")
	})

	Context("when the `bosh releases` returns successfully", func() {
		Context("when there is at least one release", func() {
			BeforeEach(func() {
				output := `{
    "Tables": [
        {
            "Content": "releases",
            "Header": {
                "commit_hash": "Commit Hash",
                "name": "Name",
                "version": "Version"
            },
            "Rows": [
                {
                    "name": "release-a",
                    "version": "0.1.0*"
                },
                {
                    "name": "release-b",
                    "version": "2.0.0*"
                }
            ]
        }
    ]
}`

				returnedReader = strings.NewReader(output)
			})

			It("returns the found releases", func() {
				Expect(actualErr).ToNot(HaveOccurred())
				Expect(actualReleases).To(ConsistOf(
					release.Release{Name: "release-a", Version: "0.1.0"},
					release.Release{Name: "release-b", Version: "2.0.0"},
				))
			})
		})

		Context("when bosh-dns is included in the list", func() {
			BeforeEach(func() {
				output := `{
    "Tables": [
        {
            "Content": "releases",
            "Header": {
                "commit_hash": "Commit Hash",
                "name": "Name",
                "version": "Version"
            },
            "Rows": [
                {
                    "name": "release-a",
                    "version": "0.1.0*"
                },
                {
                    "name": "bosh-dns",
                    "version": "0.3.0*"
                },
                {
                    "name": "release-b",
                    "version": "2.0.0*"
                }
            ]
        }
    ]
}`

				returnedReader = strings.NewReader(output)
			})

			It("excludes bosh-dns from the list of returned releases", func() {
				Expect(actualErr).ToNot(HaveOccurred())
				Expect(actualReleases).To(ConsistOf(
					release.Release{Name: "release-a", Version: "0.1.0"},
					release.Release{Name: "release-b", Version: "2.0.0"},
				))
			})
		})
	})

	Context("when the `bosh releases` fails", func() {
		BeforeEach(func() {
			returnedError = errors.New("some error")
		})

		It("returns the error", func() {
			expectedErr := errors.New("some error")
			Expect(actualErr).To(MatchError(expectedErr))
		})
	})
})

var _ = Describe("Export", func() {
	var (
		fakeCLI     *releasefakes.FakeBoshCLI
		releaseArg  release.Release
		stemcellArg stemcell.Stemcell

		returnedReader io.Reader
		returnedError  error

		actualErr error
	)

	BeforeEach(func() {
		fakeCLI = new(releasefakes.FakeBoshCLI)

		returnedReader = new(bytes.Buffer)
		returnedError = nil
	})

	JustBeforeEach(func() {
		fakeCLI.CmdReturns(returnedReader, returnedError)

		actualErr = release.Export(fakeCLI, releaseArg, stemcellArg)
	})

	Context("when a valid release and os are passed in", func() {

		BeforeEach(func() {
			releaseArg = release.Release{Name: "some-release", Version: "some-release-version"}
			stemcellArg = stemcell.Stemcell{OS: "some-os", Version: "some-os-version"}
		})

		It("should call `bosh export-release --json`", func() {
			Expect(fakeCLI.CmdCallCount()).To(Equal(1), "expected bosh cli calls")
			name, args := fakeCLI.CmdArgsForCall(0)

			Expect(name).To(Equal("export-release"), "expected command name")
			Expect(args).To(Equal([]string{
				"--json",
				"some-release/some-release-version",
				"some-os/some-os-version",
			}), "expected command args")
		})

		Context("when it successfully exports a release", func() {
			It("successfully returns", func() {
				Expect(actualErr).ToNot(HaveOccurred())
			})
		})

		Context("when the `bosh export-release` command fails", func() {
			BeforeEach(func() {
				returnedError = errors.New("some error")
			})

			It("returns the error", func() {
				Expect(actualErr).To(MatchError("some error"))
			})
		})
	})
})
