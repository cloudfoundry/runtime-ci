// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v5.26.1
// source: envoy/extensions/load_balancing_policies/common/v3/common.proto

package commonv3

import (
	reflect "reflect"
	sync "sync"

	_ "github.com/cncf/xds/go/udpa/annotations"
	v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	v31 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type LocalityLbConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to LocalityConfigSpecifier:
	//
	//	*LocalityLbConfig_ZoneAwareLbConfig_
	//	*LocalityLbConfig_LocalityWeightedLbConfig_
	LocalityConfigSpecifier isLocalityLbConfig_LocalityConfigSpecifier `protobuf_oneof:"locality_config_specifier"`
}

func (x *LocalityLbConfig) Reset() {
	*x = LocalityLbConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_envoy_extensions_load_balancing_policies_common_v3_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LocalityLbConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocalityLbConfig) ProtoMessage() {}

func (x *LocalityLbConfig) ProtoReflect() protoreflect.Message {
	mi := &file_envoy_extensions_load_balancing_policies_common_v3_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocalityLbConfig.ProtoReflect.Descriptor instead.
func (*LocalityLbConfig) Descriptor() ([]byte, []int) {
	return file_envoy_extensions_load_balancing_policies_common_v3_common_proto_rawDescGZIP(), []int{0}
}

func (m *LocalityLbConfig) GetLocalityConfigSpecifier() isLocalityLbConfig_LocalityConfigSpecifier {
	if m != nil {
		return m.LocalityConfigSpecifier
	}
	return nil
}

func (x *LocalityLbConfig) GetZoneAwareLbConfig() *LocalityLbConfig_ZoneAwareLbConfig {
	if x, ok := x.GetLocalityConfigSpecifier().(*LocalityLbConfig_ZoneAwareLbConfig_); ok {
		return x.ZoneAwareLbConfig
	}
	return nil
}

func (x *LocalityLbConfig) GetLocalityWeightedLbConfig() *LocalityLbConfig_LocalityWeightedLbConfig {
	if x, ok := x.GetLocalityConfigSpecifier().(*LocalityLbConfig_LocalityWeightedLbConfig_); ok {
		return x.LocalityWeightedLbConfig
	}
	return nil
}

type isLocalityLbConfig_LocalityConfigSpecifier interface {
	isLocalityLbConfig_LocalityConfigSpecifier()
}

type LocalityLbConfig_ZoneAwareLbConfig_ struct {
	// Configuration for local zone aware load balancing.
	ZoneAwareLbConfig *LocalityLbConfig_ZoneAwareLbConfig `protobuf:"bytes,1,opt,name=zone_aware_lb_config,json=zoneAwareLbConfig,proto3,oneof"`
}

type LocalityLbConfig_LocalityWeightedLbConfig_ struct {
	// Enable locality weighted load balancing.
	LocalityWeightedLbConfig *LocalityLbConfig_LocalityWeightedLbConfig `protobuf:"bytes,2,opt,name=locality_weighted_lb_config,json=localityWeightedLbConfig,proto3,oneof"`
}

func (*LocalityLbConfig_ZoneAwareLbConfig_) isLocalityLbConfig_LocalityConfigSpecifier() {}

func (*LocalityLbConfig_LocalityWeightedLbConfig_) isLocalityLbConfig_LocalityConfigSpecifier() {}

// Configuration for :ref:`slow start mode <arch_overview_load_balancing_slow_start>`.
type SlowStartConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Represents the size of slow start window.
	// If set, the newly created host remains in slow start mode starting from its creation time
	// for the duration of slow start window.
	SlowStartWindow *durationpb.Duration `protobuf:"bytes,1,opt,name=slow_start_window,json=slowStartWindow,proto3" json:"slow_start_window,omitempty"`
	// This parameter controls the speed of traffic increase over the slow start window. Defaults to 1.0,
	// so that endpoint would get linearly increasing amount of traffic.
	// When increasing the value for this parameter, the speed of traffic ramp-up increases non-linearly.
	// The value of aggression parameter should be greater than 0.0.
	// By tuning the parameter, is possible to achieve polynomial or exponential shape of ramp-up curve.
	//
	// During slow start window, effective weight of an endpoint would be scaled with time factor and aggression:
	// “new_weight = weight * max(min_weight_percent, time_factor ^ (1 / aggression))“,
	// where “time_factor=(time_since_start_seconds / slow_start_time_seconds)“.
	//
	// As time progresses, more and more traffic would be sent to endpoint, which is in slow start window.
	// Once host exits slow start, time_factor and aggression no longer affect its weight.
	Aggression *v3.RuntimeDouble `protobuf:"bytes,2,opt,name=aggression,proto3" json:"aggression,omitempty"`
	// Configures the minimum percentage of origin weight that avoids too small new weight,
	// which may cause endpoints in slow start mode receive no traffic in slow start window.
	// If not specified, the default is 10%.
	MinWeightPercent *v31.Percent `protobuf:"bytes,3,opt,name=min_weight_percent,json=minWeightPercent,proto3" json:"min_weight_percent,omitempty"`
}

