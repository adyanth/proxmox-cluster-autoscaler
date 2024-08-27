/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package proxmox

import (
	"fmt"
	"io"
	"os"

	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
)

func BuildProxmoxEngine(cloudConfigPath string) cloudprovider.CloudProvider {
	var configFile io.ReadCloser
	if cloudConfigPath == "" {
		panic(fmt.Errorf("config file needed, not provided"))
	}

	var err error
	configFile, err = os.Open(cloudConfigPath)
	if err != nil {
		panic(fmt.Errorf("couldn't open cloud provider configuration %s: %#v", cloudConfigPath, err))
	}
	defer configFile.Close()

	manager, err := newProxmoxManager(configFile)
	if err != nil {
		panic(fmt.Errorf("Failed to create Proxmox config: %v", err))
	}

	return newProxmoxCloudProvider(manager)
}
