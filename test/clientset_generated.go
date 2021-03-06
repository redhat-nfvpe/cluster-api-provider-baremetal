package test

import (
	"github.com/openshift/cluster-api/pkg/client/clientset_generated/clientset"
	clusterv1alpha1 "github.com/openshift/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
	fakeclusterv1alpha1 "github.com/openshift/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1/fake"
	machinev1beta1 "github.com/openshift/cluster-api/pkg/client/clientset_generated/clientset/typed/machine/v1beta1"
	fakemachinev1beta1 "github.com/openshift/cluster-api/pkg/client/clientset_generated/clientset/typed/machine/v1beta1/fake"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/testing"
)

// Clientset implements clientset.Interface. Meant to be embedded into a
// struct to get a default implementation. This makes faking out just the method
// you want to test easier.
type Clientset struct {
	*testing.Fake
	discovery *fakediscovery.FakeDiscovery
}

func NewSimpleClientset(objects ...runtime.Object) *Clientset {
	o := testing.NewObjectTracker(scheme, codecs.UniversalDecoder())
	for _, obj := range objects {
		if err := o.Add(obj); err != nil {
			panic(err)
		}
	}

	// https://github.com/kubernetes/kubernetes/pull/60584
	fakePtr := &testing.Fake{}
	fakePtr.AddReactor("*", "*", testing.ObjectReaction(o))
	fakePtr.AddWatchReactor("*", testing.DefaultWatchReactor(watch.NewFake(), nil))

	return &Clientset{fakePtr, &fakediscovery.FakeDiscovery{Fake: fakePtr}}
}

func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	return c.discovery
}

var _ clientset.Interface = &Clientset{}

// MachineV1beta1 retrieves the MachineV1beta1Client
func (c *Clientset) MachineV1beta1() machinev1beta1.MachineV1beta1Interface {
	return &fakemachinev1beta1.FakeMachineV1beta1{Fake: c.Fake}
}

// Machine retrieves the MachineV1beta1Client
func (c *Clientset) Machine() machinev1beta1.MachineV1beta1Interface {
	return c.MachineV1beta1()
}

// ClusterV1alpha1 retrieves the ClusterV1alpha1Client
func (c *Clientset) ClusterV1alpha1() clusterv1alpha1.ClusterV1alpha1Interface {
	return &fakeclusterv1alpha1.FakeClusterV1alpha1{Fake: c.Fake}
}

// Cluster retrieves the ClusterV1alpha1Client
func (c *Clientset) Cluster() clusterv1alpha1.ClusterV1alpha1Interface {
	return &fakeclusterv1alpha1.FakeClusterV1alpha1{Fake: c.Fake}
}
