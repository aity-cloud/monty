// Code generated by internal/codegen. DO NOT EDIT.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v1.0.0
// source: github.com/aity-cloud/monty/internal/cortex/config/querier/querier.proto

package querier

import (
	_ "github.com/aity-cloud/monty/internal/codegen/cli"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Config struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The maximum number of concurrent queries.
	MaxConcurrent *int32 `protobuf:"varint,1,opt,name=max_concurrent,json=maxConcurrent,proto3,oneof" json:"max_concurrent,omitempty"`
	// The timeout for a query.
	Timeout *durationpb.Duration `protobuf:"bytes,2,opt,name=timeout,proto3" json:"timeout,omitempty"`
	// Use iterators to execute query, as opposed to fully materialising the series in memory.
	Iterators *bool `protobuf:"varint,3,opt,name=iterators,proto3,oneof" json:"iterators,omitempty"`
	// Use batch iterators to execute query, as opposed to fully materialising the series in memory.  Takes precedent over the -querier.iterators flag.
	BatchIterators *bool `protobuf:"varint,4,opt,name=batch_iterators,json=batchIterators,proto3,oneof" json:"batch_iterators,omitempty"`
	// Use streaming RPCs to query ingester.
	IngesterStreaming *bool `protobuf:"varint,5,opt,name=ingester_streaming,json=ingesterStreaming,proto3,oneof" json:"ingester_streaming,omitempty"`
	// Use streaming RPCs for metadata APIs from ingester.
	IngesterMetadataStreaming *bool `protobuf:"varint,6,opt,name=ingester_metadata_streaming,json=ingesterMetadataStreaming,proto3,oneof" json:"ingester_metadata_streaming,omitempty"`
	// Maximum number of samples a single query can load into memory.
	MaxSamples *int32 `protobuf:"varint,7,opt,name=max_samples,json=maxSamples,proto3,oneof" json:"max_samples,omitempty"`
	// Maximum lookback beyond which queries are not sent to ingester. 0 means all queries are sent to ingester.
	QueryIngestersWithin *durationpb.Duration `protobuf:"bytes,8,opt,name=query_ingesters_within,json=queryIngestersWithin,proto3" json:"query_ingesters_within,omitempty"`
	// Deprecated (Querying long-term store for labels will be always enabled in the future.): Query long-term store for series, label values and label names APIs.
	QueryStoreForLabelsEnabled *bool `protobuf:"varint,9,opt,name=query_store_for_labels_enabled,json=queryStoreForLabelsEnabled,proto3,oneof" json:"query_store_for_labels_enabled,omitempty"`
	// Enable returning samples stats per steps in query response.
	PerStepStatsEnabled *bool `protobuf:"varint,10,opt,name=per_step_stats_enabled,json=perStepStatsEnabled,proto3,oneof" json:"per_step_stats_enabled,omitempty"`
	// The time after which a metric should be queried from storage and not just ingesters. 0 means all queries are sent to store. When running the blocks storage, if this option is enabled, the time range of the query sent to the store will be manipulated to ensure the query end is not more recent than 'now - query-store-after'.
	QueryStoreAfter *durationpb.Duration `protobuf:"bytes,11,opt,name=query_store_after,json=queryStoreAfter,proto3" json:"query_store_after,omitempty"`
	// Maximum duration into the future you can query. 0 to disable.
	MaxQueryIntoFuture *durationpb.Duration `protobuf:"bytes,12,opt,name=max_query_into_future,json=maxQueryIntoFuture,proto3" json:"max_query_into_future,omitempty"`
	// The default evaluation interval or step size for subqueries.
	DefaultEvaluationInterval *durationpb.Duration `protobuf:"bytes,13,opt,name=default_evaluation_interval,json=defaultEvaluationInterval,proto3" json:"default_evaluation_interval,omitempty"`
	// Time since the last sample after which a time series is considered stale and ignored by expression evaluations.
	LookbackDelta *durationpb.Duration `protobuf:"bytes,14,opt,name=lookback_delta,json=lookbackDelta,proto3" json:"lookback_delta,omitempty"`
	// When distributor's sharding strategy is shuffle-sharding and this setting is > 0, queriers fetch in-memory series from the minimum set of required ingesters, selecting only ingesters which may have received series since 'now - lookback period'. The lookback period should be greater or equal than the configured 'query store after' and 'query ingesters within'. If this setting is 0, queriers always query all ingesters (ingesters shuffle sharding on read path is disabled).
	ShuffleShardingIngestersLookbackPeriod *durationpb.Duration `protobuf:"bytes,15,opt,name=shuffle_sharding_ingesters_lookback_period,json=shuffleShardingIngestersLookbackPeriod,proto3" json:"shuffle_sharding_ingesters_lookback_period,omitempty"`
	// Experimental. Use Thanos promql engine https://github.com/thanos-io/promql-engine rather than the Prometheus promql engine.
	ThanosEngine *bool `protobuf:"varint,16,opt,name=thanos_engine,json=thanosEngine,proto3,oneof" json:"thanos_engine,omitempty"`
}

