package stemcell_test

import (
	"bytes"
	"errors"
	"io"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/runtime-ci/tasks/export-all-compiled-release-tarballs/stemcell"
	"github.com/cloudfoundry/runtime-ci/tasks/export-all-compiled-release-tarballs/stemcell/stemcellfakes"
)

var _ = Describe("List", func() {
	var (
		fakeCLI *stemcellfakes.FakeBoshCLI

		returnedReader io.Reader
		returnedError  error

		actualStemcells []stemcell.Stemcell
		actualErr      error
	)

	BeforeEach(func() {
		fakeCLI = new(stemcellfakes.FakeBoshCLI)

		returnedReader = new(bytes.Buffer)
		returnedError = nil
	})

	JustBeforeEach(func() {
		fakeCLI.CmdReturns(returnedReader, returnedError)

		actualStemcells, actualErr = stemcell.List(fakeCLI)
	})

	It("should call `bosh stemcells --json`", func() {
		Expect(fakeCLI.CmdCallCount()).To(Equal(1), "expected bosh cli calls")
		name, args := fakeCLI.CmdArgsForCall(0)

		Expect(name).To(Equal("stemcells"), "expected command name")
		Expect(args).To(ConsistOf("--json"), "expected command args")
	})

	Context("when the `bosh stemcells` returns successfully", func() {
		Context("when there is at least one stemcell", func() {
			BeforeEach(func() {
				output := `{
    "Tables": [
        {
            "Content": "stemcells",
            "Header": {
                "cid": "CID",
                "cpi": "CPI",
                "name": "Name",
                "os": "OS",
                "version": "Version"
            },
            "Rows": [
                {
                    "cid": "stemcell-id",
                    "cpi": "",
                    "name": "some-stemcell",
                    "os": "some-os",
                    "version": "1.2*"
                }
            ]
        }
    ]
}`

				returnedReader = strings.NewReader(output)
			})

			It("returns the found stemcell", func() {
				Expect(actualErr).ToNot(HaveOccurred())
				Expect(actualStemcells).To(ConsistOf(
					stemcell.Stemcell{
						Name: "some-stemcell",
						OS: "some-os",
						Version: "1.2",
					},
				))
			})
		})
	})

	Context("when the `bosh stemcells` fails", func() {
		BeforeEach(func() {
			returnedError = errors.New("some error")
		})

		It("returns the error", func() {
			expectedErr := errors.New("some error")
			Expect(actualErr).To(MatchError(expectedErr))
		})
	})
})
