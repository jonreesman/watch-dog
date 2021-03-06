// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: watchdog.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

type SentimentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tweet string `protobuf:"bytes,1,opt,name=tweet,proto3" json:"tweet,omitempty"`
}

func (x *SentimentRequest) Reset() {
	*x = SentimentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_watchdog_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SentimentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SentimentRequest) ProtoMessage() {}

func (x *SentimentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_watchdog_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SentimentRequest.ProtoReflect.Descriptor instead.
func (*SentimentRequest) Descriptor() ([]byte, []int) {
	return file_watchdog_proto_rawDescGZIP(), []int{0}
}

func (x *SentimentRequest) GetTweet() string {
	if x != nil {
		return x.Tweet
	}
	return ""
}

type SentimentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Polarity float32 `protobuf:"fixed32,1,opt,name=polarity,proto3" json:"polarity,omitempty"`
}

func (x *SentimentResponse) Reset() {
	*x = SentimentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_watchdog_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SentimentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SentimentResponse) ProtoMessage() {}

func (x *SentimentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_watchdog_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SentimentResponse.ProtoReflect.Descriptor instead.
func (*SentimentResponse) Descriptor() ([]byte, []int) {
	return file_watchdog_proto_rawDescGZIP(), []int{1}
}

func (x *SentimentResponse) GetPolarity() float32 {
	if x != nil {
		return x.Polarity
	}
	return 0
}

type QuoteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Period string `protobuf:"bytes,2,opt,name=period,proto3" json:"period,omitempty"`
}

func (x *QuoteRequest) Reset() {
	*x = QuoteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_watchdog_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuoteRequest) ProtoMessage() {}

func (x *QuoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_watchdog_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuoteRequest.ProtoReflect.Descriptor instead.
func (*QuoteRequest) Descriptor() ([]byte, []int) {
	return file_watchdog_proto_rawDescGZIP(), []int{2}
}

func (x *QuoteRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *QuoteRequest) GetPeriod() string {
	if x != nil {
		return x.Period
	}
	return ""
}

type Quote struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Time  *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=time,proto3" json:"time,omitempty"`
	Price float32                `protobuf:"fixed32,2,opt,name=price,proto3" json:"price,omitempty"`
}

