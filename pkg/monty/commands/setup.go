//go:build !minimal

package commands

import (
	"strings"
	"sync"

	managementv1 "github.com/aity-cloud/monty/pkg/apis/management/v1"
	"github.com/aity-cloud/monty/pkg/clients"
	"github.com/aity-cloud/monty/pkg/config"
	"github.com/aity-cloud/monty/pkg/config/v1beta1"
	"github.com/aity-cloud/monty/pkg/logger"
	"github.com/aity-cloud/monty/pkg/monty/cliutil"
	"github.com/aity-cloud/monty/pkg/tracing"
	"github.com/spf13/cobra"
)

var mgmtClient managementv1.ManagementClient
var managementListenAddress string
var initManagementOnce sync.Once
var lg = logger.New()

func ConfigureManagementCommand(cmd *cobra.Command) {
	if cmd.PersistentPreRunE == nil {
		cmd.PersistentPreRunE = managementPreRunE
	} else {
		oldPreRunE := cmd.PersistentPreRunE
		cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			if err := managementPreRunE(cmd, args); err != nil {
				return err
			}
			return oldPreRunE(cmd, args)
		}
	}
}

func managementPreRunE(cmd *cobra.Command, _ []string) (err error) {
	initManagementOnce.Do(func() {
		tracing.Configure("cli")
		address := cmd.Flag("address").Value.String()
		if address == "" {
			path, err := config.FindConfig()
			if err == nil {
				objects := cliutil.LoadConfigObjectsOrDie(path, lg)
				objects.Visit(func(obj *v1beta1.GatewayConfig) {
					address = strings.TrimPrefix(obj.Spec.Management.GRPCListenAddress, "tcp://")
				})
			}
		}
		if address == "" {
			address = managementv1.DefaultManagementSocket()
		}
		managementListenAddress = address
		mgmtClient, err = clients.NewManagementClient(cmd.Context(), clients.WithAddress(address))
	})
	return err
}
