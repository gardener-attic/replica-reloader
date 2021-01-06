// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"errors"
	goflag "flag"

	"github.com/gardener/replica-reloader/internal/kubeconfig"
	"github.com/gardener/replica-reloader/internal/version"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

// NewReplicaReloader creates the root command of the replica-reloader.
func NewReplicaReloader(ctx context.Context) *cobra.Command {
	opts := &options{}
	kubeConfig := &kubeconfig.KubeConfig{}
	cmd := &cobra.Command{
		Use:     "replica-reloader [flags] -- COMMAND [args...] INJECTED-REPLICA-COUNT",
		Short:   "Executes a command depending on the deployment's replica count",
		Version: version.Version,
		Example: `
    $(terminal-1) kubectl create deployment --image=nginx nginx
    deployment.apps/nginx created

and running

    $(terminal-2) replica-reloader --namespace=default --deployment-name=nginx -- sleep

would start a "sleep 1" process.
If the watched deployment is scaled, then the controller stops the previous process and
starts a new one:

    $(terminal-1) kubectl scale deployment my-dep --replicas=10
    deployment.apps/my-dep scaled

    $(terminal-1) ps | grep sleep
	61191 ttys003    0:00.00 sleep 10`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			rc, err := kubeConfig.Complete()
			if err != nil {
				return err
			}

			opts.client, err = kubernetes.NewForConfig(rc)
			if err != nil {
				return err
			}

			argsLenAtDash := cmd.ArgsLenAtDash()

			if argsLenAtDash > -1 && len(args[argsLenAtDash:]) > 0 {
				opts.command = args[argsLenAtDash:]
			} else {
				return errors.New("COMMAND should be passed e.g. replica-reloader --namespace=default -- sleep")
			}

			if opts.namespace == "" {
				return errors.New("namespace is required")
			}

			if opts.deploymentName == "" {
				return errors.New("deployment-name is required")
			}

			if opts.jitter < 0 {
				return errors.New("jitter cannot be negative")
			}

			if opts.jitterFactor < 0 {
				return errors.New("jitter-factor cannot be negative")
			}

			if opts.jitterFactor > 0 && opts.jitter == 0 {
				return errors.New("jitter must be set when specifying jitter-factor")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.start(ctx)
		},
	}

	flags := cmd.Flags()

	klog.InitFlags(goflag.CommandLine)
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	kubeConfig.AddFlags(flags)
	flags.StringVarP(&opts.namespace, "namespace", "n", "", "namespace of the deployment")
	flags.StringVar(&opts.deploymentName, "deployment-name", "kube-apiserver", "name of the deployment")
	flags.DurationVar(&opts.jitter, "jitter", 0, "duration between receiving a scale event and process restart")
	flags.Float64Var(&opts.jitterFactor, "jitter-factor", 0.0, "adds random factor to jitter. Requires jitter to be set")

	return cmd
}