func (x *Quote) Reset() {
	*x = Quote{}
	if protoimpl.UnsafeEnabled {
		mi := &file_watchdog_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Quote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Quote) ProtoMessage() {}

func (x *Quote) ProtoReflect() protoreflect.Message {
	mi := &file_watchdog_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Quote.ProtoReflect.Descriptor instead.
func (*Quote) Descriptor() ([]byte, []int) {
	return file_watchdog_proto_rawDescGZIP(), []int{3}
}

func (x *Quote) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

func (x *Quote) GetPrice() float32 {
	if x != nil {
		return x.Price
	}
	return 0
}

type QuoteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Quotes []*Quote `protobuf:"bytes,1,rep,name=quotes,proto3" json:"quotes,omitempty"`
}

func (x *QuoteResponse) Reset() {
	*x = QuoteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_watchdog_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuoteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuoteResponse) ProtoMessage() {}

func (x *QuoteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_watchdog_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuoteResponse.ProtoReflect.Descriptor instead.
func (*QuoteResponse) Descriptor() ([]byte, []int) {
	return file_watchdog_proto_rawDescGZIP(), []int{4}
}

func (x *QuoteResponse) GetQuotes() []*Quote {
	if x != nil {
		return x.Quotes
	}
	return nil
}

var File_watchdog_proto protoreflect.FileDescriptor

var file_watchdog_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x77, 0x61, 0x74, 0x63, 0x68, 0x64, 0x6f, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x02, 0x70, 0x62, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x28, 0x0a, 0x10, 0x53, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x77, 0x65,
	0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x77, 0x65, 0x65, 0x74, 0x22,
	0x2f, 0x0a, 0x11, 0x53, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x6f, 0x6c, 0x61, 0x72, 0x69, 0x74, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x08, 0x70, 0x6f, 0x6c, 0x61, 0x72, 0x69, 0x74, 0x79,
	0x22, 0x3a, 0x0a, 0x0c, 0x51, 0x75, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x22, 0x4d, 0x0a, 0x05,
	0x51, 0x75, 0x6f, 0x74, 0x65, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x22, 0x32, 0x0a, 0x0d, 0x51,
	0x75, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x06,
	0x71, 0x75, 0x6f, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x70,
	0x62, 0x2e, 0x51, 0x75, 0x6f, 0x74, 0x65, 0x52, 0x06, 0x71, 0x75, 0x6f, 0x74, 0x65, 0x73, 0x32,
	0x44, 0x0a, 0x09, 0x53, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x37, 0x0a, 0x06,
	0x44, 0x65, 0x74, 0x65, 0x63, 0x74, 0x12, 0x14, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x6e, 0x74,
	0x69, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x70,
	0x62, 0x2e, 0x53, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x32, 0x39, 0x0a, 0x06, 0x51, 0x75, 0x6f, 0x74, 0x65, 0x73, 0x12,
	0x2f, 0x0a, 0x06, 0x44, 0x65, 0x74, 0x65, 0x63, 0x74, 0x12, 0x10, 0x2e, 0x70, 0x62, 0x2e, 0x51,
	0x75, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x70, 0x62,
	0x2e, 0x51, 0x75, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x23, 0x5a, 0x21, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a,
	0x6f, 0x6e, 0x72, 0x65, 0x65, 0x73, 0x6d, 0x61, 0x6e, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x54, 0x65,
	0x73, 0x74, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_watchdog_proto_rawDescOnce sync.Once
	file_watchdog_proto_rawDescData = file_watchdog_proto_rawDesc
)

func file_watchdog_proto_rawDescGZIP() []byte {
	file_watchdog_proto_rawDescOnce.Do(func() {
		file_watchdog_proto_rawDescData = protoimpl.X.CompressGZIP(file_watchdog_proto_rawDescData)
	})
	return file_watchdog_proto_rawDescData
}

var file_watchdog_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_watchdog_proto_goTypes = []interface{}{
	(*SentimentRequest)(nil),      // 0: pb.SentimentRequest
	(*SentimentResponse)(nil),     // 1: pb.SentimentResponse
	(*QuoteRequest)(nil),          // 2: pb.QuoteRequest
	(*Quote)(nil),                 // 3: pb.Quote
	(*QuoteResponse)(nil),         // 4: pb.QuoteResponse
	(*timestamppb.Timestamp)(nil), // 5: google.protobuf.Timestamp
}
var file_watchdog_proto_depIdxs = []int32{
	5, // 0: pb.Quote.time:type_name -> google.protobuf.Timestamp
	3, // 1: pb.QuoteResponse.quotes:type_name -> pb.Quote
	0, // 2: pb.Sentiment.Detect:input_type -> pb.SentimentRequest
	2, // 3: pb.Quotes.Detect:input_type -> pb.QuoteRequest
	1, // 4: pb.Sentiment.Detect:output_type -> pb.SentimentResponse
	4, // 5: pb.Quotes.Detect:output_type -> pb.QuoteResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_watchdog_proto_init() }
func file_watchdog_proto_init() {
	if File_watchdog_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_watchdog_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SentimentRequest); i {
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
		file_watchdog_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SentimentResponse); i {
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
		file_watchdog_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QuoteRequest); i {
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
		file_watchdog_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Quote); i {
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
		file_watchdog_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QuoteResponse); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_watchdog_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_watchdog_proto_goTypes,
		DependencyIndexes: file_watchdog_proto_depIdxs,
		MessageInfos:      file_watchdog_proto_msgTypes,
	}.Build()
	File_watchdog_proto = out.File
	file_watchdog_proto_rawDesc = nil
	file_watchdog_proto_goTypes = nil
	file_watchdog_proto_depIdxs = nil
}
