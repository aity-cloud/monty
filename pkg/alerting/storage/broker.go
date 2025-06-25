package storage

import (
	"github.com/aity-cloud/monty/pkg/alerting/drivers/routing"
	"github.com/aity-cloud/monty/pkg/alerting/interfaces"
	"github.com/aity-cloud/monty/pkg/alerting/storage/jetstream"
	"github.com/aity-cloud/monty/pkg/alerting/storage/mem"
	"github.com/aity-cloud/monty/pkg/alerting/storage/opts"
	"github.com/aity-cloud/monty/pkg/alerting/storage/spec"
	alertingv1 "github.com/aity-cloud/monty/pkg/apis/alerting/v1"
)

var _ spec.RouterStorage = (*jetstream.JetstreamRouterStore[routing.MontyRouting])(nil)
var _ spec.AlertingStorage[interfaces.AlertingSecret] = (*jetstream.JetStreamAlertingStorage[interfaces.AlertingSecret])(nil)
var _ spec.AlertingStateCache[*alertingv1.CachedState] = (*jetstream.JetStreamAlertingStateCache)(nil)
var _ spec.AlertingIncidentTracker[*alertingv1.IncidentIntervals] = (*jetstream.JetStreamAlertingIncidentTracker)(nil)
var _ spec.AlertingStorage[interfaces.AlertingSecret] = (*jetstream.JetStreamAlertingStorage[interfaces.AlertingSecret])(nil)
var _ spec.RouterStorage = (*mem.InMemoryRouterStore)(nil)

type CompositeAlertingBroker struct {
	opts.ClientSetOptions
	*CompositeAlertingClientSet
}

func NewCompositeAlertingBroker(options opts.ClientSetOptions) *CompositeAlertingBroker {
	return &CompositeAlertingBroker{
		ClientSetOptions: options,
		CompositeAlertingClientSet: &CompositeAlertingClientSet{
			hashes: make(map[string]string),
			Logger: options.Logger,
		},
	}
}

var _ spec.AlertingClientSet = (*CompositeAlertingBroker)(nil)
var _ spec.AlertingStoreBroker = (*CompositeAlertingBroker)(nil)

func (c *CompositeAlertingBroker) Use(store any) {
	if cs, ok := store.(spec.ConditionStorage); ok {
		c.conds = cs
	}
	if es, ok := store.(spec.EndpointStorage); ok {
		c.endps = es
	}
	if rs, ok := store.(spec.RouterStorage); ok {
		c.routers = rs
	}
	if ss, ok := store.(spec.StateStorage); ok {
		c.states = ss
	}
	if is, ok := store.(spec.IncidentStorage); ok {
		c.incidents = is
	}
}

func (c *CompositeAlertingBroker) NewClientSet() spec.AlertingClientSet {
	return c.CompositeAlertingClientSet
}
