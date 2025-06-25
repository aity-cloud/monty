//go:build !minimal

package apis

import opnigrafanav1beta1 "github.com/aity-cloud/monty/apis/grafana/v1beta1"

func init() {
	addSchemeBuilders(opnigrafanav1beta1.AddToScheme)
}
