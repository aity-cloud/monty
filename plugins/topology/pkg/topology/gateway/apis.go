package gateway

/*
Orchestrator API implementation
*/

import (
	"context"
	"encoding/json"

	corev1 "github.com/aity-cloud/monty/pkg/apis/core/v1"
	"github.com/aity-cloud/monty/plugins/topology/apis/orchestrator"
	"github.com/aity-cloud/monty/plugins/topology/apis/representation"
	kgraph "github.com/steveteuber/kubectl-graph/pkg/graph"

	"github.com/aity-cloud/monty/pkg/topology/graph"
	"github.com/aity-cloud/monty/pkg/topology/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (p *Plugin) GetGraph(_ context.Context, _ *corev1.Reference) (*representation.TopologyGraph, error) {
	// TODO(topology) : implement me
	return nil, status.Error(codes.Unimplemented, "method not implemented")
}

func (p *Plugin) RenderGraph(ctx context.Context, clusterRef *corev1.Reference) (*representation.DOTRepresentation, error) {
	if !p.topologyRemoteWrite.Initialized() {
		return nil, status.Error(codes.Unavailable, "topology remote write not initialized")
	}
	ctxCa, cancel := context.WithCancel(ctx)
	defer cancel()
	mgr, err := p.nc.GetContext(ctxCa)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	obj, err := store.GetTopologyObjectStore(mgr)
	if err != nil {
		return nil, err
	}
	objInfo, err := obj.GetInfo(store.NewClusterKey(clusterRef))
	if err != nil {
		return nil, err
	}
	reprKey := objInfo.Headers.Get(store.ReprHeaderKey)
	if reprKey == "" {
		return nil, status.Error(
			codes.Internal,
			"no representation header found for the cluster",
		)
	}
	graphObj, err := obj.Get(store.NewClusterKey(clusterRef))
	if err != nil {
		return nil, err
	}

	switch reprKey {
	case representation.GraphRepr_KubectlGraph.String():
		var g *kgraph.Graph
		if err := json.NewDecoder(graphObj).Decode(&g); err != nil {
			return nil, err
		}
		// !! Cannot store marshalled digraph, since it will not capture information unless
		// we implement the entire gonum.Graph and its sub interfaces ourselves
		diGraph := graph.NewScientificKubeGraph()
		diGraph.FromKubectlGraph(g)
		bytes, err := diGraph.RenderDOT()
		if err != nil {
			return nil, err
		}
		return &representation.DOTRepresentation{
			RawDotFormat: string(bytes),
		}, nil
	default:
		return nil, status.Error(codes.Internal, "invalid representation key")
	}
}

func (p *Plugin) GetClusterStatus(_ context.Context, _ *emptypb.Empty) (*orchestrator.InstallStatus, error) {
	return &orchestrator.InstallStatus{
		State:   orchestrator.InstallState_Installed,
		Version: "0.1",
	}, nil
}
