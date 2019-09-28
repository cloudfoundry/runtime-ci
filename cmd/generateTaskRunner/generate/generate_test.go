package generate

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("generate", func() {
	Describe("NewPackageFromFile", func() {
		var (
			taskYAMLPathArg     string
			outputPackageDirArg string

			actualPackage Package
			actualErr     error
		)

		JustBeforeEach(func() {
			actualPackage, actualErr = NewPackage(outputPackageDirArg, taskYAMLPathArg)
		})

		Context("with a standard task.yml", func() {
			BeforeEach(func() {
				outputPackageDirArg = "example-runner"

				f, err := ioutil.TempFile("", "task.yml-")
				Expect(err).NotTo(HaveOccurred())
				taskYAMLPathArg = f.Name()

				fmt.Fprint(f, `---
platform: linux

inputs:
- name: example-input-a
- name: example-input-b

outputs:
- name: example-output-a
- name: example-output-b
`)
				Expect(f.Close()).To(Succeed())
			})

			AfterEach(func() {
				Expect(os.Remove(taskYAMLPathArg)).To(Succeed())
			})

			It("expect it to be initialized from the task yaml", func() {
				Expect(actualErr).NotTo(HaveOccurred())

				Expect(actualPackage).To(Equal(Package{
					name:     "example-runner",
					dir: "example-runner",
					inputs:  []string{"example-input-a", "example-input-b"},
					outputs: []string{"example-output-a", "example-output-b"},
				}))
			})
		})
	})

	Describe("Write", func() {
		var (
			pkgSubject Package

			nameValue     string
			dirValue     string
			inputsValue  []string
			outputsValue []string

			actualErr error
		)

		JustBeforeEach(func() {
			pkgSubject = Package{
				name:     nameValue,
				dir:     dirValue,
				inputs:  inputsValue,
				outputs: outputsValue,
			}

			actualErr = pkgSubject.Write()
		})

		Context("when all provided a full task YAML package", func() {
			BeforeEach(func() {
				nameValue = "generated-runner"

				tmp, err := ioutil.TempDir("", "generated-runner-")
				Expect(err).NotTo(HaveOccurred())
				Expect(os.RemoveAll(tmp)).To(Succeed())
				dirValue = tmp

				inputsValue = []string{"example-input-a", "example-input-b"}
				outputsValue = []string{"example-output-a", "example-output-b"}
			})

			It("creates the package directory and content", func() {
				Expect(actualErr).NotTo(HaveOccurred())
				Expect(dirValue).To(BeADirectory())
			})
		})
	})
})
