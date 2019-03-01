package main

import (
	clusterapis "github.com/openshift/cluster-api/pkg/apis"
	machineactuator "github.com/redhat-nfvpe/cluster-api-provider-baremetal/pkg/actuators/machine"
	"github.com/redhat-nfvpe/cluster-api-provider-baremetal/pkg/apis"
	"github.com/redhat-nfvpe/cluster-api-provider-baremetal/pkg/apis/baremetalproviderconfig/v1alpha1"
	"github.com/redhat-nfvpe/cluster-api-provider-baremetal/pkg/controller"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"

	"flag"

	"github.com/golang/glog"
	"k8s.io/klog"
)

func main() {
	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlags)
	flag.Parse()
	flag.VisitAll(func(f1 *flag.Flag) {
		f2 := klogFlags.Lookup(f1.Name)
		if f2 != nil {
			value := f1.Value.String()
			f2.Value.Set(value)
		}
	})

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		glog.Fatal(err)
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{})
	if err != nil {
		glog.Fatal(err)
	}

	glog.Infof("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		glog.Fatal(err)
	}

	if err := clusterapis.AddToScheme(mgr.GetScheme()); err != nil {
		glog.Fatal(err)
	}

	initActuator(mgr)
	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		glog.Fatal(err)
	}

	glog.Infof("Starting the Cmd.")

	// Start the Cmd
	glog.Fatal(mgr.Start(signals.SetupSignalHandler()))
}

func initActuator(m manager.Manager) (*machineactuator.Actuator, error) {
	config := m.GetConfig()
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Could not create kubernetes client to talk to the apiserver: %v", err)
	}

	codec, err := v1alpha1.NewCodec()
	if err != nil {
		glog.Fatal(err)
	}

	params := machineactuator.ActuatorParams{
		Client:              m.GetClient(),
		Codec:               codec,
		KubeClient:          kubeClient,
		EventRecorder:       m.GetRecorder("baremetal-controller"),
		ServerListenAddress: "localhost:8081",
	}

	machineactuator.MachineActuator, err = machineactuator.NewActuator(params)
	if err != nil {
		glog.Fatalf("Could not create Baremetal machine actuator: %v", err)
	}

	return machineactuator.MachineActuator, nil
}
