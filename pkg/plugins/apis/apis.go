// Package apis can be imported to ensure all plugin APIs are added to client schemes.
package apis

import (
	_ "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions"
	_ "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions/http"
	_ "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions/management"
	_ "github.com/aity-cloud/monty/pkg/plugins/apis/apiextensions/stream"
	_ "github.com/aity-cloud/monty/pkg/plugins/apis/capability"
	_ "github.com/aity-cloud/monty/pkg/plugins/apis/health"
	_ "github.com/aity-cloud/monty/pkg/plugins/apis/metrics"
	_ "github.com/aity-cloud/monty/pkg/plugins/apis/proxy"
	_ "github.com/aity-cloud/monty/pkg/plugins/apis/system"
)
