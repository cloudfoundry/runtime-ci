package bosh_test

import (
	"io/ioutil"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/runtime-ci/task-libs/bosh"
)

var _ = Describe("Stemcell", func() {
	var (
		stemcellDirArg string

		actualStemcell Stemcell
		actualErr      error
	)

	BeforeEach(func() {
		var err error
		stemcellDirArg, err = ioutil.TempDir("", "stemcell-input-")
		Expect(err).ToNot(HaveOccurred())
	})

	JustBeforeEach(func() {
		actualStemcell, actualErr = NewStemcellFromInput(stemcellDirArg)
	})

	Context("when the stemcell dir contains all necessary files", func() {
		BeforeEach(func() {
			Expect(ioutil.WriteFile(filepath.Join(stemcellDirArg, "version"), []byte("some-version"), 0777)).
				To(Succeed())
			Expect(ioutil.WriteFile(filepath.Join(stemcellDirArg, "url"), []byte("https://s3.amazonaws.com/some-stemcell/stuff-ubuntu-some-os-go_agent.tgz"), 0777)).
				To(Succeed())
		})

		It("sets the stemcell OS and Version", func() {
			Expect(actualErr).ToNot(HaveOccurred())

			Expect(actualStemcell).To(Equal(Stemcell{OS: "ubuntu-some-os", Version: "some-version"}))
		})
	})

	Context("when the stemcell dir is missing some files", func() {
		It("returns a missing files err", func() {
			Expect(actualErr.Error()).To(ContainSubstring("missing files"))
		})
	})
})