func (x *SlowStartConfig) Reset() {
	*x = SlowStartConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_envoy_extensions_load_balancing_policies_common_v3_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SlowStartConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SlowStartConfig) ProtoMessage() {}

func (x *SlowStartConfig) ProtoReflect() protoreflect.Message {
	mi := &file_envoy_extensions_load_balancing_policies_common_v3_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SlowStartConfig.ProtoReflect.Descriptor instead.
func (*SlowStartConfig) Descriptor() ([]byte, []int) {
	return file_envoy_extensions_load_balancing_policies_common_v3_common_proto_rawDescGZIP(), []int{1}
}

func (x *SlowStartConfig) GetSlowStartWindow() *durationpb.Duration {
	if x != nil {
		return x.SlowStartWindow
	}
	return nil
}

func (x *SlowStartConfig) GetAggression() *v3.RuntimeDouble {
	if x != nil {
		return x.Aggression
	}
	return nil
}

func (x *SlowStartConfig) GetMinWeightPercent() *v31.Percent {
	if x != nil {
		return x.MinWeightPercent
	}
	return nil
}

// Common Configuration for all consistent hashing load balancers (MaglevLb, RingHashLb, etc.)
type ConsistentHashingLbConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// If set to “true“, the cluster will use hostname instead of the resolved
	// address as the key to consistently hash to an upstream host. Only valid for StrictDNS clusters with hostnames which resolve to a single IP address.
	UseHostnameForHashing bool `protobuf:"varint,1,opt,name=use_hostname_for_hashing,json=useHostnameForHashing,proto3" json:"use_hostname_for_hashing,omitempty"`
	// Configures percentage of average cluster load to bound per upstream host. For example, with a value of 150
	// no upstream host will get a load more than 1.5 times the average load of all the hosts in the cluster.
	// If not specified, the load is not bounded for any upstream host. Typical value for this parameter is between 120 and 200.
	// Minimum is 100.
	//
	// Applies to both Ring Hash and Maglev load balancers.
	//
	// This is implemented based on the method described in the paper https://arxiv.org/abs/1608.01350. For the specified
	// “hash_balance_factor“, requests to any upstream host are capped at “hash_balance_factor/100“ times the average number of requests
	// across the cluster. When a request arrives for an upstream host that is currently serving at its max capacity, linear probing
	// is used to identify an eligible host. Further, the linear probe is implemented using a random jump in hosts ring/table to identify
	// the eligible host (this technique is as described in the paper https://arxiv.org/abs/1908.08762 - the random jump avoids the
	// cascading overflow effect when choosing the next host in the ring/table).
	//
	// If weights are specified on the hosts, they are respected.
	//
	// This is an O(N) algorithm, unlike other load balancers. Using a lower “hash_balance_factor“ results in more hosts
	// being probed, so use a higher value if you require better performance.
	HashBalanceFactor *wrapperspb.UInt32Value `protobuf:"bytes,2,opt,name=hash_balance_factor,json=hashBalanceFactor,proto3" json:"hash_balance_factor,omitempty"`
}

func (x *ConsistentHashingLbConfig) Reset() {
	*x = ConsistentHashingLbConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_envoy_extensions_load_balancing_policies_common_v3_common_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsistentHashingLbConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsistentHashingLbConfig) ProtoMessage() {}

