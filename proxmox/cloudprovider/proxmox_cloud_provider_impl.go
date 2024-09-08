package cloudprovider

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/adyanth/proxmox-cluster-autoscaler/proxmox"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/autoscaler/cluster-autoscaler/utils/errors"
)

func trace() {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	fmt.Printf("==> %s:%d %s\n", file, line, f.Name())
}

// cloudProvider.CloudProvider interface implementation on ProxmoxCloudProvider

// Name returns name of the cloud provider.
func (p *ProxmoxCloudProvider) Name() string {
	trace()
	return "ProxmoxCloudProvider"
}

// NodeGroups returns all node groups configured for this cloud provider.
func (p *ProxmoxCloudProvider) NodeGroups() []cloudprovider.NodeGroup {
	trace()
	ngs := make([]cloudprovider.NodeGroup, 0, len(p.manager.NodeGroupManagers))
	for _, ng := range p.manager.NodeGroupManagers {
		ngs = append(ngs, ng)
	}
	return ngs
}

// NodeGroupForNode returns the node group for the given node, nil if the node
// should not be processed by cluster autoscaler, or non-nil error if such
// occurred. Must be implemented.
func (p *ProxmoxCloudProvider) NodeGroupForNode(node *apiv1.Node) (ng cloudprovider.NodeGroup, err error) {
	trace()
	log.Println("Getting nodegroup for node " + node.Name + " with spec.providerId: " + node.Spec.ProviderID)
	ng, _, err = p.manager.GetDetailsFromNode(node)
	return ng, err
}

// HasInstance returns whether the node has corresponding instance in cloud provider,
// true if the node has an instance, false if it no longer exists
func (p *ProxmoxCloudProvider) HasInstance(node *apiv1.Node) (bool, error) {
	trace()
	ngm, offset, err := p.manager.GetDetailsFromNode(node)
	if err != nil || ngm == nil {
		return false, err
	}
	return ngm.HasNodeInstance(offset), nil
}

// Pricing returns pricing model for this cloud provider or error if not available.
// Implementation optional.
func (p *ProxmoxCloudProvider) Pricing() (cloudprovider.PricingModel, errors.AutoscalerError) {
	trace()
	return nil, cloudprovider.ErrNotImplemented
}

// GetAvailableMachineTypes get all machine types that can be requested from the cloud provider.
// Implementation optional.
func (p *ProxmoxCloudProvider) GetAvailableMachineTypes() ([]string, error) {
	trace()
	return []string{}, cloudprovider.ErrNotImplemented
}

// NewNodeGroup builds a theoretical node group based on the node definition provided. The node group is not automatically
// created on the cloud provider side. The node group is not returned by NodeGroups() until it is created.
// Implementation optional.
func (p *ProxmoxCloudProvider) NewNodeGroup(machineType string, labels map[string]string, systemLabels map[string]string,
	taints []apiv1.Taint, extraResources map[string]resource.Quantity) (cloudprovider.NodeGroup, error) {
	trace()
	return nil, cloudprovider.ErrNotImplemented
}

// GetResourceLimiter returns struct containing limits (max, min) for resources (cores, memory etc.).
func (p *ProxmoxCloudProvider) GetResourceLimiter() (*cloudprovider.ResourceLimiter, error) {
	trace()
	return nil, nil
}

// GPULabel returns the label added to nodes with GPU resource.
func (p *ProxmoxCloudProvider) GPULabel() string {
	trace()
	return proxmox.GPULabel
}

// GetAvailableGPUTypes return all available GPU types cloud provider supports.
func (p *ProxmoxCloudProvider) GetAvailableGPUTypes() map[string]struct{} {
	trace()
	return nil
}

// GetNodeGpuConfig returns the label, type and resource name for the GPU added to node. If node doesn't have
// any GPUs, it returns nil.
func (p *ProxmoxCloudProvider) GetNodeGpuConfig(node *apiv1.Node) *cloudprovider.GpuConfig {
	trace()
	return nil
}

// Cleanup cleans up open resources before the cloud provider is destroyed, i.e. go routines etc.
func (p *ProxmoxCloudProvider) Cleanup() error {
	trace()
	return nil
}

// Refresh is called before every main loop and can be used to dynamically update cloud provider state.
// In particular the list of node groups returned by NodeGroups can change as a result of CloudProvider.Refresh().
func (p *ProxmoxCloudProvider) Refresh() error {
	trace()
	ctx := context.Background()
	for _, ngm := range p.manager.NodeGroupManagers {
		ngm.FillCurrentSize(ctx)
	}
	return nil
}
