package drivers

import (
	"context"
	"fmt"

	"github.com/aity-cloud/monty/pkg/config/v1beta1"
	"github.com/aity-cloud/monty/pkg/plugins/driverutil"
	"github.com/aity-cloud/monty/pkg/rules"
	"github.com/aity-cloud/monty/pkg/util/notifier"
	"github.com/aity-cloud/monty/plugins/metrics/apis/node"
	"github.com/aity-cloud/monty/plugins/metrics/apis/remoteread"
)

type MetricsNodeConfigurator interface {
	ConfigureNode(nodeId string, conf *node.MetricsCapabilityConfig) error
}

type MetricsNodeDriver interface {
	MetricsNodeConfigurator
	DiscoverPrometheuses(context.Context, string) ([]*remoteread.DiscoveryEntry, error)
	ConfigureRuleGroupFinder(config *v1beta1.RulesSpec) notifier.Finder[rules.RuleGroup]
}

var NodeDrivers = driverutil.NewDriverCache[MetricsNodeDriver]()

type ConfigurationError struct {
	NodeId string
	Cfg    *node.MetricsCapabilityConfig
	Err    error
}

func (e *ConfigurationError) Error() string {
	return fmt.Errorf("error configuring node %s: %w", e.NodeId, e.Err).Error()
}
