// Package ephemeral License-Identifier: Apache-2.0

package ephemeral

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/ptr"
)

func NewEphemeralContainerCmd(streams genericiooptions.IOStreams) *cobra.Command {
	o := NewEphemeralContainerOptions(streams)

	cmd := &cobra.Command{
		Use:          "ephemeral [pod name] [flags]",
		Short:        "Create an ephemeral container in a target pod from YAML specification.",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(c, args); err != nil {
				return err
			}

			if err := o.Validate(); err != nil {
				return err
			}

			if err := Run(o, c.Context()); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&o.ContainerFilePath, "file", "f", "", "file containing a YAML container spec to create an ephemeral container from")
	cmd.Flags().StringVarP(&o.TargetContainerName, "container", "c", "", "container within the target pod that the ephemeral container will attach itself to")
	o.ConfigFlags.AddFlags(cmd.Flags())

	return cmd
}

func Run(opts *ContainerOptions, ctx context.Context) error {
	container, err := getContainerFromFile(opts.ContainerFilePath)
	if err != nil {
		return fmt.Errorf("failed to read ephemeral container spec from file: %w", err)
	}
	kubeConfigPath := os.ExpandEnv("$HOME/.kube/config")
	if _, err2 := os.Stat(kubeConfigPath); err2 != nil {
		kubeConfigPath = ""
	}
	conf, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return fmt.Errorf("Couldn't get kubeconfig. %w", err)
	}
	client, err := corev1.NewForConfig(conf)
	if err != nil {
		return fmt.Errorf("failed to construct client. %w", err)
	}

	namespace := opts.ConfigFlags.Namespace
	if namespace == nil || *namespace == "" {
		namespace = ptr.To(metav1.NamespaceDefault)
	}

	pod, err := client.Pods(*namespace).Get(ctx, opts.TargetPodName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("target pod not found: %w", err)
	}

	if opts.TargetContainerName != "" {
		containerValid := false
		for _, container := range pod.Spec.Containers {
			if container.Name == opts.TargetContainerName {
				containerValid = true
				break
			}
		}

		if !containerValid {
			return fmt.Errorf("container %s not found in target pod %s", opts.TargetContainerName, opts.TargetPodName)
		}

		container.TargetContainerName = opts.TargetContainerName
	}

	pod.Spec.EphemeralContainers = append(pod.Spec.EphemeralContainers, *container)

	_, err = client.Pods(*namespace).UpdateEphemeralContainers(ctx, pod.Name, pod, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create ephemeral container: %w", err)
	}

	fmt.Printf("EphemeralContainer/%s created\n", container.Name)

	return nil
}
