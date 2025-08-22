//go:build !noevents && !cli

package commands

import (
	"github.com/aity-cloud/monty/pkg/events"
	"github.com/spf13/cobra"
)

var (
	shipperEndpoint string
)

func BuildEventsCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "events",
		Short: "Run the Kubernetes event collector",
		Long:  "The event collector ships Kubernetes events to an monty-shipper endpoint.",
		RunE:  doEvents,
	}

	command.Flags().StringVar(&shipperEndpoint, "endpoint", "http://monty-shipper:2021/log/ingest", "endpoint to post events to")
	return command
}

func doEvents(cmd *cobra.Command, _ []string) error {
	collector := events.NewEventCollector(cmd.Context(), shipperEndpoint)
	return collector.Run(cmd.Context().Done())
}

func init() {
	AddCommandsToGroup(MontyComponents, BuildEventsCmd())
}
