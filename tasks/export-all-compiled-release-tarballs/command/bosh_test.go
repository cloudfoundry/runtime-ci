package command

import (
	"errors"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("parseErr", func() {
	var (
		readerArg *strings.Reader
		errArg    error

		actualErr error
	)

	BeforeEach(func() {
		readerArg = new(strings.Reader)
		errArg = errors.New("")
	})

	JustBeforeEach(func() {
		actualErr = parseErr(readerArg, errArg)
	})

	Context("with a valid block error", func() {
		BeforeEach(func() {
			readerArg = strings.NewReader(`{
    "Tables": null,
    "Blocks": [
        "\nTask 36835 | 23:50:15 | ",
        "Error: Release 'bad' doesn't exist"
    ]
}`)
		})

		It("returns the error message", func() {
			Expect(actualErr).To(MatchError("Error: Release 'bad' doesn't exist"))
		})
	})
})
