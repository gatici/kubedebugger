// Package main License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/gatici/KubeDebugger/kubedebugger/ephemeral"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericiooptions"

	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

func main() {
	flags := pflag.NewFlagSet("kubedb", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := ephemeral.NewEphemeralContainerCmd(genericiooptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
