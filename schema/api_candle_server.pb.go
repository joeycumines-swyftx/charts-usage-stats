// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.5
// source: schema/api_candle_server.proto

package schema

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ApiEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Data:
	//
	//	*ApiEvent_AccessLog
	Data isApiEvent_Data `protobuf_oneof:"data"`
}

func (x *ApiEvent) Reset() {
	*x = ApiEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_schema_api_candle_server_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApiEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApiEvent) ProtoMessage() {}

func (x *ApiEvent) ProtoReflect() protoreflect.Message {
	mi := &file_schema_api_candle_server_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApiEvent.ProtoReflect.Descriptor instead.
func (*ApiEvent) Descriptor() ([]byte, []int) {
	return file_schema_api_candle_server_proto_rawDescGZIP(), []int{0}
}

func (m *ApiEvent) GetData() isApiEvent_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *ApiEvent) GetAccessLog() *ApiAccessLog {
	if x, ok := x.GetData().(*ApiEvent_AccessLog); ok {
		return x.AccessLog
	}
	return nil
}

type isApiEvent_Data interface {
	isApiEvent_Data()
}

type ApiEvent_AccessLog struct {
	AccessLog *ApiAccessLog `protobuf:"bytes,6,opt,name=access_log,json=accessLog,proto3,oneof"`
}

func (*ApiEvent_AccessLog) isApiEvent_Data() {}

