package shared

import (
	"fmt"

	"github.com/samber/lo"
)

type AlertingClusterNotification lo.Tuple2[bool, *AlertingClusterOptions]

type AlertingClusterOptions struct {
	Namespace             string
	WorkerNodesService    string
	WorkerNodePort        int
	WorkerStatefulSet     string
	ControllerNodeService string
	ControllerNodePort    int
	ControllerClusterPort int
	ControllerStatefulSet string
	MontyPort             int
	ConfigMap             string
	ManagementHookHandler string
}

func (a *AlertingClusterOptions) GetInternalControllerMontyEndpoint() string {
	return fmt.Sprintf("%s:%d", a.ControllerNodeService, AlertingDefaultHookPort)
}

func (a *AlertingClusterOptions) GetInternalWorkerMontyEndpoint() string {
	return fmt.Sprintf("%s:%d", a.WorkerNodesService, AlertingDefaultHookPort)
}

func (a *AlertingClusterOptions) GetControllerEndpoint() string {
	return fmt.Sprintf("%s:%d", a.ControllerNodeService, a.ControllerNodePort)
}

func (a *AlertingClusterOptions) GetWorkerEndpoint() string {
	return fmt.Sprintf("%s:%d", a.WorkerNodesService, a.WorkerNodePort)
}
