package actuators

import (
	clusterv1 "github.com/openshift/cluster-api/pkg/apis/cluster/v1alpha1"
	client "github.com/openshift/cluster-api/pkg/client/clientset_generated/clientset/typed/cluster/v1alpha1"
	"github.com/pkg/errors"
	v1alpha1 "github.com/redhat-nfvpe/cluster-api-provider-baremetal/pkg/apis/baremetalproviderconfig/v1alpha1"
	"k8s.io/klog"
)

// ScopeParams defines the input parameters used to create a new Scope.
type ScopeParams struct {
	Cluster *clusterv1.Cluster
	Client  client.ClusterV1alpha1Interface
}

// Scope defines the basic context for an actuator to operate upon.
type Scope struct {
	Cluster       *clusterv1.Cluster
	ClusterClient client.ClusterInterface
	ClusterConfig *v1alpha1.BaremetalClusterProviderSpec
	ClusterStatus *v1alpha1.BaremetalClusterProviderStatus
}

// NewScope creates a new Scope from the supplied parameters.
// This is meant to be called for each different actuator iteration.
func NewScope(params ScopeParams) (*Scope, error) {
	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil cluster")
	}

	clusterConfig, err := v1alpha1.ClusterConfigFromProviderSpec(params.Cluster.Spec.ProviderSpec)
	if err != nil {
		return nil, errors.Errorf("failed to load cluster provider config: %v", err)
	}

	clusterStatus, err := v1alpha1.ClusterStatusFromProviderStatus(params.Cluster.Status.ProviderStatus)
	if err != nil {
		return nil, errors.Errorf("failed to load cluster provider status: %v", err)
	}

	var clusterClient client.ClusterInterface
	if params.Client != nil {
		clusterClient = params.Client.Clusters(params.Cluster.Namespace)
	}

	return &Scope{
		Cluster:       params.Cluster,
		ClusterClient: clusterClient,
		ClusterConfig: clusterConfig,
		ClusterStatus: clusterStatus,
	}, nil
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *Scope) Close() {
	if s.ClusterClient == nil {
		return
	}

	latestCluster, err := s.storeClusterConfig(s.Cluster)
	if err != nil {
		klog.Errorf("[scope] failed to store provider config for cluster %q in namespace %q: %v", s.Cluster.Name, s.Cluster.Namespace, err)
		return
	}

	_, err = s.storeClusterStatus(latestCluster)
	if err != nil {
		klog.Errorf("[scope] failed to store provider status for cluster %q in namespace %q: %v", s.Cluster.Name, s.Cluster.Namespace, err)
	}
}

func (s *Scope) storeClusterConfig(cluster *clusterv1.Cluster) (*clusterv1.Cluster, error) {
	ext, err := v1alpha1.EncodeClusterSpec(s.ClusterConfig)
	if err != nil {
		return nil, err
	}

	cluster.Spec.ProviderSpec.Value = ext
	return s.ClusterClient.Update(cluster)
}

func (s *Scope) storeClusterStatus(cluster *clusterv1.Cluster) (*clusterv1.Cluster, error) {
	ext, err := v1alpha1.EncodeClusterStatus(s.ClusterStatus)
	if err != nil {
		return nil, err
	}

	cluster.Status.ProviderStatus = ext
	return s.ClusterClient.UpdateStatus(cluster)
}
