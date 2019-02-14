package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	goflag "flag"

	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/openshift/cluster-api-actuator-pkg/pkg/e2e/framework"

	"github.com/redhat-nfvpe/cluster-api-provider-baremetal/cmd/baremetal-actuator/utils"
)

func usage() {
	fmt.Printf("Usage: %s\n\n", os.Args[0])
}

var rootCmd = &cobra.Command{
	Use:   "baremetal-actuator-test",
	Short: "Test for Cluster API Baremetal actuator",
}

func cmdRun(binaryPath string, args ...string) ([]byte, error) {
	cmd := exec.Command(binaryPath, args...)
	return cmd.CombinedOutput()
}

func BuildPKSecret(secretName, namespace, pkLoc string) (*apiv1.Secret, error) {
	pkBytes, err := ioutil.ReadFile(pkLoc)
	if err != nil {
		return nil, fmt.Errorf("unable to read %v: %v", pkLoc, err)
	}

	return &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"privatekey": pkBytes,
		},
	}, nil
}

func createSecretAndWait(f *framework.Framework, secret *apiv1.Secret) error {
	_, err := f.KubeClient.CoreV1().Secrets(secret.Namespace).Create(secret)
	if err != nil {
		return err
	}

	err = wait.Poll(framework.PollInterval, framework.PoolTimeout, func() (bool, error) {
		_, err := f.KubeClient.CoreV1().Secrets(secret.Namespace).Get(secret.Name, metav1.GetOptions{})
		return err == nil, nil
	})
	return err
}

func createCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create machine instance for specified cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkFlags(cmd); err != nil {
				return err
			}
			cluster, machine, userData, err := utils.ReadClusterResources(
				cmd.Flag("cluster").Value.String(),
				cmd.Flag("machine").Value.String(),
				cmd.Flag("userdata").Value.String(),
			)
			if err != nil {
				return err
			}

			actuator := utils.CreateActuator(machine, userData)
			err = actuator.Create(context.TODO(), cluster, machine)
			if err != nil {
				return err
			}
			fmt.Printf("Machine creation was successful!\n")
			return nil
		},
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("machine", "m", "", "Machine manifest")
	rootCmd.PersistentFlags().StringP("cluster", "c", "", "Cluster manifest")
	rootCmd.PersistentFlags().StringP("userdata", "u", "", "User data manifest")

	cUser, err := user.Current()
	if err != nil {
		rootCmd.PersistentFlags().StringP("environment-id", "p", "", "Directory with bootstrapping manifests")
	} else {
		rootCmd.PersistentFlags().StringP("environment-id", "p", cUser.Username, "Machine prefix, by default set to the current user")
	}

	rootCmd.AddCommand(createCommand())

	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	// the following line exists to make glog happy, for more information, see: https://github.com/kubernetes/kubernetes/issues/17162
	flag.CommandLine.Parse([]string{})
}

func checkFlags(cmd *cobra.Command) error {
	if cmd.Flag("cluster").Value.String() == "" {
		return fmt.Errorf("--%v/-%v flag is required", cmd.Flag("cluster").Name, cmd.Flag("cluster").Shorthand)
	}
	if cmd.Flag("machine").Value.String() == "" {
		return fmt.Errorf("--%v/-%v flag is required", cmd.Flag("machine").Name, cmd.Flag("machine").Shorthand)
	}
	return nil
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occurred: %v\n", err)
		os.Exit(1)
	}
}
