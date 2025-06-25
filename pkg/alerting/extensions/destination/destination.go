package destination

import (
	"context"

	"github.com/aity-cloud/monty/pkg/alerting/drivers/config"
	alertingv1 "github.com/aity-cloud/monty/pkg/apis/alerting/v1"
)

var (
	defaultSeverity = alertingv1.MontySeverity_Info.String()
)

const (
	missingTitle = "missing alert title"
	missingBody  = "missing alert body"
)

type Destination interface {
	Push(context.Context, config.WebhookMessage) error
	Name() string
}