type ApiAccessLog struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StatusCode    int32                `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	ContentLength int64                `protobuf:"varint,2,opt,name=content_length,json=contentLength,proto3" json:"content_length,omitempty"`
	Duration      *durationpb.Duration `protobuf:"bytes,3,opt,name=duration,proto3" json:"duration,omitempty"`
	// Types that are assignable to Data:
	//
	//	*ApiAccessLog_GetBars_
	//	*ApiAccessLog_LastKnownPrice_
	//	*ApiAccessLog_Rate_
	Data isApiAccessLog_Data `protobuf_oneof:"data"`
}

func (x *ApiAccessLog) Reset() {
	*x = ApiAccessLog{}
	if protoimpl.UnsafeEnabled {
		mi := &file_schema_api_candle_server_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApiAccessLog) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApiAccessLog) ProtoMessage() {}

func (x *ApiAccessLog) ProtoReflect() protoreflect.Message {
	mi := &file_schema_api_candle_server_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApiAccessLog.ProtoReflect.Descriptor instead.
func (*ApiAccessLog) Descriptor() ([]byte, []int) {
	return file_schema_api_candle_server_proto_rawDescGZIP(), []int{1}
}

func (x *ApiAccessLog) GetStatusCode() int32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *ApiAccessLog) GetContentLength() int64 {
	if x != nil {
		return x.ContentLength
	}
	return 0
}

func (x *ApiAccessLog) GetDuration() *durationpb.Duration {
	if x != nil {
		return x.Duration
	}
	return nil
}

func (m *ApiAccessLog) GetData() isApiAccessLog_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *ApiAccessLog) GetGetBars() *ApiAccessLog_GetBars {
	if x, ok := x.GetData().(*ApiAccessLog_GetBars_); ok {
		return x.GetBars
	}
	return nil
}

func (x *ApiAccessLog) GetLastKnownPrice() *ApiAccessLog_LastKnownPrice {
	if x, ok := x.GetData().(*ApiAccessLog_LastKnownPrice_); ok {
		return x.LastKnownPrice
	}
	return nil
}

func (x *ApiAccessLog) GetRate() *ApiAccessLog_Rate {
	if x, ok := x.GetData().(*ApiAccessLog_Rate_); ok {
		return x.Rate
	}
	return nil
}

type isApiAccessLog_Data interface {
	isApiAccessLog_Data()
}

type ApiAccessLog_GetBars_ struct {
	GetBars *ApiAccessLog_GetBars `protobuf:"bytes,7,opt,name=get_bars,json=getBars,proto3,oneof"`
}

type ApiAccessLog_LastKnownPrice_ struct {
	LastKnownPrice *ApiAccessLog_LastKnownPrice `protobuf:"bytes,8,opt,name=last_known_price,json=lastKnownPrice,proto3,oneof"`
}

type ApiAccessLog_Rate_ struct {
	Rate *ApiAccessLog_Rate `protobuf:"bytes,9,opt,name=rate,proto3,oneof"`
}

func (*ApiAccessLog_GetBars_) isApiAccessLog_Data() {}

func (*ApiAccessLog_LastKnownPrice_) isApiAccessLog_Data() {}

func (*ApiAccessLog_Rate_) isApiAccessLog_Data() {}

type ApiAccessLog_GetBars struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PrimaryAsset   string                 `protobuf:"bytes,1,opt,name=primary_asset,json=primaryAsset,proto3" json:"primary_asset,omitempty"`
	SecondaryAsset string                 `protobuf:"bytes,2,opt,name=secondary_asset,json=secondaryAsset,proto3" json:"secondary_asset,omitempty"`
	MarketSide     MarketSide             `protobuf:"varint,3,opt,name=market_side,json=marketSide,proto3,enum=charts.MarketSide" json:"market_side,omitempty"`
	StartTime      *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime        *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	Resolution     *durationpb.Duration   `protobuf:"bytes,6,opt,name=resolution,proto3" json:"resolution,omitempty"`
}

func (x *ApiAccessLog_GetBars) Reset() {
	*x = ApiAccessLog_GetBars{}
	if protoimpl.UnsafeEnabled {
		mi := &file_schema_api_candle_server_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApiAccessLog_GetBars) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApiAccessLog_GetBars) ProtoMessage() {}

func (x *ApiAccessLog_GetBars) ProtoReflect() protoreflect.Message {
	mi := &file_schema_api_candle_server_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApiAccessLog_GetBars.ProtoReflect.Descriptor instead.
func (*ApiAccessLog_GetBars) Descriptor() ([]byte, []int) {
	return file_schema_api_candle_server_proto_rawDescGZIP(), []int{1, 0}
}

func (x *ApiAccessLog_GetBars) GetPrimaryAsset() string {
	if x != nil {
		return x.PrimaryAsset
	}
	return ""
}

func (x *ApiAccessLog_GetBars) GetSecondaryAsset() string {
	if x != nil {
		return x.SecondaryAsset
	}
	return ""
}

func (x *ApiAccessLog_GetBars) GetMarketSide() MarketSide {
	if x != nil {
		return x.MarketSide
	}
	return MarketSide_MARKET_SIDE_UNSPECIFIED
}

func (x *ApiAccessLog_GetBars) GetStartTime() *timestamppb.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *ApiAccessLog_GetBars) GetEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

func (x *ApiAccessLog_GetBars) GetResolution() *durationpb.Duration {
	if x != nil {
		return x.Resolution
	}
	return nil
}

type ApiAccessLog_LastKnownPrice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PrimaryAsset   string `protobuf:"bytes,1,opt,name=primary_asset,json=primaryAsset,proto3" json:"primary_asset,omitempty"`
	SecondaryAsset string `protobuf:"bytes,2,opt,name=secondary_asset,json=secondaryAsset,proto3" json:"secondary_asset,omitempty"`
}

func (x *ApiAccessLog_LastKnownPrice) Reset() {
	*x = ApiAccessLog_LastKnownPrice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_schema_api_candle_server_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApiAccessLog_LastKnownPrice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApiAccessLog_LastKnownPrice) ProtoMessage() {}

func (x *ApiAccessLog_LastKnownPrice) ProtoReflect() protoreflect.Message {
	mi := &file_schema_api_candle_server_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApiAccessLog_LastKnownPrice.ProtoReflect.Descriptor instead.
func (*ApiAccessLog_LastKnownPrice) Descriptor() ([]byte, []int) {
	return file_schema_api_candle_server_proto_rawDescGZIP(), []int{1, 1}
}

func (x *ApiAccessLog_LastKnownPrice) GetPrimaryAsset() string {
	if x != nil {
		return x.PrimaryAsset
	}
	return ""
}

func (x *ApiAccessLog_LastKnownPrice) GetSecondaryAsset() string {
	if x != nil {
		return x.SecondaryAsset
	}
	return ""
}

type ApiAccessLog_Rate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PrimaryAsset   string                 `protobuf:"bytes,1,opt,name=primary_asset,json=primaryAsset,proto3" json:"primary_asset,omitempty"`
	SecondaryAsset string                 `protobuf:"bytes,2,opt,name=secondary_asset,json=secondaryAsset,proto3" json:"secondary_asset,omitempty"`
	Timestamp      *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *ApiAccessLog_Rate) Reset() {
	*x = ApiAccessLog_Rate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_schema_api_candle_server_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApiAccessLog_Rate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApiAccessLog_Rate) ProtoMessage() {}

func (x *ApiAccessLog_Rate) ProtoReflect() protoreflect.Message {
	mi := &file_schema_api_candle_server_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApiAccessLog_Rate.ProtoReflect.Descriptor instead.
func (*ApiAccessLog_Rate) Descriptor() ([]byte, []int) {
	return file_schema_api_candle_server_proto_rawDescGZIP(), []int{1, 2}
}

func (x *ApiAccessLog_Rate) GetPrimaryAsset() string {
	if x != nil {
		return x.PrimaryAsset
	}
	return ""
}

func (x *ApiAccessLog_Rate) GetSecondaryAsset() string {
	if x != nil {
		return x.SecondaryAsset
	}
	return ""
}

func (x *ApiAccessLog_Rate) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

var File_schema_api_candle_server_proto protoreflect.FileDescriptor

var file_schema_api_candle_server_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2f, 0x61, 0x70, 0x69, 0x5f, 0x63, 0x61, 0x6e,
	0x64, 0x6c, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x06, 0x63, 0x68, 0x61, 0x72, 0x74, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x18, 0x73, 0x63, 0x68, 0x65, 0x6d,
	0x61, 0x2f, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x5f, 0x73, 0x69, 0x64, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x49, 0x0a, 0x08, 0x41, 0x70, 0x69, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12,
	0x35, 0x0a, 0x0a, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x6c, 0x6f, 0x67, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63, 0x68, 0x61, 0x72, 0x74, 0x73, 0x2e, 0x41, 0x70, 0x69,
	0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4c, 0x6f, 0x67, 0x48, 0x00, 0x52, 0x09, 0x61, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x4c, 0x6f, 0x67, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xff,
	0x06, 0x0a, 0x0c, 0x41, 0x70, 0x69, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4c, 0x6f, 0x67, 0x12,
	0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x25, 0x0a, 0x0e, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x5f, 0x6c, 0x65, 0x6e, 0x67,
	0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12, 0x35, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x39,
	0x0a, 0x08, 0x67, 0x65, 0x74, 0x5f, 0x62, 0x61, 0x72, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1c, 0x2e, 0x63, 0x68, 0x61, 0x72, 0x74, 0x73, 0x2e, 0x41, 0x70, 0x69, 0x41, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x4c, 0x6f, 0x67, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x61, 0x72, 0x73, 0x48, 0x00,
	0x52, 0x07, 0x67, 0x65, 0x74, 0x42, 0x61, 0x72, 0x73, 0x12, 0x4f, 0x0a, 0x10, 0x6c, 0x61, 0x73,
	0x74, 0x5f, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x63, 0x68, 0x61, 0x72, 0x74, 0x73, 0x2e, 0x41, 0x70, 0x69,
	0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4c, 0x6f, 0x67, 0x2e, 0x4c, 0x61, 0x73, 0x74, 0x4b, 0x6e,
	0x6f, 0x77, 0x6e, 0x50, 0x72, 0x69, 0x63, 0x65, 0x48, 0x00, 0x52, 0x0e, 0x6c, 0x61, 0x73, 0x74,
	0x4b, 0x6e, 0x6f, 0x77, 0x6e, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x2f, 0x0a, 0x04, 0x72, 0x61,
	0x74, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x68, 0x61, 0x72, 0x74,
	0x73, 0x2e, 0x41, 0x70, 0x69, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4c, 0x6f, 0x67, 0x2e, 0x52,
	0x61, 0x74, 0x65, 0x48, 0x00, 0x52, 0x04, 0x72, 0x61, 0x74, 0x65, 0x1a, 0xb9, 0x02, 0x0a, 0x07,
	0x47, 0x65, 0x74, 0x42, 0x61, 0x72, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x72, 0x69, 0x6d, 0x61,
	0x72, 0x79, 0x5f, 0x61, 0x73, 0x73, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x41, 0x73, 0x73, 0x65, 0x74, 0x12, 0x27, 0x0a, 0x0f,
	0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x61, 0x72, 0x79, 0x5f, 0x61, 0x73, 0x73, 0x65, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x61, 0x72, 0x79,
	0x41, 0x73, 0x73, 0x65, 0x74, 0x12, 0x33, 0x0a, 0x0b, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x5f,
	0x73, 0x69, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x63, 0x68, 0x61,
	0x72, 0x74, 0x73, 0x2e, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x53, 0x69, 0x64, 0x65, 0x52, 0x0a,
	0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x53, 0x69, 0x64, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x39, 0x0a, 0x0a,
	0x72, 0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x72, 0x65, 0x73,
	0x6f, 0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x5e, 0x0a, 0x0e, 0x4c, 0x61, 0x73, 0x74, 0x4b,
	0x6e, 0x6f, 0x77, 0x6e, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x72, 0x69,
	0x6d, 0x61, 0x72, 0x79, 0x5f, 0x61, 0x73, 0x73, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x41, 0x73, 0x73, 0x65, 0x74, 0x12, 0x27,
	0x0a, 0x0f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x61, 0x72, 0x79, 0x5f, 0x61, 0x73, 0x73, 0x65,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x61,
	0x72, 0x79, 0x41, 0x73, 0x73, 0x65, 0x74, 0x1a, 0x8e, 0x01, 0x0a, 0x04, 0x52, 0x61, 0x74, 0x65,
	0x12, 0x23, 0x0a, 0x0d, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x5f, 0x61, 0x73, 0x73, 0x65,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79,
	0x41, 0x73, 0x73, 0x65, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x61,
	0x72, 0x79, 0x5f, 0x61, 0x73, 0x73, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e,
	0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x61, 0x72, 0x79, 0x41, 0x73, 0x73, 0x65, 0x74, 0x12, 0x38,
	0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x42, 0x39, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a,
	0x6f, 0x65, 0x79, 0x63, 0x75, 0x6d, 0x69, 0x6e, 0x65, 0x73, 0x2d, 0x73, 0x77, 0x79, 0x66, 0x74,
	0x78, 0x2f, 0x63, 0x68, 0x61, 0x72, 0x74, 0x73, 0x2d, 0x75, 0x73, 0x61, 0x67, 0x65, 0x2d, 0x73,
	0x74, 0x61, 0x74, 0x73, 0x2f, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_schema_api_candle_server_proto_rawDescOnce sync.Once
	file_schema_api_candle_server_proto_rawDescData = file_schema_api_candle_server_proto_rawDesc
)

func file_schema_api_candle_server_proto_rawDescGZIP() []byte {
	file_schema_api_candle_server_proto_rawDescOnce.Do(func() {
		file_schema_api_candle_server_proto_rawDescData = protoimpl.X.CompressGZIP(file_schema_api_candle_server_proto_rawDescData)
	})
	return file_schema_api_candle_server_proto_rawDescData
}

var file_schema_api_candle_server_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_schema_api_candle_server_proto_goTypes = []interface{}{
	(*ApiEvent)(nil),                    // 0: charts.ApiEvent
	(*ApiAccessLog)(nil),                // 1: charts.ApiAccessLog
	(*ApiAccessLog_GetBars)(nil),        // 2: charts.ApiAccessLog.GetBars
	(*ApiAccessLog_LastKnownPrice)(nil), // 3: charts.ApiAccessLog.LastKnownPrice
	(*ApiAccessLog_Rate)(nil),           // 4: charts.ApiAccessLog.Rate
	(*durationpb.Duration)(nil),         // 5: google.protobuf.Duration
	(MarketSide)(0),                     // 6: charts.MarketSide
	(*timestamppb.Timestamp)(nil),       // 7: google.protobuf.Timestamp
}
var file_schema_api_candle_server_proto_depIdxs = []int32{
	1,  // 0: charts.ApiEvent.access_log:type_name -> charts.ApiAccessLog
	5,  // 1: charts.ApiAccessLog.duration:type_name -> google.protobuf.Duration
	2,  // 2: charts.ApiAccessLog.get_bars:type_name -> charts.ApiAccessLog.GetBars
	3,  // 3: charts.ApiAccessLog.last_known_price:type_name -> charts.ApiAccessLog.LastKnownPrice
	4,  // 4: charts.ApiAccessLog.rate:type_name -> charts.ApiAccessLog.Rate
	6,  // 5: charts.ApiAccessLog.GetBars.market_side:type_name -> charts.MarketSide
	7,  // 6: charts.ApiAccessLog.GetBars.start_time:type_name -> google.protobuf.Timestamp
	7,  // 7: charts.ApiAccessLog.GetBars.end_time:type_name -> google.protobuf.Timestamp
	5,  // 8: charts.ApiAccessLog.GetBars.resolution:type_name -> google.protobuf.Duration
	7,  // 9: charts.ApiAccessLog.Rate.timestamp:type_name -> google.protobuf.Timestamp
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_schema_api_candle_server_proto_init() }
func file_schema_api_candle_server_proto_init() {
	if File_schema_api_candle_server_proto != nil {
		return
	}
	file_schema_market_side_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_schema_api_candle_server_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApiEvent); i {
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
		file_schema_api_candle_server_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApiAccessLog); i {
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
		file_schema_api_candle_server_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApiAccessLog_GetBars); i {
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
		file_schema_api_candle_server_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApiAccessLog_LastKnownPrice); i {
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
		file_schema_api_candle_server_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApiAccessLog_Rate); i {
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
	file_schema_api_candle_server_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*ApiEvent_AccessLog)(nil),
	}
	file_schema_api_candle_server_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*ApiAccessLog_GetBars_)(nil),
		(*ApiAccessLog_LastKnownPrice_)(nil),
		(*ApiAccessLog_Rate_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_schema_api_candle_server_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_schema_api_candle_server_proto_goTypes,
		DependencyIndexes: file_schema_api_candle_server_proto_depIdxs,
		MessageInfos:      file_schema_api_candle_server_proto_msgTypes,
	}.Build()
	File_schema_api_candle_server_proto = out.File
	file_schema_api_candle_server_proto_rawDesc = nil
	file_schema_api_candle_server_proto_goTypes = nil
	file_schema_api_candle_server_proto_depIdxs = nil
}