func (x *ConsistentHashingLbConfig) ProtoReflect() protoreflect.Message {
	mi := &file_envoy_extensions_load_balancing_policies_common_v3_common_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsistentHashingLbConfig.ProtoReflect.Descriptor instead.
func (*ConsistentHashingLbConfig) Descriptor() ([]byte, []int) {
	return file_envoy_extensions_load_balancing_policies_common_v3_common_proto_rawDescGZIP(), []int{2}
}

func (x *ConsistentHashingLbConfig) GetUseHostnameForHashing() bool {
	if x != nil {
		return x.UseHostnameForHashing
	}
	return false
}

func (x *ConsistentHashingLbConfig) GetHashBalanceFactor() *wrapperspb.UInt32Value {
	if x != nil {
		return x.HashBalanceFactor
	}
	return nil
}

// Configuration for :ref:`zone aware routing
// <arch_overview_load_balancing_zone_aware_routing>`.
type LocalityLbConfig_ZoneAwareLbConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Configures percentage of requests that will be considered for zone aware routing
	// if zone aware routing is configured. If not specified, the default is 100%.
	// * :ref:`runtime values <config_cluster_manager_cluster_runtime_zone_routing>`.
	// * :ref:`Zone aware routing support <arch_overview_load_balancing_zone_aware_routing>`.
	RoutingEnabled *v31.Percent `protobuf:"bytes,1,opt,name=routing_enabled,json=routingEnabled,proto3" json:"routing_enabled,omitempty"`
	// Configures minimum upstream cluster size required for zone aware routing
	// If upstream cluster size is less than specified, zone aware routing is not performed
	// even if zone aware routing is configured. If not specified, the default is 6.
	// * :ref:`runtime values <config_cluster_manager_cluster_runtime_zone_routing>`.
	// * :ref:`Zone aware routing support <arch_overview_load_balancing_zone_aware_routing>`.
	MinClusterSize *wrapperspb.UInt64Value `protobuf:"bytes,2,opt,name=min_cluster_size,json=minClusterSize,proto3" json:"min_cluster_size,omitempty"`
	// If set to true, Envoy will not consider any hosts when the cluster is in :ref:`panic
	// mode<arch_overview_load_balancing_panic_threshold>`. Instead, the cluster will fail all
	// requests as if all hosts are unhealthy. This can help avoid potentially overwhelming a
	// failing service.
	FailTrafficOnPanic bool `protobuf:"varint,3,opt,name=fail_traffic_on_panic,json=failTrafficOnPanic,proto3" json:"fail_traffic_on_panic,omitempty"`
}

func (x *LocalityLbConfig_ZoneAwareLbConfig) Reset() {
	*x = LocalityLbConfig_ZoneAwareLbConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_envoy_extensions_load_balancing_policies_common_v3_common_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LocalityLbConfig_ZoneAwareLbConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocalityLbConfig_ZoneAwareLbConfig) ProtoMessage() {}

func (x *LocalityLbConfig_ZoneAwareLbConfig) ProtoReflect() protoreflect.Message {
	mi := &file_envoy_extensions_load_balancing_policies_common_v3_common_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocalityLbConfig_ZoneAwareLbConfig.ProtoReflect.Descriptor instead.
func (*LocalityLbConfig_ZoneAwareLbConfig) Descriptor() ([]byte, []int) {
	return file_envoy_extensions_load_balancing_policies_common_v3_common_proto_rawDescGZIP(), []int{0, 0}
}

func (x *LocalityLbConfig_ZoneAwareLbConfig) GetRoutingEnabled() *v31.Percent {
	if x != nil {
		return x.RoutingEnabled
	}
	return nil
}

func (x *LocalityLbConfig_ZoneAwareLbConfig) GetMinClusterSize() *wrapperspb.UInt64Value {
	if x != nil {
		return x.MinClusterSize
	}
	return nil
}

func (x *LocalityLbConfig_ZoneAwareLbConfig) GetFailTrafficOnPanic() bool {
	if x != nil {
		return x.FailTrafficOnPanic
	}
	return false
}

// Configuration for :ref:`locality weighted load balancing
// <arch_overview_load_balancing_locality_weighted_lb>`
type LocalityLbConfig_LocalityWeightedLbConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *LocalityLbConfig_LocalityWeightedLbConfig) Reset() {
	*x = LocalityLbConfig_LocalityWeightedLbConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_envoy_extensions_load_balancing_policies_common_v3_common_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LocalityLbConfig_LocalityWeightedLbConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LocalityLbConfig_LocalityWeightedLbConfig) ProtoMessage() {}

