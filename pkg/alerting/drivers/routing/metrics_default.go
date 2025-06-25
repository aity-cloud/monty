package routing

import (
	"github.com/aity-cloud/monty/pkg/alerting/drivers/config"
	"github.com/prometheus/alertmanager/pkg/labels"
	"github.com/prometheus/common/model"
)

var MontyMetricsSubRoutingTreeId config.Matchers = []*labels.Matcher{
	MontyMetricsSubRoutingTreeMatcher,
}

// this subtree connects to cortex
func NewMontyMetricsSubtree() *config.Route {
	metricsRoute := &config.Route{
		// expand all labels in case our default grouping overshadows some of the user's configs
		GroupBy: []model.LabelName{
			"...",
		},
		Matchers: MontyMetricsSubRoutingTreeId,
		Routes:   []*config.Route{},
		Continue: true, // want to expand the subtree
	}
	setDefaultRateLimitingFromProto(metricsRoute)
	return metricsRoute
}
