package metrics

import (
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/model/relabel"
)

const (
	LabelImpersonateAs = "zz_monty_impersonate_as"
)

var matchNonEmptyRegex = relabel.MustNewRegexp(".+")

// Drops metrics containing any monty internal labels.
func MontyInternalLabelFilter() *relabel.Config {
	return &relabel.Config{
		SourceLabels: model.LabelNames{
			LabelImpersonateAs,
		},
		Regex:  matchNonEmptyRegex,
		Action: relabel.Drop,
	}
}
