package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/runtime-ci/tasks/update-stemcell/manifest"
)

var _ = Describe("Update Manifest", func() {
	var (
		contentArg []byte
		stemcellArg manifest.Stemcell

		actualContent []byte
		actualError error
	)

	BeforeEach(func() {
		contentArg = nil
		stemcellArg = manifest.Stemcell{}
	})

	JustBeforeEach(func() {
		actualContent, actualError = manifest.Update(contentArg, stemcellArg)
	})

	Context("when provided a valid manifest and stemcell", func() {
		BeforeEach(func() {
			contentArg = []byte(`---
name: cf
manifest_version: v9.9.9
update:
	some-update-values: 0
addons:
- name: some-addon
instance_groups:
- name: some-instance
  jobs:
	- name: some-job
releases:
- name: some-release
stemcells:
- alias: some-stemcell
  os: some-old-os
	version: some-old-version
`)

			stemcellArg.OS = "some-new-os"
			stemcellArg.Version = "some-new-version"
		})

		It("should update the stemcell without modifying the shape", func(){
			Expect(actualError).ToNot(HaveOccurred())

			Expect(string(actualContent)).To(Equal(`---
name: cf
manifest_version: v9.9.9
update:
	some-update-values: 0
addons:
- name: some-addon
instance_groups:
- name: some-instance
  jobs:
	- name: some-job
releases:
- name: some-release
stemcells:
- alias: some-stemcell
  os: some-new-os
	version: some-new-version
`))
		})
	})
})
