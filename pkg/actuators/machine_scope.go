package actuators

import (
	clusterv1 "github.com/openshift/cluster-api/pkg/apis/cluster/v1alpha1"
	client "github.com/openshift/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
	"github.com/pkg/errors"
	"github.com/redhat-nfvpe/cluster-api-provider-baremetal/pkg/apis/baremetalproviderconfig/v1alpha1"
	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
)

// MachineScope defines a scope defined around a machine and its cluster.
type MachineScope struct {
	*Scope

	Machine       *clusterv1.Machine
	MachineClient client.MachineInterface
	MachineConfig *v1alpha1.BaremetalMachineProviderSpec
	MachineStatus *v1alpha1.BaremetalMachineProviderStatus
}

// MachineScopeParams defines the input parameters used to create a new MachineScope.
type MachineScopeParams struct {
	Cluster *clusterv1.Cluster
	Machine *clusterv1.Machine
	Client  client.ClusterV1alpha1Interface
}

// NewMachineScope creates a new MachineScope from the supplied parameters.
// This is meant to be called for each machine actuator operation.
func NewMachineScope(params MachineScopeParams) (*MachineScope, error) {
	scope, err := NewScope(ScopeParams{Client: params.Client, Cluster: params.Cluster})
	if err != nil {
		return nil, err
	}

	machineConfig, err := MachineConfigFromProviderSpec(params.Client, params.Machine.Spec.ProviderSpec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get machine config")
	}

	machineStatus, err := v1alpha1.MachineStatusFromProviderStatus(params.Machine.Status.ProviderStatus)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get machine provider status")
	}

	var machineClient client.MachineInterface
	if params.Client != nil {
		machineClient = params.Client.Machines(params.Machine.Namespace)
	}

	return &MachineScope{
		Scope:         scope,
		Machine:       params.Machine,
		MachineClient: machineClient,
		MachineConfig: machineConfig,
		MachineStatus: machineStatus,
	}, nil
}

// MachineConfigFromProviderSpec tries to decode the JSON-encoded spec, falling back on getting a MachineClass if the value is absent.
func MachineConfigFromProviderSpec(clusterClient client.MachineClassesGetter, providerConfig clusterv1.ProviderSpec) (*v1alpha1.BaremetalMachineProviderSpec, error) {
	var config v1alpha1.BaremetalMachineProviderSpec
	if providerConfig.Value != nil {
		klog.V(4).Info("Decoding ProviderConfig from Value")
		return unmarshalProviderSpec(providerConfig.Value)
	}

	if providerConfig.ValueFrom != nil && providerConfig.ValueFrom.MachineClass != nil {
		ref := providerConfig.ValueFrom.MachineClass
		klog.V(4).Info("Decoding ProviderConfig from MachineClass")
		klog.V(6).Infof("ref: %v", ref)
		if ref.Provider != "" && ref.Provider != "aws" {
			return nil, errors.Errorf("Unsupported provider: %q", ref.Provider)
		}

		if len(ref.Namespace) > 0 && len(ref.Name) > 0 {
			klog.V(4).Infof("Getting MachineClass: %s/%s", ref.Namespace, ref.Name)
			mc, err := clusterClient.MachineClasses(ref.Namespace).Get(ref.Name, metav1.GetOptions{})
			klog.V(6).Infof("Retrieved MachineClass: %+v", mc)
			if err != nil {
				return nil, err
			}
			providerConfig.Value = &mc.ProviderSpec
			return unmarshalProviderSpec(&mc.ProviderSpec)
		}
	}

	return &config, nil
}

func unmarshalProviderSpec(spec *runtime.RawExtension) (*v1alpha1.BaremetalMachineProviderSpec, error) {
	var config v1alpha1.BaremetalMachineProviderSpec
	if spec != nil {
		if err := yaml.Unmarshal(spec.Raw, &config); err != nil {
			return nil, err
		}
	}
	klog.V(6).Infof("Found ProviderSpec: %+v", config)
	return &config, nil
}
