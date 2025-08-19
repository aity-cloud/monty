package drivers

import (
	"context"

	"github.com/aity-cloud/monty/pkg/plugins/driverutil"
	"github.com/aity-cloud/monty/plugins/alerting/pkg/apis/node"
	"github.com/aity-cloud/monty/plugins/alerting/pkg/apis/rules"
)

type ConfigPropagator interface {
	ConfigureNode(nodeId string, conf *node.AlertingCapabilityConfig) error
}

type NodeDriver interface {
	ConfigPropagator
	DiscoverRules(ctx context.Context) (*rules.RuleManifest, error)
}

var NodeDrivers = driverutil.NewCache[NodeDriver]()
