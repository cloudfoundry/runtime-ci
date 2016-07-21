package validationerrors_test

import (
	"fmt"

	. "github.com/cloudfoundry/runtime-ci/scripts/ci/run-cats/validationerrors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ValidationErrors", func() {
	It("Aggregates errors", func() {
		ve := Errors{}

		ve.Add(fmt.Errorf("error!!!!"))
		ve.Add(fmt.Errorf("error 2!!!!"))

		Expect(ve.Error()).To(And(
			ContainSubstring("error!!!!"),
			ContainSubstring("error 2!!!!"),
		))
	})

	It("Knows if it is empty", func() {
		ve := Errors{}
		Expect(ve.Empty()).To(BeTrue())
	})
})
