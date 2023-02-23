package command

import (
	"errors"
	"strings"

	. "github.com/onsi/ginkgo/v2"
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

	Context("with an empty block error", func() {
		BeforeEach(func() {
			readerArg = strings.NewReader(`{
		"Tables": null,
		"Blocks": null,
		"Lines": [
			"Using environment 'https://10.0.0.6:25555' as client 'admin'",
			"Expected non-empty deployment name",
			"Exit code 1"
		]
}`)
		})

		It("returns the filtered error message", func() {
			Expect(actualErr).To(MatchError("Expected non-empty deployment name"))
		})
	})
	{
	}
})
