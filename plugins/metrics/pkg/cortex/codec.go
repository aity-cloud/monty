package cortex

import "github.com/aity-cloud/monty/pkg/util"

var (
	orgIDCodec = util.NewDelimiterCodec("X-Scope-OrgID", "|")
)