func (x *LocalityLbConfig_LocalityWeightedLbConfig) ProtoReflect() protoreflect.Message {
	mi := &file_envoy_extensions_load_balancing_policies_common_v3_common_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LocalityLbConfig_LocalityWeightedLbConfig.ProtoReflect.Descriptor instead.
func (*LocalityLbConfig_LocalityWeightedLbConfig) Descriptor() ([]byte, []int) {
	return file_envoy_extensions_load_balancing_policies_common_v3_common_proto_rawDescGZIP(), []int{0, 1}
}

var File_envoy_extensions_load_balancing_policies_common_v3_common_proto protoreflect.FileDescriptor

var file_envoy_extensions_load_balancing_policies_common_v3_common_proto_rawDesc = []byte{
	0x0a, 0x3f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x69, 0x6e,
	0x67, 0x5f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65, 0x73, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2f, 0x76, 0x33, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x32, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x69,
	0x6e, 0x67, 0x5f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x76, 0x33, 0x1a, 0x1f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x33, 0x2f, 0x62, 0x61, 0x73, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2f, 0x74, 0x79,
	0x70, 0x65, 0x2f, 0x76, 0x33, 0x2f, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x75, 0x64, 0x70, 0x61, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcf, 0x04, 0x0a, 0x10,
	0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x4c, 0x62, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x12, 0x89, 0x01, 0x0a, 0x14, 0x7a, 0x6f, 0x6e, 0x65, 0x5f, 0x61, 0x77, 0x61, 0x72, 0x65, 0x5f,
	0x6c, 0x62, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x56, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x69, 0x6e,
	0x67, 0x5f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x76, 0x33, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x4c, 0x62, 0x43,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x5a, 0x6f, 0x6e, 0x65, 0x41, 0x77, 0x61, 0x72, 0x65, 0x4c,
	0x62, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x48, 0x00, 0x52, 0x11, 0x7a, 0x6f, 0x6e, 0x65, 0x41,
	0x77, 0x61, 0x72, 0x65, 0x4c, 0x62, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x9e, 0x01, 0x0a,
	0x1b, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74,
	0x65, 0x64, 0x5f, 0x6c, 0x62, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x5d, 0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e,
	0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x62, 0x61, 0x6c, 0x61, 0x6e,
	0x63, 0x69, 0x6e, 0x67, 0x5f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65, 0x73, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x33, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x69, 0x74, 0x79,
	0x4c, 0x62, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x69, 0x74,
	0x79, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x65, 0x64, 0x4c, 0x62, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x48, 0x00, 0x52, 0x18, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x57, 0x65, 0x69,
	0x67, 0x68, 0x74, 0x65, 0x64, 0x4c, 0x62, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x1a, 0xcf, 0x01,
	0x0a, 0x11, 0x5a, 0x6f, 0x6e, 0x65, 0x41, 0x77, 0x61, 0x72, 0x65, 0x4c, 0x62, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x12, 0x3f, 0x0a, 0x0f, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x65,
	0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x65,
	0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x76, 0x33, 0x2e, 0x50, 0x65, 0x72,
	0x63, 0x65, 0x6e, 0x74, 0x52, 0x0e, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x45, 0x6e, 0x61,
	0x62, 0x6c, 0x65, 0x64, 0x12, 0x46, 0x0a, 0x10, 0x6d, 0x69, 0x6e, 0x5f, 0x63, 0x6c, 0x75, 0x73,
	0x74, 0x65, 0x72, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x55, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0e, 0x6d, 0x69,
	0x6e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x31, 0x0a, 0x15,
	0x66, 0x61, 0x69, 0x6c, 0x5f, 0x74, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x5f, 0x6f, 0x6e, 0x5f,
	0x70, 0x61, 0x6e, 0x69, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x12, 0x66, 0x61, 0x69,
	0x6c, 0x54, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x4f, 0x6e, 0x50, 0x61, 0x6e, 0x69, 0x63, 0x1a,
	0x1a, 0x0a, 0x18, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x57, 0x65, 0x69, 0x67, 0x68,
	0x74, 0x65, 0x64, 0x4c, 0x62, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x20, 0x0a, 0x19, 0x6c,
	0x6f, 0x63, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x73,
	0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x03, 0xf8, 0x42, 0x01, 0x22, 0xe3, 0x01,
	0x0a, 0x0f, 0x53, 0x6c, 0x6f, 0x77, 0x53, 0x74, 0x61, 0x72, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x12, 0x45, 0x0a, 0x11, 0x73, 0x6c, 0x6f, 0x77, 0x5f, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f,
	0x77, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0f, 0x73, 0x6c, 0x6f, 0x77, 0x53, 0x74, 0x61,
	0x72, 0x74, 0x57, 0x69, 0x6e, 0x64, 0x6f, 0x77, 0x12, 0x43, 0x0a, 0x0a, 0x61, 0x67, 0x67, 0x72,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x65,
	0x6e, 0x76, 0x6f, 0x79, 0x2e, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x76, 0x33, 0x2e, 0x52, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x44, 0x6f, 0x75, 0x62, 0x6c,
	0x65, 0x52, 0x0a, 0x61, 0x67, 0x67, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x44, 0x0a,
	0x12, 0x6d, 0x69, 0x6e, 0x5f, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x5f, 0x70, 0x65, 0x72, 0x63,
	0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x65, 0x6e, 0x76, 0x6f,
	0x79, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x76, 0x33, 0x2e, 0x50, 0x65, 0x72, 0x63, 0x65, 0x6e,
	0x74, 0x52, 0x10, 0x6d, 0x69, 0x6e, 0x57, 0x65, 0x69, 0x67, 0x68, 0x74, 0x50, 0x65, 0x72, 0x63,
	0x65, 0x6e, 0x74, 0x22, 0xab, 0x01, 0x0a, 0x19, 0x43, 0x6f, 0x6e, 0x73, 0x69, 0x73, 0x74, 0x65,
	0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x69, 0x6e, 0x67, 0x4c, 0x62, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x12, 0x37, 0x0a, 0x18, 0x75, 0x73, 0x65, 0x5f, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d,
	0x65, 0x5f, 0x66, 0x6f, 0x72, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x15, 0x75, 0x73, 0x65, 0x48, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65,
	0x46, 0x6f, 0x72, 0x48, 0x61, 0x73, 0x68, 0x69, 0x6e, 0x67, 0x12, 0x55, 0x0a, 0x13, 0x68, 0x61,
	0x73, 0x68, 0x5f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x66, 0x61, 0x63, 0x74, 0x6f,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e, 0x74, 0x33, 0x32,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x2a, 0x02, 0x28, 0x64, 0x52, 0x11,
	0x68, 0x61, 0x73, 0x68, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x46, 0x61, 0x63, 0x74, 0x6f,
	0x72, 0x42, 0xbd, 0x01, 0xba, 0x80, 0xc8, 0xd1, 0x06, 0x02, 0x10, 0x02, 0x0a, 0x40, 0x69, 0x6f,
	0x2e, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x65, 0x6e, 0x76, 0x6f,
	0x79, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x6c, 0x6f, 0x61,
	0x64, 0x5f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x69, 0x6e, 0x67, 0x5f, 0x70, 0x6f, 0x6c, 0x69,
	0x63, 0x69, 0x65, 0x73, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x33, 0x42, 0x0b,
	0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x62, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x70,
	0x72, 0x6f, 0x78, 0x79, 0x2f, 0x67, 0x6f, 0x2d, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2d,
	0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x65, 0x6e, 0x76, 0x6f, 0x79, 0x2f, 0x65, 0x78, 0x74, 0x65,
	0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x62, 0x61, 0x6c, 0x61,
	0x6e, 0x63, 0x69, 0x6e, 0x67, 0x5f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65, 0x73, 0x2f, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x33, 0x3b, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x76,
	0x33, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_envoy_extensions_load_balancing_policies_common_v3_common_proto_rawDescOnce sync.Once
	file_envoy_extensions_load_balancing_policies_common_v3_common_proto_rawDescData = file_envoy_extensions_load_balancing_policies_common_v3_common_proto_rawDesc
)

func file_envoy_extensions_load_balancing_policies_common_v3_common_proto_rawDescGZIP() []byte {
	file_envoy_extensions_load_balancing_policies_common_v3_common_proto_rawDescOnce.Do(func() {
		file_envoy_extensions_load_balancing_policies_common_v3_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_envoy_extensions_load_balancing_policies_common_v3_common_proto_rawDescData)
	})
	return file_envoy_extensions_load_balancing_policies_common_v3_common_proto_rawDescData
}

