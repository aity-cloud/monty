//go:build !minimal

package apis

import montygrafanav1beta1 "github.com/aity-cloud/monty/apis/grafana/v1beta1"

func init() {
	addSchemeBuilders(montygrafanav1beta1.AddToScheme)
}
