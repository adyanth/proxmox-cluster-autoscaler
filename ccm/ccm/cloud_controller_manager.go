package ccm

import (
	"github.com/adyanth/proxmox-cluster-autoscaler/proxmox"
	cloudprovider "k8s.io/cloud-provider"
)

// Make sure struct implements the interface
var _ cloudprovider.Interface = &ProxmoxCloudProvider{}

type ProxmoxCloudProvider struct {
	Name string
}

func (p *ProxmoxCloudProvider) Initialize(clientBulder cloudprovider.ControllerClientBuilder, stop <-chan struct{}) {
}

func (p *ProxmoxCloudProvider) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return nil, false
}

func (p *ProxmoxCloudProvider) Instances() (cloudprovider.Instances, bool) {
	return nil, false
}

func (p *ProxmoxCloudProvider) InstancesV2() (cloudprovider.InstancesV2, bool) {
	return nil, false
}

func (p *ProxmoxCloudProvider) Zones() (cloudprovider.Zones, bool) {
	return nil, false
}

func (p *ProxmoxCloudProvider) Clusters() (cloudprovider.Clusters, bool) {
	return nil, false
}

func (p *ProxmoxCloudProvider) Routes() (cloudprovider.Routes, bool) {
	return nil, false
}

func (p *ProxmoxCloudProvider) ProviderName() string {
	return proxmox.ProviderId
}

func (p *ProxmoxCloudProvider) HasClusterID() bool {
	return false
}
