// Package ephemeral License-Identifier: Apache-2.0

package ephemeral

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericiooptions"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

type ContainerOptions struct {
	genericiooptions.IOStreams
	ConfigFlags         *genericclioptions.ConfigFlags
	ContainerFilePath   string
	TargetPodName       string
	TargetContainerName string
	args                []string
}

func NewEphemeralContainerOptions(streams genericiooptions.IOStreams) *ContainerOptions {
	return &ContainerOptions{
		ConfigFlags: genericclioptions.NewConfigFlags(true),

		IOStreams: streams,
	}
}

// Complete sets all information required for updating the current context
func (o *ContainerOptions) Complete(cmd *cobra.Command, args []string) error {
	o.args = args

	if len(o.args) > 0 {
		o.TargetPodName = args[0]
	}

	return nil
}

// Validate ensures that all required arguments and flag values are provided
func (o *ContainerOptions) Validate() error {
	if len(o.args) != 1 {
		return fmt.Errorf("missing a target pod")
	}

	if o.ContainerFilePath == "" {
		return fmt.Errorf("must pass a container file path")
	}

	return nil
}
