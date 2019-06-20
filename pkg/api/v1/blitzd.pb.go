// Code generated by protoc-gen-go. DO NOT EDIT.
// source: blitzd.proto

package v1

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Style int32

const (
	Style_UNKNOWN_STYLE Style = 0
	Style_NORMAL        Style = 1
	Style_HIGHLIGHT     Style = 2
	Style_WARNING       Style = 3
	Style_ERROR         Style = 4
	Style_CRITICAL      Style = 5
)

var Style_name = map[int32]string{
	0: "UNKNOWN_STYLE",
	1: "NORMAL",
	2: "HIGHLIGHT",
	3: "WARNING",
	4: "ERROR",
	5: "CRITICAL",
}

var Style_value = map[string]int32{
	"UNKNOWN_STYLE": 0,
	"NORMAL":        1,
	"HIGHLIGHT":     2,
	"WARNING":       3,
	"ERROR":         4,
	"CRITICAL":      5,
}

func (x Style) String() string {
	return proto.EnumName(Style_name, int32(x))
}

func (Style) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_84366feb8d987a6e, []int{0}
}

type Kind int32

const (
	Kind_UNKNOWN_KIND Kind = 0
	Kind_STATIC       Kind = 1
	Kind_TIME_BASED   Kind = 2
	Kind_EVENT_BASED  Kind = 3
)

var Kind_name = map[int32]string{
	0: "UNKNOWN_KIND",
	1: "STATIC",
	2: "TIME_BASED",
	3: "EVENT_BASED",
}

var Kind_value = map[string]int32{
	"UNKNOWN_KIND": 0,
	"STATIC":       1,
	"TIME_BASED":   2,
	"EVENT_BASED":  3,
}

func (x Kind) String() string {
	return proto.EnumName(Kind_name, int32(x))
}

func (Kind) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_84366feb8d987a6e, []int{1}
}

// The request message
type ShutdownRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShutdownRequest) Reset()         { *m = ShutdownRequest{} }
func (m *ShutdownRequest) String() string { return proto.CompactTextString(m) }
func (*ShutdownRequest) ProtoMessage()    {}
func (*ShutdownRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_84366feb8d987a6e, []int{0}
}

func (m *ShutdownRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShutdownRequest.Unmarshal(m, b)
}
func (m *ShutdownRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShutdownRequest.Marshal(b, m, deterministic)
}
func (m *ShutdownRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShutdownRequest.Merge(m, src)
}
func (m *ShutdownRequest) XXX_Size() int {
	return xxx_messageInfo_ShutdownRequest.Size(m)
}
func (m *ShutdownRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ShutdownRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ShutdownRequest proto.InternalMessageInfo

// The response message
type ShutdownResponse struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShutdownResponse) Reset()         { *m = ShutdownResponse{} }
func (m *ShutdownResponse) String() string { return proto.CompactTextString(m) }
func (*ShutdownResponse) ProtoMessage()    {}
func (*ShutdownResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_84366feb8d987a6e, []int{1}
}

