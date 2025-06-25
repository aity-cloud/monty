//go:build !minimal

package apis

import (
	montyaiv1beta1 "github.com/aity-cloud/monty/apis/ai/v1beta1"
	montyloggingv1beta1 "github.com/aity-cloud/monty/apis/logging/v1beta1"
	montymonitoringv1beta1 "github.com/aity-cloud/monty/apis/monitoring/v1beta1"
)

func init() {
	addSchemeBuilders(
		montyaiv1beta1.AddToScheme,
		montyloggingv1beta1.AddToScheme,
		montymonitoringv1beta1.AddToScheme,
	)
}
