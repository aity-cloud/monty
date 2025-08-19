package gateway

import (
	"context"

	corev1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
	managementext "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions/management"
	"github.com/aity-cloud/monty/pkg/util"
)

// ManagementServices implements managementext.ManagementAPIExtension.
func (p *Plugin) ManagementServices(s managementext.ServiceController) []util.ServicePackInterface {
	p.serviceCtrl.C() <- s
	return p.managementServices
}

// Authorized checks whether a given set of roles is allowed to access a given request
func (p *Plugin) CheckAuthz(_ context.Context, _ *corev1.ReferenceList, _, _ string) bool {
	return true
}
