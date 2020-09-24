package cli

import (
	"fmt"
	"github.com/myonlyzzy/kubectl-resource-view/pkg/logger"
	"github.com/myonlyzzy/kubectl-resource-view/pkg/plugin"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
	"strings"
)

var (
	KubernetesConfigFlags *genericclioptions.ConfigFlags
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubectl-resource-view  NODENAME ) [flags]",
		Short: "Display Node resources",
		Long: `Display Node resources
Examples:
  # Show all nodes all resources
  kubectl resource-view
  
  # Show all resources of node marked as master
  kubectl resource-view -l node-role.kubernetes.io/master

  # Show cpu of node 
  kubectl resource-view -r cpu
  
  # show node host192.0.0.1 resources
  kubectl resource-view -n host192.0.0.1
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			log := logger.NewLogger()
			log.Info("")

			if err := plugin.RunPlugin(KubernetesConfigFlags, cmd); err != nil {
				return errors.Cause(err)
			}

			log.Info("")

			return nil
		},
	}

	cobra.OnInitialize(initConfig)
	var nodeName string
	KubernetesConfigFlags = genericclioptions.NewConfigFlags(false)
	ResourceBuilderFlags := genericclioptions.NewResourceBuilderFlags().WithAllNamespaces(false).
		WithFieldSelector("").
		WithLabelSelector("").
		WithLatest()

	KubernetesConfigFlags.AddFlags(cmd.Flags())
	ResourceBuilderFlags.AddFlags(cmd.Flags())
	cmd.Flags().StringVar(&nodeName, "node", "", "node name")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return cmd
}

func InitAndExecute() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.AutomaticEnv()
}
