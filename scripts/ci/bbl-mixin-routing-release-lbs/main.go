package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type CloudConfig struct {
	AZs          interface{}              `yaml:"azs,omitempty"`
	VMTypes      interface{}              `yaml:"vm_types,omitempty"`
	DiskTypes    interface{}              `yaml:"disk_types,omitempty"`
	Compilation  interface{}              `yaml:"compilation,omitempty"`
	Networks     interface{}              `yaml:"networks,omitempty"`
	VMExtensions []map[string]interface{} `yaml:"vm_extensions"`
}

func main() {
	var (
		cloudConfigPath          string
		tcpRouterELBID           string
		tcpRouterSecurityGroupID string
		internalSecurityGroupID  string
	)

	flag.StringVar(&cloudConfigPath, "cloud-config", "", "path to the cloud config file to mixin")
	flag.StringVar(&tcpRouterELBID, "tcp-router-elb-id", "", "the tcp router id to mixin to the cloud config")
	flag.StringVar(&tcpRouterSecurityGroupID, "tcp-router-security-group-id", "", "the tcp router security group id to mixin to the cloud config")
	flag.StringVar(&internalSecurityGroupID, "internal-security-group-id", "", "the internal security group id to mixin to the cloud config")
	flag.Parse()

	inputCloudConfigContents, err := ioutil.ReadFile(cloudConfigPath)
	if err != nil {
		fmt.Fprint(os.Stdout, err)
		os.Exit(1)
	}

	var cloudConfig CloudConfig
	if err := yaml.Unmarshal(inputCloudConfigContents, &cloudConfig); err != nil {
		fmt.Printf("error parsing cloud-config: %q", err)
		os.Exit(1)
	}

	newVMExtensions := []map[string]interface{}{}
	for _, vmExtension := range cloudConfig.VMExtensions {
		if vmExtension["name"] != "tcp-router-lb" {
			newVMExtensions = append(newVMExtensions, vmExtension)
		}
	}

	tcpRouterVMExtension := map[string]interface{}{
		"name": "tcp-router-lb",
		"cloud_properties": map[string]interface{}{
			"elbs":            []string{tcpRouterELBID},
			"security_groups": []string{tcpRouterSecurityGroupID, internalSecurityGroupID},
		},
	}

	cloudConfig.VMExtensions = append(newVMExtensions, tcpRouterVMExtension)

	mixedInCloudConfig, err := yaml.Marshal(cloudConfig)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(os.Stdout, string(mixedInCloudConfig))
}
