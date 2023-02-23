package deployment_test

import (
	"bytes"
	"errors"
	"io"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/runtime-ci/tasks/export-all-compiled-release-tarballs/deployment"
	"github.com/cloudfoundry/runtime-ci/tasks/export-all-compiled-release-tarballs/deployment/deploymentfakes"
	"github.com/cloudfoundry/runtime-ci/tasks/export-all-compiled-release-tarballs/stemcell"
)

var _ = Describe("List", func() {
	var (
		fakeCLI *deploymentfakes.FakeBoshCLI

		returnedReader io.Reader
		returnedError  error

		actualDeployments []deployment.Deployment
		actualErr         error
	)

	BeforeEach(func() {
		fakeCLI = new(deploymentfakes.FakeBoshCLI)

		returnedReader = new(bytes.Buffer)
		returnedError = nil
	})

	JustBeforeEach(func() {
		fakeCLI.CmdReturns(returnedReader, returnedError)

		actualDeployments, actualErr = deployment.List(fakeCLI, []stemcell.Stemcell{{Name: "stemcell-name", OS: "stemcell-os", Version: "1.2"}})
	})

	It("should call `bosh deployments --json`", func() {
		Expect(fakeCLI.CmdCallCount()).To(Equal(1), "expected bosh cli calls")
		name, args := fakeCLI.CmdArgsForCall(0)

		Expect(name).To(Equal("deployments"), "expected command name")
		Expect(args).To(ConsistOf("--json"), "expected command args")
	})

	Context("when the `bosh deployments` returns successfully", func() {
		Context("when there is at least one release", func() {
			BeforeEach(func() {
				output := `{
    "Tables": [
        {
            "Content": "deployments",
            "Header": {
                "name": "Name",
                "release_s": "Release(s)",
                "stemcell_s": "Stemcell(s)",
                "team_s": "Team(s)"
            },
            "Rows": [
                {
                    "name": "cf-compilation-release-a",
                    "release_s": "release-a/0.1.0",
                    "stemcell_s": "stemcell-name/1.2",
                    "team_s": ""
                },
                {
                    "name": "cf-compilation-release-b",
                    "release_s": "release-b/2.0.0",
                    "stemcell_s": "stemcell-name/1.2",
                    "team_s": ""
                }
            ],
            "Notes": null
        }
    ],
    "Blocks": null,
    "Lines": [
        "Using environment 'https://10.0.0.6:25555' as client 'admin'",
        "Succeeded"
    ]
}`

				returnedReader = strings.NewReader(output)
			})

			It("returns the found releases", func() {
				Expect(actualErr).ToNot(HaveOccurred())
				Expect(actualDeployments).To(ConsistOf(
					deployment.Deployment{
						Name: "cf-compilation-release-a",
						Releases: []deployment.Release{
							{Name: "release-a", Version: "0.1.0"},
						},
						Stemcell: stemcell.Stemcell{OS: "stemcell-os", Version: "1.2"},
					},
					deployment.Deployment{
						Name: "cf-compilation-release-b",
						Releases: []deployment.Release{
							{Name: "release-b", Version: "2.0.0"},
						},
						Stemcell: stemcell.Stemcell{OS: "stemcell-os", Version: "1.2"},
					},
				))
			})
		})

		Context("when there are multiple releases in a single deployment", func() {
			BeforeEach(func() {
				output := `{
    "Tables": [
        {
            "Content": "deployments",
            "Header": {
                "name": "Name",
                "release_s": "Release(s)",
                "stemcell_s": "Stemcell(s)",
                "team_s": "Team(s)"
            },
            "Rows": [
                {
                    "name": "cf-compilation-releases",
                    "release_s": "release-a/0.1.0\nrelease-b/2.0.0",
                    "stemcell_s": "stemcell-name/1.2",
                    "team_s": ""
                }
            ],
            "Notes": null
        }
    ],
    "Blocks": null,
    "Lines": [
        "Using environment 'https://10.0.0.6:25555' as client 'admin'",
        "Succeeded"
    ]
}`

				returnedReader = strings.NewReader(output)
			})

			It("returns the found releases", func() {
				Expect(actualErr).ToNot(HaveOccurred())
				Expect(actualDeployments).To(ConsistOf(
					deployment.Deployment{
						Name: "cf-compilation-releases",
						Releases: []deployment.Release{
							{Name: "release-a", Version: "0.1.0"},
							{Name: "release-b", Version: "2.0.0"},
						},
						Stemcell: stemcell.Stemcell{OS: "stemcell-os", Version: "1.2"},
					},
				))
			})
		})

		Context("when bosh-dns is included in the list", func() {
			BeforeEach(func() {
				output := `{
    "Tables": [
        {
            "Content": "deployments",
            "Header": {
                "name": "Name",
                "release_s": "Release(s)",
                "stemcell_s": "Stemcell(s)",
                "team_s": "Team(s)"
            },
            "Rows": [
                {
                    "name": "cf-compilation-releases",
                    "release_s": "release-a/0.1.0\nrelease-b/2.0.0\nbosh-dns/0.3.0",
                    "stemcell_s": "stemcell-name/1.2",
                    "team_s": ""
                }
            ],
            "Notes": null
        }
    ],
    "Blocks": null,
    "Lines": [
        "Using environment 'https://10.0.0.6:25555' as client 'admin'",
        "Succeeded"
    ]
}`

				returnedReader = strings.NewReader(output)
			})

			It("excludes bosh-dns from the list of returned releases", func() {
				Expect(actualErr).ToNot(HaveOccurred())
				Expect(actualDeployments).To(ConsistOf(
					deployment.Deployment{
						Name: "cf-compilation-releases",
						Releases: []deployment.Release{
							{Name: "release-a", Version: "0.1.0"},
							{Name: "release-b", Version: "2.0.0"},
						},
						Stemcell: stemcell.Stemcell{OS: "stemcell-os", Version: "1.2"},
					},
				))
			})
		})
	})

	Context("when the `bosh deployments` fails", func() {
		BeforeEach(func() {
			returnedError = errors.New("some error")
		})

		It("returns the error", func() {
			expectedErr := errors.New("some error")
			Expect(actualErr).To(MatchError(expectedErr))
		})
	})
})

var _ = Describe("ExportRelease", func() {
	var (
		fakeCLI       *deploymentfakes.FakeBoshCLI
		deploymentArg deployment.Deployment
		releaseArg    deployment.Release
		stemcellArg   stemcell.Stemcell

		returnedReader io.Reader
		returnedError  error

		actualErr error
	)

	BeforeEach(func() {
		fakeCLI = new(deploymentfakes.FakeBoshCLI)

		returnedReader = new(bytes.Buffer)
		returnedError = nil
	})

	JustBeforeEach(func() {
		fakeCLI.CmdReturns(returnedReader, returnedError)

		actualErr = deployment.ExportRelease(fakeCLI, releaseArg, stemcellArg, deploymentArg)
	})

	Context("when a valid release and os are passed in", func() {

		BeforeEach(func() {
			deploymentArg = deployment.Deployment{Name: "some-deployment"}
			releaseArg = deployment.Release{Name: "some-release", Version: "some-release-version"}
			stemcellArg = stemcell.Stemcell{OS: "some-os", Version: "some-os-version"}
		})

		It("should call `bosh export-release --json`", func() {
			Expect(fakeCLI.CmdCallCount()).To(Equal(1), "expected bosh cli calls")
			name, args := fakeCLI.CmdArgsForCall(0)

			Expect(name).To(Equal("export-release"), "expected command name")
			Expect(args).To(Equal([]string{
				"-d",
				"some-deployment",
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
