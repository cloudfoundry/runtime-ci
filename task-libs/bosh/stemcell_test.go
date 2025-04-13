package bosh_test

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/runtime-ci/task-libs/bosh"
)

var _ = Describe("Stemcell", func() {
	Describe("NewStemcellFromInput", func() {
		var (
			stemcellDirArg string

			actualStemcell Stemcell
			actualErr      error
		)

		BeforeEach(func() {
			var err error
			stemcellDirArg, err = os.MkdirTemp("", "stemcell-input-")
			Expect(err).ToNot(HaveOccurred())
		})

		JustBeforeEach(func() {
			actualStemcell, actualErr = NewStemcellFromInput(stemcellDirArg)
		})

		Context("when the stemcell dir contains all necessary files", func() {
			BeforeEach(func() {
				Expect(os.WriteFile(filepath.Join(stemcellDirArg, "version"), []byte("some-version"), 0777)).
					To(Succeed())
				Expect(os.WriteFile(filepath.Join(stemcellDirArg, "url"), []byte("https://s3.amazonaws.com/some-stemcell/stuff-ubuntu-some-os-go_agent.tgz"), 0777)).
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

	Describe("CompareVersion", func() {
		var (
			stemcell1 Stemcell
			stemcell2 Stemcell
		)

		It("returns -1 when stemcell1 is smaller", func() {
			stemcell1 = Stemcell{
				OS:      "whatever",
				Version: "456.1",
			}
			stemcell2 = Stemcell{
				OS:      "whatever",
				Version: "456.2",
			}

			actualResult, actualErr := stemcell1.CompareVersion(stemcell2)

			Expect(actualErr).NotTo(HaveOccurred())
			Expect(actualResult).To(Equal(-1))
		})

		It("returns 0 when the versions are the same", func() {
			stemcell1 = Stemcell{
				OS:      "whatever",
				Version: "1.2",
			}
			stemcell2 = Stemcell{
				OS:      "whatever",
				Version: "1.2",
			}

			actualResult, actualErr := stemcell1.CompareVersion(stemcell2)

			Expect(actualErr).NotTo(HaveOccurred())
			Expect(actualResult).To(Equal(0))
		})

		It("returns 1 when stemcell1 is larger", func() {
			stemcell1 = Stemcell{
				OS:      "whatever",
				Version: "5.1",
			}
			stemcell2 = Stemcell{
				OS:      "whatever",
				Version: "1.2",
			}

			actualResult, actualErr := stemcell1.CompareVersion(stemcell2)

			Expect(actualErr).NotTo(HaveOccurred())
			Expect(actualResult).To(Equal(1))
		})

		It("returns an error when the stemcell1 version is invalid", func() {
			stemcell1 = Stemcell{
				OS:      "whatever",
				Version: "5",
			}
			stemcell2 = Stemcell{
				OS:      "whatever",
				Version: "1.2",
			}

			_, actualErr := stemcell1.CompareVersion(stemcell2)

			Expect(actualErr).To(MatchError("failed to parse stemcell version \"5\": No Major.Minor.Patch elements found"))
		})

		It("returns an error when the stemcell2 version is invalid", func() {
			stemcell1 = Stemcell{
				OS:      "whatever",
				Version: "5.1",
			}
			stemcell2 = Stemcell{
				OS:      "whatever",
				Version: "",
			}

			_, actualErr := stemcell1.CompareVersion(stemcell2)

			Expect(actualErr).To(MatchError("failed to parse stemcell version \"\": No Major.Minor.Patch elements found"))
		})

		It("returns an error when the stemcell OS do not match", func() {
			stemcell1 = Stemcell{
				OS:      "os1",
				Version: "5",
			}
			stemcell2 = Stemcell{
				OS:      "os2",
				Version: "1.2",
			}

			_, actualErr := stemcell1.CompareVersion(stemcell2)

			Expect(actualErr).To(MatchError("stemcell OS mismatch: \"os1\" vs \"os2\""))
		})
	})

	Describe("DetectBumpTypeFrom", func() {
		var (
			targetStemcell Stemcell
			baseStemcell   Stemcell
		)

		Context("is a forward bump for the same OS", func() {
			It("returns \"minor\" when the stemcells have the same major version but the target's minor version is greater", func() {
				targetStemcell = Stemcell{
					OS:      "whatever",
					Version: "456.2",
				}
				baseStemcell = Stemcell{
					OS:      "whatever",
					Version: "456.1",
				}

				actualResult, actualErr := targetStemcell.DetectBumpTypeFrom(baseStemcell)

				Expect(actualErr).NotTo(HaveOccurred())
				Expect(actualResult).To(Equal("minor"))
			})

			It("returns \"major\" when target stemcell's major version is greater", func() {
				targetStemcell = Stemcell{
					OS:      "whatever",
					Version: "460.0",
				}
				baseStemcell = Stemcell{
					OS:      "whatever",
					Version: "456.1",
				}

				actualResult, actualErr := targetStemcell.DetectBumpTypeFrom(baseStemcell)

				Expect(actualErr).NotTo(HaveOccurred())
				Expect(actualResult).To(Equal("major"))
			})
		})

		Context("is NOT a forward bump for the same OS", func() {
			It("returns a non-forward bump error", func() {
				targetStemcell = Stemcell{
					OS:      "whatever",
					Version: "456.1",
				}
				baseStemcell = Stemcell{
					OS:      "whatever",
					Version: "456.2",
				}

				_, actualErr := targetStemcell.DetectBumpTypeFrom(baseStemcell)
				Expect(actualErr).To(MatchError(fmt.Errorf("change from 456.2 to 456.1 is not a forward bump")))
			})
		})

		Context("is a bump to a new OS", func() {
			It("returns \"major\" when target stemcell's OS changes", func() {
				targetStemcell = Stemcell{
					OS:      "new-os",
					Version: "1.1",
				}
				baseStemcell = Stemcell{
					OS:      "whatever",
					Version: "460.0",
				}

				actualResult, actualErr := targetStemcell.DetectBumpTypeFrom(baseStemcell)

				Expect(actualErr).NotTo(HaveOccurred())
				Expect(actualResult).To(Equal("major"))
			})
		})
	})
})