var file_envoy_extensions_load_balancing_policies_common_v3_common_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_envoy_extensions_load_balancing_policies_common_v3_common_proto_goTypes = []interface{}{
	(*LocalityLbConfig)(nil),                          // 0: envoy.extensions.load_balancing_policies.common.v3.LocalityLbConfig
	(*SlowStartConfig)(nil),                           // 1: envoy.extensions.load_balancing_policies.common.v3.SlowStartConfig
	(*ConsistentHashingLbConfig)(nil),                 // 2: envoy.extensions.load_balancing_policies.common.v3.ConsistentHashingLbConfig
	(*LocalityLbConfig_ZoneAwareLbConfig)(nil),        // 3: envoy.extensions.load_balancing_policies.common.v3.LocalityLbConfig.ZoneAwareLbConfig
	(*LocalityLbConfig_LocalityWeightedLbConfig)(nil), // 4: envoy.extensions.load_balancing_policies.common.v3.LocalityLbConfig.LocalityWeightedLbConfig
	(*durationpb.Duration)(nil),                       // 5: google.protobuf.Duration
	(*v3.RuntimeDouble)(nil),                          // 6: envoy.config.core.v3.RuntimeDouble
	(*v31.Percent)(nil),                               // 7: envoy.type.v3.Percent
	(*wrapperspb.UInt32Value)(nil),                    // 8: google.protobuf.UInt32Value
	(*wrapperspb.UInt64Value)(nil),                    // 9: google.protobuf.UInt64Value
}
var file_envoy_extensions_load_balancing_policies_common_v3_common_proto_depIdxs = []int32{
	3, // 0: envoy.extensions.load_balancing_policies.common.v3.LocalityLbConfig.zone_aware_lb_config:type_name -> envoy.extensions.load_balancing_policies.common.v3.LocalityLbConfig.ZoneAwareLbConfig
	4, // 1: envoy.extensions.load_balancing_policies.common.v3.LocalityLbConfig.locality_weighted_lb_config:type_name -> envoy.extensions.load_balancing_policies.common.v3.LocalityLbConfig.LocalityWeightedLbConfig
	5, // 2: envoy.extensions.load_balancing_policies.common.v3.SlowStartConfig.slow_start_window:type_name -> google.protobuf.Duration
	6, // 3: envoy.extensions.load_balancing_policies.common.v3.SlowStartConfig.aggression:type_name -> envoy.config.core.v3.RuntimeDouble
	7, // 4: envoy.extensions.load_balancing_policies.common.v3.SlowStartConfig.min_weight_percent:type_name -> envoy.type.v3.Percent
	8, // 5: envoy.extensions.load_balancing_policies.common.v3.ConsistentHashingLbConfig.hash_balance_factor:type_name -> google.protobuf.UInt32Value
	7, // 6: envoy.extensions.load_balancing_policies.common.v3.LocalityLbConfig.ZoneAwareLbConfig.routing_enabled:type_name -> envoy.type.v3.Percent
	9, // 7: envoy.extensions.load_balancing_policies.common.v3.LocalityLbConfig.ZoneAwareLbConfig.min_cluster_size:type_name -> google.protobuf.UInt64Value
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_envoy_extensions_load_balancing_policies_common_v3_common_proto_init() }
func file_envoy_extensions_load_balancing_policies_common_v3_common_proto_init() {
	if File_envoy_extensions_load_balancing_policies_common_v3_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_envoy_extensions_load_balancing_policies_common_v3_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LocalityLbConfig); i {
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
		file_envoy_extensions_load_balancing_policies_common_v3_common_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SlowStartConfig); i {
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
		file_envoy_extensions_load_balancing_policies_common_v3_common_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConsistentHashingLbConfig); i {
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
		file_envoy_extensions_load_balancing_policies_common_v3_common_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LocalityLbConfig_ZoneAwareLbConfig); i {
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
		file_envoy_extensions_load_balancing_policies_common_v3_common_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LocalityLbConfig_LocalityWeightedLbConfig); i {
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
	file_envoy_extensions_load_balancing_policies_common_v3_common_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*LocalityLbConfig_ZoneAwareLbConfig_)(nil),
		(*LocalityLbConfig_LocalityWeightedLbConfig_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_envoy_extensions_load_balancing_policies_common_v3_common_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_envoy_extensions_load_balancing_policies_common_v3_common_proto_goTypes,
		DependencyIndexes: file_envoy_extensions_load_balancing_policies_common_v3_common_proto_depIdxs,
		MessageInfos:      file_envoy_extensions_load_balancing_policies_common_v3_common_proto_msgTypes,
	}.Build()
	File_envoy_extensions_load_balancing_policies_common_v3_common_proto = out.File
	file_envoy_extensions_load_balancing_policies_common_v3_common_proto_rawDesc = nil
	file_envoy_extensions_load_balancing_policies_common_v3_common_proto_goTypes = nil
	file_envoy_extensions_load_balancing_policies_common_v3_common_proto_depIdxs = nil
}
