package cortex

import "github.com/aity-cloud/monty/pkg/util"

var (
	OrgIDCodec = util.NewDelimiterCodec("X-Scope-OrgID", "|")
)
