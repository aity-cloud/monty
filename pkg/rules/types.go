package rules

import (
	"maps"

	"github.com/prometheus/prometheus/model/rulefmt"
)

// Alias so we can implement the Finder[T] interface
// on Prometheus' rulefmt.RuleGroup
type RuleGroup rulefmt.RuleGroup
type RuleGroupList []RuleGroup

func (r RuleGroup) Clone() RuleGroup {
	cloned := RuleGroup{
		Name:     r.Name,
		Interval: r.Interval,
		Rules:    make([]rulefmt.Rule, len(r.Rules)),
	}

	for i, r := range r.Rules {
		cloned.Rules[i] = rulefmt.Rule{
			Record:      r.Record,
			Alert:       r.Alert,
			Expr:        r.Expr,
			For:         r.For,
			Labels:      maps.Clone(r.Labels),
			Annotations: maps.Clone(r.Annotations),
		}
	}
	return cloned
}
