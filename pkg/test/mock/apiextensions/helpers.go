package mock_apiextensions

import apiextensions "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions"

type MockManagementAPIExtensionServerImpl struct {
	apiextensions.UnsafeManagementAPIExtensionServer
	*MockManagementAPIExtensionServer
}