func (m *ShutdownResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShutdownResponse.Unmarshal(m, b)
}
func (m *ShutdownResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShutdownResponse.Marshal(b, m, deterministic)
}
func (m *ShutdownResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShutdownResponse.Merge(m, src)
}
func (m *ShutdownResponse) XXX_Size() int {
	return xxx_messageInfo_ShutdownResponse.Size(m)
}
func (m *ShutdownResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ShutdownResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ShutdownResponse proto.InternalMessageInfo

func (m *ShutdownResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type Metric struct {
	Kind                 Kind                 `protobuf:"varint,1,opt,name=kind,proto3,enum=v1.Kind" json:"kind,omitempty"`
	Module               string               `protobuf:"bytes,2,opt,name=module,proto3" json:"module,omitempty"`
	Title                string               `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Interval             float64              `protobuf:"fixed64,4,opt,name=interval,proto3" json:"interval,omitempty"`
	Timeout              float64              `protobuf:"fixed64,5,opt,name=timeout,proto3" json:"timeout,omitempty"`
	UpdatedAt            *timestamp.Timestamp `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	ExpiredAfter         *timestamp.Timestamp `protobuf:"bytes,7,opt,name=expired_after,json=expiredAfter,proto3" json:"expired_after,omitempty"`
	Expired              bool                 `protobuf:"varint,8,opt,name=expired,proto3" json:"expired,omitempty"`
	Value                string               `protobuf:"bytes,9,opt,name=value,proto3" json:"value,omitempty"`
	Prefix               string               `protobuf:"bytes,10,opt,name=prefix,proto3" json:"prefix,omitempty"`
	Suffix               string               `protobuf:"bytes,11,opt,name=suffix,proto3" json:"suffix,omitempty"`
	Style                Style                `protobuf:"varint,12,opt,name=style,proto3,enum=v1.Style" json:"style,omitempty"`
	Text                 string               `protobuf:"bytes,13,opt,name=text,proto3" json:"text,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Metric) Reset()         { *m = Metric{} }
func (m *Metric) String() string { return proto.CompactTextString(m) }
func (*Metric) ProtoMessage()    {}
func (*Metric) Descriptor() ([]byte, []int) {
	return fileDescriptor_84366feb8d987a6e, []int{2}
}

func (m *Metric) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Metric.Unmarshal(m, b)
}
func (m *Metric) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Metric.Marshal(b, m, deterministic)
}
func (m *Metric) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Metric.Merge(m, src)
}
func (m *Metric) XXX_Size() int {
	return xxx_messageInfo_Metric.Size(m)
}
func (m *Metric) XXX_DiscardUnknown() {
	xxx_messageInfo_Metric.DiscardUnknown(m)
}

var xxx_messageInfo_Metric proto.InternalMessageInfo

func (m *Metric) GetKind() Kind {
	if m != nil {
		return m.Kind
	}
	return Kind_UNKNOWN_KIND
}

func (m *Metric) GetModule() string {
	if m != nil {
		return m.Module
	}
	return ""
}

func (m *Metric) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Metric) GetInterval() float64 {
	if m != nil {
		return m.Interval
	}
	return 0
}

func (m *Metric) GetTimeout() float64 {
	if m != nil {
		return m.Timeout
	}
	return 0
}

func (m *Metric) GetUpdatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.UpdatedAt
	}
	return nil
}

func (m *Metric) GetExpiredAfter() *timestamp.Timestamp {
	if m != nil {
		return m.ExpiredAfter
	}
	return nil
}

func (m *Metric) GetExpired() bool {
	if m != nil {
		return m.Expired
	}
	return false
}

func (m *Metric) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *Metric) GetPrefix() string {
	if m != nil {
		return m.Prefix
	}
	return ""
}

func (m *Metric) GetSuffix() string {
	if m != nil {
		return m.Suffix
	}
	return ""
}

func (m *Metric) GetStyle() Style {
	if m != nil {
		return m.Style
	}
	return Style_UNKNOWN_STYLE
}

func (m *Metric) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

// The request message
type GetMetricByPathRequest struct {
	Path                 string   `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetMetricByPathRequest) Reset()         { *m = GetMetricByPathRequest{} }
func (m *GetMetricByPathRequest) String() string { return proto.CompactTextString(m) }
func (*GetMetricByPathRequest) ProtoMessage()    {}
func (*GetMetricByPathRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_84366feb8d987a6e, []int{3}
}

func (m *GetMetricByPathRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetMetricByPathRequest.Unmarshal(m, b)
}
func (m *GetMetricByPathRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetMetricByPathRequest.Marshal(b, m, deterministic)
}
func (m *GetMetricByPathRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetMetricByPathRequest.Merge(m, src)
}
func (m *GetMetricByPathRequest) XXX_Size() int {
	return xxx_messageInfo_GetMetricByPathRequest.Size(m)
}
func (m *GetMetricByPathRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetMetricByPathRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetMetricByPathRequest proto.InternalMessageInfo

func (m *GetMetricByPathRequest) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

// The request message
type GetMetricRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetMetricRequest) Reset()         { *m = GetMetricRequest{} }
func (m *GetMetricRequest) String() string { return proto.CompactTextString(m) }
func (*GetMetricRequest) ProtoMessage()    {}
func (*GetMetricRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_84366feb8d987a6e, []int{4}
}

func (m *GetMetricRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetMetricRequest.Unmarshal(m, b)
}
func (m *GetMetricRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetMetricRequest.Marshal(b, m, deterministic)
}
func (m *GetMetricRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetMetricRequest.Merge(m, src)
}
func (m *GetMetricRequest) XXX_Size() int {
	return xxx_messageInfo_GetMetricRequest.Size(m)
}
func (m *GetMetricRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetMetricRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetMetricRequest proto.InternalMessageInfo

// The response message
type GetMetricResponse struct {
	// API versioning: it is my best practice to specify version explicitly
	Api string `protobuf:"bytes,1,opt,name=api,proto3" json:"api,omitempty"`
	// Task entity to add
	Metric               *Metric  `protobuf:"bytes,2,opt,name=metric,proto3" json:"metric,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetMetricResponse) Reset()         { *m = GetMetricResponse{} }
func (m *GetMetricResponse) String() string { return proto.CompactTextString(m) }
func (*GetMetricResponse) ProtoMessage()    {}
func (*GetMetricResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_84366feb8d987a6e, []int{5}
}

func (m *GetMetricResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetMetricResponse.Unmarshal(m, b)
}
func (m *GetMetricResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetMetricResponse.Marshal(b, m, deterministic)
}
func (m *GetMetricResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetMetricResponse.Merge(m, src)
}
func (m *GetMetricResponse) XXX_Size() int {
	return xxx_messageInfo_GetMetricResponse.Size(m)
}
func (m *GetMetricResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetMetricResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetMetricResponse proto.InternalMessageInfo

func (m *GetMetricResponse) GetApi() string {
	if m != nil {
		return m.Api
	}
	return ""
}

func (m *GetMetricResponse) GetMetric() *Metric {
	if m != nil {
		return m.Metric
	}
	return nil
}

func init() {
	proto.RegisterEnum("v1.Style", Style_name, Style_value)
	proto.RegisterEnum("v1.Kind", Kind_name, Kind_value)
	proto.RegisterType((*ShutdownRequest)(nil), "v1.ShutdownRequest")
	proto.RegisterType((*ShutdownResponse)(nil), "v1.ShutdownResponse")
	proto.RegisterType((*Metric)(nil), "v1.Metric")
	proto.RegisterType((*GetMetricByPathRequest)(nil), "v1.GetMetricByPathRequest")
	proto.RegisterType((*GetMetricRequest)(nil), "v1.GetMetricRequest")
	proto.RegisterType((*GetMetricResponse)(nil), "v1.GetMetricResponse")
}

func init() { proto.RegisterFile("blitzd.proto", fileDescriptor_84366feb8d987a6e) }

var fileDescriptor_84366feb8d987a6e = []byte{
	// 605 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x53, 0x5d, 0x4f, 0xdb, 0x30,
	0x14, 0x25, 0xfd, 0xa2, 0xbd, 0x6d, 0x21, 0xdc, 0x31, 0x64, 0x55, 0x93, 0xa8, 0xf2, 0x54, 0x21,
	0x54, 0x44, 0xa7, 0x3d, 0x4c, 0x9a, 0x34, 0x05, 0x9a, 0x95, 0x88, 0x12, 0x26, 0x27, 0x1b, 0x9a,
	0xf6, 0x80, 0x02, 0x71, 0x21, 0x5a, 0xda, 0x64, 0x89, 0xd3, 0xc1, 0xfe, 0xc6, 0x1e, 0xf7, 0x67,
	0x27, 0xdb, 0x49, 0xb7, 0xa2, 0x69, 0x7b, 0xf3, 0x39, 0x3e, 0xc7, 0xba, 0xf7, 0xdc, 0x6b, 0xe8,
	0xdc, 0x44, 0x21, 0xff, 0x1e, 0x0c, 0x93, 0x34, 0xe6, 0x31, 0x56, 0x96, 0xc7, 0xbd, 0xfd, 0xbb,
	0x38, 0xbe, 0x8b, 0xd8, 0x91, 0x64, 0x6e, 0xf2, 0xd9, 0x11, 0x0f, 0xe7, 0x2c, 0xe3, 0xfe, 0x3c,
	0x51, 0x22, 0x63, 0x07, 0xb6, 0xdd, 0xfb, 0x9c, 0x07, 0xf1, 0xb7, 0x05, 0x65, 0x5f, 0x73, 0x96,
	0x71, 0xe3, 0x10, 0xf4, 0xdf, 0x54, 0x96, 0xc4, 0x8b, 0x8c, 0x21, 0x81, 0xcd, 0x39, 0xcb, 0x32,
	0xff, 0x8e, 0x11, 0xad, 0xaf, 0x0d, 0x5a, 0xb4, 0x84, 0xc6, 0xcf, 0x2a, 0x34, 0x2e, 0x18, 0x4f,
	0xc3, 0x5b, 0x7c, 0x01, 0xb5, 0x2f, 0xe1, 0x22, 0x90, 0x8a, 0xad, 0x51, 0x73, 0xb8, 0x3c, 0x1e,
	0x9e, 0x87, 0x8b, 0x80, 0x4a, 0x16, 0xf7, 0xa0, 0x31, 0x8f, 0x83, 0x3c, 0x62, 0xa4, 0x22, 0x5f,
	0x28, 0x10, 0xee, 0x42, 0x9d, 0x87, 0x3c, 0x62, 0xa4, 0x2a, 0x69, 0x05, 0xb0, 0x07, 0xcd, 0x70,
	0xc1, 0x59, 0xba, 0xf4, 0x23, 0x52, 0xeb, 0x6b, 0x03, 0x8d, 0xae, 0xb0, 0x28, 0x46, 0xb4, 0x11,
	0xe7, 0x9c, 0xd4, 0xe5, 0x55, 0x09, 0xf1, 0x35, 0x40, 0x9e, 0x04, 0x3e, 0x67, 0xc1, 0xb5, 0xcf,
	0x49, 0xa3, 0xaf, 0x0d, 0xda, 0xa3, 0xde, 0x50, 0x65, 0x30, 0x2c, 0x33, 0x18, 0x7a, 0x65, 0x06,
	0xb4, 0x55, 0xa8, 0x4d, 0x8e, 0x6f, 0xa1, 0xcb, 0x1e, 0x92, 0x30, 0x15, 0xd6, 0x19, 0x67, 0x29,
	0xd9, 0xfc, 0xaf, 0xbb, 0x53, 0x18, 0x4c, 0xa1, 0x17, 0x55, 0x15, 0x98, 0x34, 0xfb, 0xda, 0xa0,
	0x49, 0x4b, 0x28, 0x3a, 0x5c, 0xfa, 0x51, 0xce, 0x48, 0x4b, 0x75, 0x28, 0x81, 0xc8, 0x23, 0x49,
	0xd9, 0x2c, 0x7c, 0x20, 0xa0, 0xf2, 0x50, 0x48, 0xf0, 0x59, 0x3e, 0x13, 0x7c, 0x5b, 0xf1, 0x0a,
	0xe1, 0x3e, 0xd4, 0x33, 0xfe, 0x18, 0x31, 0xd2, 0x91, 0xf1, 0xb6, 0x44, 0xbc, 0xae, 0x20, 0xa8,
	0xe2, 0x11, 0xa1, 0xc6, 0xd9, 0x03, 0x27, 0x5d, 0x69, 0x93, 0x67, 0xe3, 0x10, 0xf6, 0x26, 0x8c,
	0xab, 0xf9, 0x9c, 0x3c, 0xbe, 0xf7, 0xf9, 0x7d, 0x31, 0x65, 0xa1, 0x4e, 0x7c, 0x7e, 0x5f, 0x8c,
	0x53, 0x9e, 0x0d, 0x04, 0x7d, 0xa5, 0x2e, 0xb7, 0xc1, 0x86, 0x9d, 0x3f, 0xb8, 0x62, 0x1d, 0x74,
	0xa8, 0xfa, 0x49, 0x58, 0x78, 0xc5, 0x11, 0x0d, 0x68, 0xcc, 0xa5, 0x46, 0x4e, 0xb7, 0x3d, 0x02,
	0x51, 0x5e, 0xe1, 0x2a, 0x6e, 0x0e, 0x3e, 0x43, 0x5d, 0x16, 0x8c, 0x3b, 0xd0, 0xfd, 0xe0, 0x9c,
	0x3b, 0x97, 0x57, 0xce, 0xb5, 0xeb, 0x7d, 0x9a, 0x5a, 0xfa, 0x06, 0x02, 0x34, 0x9c, 0x4b, 0x7a,
	0x61, 0x4e, 0x75, 0x0d, 0xbb, 0xd0, 0x3a, 0xb3, 0x27, 0x67, 0x53, 0x7b, 0x72, 0xe6, 0xe9, 0x15,
	0x6c, 0xc3, 0xe6, 0x95, 0x49, 0x1d, 0xdb, 0x99, 0xe8, 0x55, 0x6c, 0x41, 0xdd, 0xa2, 0xf4, 0x92,
	0xea, 0x35, 0xec, 0x40, 0xf3, 0x94, 0xda, 0x9e, 0x7d, 0x6a, 0x4e, 0xf5, 0xfa, 0x81, 0x05, 0x35,
	0xb1, 0x6c, 0xa8, 0x43, 0xa7, 0x7c, 0xfb, 0xdc, 0x76, 0xc6, 0xea, 0x69, 0xd7, 0x33, 0x3d, 0xfb,
	0x54, 0xd7, 0x70, 0x0b, 0xc0, 0xb3, 0x2f, 0xac, 0xeb, 0x13, 0xd3, 0xb5, 0xc6, 0x7a, 0x05, 0xb7,
	0xa1, 0x6d, 0x7d, 0xb4, 0x1c, 0xaf, 0x20, 0xaa, 0x23, 0x0b, 0x9a, 0xe5, 0xf2, 0x8b, 0x6d, 0x1a,
	0xc7, 0x2b, 0xf4, 0x4c, 0x06, 0xbe, 0xfe, 0x57, 0x7a, 0xbb, 0xeb, 0xa4, 0x8a, 0xc7, 0xd8, 0x18,
	0xfd, 0xd0, 0xa0, 0xab, 0xba, 0x77, 0x59, 0xba, 0x0c, 0x6f, 0x19, 0x8e, 0x61, 0xfb, 0xc9, 0x24,
	0xb0, 0x27, 0xcc, 0x7f, 0x1f, 0x4f, 0xef, 0xf9, 0xda, 0xdd, 0x2a, 0xf8, 0x37, 0xd0, 0x5d, 0x91,
	0xef, 0xe2, 0xf8, 0x15, 0xee, 0x3e, 0xd1, 0xfd, 0xcb, 0x7d, 0xd3, 0x90, 0x4b, 0xfc, 0xf2, 0x57,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x02, 0x94, 0x63, 0xb6, 0x28, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ShutdownClient is the client API for Shutdown service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ShutdownClient interface {
	// Sends a greeting
	DoShutdown(ctx context.Context, in *ShutdownRequest, opts ...grpc.CallOption) (*ShutdownResponse, error)
}

type shutdownClient struct {
	cc *grpc.ClientConn
}

func NewShutdownClient(cc *grpc.ClientConn) ShutdownClient {
	return &shutdownClient{cc}
}

func (c *shutdownClient) DoShutdown(ctx context.Context, in *ShutdownRequest, opts ...grpc.CallOption) (*ShutdownResponse, error) {
	out := new(ShutdownResponse)
	err := c.cc.Invoke(ctx, "/v1.Shutdown/DoShutdown", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShutdownServer is the server API for Shutdown service.
type ShutdownServer interface {
	// Sends a greeting
	DoShutdown(context.Context, *ShutdownRequest) (*ShutdownResponse, error)
}

// UnimplementedShutdownServer can be embedded to have forward compatible implementations.
type UnimplementedShutdownServer struct {
}

func (*UnimplementedShutdownServer) DoShutdown(ctx context.Context, req *ShutdownRequest) (*ShutdownResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DoShutdown not implemented")
}

func RegisterShutdownServer(s *grpc.Server, srv ShutdownServer) {
	s.RegisterService(&_Shutdown_serviceDesc, srv)
}

func _Shutdown_DoShutdown_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShutdownRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShutdownServer).DoShutdown(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.Shutdown/DoShutdown",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShutdownServer).DoShutdown(ctx, req.(*ShutdownRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Shutdown_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v1.Shutdown",
	HandlerType: (*ShutdownServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DoShutdown",
			Handler:    _Shutdown_DoShutdown_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "blitzd.proto",
}

// MetricServiceClient is the client API for MetricService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MetricServiceClient interface {
	GetMetricByPath(ctx context.Context, in *GetMetricByPathRequest, opts ...grpc.CallOption) (*GetMetricResponse, error)
	GetMetricFoo5(ctx context.Context, in *GetMetricRequest, opts ...grpc.CallOption) (*GetMetricResponse, error)
}

type metricServiceClient struct {
	cc *grpc.ClientConn
}

func NewMetricServiceClient(cc *grpc.ClientConn) MetricServiceClient {
	return &metricServiceClient{cc}
}

func (c *metricServiceClient) GetMetricByPath(ctx context.Context, in *GetMetricByPathRequest, opts ...grpc.CallOption) (*GetMetricResponse, error) {
	out := new(GetMetricResponse)
	err := c.cc.Invoke(ctx, "/v1.MetricService/GetMetricByPath", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *metricServiceClient) GetMetricFoo5(ctx context.Context, in *GetMetricRequest, opts ...grpc.CallOption) (*GetMetricResponse, error) {
	out := new(GetMetricResponse)
	err := c.cc.Invoke(ctx, "/v1.MetricService/GetMetricFoo5", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MetricServiceServer is the server API for MetricService service.
type MetricServiceServer interface {
	GetMetricByPath(context.Context, *GetMetricByPathRequest) (*GetMetricResponse, error)
	GetMetricFoo5(context.Context, *GetMetricRequest) (*GetMetricResponse, error)
}

// UnimplementedMetricServiceServer can be embedded to have forward compatible implementations.
type UnimplementedMetricServiceServer struct {
}

func (*UnimplementedMetricServiceServer) GetMetricByPath(ctx context.Context, req *GetMetricByPathRequest) (*GetMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetricByPath not implemented")
}
func (*UnimplementedMetricServiceServer) GetMetricFoo5(ctx context.Context, req *GetMetricRequest) (*GetMetricResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetricFoo5 not implemented")
}

func RegisterMetricServiceServer(s *grpc.Server, srv MetricServiceServer) {
	s.RegisterService(&_MetricService_serviceDesc, srv)
}

func _MetricService_GetMetricByPath_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMetricByPathRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricServiceServer).GetMetricByPath(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.MetricService/GetMetricByPath",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricServiceServer).GetMetricByPath(ctx, req.(*GetMetricByPathRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MetricService_GetMetricFoo5_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MetricServiceServer).GetMetricFoo5(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.MetricService/GetMetricFoo5",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MetricServiceServer).GetMetricFoo5(ctx, req.(*GetMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _MetricService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v1.MetricService",
	HandlerType: (*MetricServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetMetricByPath",
			Handler:    _MetricService_GetMetricByPath_Handler,
		},
		{
			MethodName: "GetMetricFoo5",
			Handler:    _MetricService_GetMetricFoo5_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "blitzd.proto",
}
