package main_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

const (
	expectedOpsFile string = `
---
- type: replace
path: /releases/-
value:
	name: release1
	version: new-release1-version
	url: new-release1-url
	sha1: new-release1-sha

- type: replace
path: /releases/-
value:
	name: release2
	version: original-release2-version
	url: original-release2-url
	sha1: original-release2-sha
	`

	opsFileStub string = `
---
- type: replace
path: /releases/-
value:
	name: release1
	version: original-release1-version
	url: original-release1-url
	sha1: original-release1-sha

- type: replace
path: /releases/-
value:
	name: release2
	version: original-release2-version
	url: original-release2-url
	sha1: original-release2-sha
	`
)

var _ = Describe("main", func() {
	var (
		pathToBinary string
		buildDir     string
	)

	BeforeEach(func() {
		var err error

		pathToBinary, err = gexec.Build("github.com/cloudfoundry/runtime-ci/util/update-ops-file-releases")
		Expect(err).NotTo(HaveOccurred())

		buildDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())
	})
})
