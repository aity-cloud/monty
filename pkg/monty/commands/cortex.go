//go:build !minimal && !cli

package commands

import (
	"flag"

	cortex_internal "github.com/aity-cloud/monty/internal/cortex"
	"github.com/spf13/cobra"
)

func BuildCortexCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "cortex",
		Short:              "Run the embedded Cortex server",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			flag.CommandLine = flag.NewFlagSet("cortex", flag.ExitOnError)
			cortex_internal.Main(append([]string{"cortex"}, args...))
		},
	}

	return cmd
}

func init() {
	AddCommandsToGroup(MontyComponents, BuildCortexCmd())
}
