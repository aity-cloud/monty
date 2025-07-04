//go:build !minimal

package commands

import (
	"github.com/aity-cloud/monty/plugins/metrics/apis/cortexops"
	"github.com/spf13/cobra"
)

func BuildMetricsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Interact with metrics plugin APIs",
	}

	cmd.AddCommand(BuildCortexAdminCmd())
	cmd.AddCommand(cortexops.BuildCortexOpsCmd())
	cmd.AddCommand(BuildMetricsConfigCmd())

	ConfigureManagementCommand(cmd)
	ConfigureCortexAdminCommand(cmd)
	return cmd
}

func init() {
	AddCommandsToGroup(PluginAPIs, BuildMetricsCmd())
}
