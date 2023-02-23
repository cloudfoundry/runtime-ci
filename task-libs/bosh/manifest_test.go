package bosh_test

import (
	"bytes"
	"io"
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/runtime-ci/task-libs/bosh"
	"github.com/cloudfoundry/runtime-ci/task-libs/bosh/boshfakes"
)

var _ = Describe("Manifest", func() {
	Describe("NewManifestFromFile", func() {
		var (
			fileArg []byte

			actualManifest Manifest
			actualError    error
		)

		JustBeforeEach(func() {
			actualManifest, actualError = NewManifestFromFile(fileArg)
		})

		BeforeEach(func() {
			fileArg = []byte(`---
name: some-name
stemcells:
- alias: default
  os: some-os
  version: some-version
`)
		})

		It("returns the manifest", func() {
			Expect(actualError).ToNot(HaveOccurred())

			Expect(actualManifest).To(Equal(Manifest{
				Name: "some-name",
				Stemcells: []Stemcell{
					{
						Alias:   "default",
						OS:      "some-os",
						Version: "some-version",
					},
				},
			}))
		})
	})

	Describe("Deploy", func() {
		var (
			manifestArg Manifest

			fakeBoshCLI *boshfakes.FakeBoshCLI

			actualManifest []byte
			actualError    error
		)

		BeforeEach(func() {
			fakeBoshCLI = new(boshfakes.FakeBoshCLI)
			fakeBoshCLI.CmdStub = func(cmd string, args ...string) (io.Reader, error) {
				path := args[0]

				var err error
				actualManifest, err = os.ReadFile(path)
				return new(bytes.Buffer), err
			}
		})

		JustBeforeEach(func() {
			actualError = manifestArg.Deploy(fakeBoshCLI)
		})

		Context("when the manifest is partially filled", func() {
			BeforeEach(func() {
				manifestArg = Manifest{
					Name:      "cf-compilation",
					Releases:  []Release{{Name: "release-a"}, {Name: "release-b"}},
					Stemcells: []Stemcell{{OS: "some-os", Version: "1.2.3"}},
				}
			})

			It("runs the bosh deploy command with the manifest file", func() {
				Expect(actualError).ToNot(HaveOccurred())

				Expect(fakeBoshCLI.CmdCallCount()).To(Equal(1), "expected boshCLI call count")

				cmd, args := fakeBoshCLI.CmdArgsForCall(0)
				Expect(cmd).To(Equal("deploy"), "expected boshCLI command")
				Expect(args).To(HaveLen(5), "expected boshCLI arg len")
				Expect(strings.Join(args[1:], " ")).To(Equal("-d cf-compilation -n --json"))

				Expect(string(actualManifest)).To(Equal(`name: cf-compilation
update:
    canaries: 1
    canary_watch_time: 1
    max_in_flight: 1
    update_watch_time: 1
releases:
    - name: release-a
      sha1: ""
      url: ""
      version: ""
    - name: release-b
      sha1: ""
      url: ""
      version: ""
stemcells:
    - alias: default
      os: some-os
      version: 1.2.3
`))
			})
		})
	})
})
