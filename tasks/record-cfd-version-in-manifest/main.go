package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

func main() {
	manifest_version, _ := boshExtract("cf-deployment.yml", "/manifest_version")
	fmt.Print(string(manifest_version))

	opsFile, _ := ioutil.TempFile("/tmp", "manifest-version-ops.yml")
	opsFile.WriteString("---\n- type: replace\n  path: /manifest_version\n  value: bar")
	fmt.Println("Your file is written to: " + opsFile.Name())

	newManifestBytes, _ := boshIncept("cf-deployment.yml", opsFile.Name())
	ioutil.WriteFile("cf-deployment-hotness.yml", newManifestBytes, 0666)

}

func boshExtract(filePath, opPath string) ([]byte, error) {
	cmd := exec.Command("bosh", "int", filePath, fmt.Sprintf("--path=%s", opPath))

	b, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return b, nil
}

func boshIncept(manifestPath, opsFilePath string) ([]byte, error) {
	cmd := exec.Command("bosh", "int", manifestPath, "-o", opsFilePath)

	b, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	return b, nil
}
