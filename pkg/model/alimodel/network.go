/*
Copyright 2018 The Kubernetes Authors.

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

package alimodel

import (
	"k8s.io/kops/upup/pkg/fi"
	"k8s.io/kops/upup/pkg/fi/cloudup/alitasks"
)

// NetWorkModelBuilder configures VPC network objects
type NetWorkModelBuilder struct {
	*ALIModelContext
	Lifecycle *fi.Lifecycle
}

var _ fi.ModelBuilder = &NetWorkModelBuilder{}

func (b *NetWorkModelBuilder) Build(c *fi.ModelBuilderContext) error {
	sharedVPC := b.Cluster.SharedVPC()

	// VPC that holds everything for the cluster
	vpc := &alitasks.VPC{}
	{
		vpcName := b.GetNameForVPC()
		vpc.Name = s(vpcName)
		vpc.Lifecycle = b.Lifecycle
		vpc.Shared = fi.Bool(sharedVPC)

		if b.Cluster.Spec.NetworkID != "" {
			vpc.ID = s(b.Cluster.Spec.NetworkID)
		}

		if b.Cluster.Spec.NetworkCIDR != "" {
			vpc.CIDR = s(b.Cluster.Spec.NetworkCIDR)
		}
		c.AddTask(vpc)
	}

	for i := range b.Cluster.Spec.Subnets {
		subnetSpec := &b.Cluster.Spec.Subnets[i]

		vswitch := &alitasks.VSwitch{
			Name:      s(b.GetNameForVSwitch(subnetSpec.Name)),
			Lifecycle: b.Lifecycle,
			VPC:       b.LinkToVPC(),
			ZoneId:    s(subnetSpec.Zone),
			CidrBlock: s(subnetSpec.CIDR),
			Shared:    fi.Bool(false),
		}

		if subnetSpec.ProviderID != "" {
			vswitch.VSwitchId = s(subnetSpec.ProviderID)
			vswitch.Shared = fi.Bool(true)
		}

		c.AddTask(vswitch)

	}

	return nil
}
