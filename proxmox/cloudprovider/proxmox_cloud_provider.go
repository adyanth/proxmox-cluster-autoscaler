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

package cloudprovider

import (
	"fmt"
	"io"
	"os"

	"github.com/adyanth/proxmox-cluster-autoscaler/proxmox"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
)

type ProxmoxCloudProvider struct {
	manager *proxmox.ProxmoxManager
}

func newProxmoxCloudProvider(manager *proxmox.ProxmoxManager) *ProxmoxCloudProvider {
	return &ProxmoxCloudProvider{
		manager: manager,
	}
}

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

	manager, err := proxmox.NewProxmoxManager(configFile)
	if err != nil {
		panic(fmt.Errorf("failed to create Proxmox config: %v", err))
	}

	return newProxmoxCloudProvider(manager)
}
