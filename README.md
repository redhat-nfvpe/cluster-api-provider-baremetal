# cluster-api-provider-baremetal

When using the k8s code-generator against this, make sure you have installed k8s.io/apimachinery!

`go get -u k8s.io/code-generator`
`go get -u k8s.io/apimachinery`

`k8s.io/code-generator/generate-groups.sh all github.com/redhat-nfvpe/cluster-api-provider-baremetal/pkg/client github.com/redhat-nfvpe/cluster-api-provider-baremetal/pkg/apis baremetalproviderconfig:v1alpha1`

https://github.com/kubernetes/code-generator/issues/21

Test it:

`go run cmd/baremetal-actuator/main.go create -c examples/cluster.yaml -m examples/machine.yaml`
