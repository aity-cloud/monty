package cortex

import (
	"fmt"
	"strings"
	"time"

	"github.com/aity-cloud/monty/pkg/alerting/message"
	"github.com/aity-cloud/monty/pkg/alerting/metrics"
	"github.com/aity-cloud/monty/pkg/alerting/shared"
	alertingv1 "github.com/aity-cloud/monty/pkg/apis/alerting/v1"
	prommodel "github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/model/rulefmt"
	"github.com/samber/lo"
)

/*
Contains the struct/function adapters required for monty alerting to
communicate with cortex.
*/

const (
	MetadataCortexNamespace = "monty.io/cortex-rule-namespace"
	MetadataCortexGroup     = "monty.io/cortex-rule-group"
	MetadataCortexRuleName  = "monty.io/cortex-rule-name"
)

const alertingSuffix = "-monty-alerting"

// this enforces whatever default the remote prometheus instance has
const defaultAlertingInterval = prommodel.Duration(0 * time.Minute)

func RuleIdFromUuid(id string) string {
	return id + alertingSuffix
}

func TimeDurationToPromStr(t time.Duration) string {
	return prommodel.Duration(t).String()
}

func ConstructRecordingRuleName(prefix, typeName string) string {
	return fmt.Sprintf("monty:%s:%s", prefix, typeName)
}

func ConstructIdLabelsForRecordingRule(alertId string) map[string]string {
	return map[string]string{
		message.NotificationPropertyMontyUuid: alertId,
	}
}

func ConstructFiltersFromMap(in map[string]string) string {
	var filters []string
	for k, v := range in {
		filters = append(filters, fmt.Sprintf("%s=\"%s\"", k, v))
	}
	return strings.Join(filters, ",")
}

func NewPrometheusAlertingRule(
	alertId,
	_ string,
	montyLabels,
	montyAnnotations map[string]string,
	info alertingv1.IndexableMetric,
	interval *time.Duration,
	rule metrics.AlertRuleBuilder,
) (ruleGroup *rulefmt.RuleGroup, metadata map[string]string, err error) {
	idLabels := ConstructIdLabelsForRecordingRule(alertId)
	alertingRule, err := rule.Build(alertId)
	if err != nil {
		return nil, nil, err
	}
	recordingRuleFmt := &rulefmt.Rule{
		Record:      ConstructRecordingRuleName(info.GoldenSignal(), info.AlertType()),
		Expr:        alertingRule.Expr,
		Labels:      idLabels,
		Annotations: map[string]string{},
	}
	// have the alerting rule instead point to the recording rule(s)
	alertingRule.Expr = fmt.Sprintf("%s{%s}", recordingRuleFmt.Record, ConstructFiltersFromMap(idLabels))
	alertingRule.Labels = lo.Assign(alertingRule.Labels, montyLabels)
	alertingRule.Annotations = lo.Assign(alertingRule.Annotations, montyAnnotations)

	var promInterval prommodel.Duration
	if interval == nil {
		promInterval = defaultAlertingInterval
	} else {
		promInterval = prommodel.Duration(*interval)
	}

	rg := &rulefmt.RuleGroup{
		Name:     RuleIdFromUuid(alertId),
		Interval: promInterval,
		Rules:    []rulefmt.Rule{*alertingRule, *recordingRuleFmt},
	}

	return rg, map[string]string{
		MetadataCortexNamespace: shared.MontyAlertingCortexNamespace,
		MetadataCortexGroup:     rg.Name,
		MetadataCortexRuleName:  alertingRule.Alert,
	}, nil
}