func (x *Config) Reset() {
	*x = Config{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_rawDescGZIP(), []int{0}
}

func (x *Config) GetMaxConcurrent() int32 {
	if x != nil && x.MaxConcurrent != nil {
		return *x.MaxConcurrent
	}
	return 0
}

func (x *Config) GetTimeout() *durationpb.Duration {
	if x != nil {
		return x.Timeout
	}
	return nil
}

func (x *Config) GetIterators() bool {
	if x != nil && x.Iterators != nil {
		return *x.Iterators
	}
	return false
}

func (x *Config) GetBatchIterators() bool {
	if x != nil && x.BatchIterators != nil {
		return *x.BatchIterators
	}
	return false
}

func (x *Config) GetIngesterStreaming() bool {
	if x != nil && x.IngesterStreaming != nil {
		return *x.IngesterStreaming
	}
	return false
}

func (x *Config) GetIngesterMetadataStreaming() bool {
	if x != nil && x.IngesterMetadataStreaming != nil {
		return *x.IngesterMetadataStreaming
	}
	return false
}

func (x *Config) GetMaxSamples() int32 {
	if x != nil && x.MaxSamples != nil {
		return *x.MaxSamples
	}
	return 0
}

func (x *Config) GetQueryIngestersWithin() *durationpb.Duration {
	if x != nil {
		return x.QueryIngestersWithin
	}
	return nil
}

func (x *Config) GetQueryStoreForLabelsEnabled() bool {
	if x != nil && x.QueryStoreForLabelsEnabled != nil {
		return *x.QueryStoreForLabelsEnabled
	}
	return false
}

func (x *Config) GetPerStepStatsEnabled() bool {
	if x != nil && x.PerStepStatsEnabled != nil {
		return *x.PerStepStatsEnabled
	}
	return false
}

func (x *Config) GetQueryStoreAfter() *durationpb.Duration {
	if x != nil {
		return x.QueryStoreAfter
	}
	return nil
}

func (x *Config) GetMaxQueryIntoFuture() *durationpb.Duration {
	if x != nil {
		return x.MaxQueryIntoFuture
	}
	return nil
}

func (x *Config) GetDefaultEvaluationInterval() *durationpb.Duration {
	if x != nil {
		return x.DefaultEvaluationInterval
	}
	return nil
}

func (x *Config) GetLookbackDelta() *durationpb.Duration {
	if x != nil {
		return x.LookbackDelta
	}
	return nil
}

func (x *Config) GetShuffleShardingIngestersLookbackPeriod() *durationpb.Duration {
	if x != nil {
		return x.ShuffleShardingIngestersLookbackPeriod
	}
	return nil
}

func (x *Config) GetThanosEngine() bool {
	if x != nil && x.ThanosEngine != nil {
		return *x.ThanosEngine
	}
	return false
}

var File_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto protoreflect.FileDescriptor

var file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_rawDesc = []byte{
	0x0a, 0x48, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69, 0x74,
	0x79, 0x2d, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x6d, 0x6f, 0x6e, 0x74, 0x79, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x72, 0x74, 0x65, 0x78, 0x2f, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x2f, 0x71, 0x75, 0x65, 0x72, 0x69, 0x65, 0x72, 0x2f, 0x71, 0x75, 0x65,
	0x72, 0x69, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x71, 0x75, 0x65, 0x72,
	0x69, 0x65, 0x72, 0x1a, 0x3a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x61, 0x69, 0x74, 0x79, 0x2d, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x6d, 0x6f, 0x6e, 0x74, 0x79,
	0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x67, 0x65,
	0x6e, 0x2f, 0x63, 0x6c, 0x69, 0x2f, 0x63, 0x6c, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x8e, 0x0b, 0x0a, 0x06, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x34, 0x0a, 0x0e, 0x6d, 0x61,
	0x78, 0x5f, 0x63, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x42, 0x08, 0x8a, 0xc0, 0x0c, 0x04, 0x0a, 0x02, 0x32, 0x30, 0x48, 0x00, 0x52, 0x0d,
	0x6d, 0x61, 0x78, 0x43, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x88, 0x01, 0x01,
	0x12, 0x3f, 0x0a, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x0a, 0x8a, 0xc0,
	0x0c, 0x06, 0x0a, 0x04, 0x32, 0x6d, 0x30, 0x73, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75,
	0x74, 0x12, 0x2e, 0x0a, 0x09, 0x69, 0x74, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x08, 0x42, 0x0b, 0x8a, 0xc0, 0x0c, 0x07, 0x0a, 0x05, 0x66, 0x61, 0x6c, 0x73,
	0x65, 0x48, 0x01, 0x52, 0x09, 0x69, 0x74, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x88, 0x01,
	0x01, 0x12, 0x38, 0x0a, 0x0f, 0x62, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x69, 0x74, 0x65, 0x72, 0x61,
	0x74, 0x6f, 0x72, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x42, 0x0a, 0x8a, 0xc0, 0x0c, 0x06,
	0x0a, 0x04, 0x74, 0x72, 0x75, 0x65, 0x48, 0x02, 0x52, 0x0e, 0x62, 0x61, 0x74, 0x63, 0x68, 0x49,
	0x74, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x88, 0x01, 0x01, 0x12, 0x3e, 0x0a, 0x12, 0x69,
	0x6e, 0x67, 0x65, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e,
	0x67, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x42, 0x0a, 0x8a, 0xc0, 0x0c, 0x06, 0x0a, 0x04, 0x74,
	0x72, 0x75, 0x65, 0x48, 0x03, 0x52, 0x11, 0x69, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x65, 0x72, 0x53,
	0x74, 0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e, 0x67, 0x88, 0x01, 0x01, 0x12, 0x50, 0x0a, 0x1b, 0x69,
	0x6e, 0x67, 0x65, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e, 0x67, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08,
	0x42, 0x0b, 0x8a, 0xc0, 0x0c, 0x07, 0x0a, 0x05, 0x66, 0x61, 0x6c, 0x73, 0x65, 0x48, 0x04, 0x52,
	0x19, 0x69, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e, 0x67, 0x88, 0x01, 0x01, 0x12, 0x34, 0x0a,
	0x0b, 0x6d, 0x61, 0x78, 0x5f, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x05, 0x42, 0x0e, 0x8a, 0xc0, 0x0c, 0x0a, 0x0a, 0x08, 0x35, 0x30, 0x30, 0x30, 0x30, 0x30,
	0x30, 0x30, 0x48, 0x05, 0x52, 0x0a, 0x6d, 0x61, 0x78, 0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73,
	0x88, 0x01, 0x01, 0x12, 0x59, 0x0a, 0x16, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f, 0x69, 0x6e, 0x67,
	0x65, 0x73, 0x74, 0x65, 0x72, 0x73, 0x5f, 0x77, 0x69, 0x74, 0x68, 0x69, 0x6e, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x08,
	0x8a, 0xc0, 0x0c, 0x04, 0x0a, 0x02, 0x30, 0x73, 0x52, 0x14, 0x71, 0x75, 0x65, 0x72, 0x79, 0x49,
	0x6e, 0x67, 0x65, 0x73, 0x74, 0x65, 0x72, 0x73, 0x57, 0x69, 0x74, 0x68, 0x69, 0x6e, 0x12, 0x54,
	0x0a, 0x1e, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x5f, 0x66, 0x6f,
	0x72, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x5f, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x08, 0x42, 0x0b, 0x8a, 0xc0, 0x0c, 0x07, 0x0a, 0x05, 0x66, 0x61,
	0x6c, 0x73, 0x65, 0x48, 0x06, 0x52, 0x1a, 0x71, 0x75, 0x65, 0x72, 0x79, 0x53, 0x74, 0x6f, 0x72,
	0x65, 0x46, 0x6f, 0x72, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x61, 0x62, 0x6c, 0x65,
	0x64, 0x88, 0x01, 0x01, 0x12, 0x45, 0x0a, 0x16, 0x70, 0x65, 0x72, 0x5f, 0x73, 0x74, 0x65, 0x70,
	0x5f, 0x73, 0x74, 0x61, 0x74, 0x73, 0x5f, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x08, 0x42, 0x0b, 0x8a, 0xc0, 0x0c, 0x07, 0x0a, 0x05, 0x66, 0x61, 0x6c, 0x73,
	0x65, 0x48, 0x07, 0x52, 0x13, 0x70, 0x65, 0x72, 0x53, 0x74, 0x65, 0x70, 0x53, 0x74, 0x61, 0x74,
	0x73, 0x45, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x88, 0x01, 0x01, 0x12, 0x4f, 0x0a, 0x11, 0x71,
	0x75, 0x65, 0x72, 0x79, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x5f, 0x61, 0x66, 0x74, 0x65, 0x72,
	0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x42, 0x08, 0x8a, 0xc0, 0x0c, 0x04, 0x0a, 0x02, 0x30, 0x73, 0x52, 0x0f, 0x71, 0x75, 0x65,
	0x72, 0x79, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x41, 0x66, 0x74, 0x65, 0x72, 0x12, 0x59, 0x0a, 0x15,
	0x6d, 0x61, 0x78, 0x5f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f, 0x69, 0x6e, 0x74, 0x6f, 0x5f, 0x66,
	0x75, 0x74, 0x75, 0x72, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x0b, 0x8a, 0xc0, 0x0c, 0x07, 0x0a, 0x05, 0x31, 0x30,
	0x6d, 0x30, 0x73, 0x52, 0x12, 0x6d, 0x61, 0x78, 0x51, 0x75, 0x65, 0x72, 0x79, 0x49, 0x6e, 0x74,
	0x6f, 0x46, 0x75, 0x74, 0x75, 0x72, 0x65, 0x12, 0x65, 0x0a, 0x1b, 0x64, 0x65, 0x66, 0x61, 0x75,
	0x6c, 0x74, 0x5f, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x0a, 0x8a, 0xc0, 0x0c, 0x06, 0x0a, 0x04, 0x31,
	0x6d, 0x30, 0x73, 0x52, 0x19, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x45, 0x76, 0x61, 0x6c,
	0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x4c,
	0x0a, 0x0e, 0x6c, 0x6f, 0x6f, 0x6b, 0x62, 0x61, 0x63, 0x6b, 0x5f, 0x64, 0x65, 0x6c, 0x74, 0x61,
	0x18, 0x0e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x42, 0x0a, 0x8a, 0xc0, 0x0c, 0x06, 0x0a, 0x04, 0x35, 0x6d, 0x30, 0x73, 0x52, 0x0d, 0x6c,
	0x6f, 0x6f, 0x6b, 0x62, 0x61, 0x63, 0x6b, 0x44, 0x65, 0x6c, 0x74, 0x61, 0x12, 0x7f, 0x0a, 0x2a,
	0x73, 0x68, 0x75, 0x66, 0x66, 0x6c, 0x65, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67,
	0x5f, 0x69, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x65, 0x72, 0x73, 0x5f, 0x6c, 0x6f, 0x6f, 0x6b, 0x62,
	0x61, 0x63, 0x6b, 0x5f, 0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x08, 0x8a, 0xc0, 0x0c,
	0x04, 0x0a, 0x02, 0x30, 0x73, 0x52, 0x26, 0x73, 0x68, 0x75, 0x66, 0x66, 0x6c, 0x65, 0x53, 0x68,
	0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x49, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x65, 0x72, 0x73, 0x4c,
	0x6f, 0x6f, 0x6b, 0x62, 0x61, 0x63, 0x6b, 0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x12, 0x35, 0x0a,
	0x0d, 0x74, 0x68, 0x61, 0x6e, 0x6f, 0x73, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x18, 0x10,
	0x20, 0x01, 0x28, 0x08, 0x42, 0x0b, 0x8a, 0xc0, 0x0c, 0x07, 0x0a, 0x05, 0x66, 0x61, 0x6c, 0x73,
	0x65, 0x48, 0x08, 0x52, 0x0c, 0x74, 0x68, 0x61, 0x6e, 0x6f, 0x73, 0x45, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x88, 0x01, 0x01, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x6d, 0x61, 0x78, 0x5f, 0x63, 0x6f, 0x6e,
	0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x69, 0x74, 0x65, 0x72,
	0x61, 0x74, 0x6f, 0x72, 0x73, 0x42, 0x12, 0x0a, 0x10, 0x5f, 0x62, 0x61, 0x74, 0x63, 0x68, 0x5f,
	0x69, 0x74, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x73, 0x42, 0x15, 0x0a, 0x13, 0x5f, 0x69, 0x6e,
	0x67, 0x65, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e, 0x67,
	0x42, 0x1e, 0x0a, 0x1c, 0x5f, 0x69, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x6d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e, 0x67,
	0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x6d, 0x61, 0x78, 0x5f, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73,
	0x42, 0x21, 0x0a, 0x1f, 0x5f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x5f, 0x66, 0x6f, 0x72, 0x5f, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x5f, 0x65, 0x6e, 0x61, 0x62,
	0x6c, 0x65, 0x64, 0x42, 0x19, 0x0a, 0x17, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x73, 0x74, 0x65, 0x70,
	0x5f, 0x73, 0x74, 0x61, 0x74, 0x73, 0x5f, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x42, 0x10,
	0x0a, 0x0e, 0x5f, 0x74, 0x68, 0x61, 0x6e, 0x6f, 0x73, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65,
	0x42, 0x44, 0x82, 0xc0, 0x0c, 0x04, 0x08, 0x01, 0x10, 0x01, 0x5a, 0x3a, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x69, 0x74, 0x79, 0x2d, 0x63, 0x6c, 0x6f, 0x75,
	0x64, 0x2f, 0x6d, 0x6f, 0x6e, 0x74, 0x79, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x2f, 0x63, 0x6f, 0x72, 0x74, 0x65, 0x78, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x71,
	0x75, 0x65, 0x72, 0x69, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_rawDescOnce sync.Once
	file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_rawDescData = file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_rawDesc
)

func file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_rawDescGZIP() []byte {
	file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_rawDescOnce.Do(func() {
		file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_rawDescData)
	})
	return file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_rawDescData
}

var file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_goTypes = []interface{}{
	(*Config)(nil),              // 0: querier.Config
	(*durationpb.Duration)(nil), // 1: google.protobuf.Duration
}
var file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_depIdxs = []int32{
	1, // 0: querier.Config.timeout:type_name -> google.protobuf.Duration
	1, // 1: querier.Config.query_ingesters_within:type_name -> google.protobuf.Duration
	1, // 2: querier.Config.query_store_after:type_name -> google.protobuf.Duration
	1, // 3: querier.Config.max_query_into_future:type_name -> google.protobuf.Duration
	1, // 4: querier.Config.default_evaluation_interval:type_name -> google.protobuf.Duration
	1, // 5: querier.Config.lookback_delta:type_name -> google.protobuf.Duration
	1, // 6: querier.Config.shuffle_sharding_ingesters_lookback_period:type_name -> google.protobuf.Duration
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_init() }
func file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_init() {
	if File_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Config); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_goTypes,
		DependencyIndexes: file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_depIdxs,
		MessageInfos:      file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_msgTypes,
	}.Build()
	File_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto = out.File
	file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_rawDesc = nil
	file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_goTypes = nil
	file_github_com_aity_cloud_monty_internal_cortex_config_querier_querier_proto_depIdxs = nil
}
