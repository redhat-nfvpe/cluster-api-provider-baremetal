package operators

import (
	"context"
	"time"

	"github.com/golang/glog"
	osconfigv1 "github.com/openshift/api/config/v1"
	e2e "github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/framework"
	cov1helpers "github.com/openshift/library-go/pkg/config/clusteroperator/v1helpers"
	kappsapi "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func isDeploymentAvailable(client runtimeclient.Client, name string) bool {
	key := types.NamespacedName{
		Namespace: e2e.TestContext.MachineApiNamespace,
		Name:      name,
	}
	d := &kappsapi.Deployment{}

	if err := wait.PollImmediate(1*time.Second, e2e.WaitShort, func() (bool, error) {
		if err := client.Get(context.TODO(), key, d); err != nil {
			glog.Errorf("Error querying api for Deployment object: %v, retrying...", err)
			return false, nil
		}
		if d.Status.AvailableReplicas < 1 {
			glog.Errorf("Deployment %q is not available. Status: (replicas: %d, updated: %d, ready: %d, available: %d, unavailable: %d)", d.Name, d.Status.Replicas, d.Status.UpdatedReplicas, d.Status.ReadyReplicas, d.Status.AvailableReplicas, d.Status.UnavailableReplicas)
			return false, nil
		}
		return true, nil
	}); err != nil {
		glog.Errorf("Error checking isDeploymentAvailable: %v", err)
		return false
	}
	glog.Infof("Deployment %q is available. Status: (replicas: %d, updated: %d, ready: %d, available: %d, unavailable: %d)", d.Name, d.Status.Replicas, d.Status.UpdatedReplicas, d.Status.ReadyReplicas, d.Status.AvailableReplicas, d.Status.UnavailableReplicas)
	return true
}

func isStatusAvailable(client runtimeclient.Client, name string) bool {
	key := types.NamespacedName{
		Namespace: e2e.TestContext.MachineApiNamespace,
		Name:      name,
	}
	clusterOperator := &osconfigv1.ClusterOperator{}

	if err := wait.PollImmediate(1*time.Second, e2e.WaitShort, func() (bool, error) {
		if err := client.Get(context.TODO(), key, clusterOperator); err != nil {
			glog.Errorf("error querying api for OperatorStatus object: %v, retrying...", err)
			return false, nil
		}
		if cov1helpers.IsStatusConditionFalse(clusterOperator.Status.Conditions, osconfigv1.OperatorAvailable) {
			glog.Errorf("Condition: %q is false", osconfigv1.OperatorAvailable)
			return false, nil
		}
		if cov1helpers.IsStatusConditionTrue(clusterOperator.Status.Conditions, osconfigv1.OperatorProgressing) {
			glog.Errorf("Condition: %q is true", osconfigv1.OperatorProgressing)
			return false, nil
		}
		if cov1helpers.IsStatusConditionTrue(clusterOperator.Status.Conditions, osconfigv1.OperatorFailing) {
			glog.Errorf("Condition: %q is true", osconfigv1.OperatorFailing)
			return false, nil
		}
		return true, nil
	}); err != nil {
		glog.Errorf("Error checking isStatusAvailable: %v", err)
		return false
	}
	return true

}
