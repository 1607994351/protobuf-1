// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/protobuf/struct.proto

package structpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	prototype "google.golang.org/protobuf/reflect/prototype"
	protoiface "google.golang.org/protobuf/runtime/protoiface"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	sync "sync"
)

const (
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 0)
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(0 - protoimpl.MinVersion)
)

// `NullValue` is a singleton enumeration to represent the null value for the
// `Value` type union.
//
//  The JSON representation for `NullValue` is JSON `null`.
type NullValue int32

const (
	// Null value.
	NullValue_NULL_VALUE NullValue = 0
)

var NullValue_name = map[int32]string{
	0: "NULL_VALUE",
}

var NullValue_value = map[string]int32{
	"NULL_VALUE": 0,
}

func (x NullValue) Enum() *NullValue {
	p := new(NullValue)
	*p = x
	return p
}

func (x NullValue) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NullValue) Descriptor() protoreflect.EnumDescriptor {
	return file_google_protobuf_struct_proto_enumTypes[0].EnumDescriptor
}

func (x NullValue) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NullValue.Type instead.
func (NullValue) EnumDescriptor() ([]byte, []int) {
	return file_google_protobuf_struct_proto_rawDescGZIP(), []int{0}
}

func (NullValue) XXX_WellKnownType() string { return "NullValue" }

// `Struct` represents a structured data value, consisting of fields
// which map to dynamically typed values. In some languages, `Struct`
// might be supported by a native representation. For example, in
// scripting languages like JS a struct is represented as an
// object. The details of that representation are described together
// with the proto support for the language.
//
// The JSON representation for `Struct` is JSON object.
type Struct struct {
	// Unordered map of dynamically typed values.
	Fields               map[string]*Value       `protobuf:"bytes,1,rep,name=fields,proto3" json:"fields,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     protoimpl.UnknownFields `json:"-"`
	XXX_sizecache        protoimpl.SizeCache     `json:"-"`
}

func (x *Struct) Reset() {
	*x = Struct{}
}

func (x *Struct) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Struct) ProtoMessage() {}

func (x *Struct) ProtoReflect() protoreflect.Message {
	return file_google_protobuf_struct_proto_msgTypes[0].MessageOf(x)
}

func (m *Struct) XXX_Methods() *protoiface.Methods {
	return file_google_protobuf_struct_proto_msgTypes[0].Methods()
}

// Deprecated: Use Struct.ProtoReflect.Type instead.
func (*Struct) Descriptor() ([]byte, []int) {
	return file_google_protobuf_struct_proto_rawDescGZIP(), []int{0}
}

func (*Struct) XXX_WellKnownType() string { return "Struct" }

func (x *Struct) GetFields() map[string]*Value {
	if x != nil {
		return x.Fields
	}
	return nil
}

// `Value` represents a dynamically typed value which can be either
// null, a number, a string, a boolean, a recursive struct value, or a
// list of values. A producer of value is expected to set one of that
// variants, absence of any variant indicates an error.
//
// The JSON representation for `Value` is JSON value.
type Value struct {
	// The kind of value.
	//
	// Types that are valid to be assigned to Kind:
	// Represents a null value.
	//	*Value_NullValue
	// Represents a double value.
	//	*Value_NumberValue
	// Represents a string value.
	//	*Value_StringValue
	// Represents a boolean value.
	//	*Value_BoolValue
	// Represents a structured value.
	//	*Value_StructValue
	// Represents a repeated `Value`.
	//	*Value_ListValue
	Kind                 isValue_Kind            `protobuf_oneof:"kind"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     protoimpl.UnknownFields `json:"-"`
	XXX_sizecache        protoimpl.SizeCache     `json:"-"`
}

func (x *Value) Reset() {
	*x = Value{}
}

func (x *Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Value) ProtoMessage() {}

func (x *Value) ProtoReflect() protoreflect.Message {
	return file_google_protobuf_struct_proto_msgTypes[1].MessageOf(x)
}

func (m *Value) XXX_Methods() *protoiface.Methods {
	return file_google_protobuf_struct_proto_msgTypes[1].Methods()
}

// Deprecated: Use Value.ProtoReflect.Type instead.
func (*Value) Descriptor() ([]byte, []int) {
	return file_google_protobuf_struct_proto_rawDescGZIP(), []int{1}
}

func (*Value) XXX_WellKnownType() string { return "Value" }

func (m *Value) GetKind() isValue_Kind {
	if m != nil {
		return m.Kind
	}
	return nil
}

func (x *Value) GetNullValue() NullValue {
	if x, ok := x.GetKind().(*Value_NullValue); ok {
		return x.NullValue
	}
	return NullValue_NULL_VALUE
}

func (x *Value) GetNumberValue() float64 {
	if x, ok := x.GetKind().(*Value_NumberValue); ok {
		return x.NumberValue
	}
	return 0
}

func (x *Value) GetStringValue() string {
	if x, ok := x.GetKind().(*Value_StringValue); ok {
		return x.StringValue
	}
	return ""
}

func (x *Value) GetBoolValue() bool {
	if x, ok := x.GetKind().(*Value_BoolValue); ok {
		return x.BoolValue
	}
	return false
}

func (x *Value) GetStructValue() *Struct {
	if x, ok := x.GetKind().(*Value_StructValue); ok {
		return x.StructValue
	}
	return nil
}

func (x *Value) GetListValue() *ListValue {
	if x, ok := x.GetKind().(*Value_ListValue); ok {
		return x.ListValue
	}
	return nil
}

type isValue_Kind interface {
	isValue_Kind()
}

type Value_NullValue struct {
	NullValue NullValue `protobuf:"varint,1,opt,name=null_value,json=nullValue,proto3,enum=google.protobuf.NullValue,oneof"`
}

type Value_NumberValue struct {
	NumberValue float64 `protobuf:"fixed64,2,opt,name=number_value,json=numberValue,proto3,oneof"`
}

type Value_StringValue struct {
	StringValue string `protobuf:"bytes,3,opt,name=string_value,json=stringValue,proto3,oneof"`
}

type Value_BoolValue struct {
	BoolValue bool `protobuf:"varint,4,opt,name=bool_value,json=boolValue,proto3,oneof"`
}

type Value_StructValue struct {
	StructValue *Struct `protobuf:"bytes,5,opt,name=struct_value,json=structValue,proto3,oneof"`
}

type Value_ListValue struct {
	ListValue *ListValue `protobuf:"bytes,6,opt,name=list_value,json=listValue,proto3,oneof"`
}

func (*Value_NullValue) isValue_Kind() {}

func (*Value_NumberValue) isValue_Kind() {}

func (*Value_StringValue) isValue_Kind() {}

func (*Value_BoolValue) isValue_Kind() {}

func (*Value_StructValue) isValue_Kind() {}

func (*Value_ListValue) isValue_Kind() {}

// `ListValue` is a wrapper around a repeated field of values.
//
// The JSON representation for `ListValue` is JSON array.
type ListValue struct {
	// Repeated field of dynamically typed values.
	Values               []*Value                `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     protoimpl.UnknownFields `json:"-"`
	XXX_sizecache        protoimpl.SizeCache     `json:"-"`
}

func (x *ListValue) Reset() {
	*x = ListValue{}
}

func (x *ListValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListValue) ProtoMessage() {}

func (x *ListValue) ProtoReflect() protoreflect.Message {
	return file_google_protobuf_struct_proto_msgTypes[2].MessageOf(x)
}

func (m *ListValue) XXX_Methods() *protoiface.Methods {
	return file_google_protobuf_struct_proto_msgTypes[2].Methods()
}

// Deprecated: Use ListValue.ProtoReflect.Type instead.
func (*ListValue) Descriptor() ([]byte, []int) {
	return file_google_protobuf_struct_proto_rawDescGZIP(), []int{2}
}

func (*ListValue) XXX_WellKnownType() string { return "ListValue" }

func (x *ListValue) GetValues() []*Value {
	if x != nil {
		return x.Values
	}
	return nil
}

var File_google_protobuf_struct_proto protoreflect.FileDescriptor

var file_google_protobuf_struct_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x22,
	0x98, 0x01, 0x0a, 0x06, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x12, 0x3b, 0x0a, 0x06, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72,
	0x75, 0x63, 0x74, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x1a, 0x51, 0x0a, 0x0b, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2c, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xb2, 0x02, 0x0a, 0x05, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x3b, 0x0a, 0x0a, 0x6e, 0x75, 0x6c, 0x6c, 0x5f, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4e, 0x75, 0x6c, 0x6c, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x48, 0x00, 0x52, 0x09, 0x6e, 0x75, 0x6c, 0x6c, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x23, 0x0a, 0x0c, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x5f, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x48, 0x00, 0x52, 0x0b, 0x6e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x23, 0x0a, 0x0c, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0b,
	0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1f, 0x0a, 0x0a, 0x62,
	0x6f, 0x6f, 0x6c, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x48,
	0x00, 0x52, 0x09, 0x62, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x3c, 0x0a, 0x0c,
	0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x48, 0x00, 0x52, 0x0b, 0x73,
	0x74, 0x72, 0x75, 0x63, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x3b, 0x0a, 0x0a, 0x6c, 0x69,
	0x73, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x48, 0x00, 0x52, 0x09, 0x6c, 0x69,
	0x73, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x06, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x22,
	0x3b, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x2e, 0x0a, 0x06,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x2a, 0x1b, 0x0a, 0x09,
	0x4e, 0x75, 0x6c, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x0e, 0x0a, 0x0a, 0x4e, 0x55, 0x4c,
	0x4c, 0x5f, 0x56, 0x41, 0x4c, 0x55, 0x45, 0x10, 0x00, 0x42, 0x7f, 0x0a, 0x13, 0x63, 0x6f, 0x6d,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x42, 0x0b, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2e, 0x6f,
	0x72, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x2f, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x70, 0x62,
	0xf8, 0x01, 0x01, 0xa2, 0x02, 0x03, 0x47, 0x50, 0x42, 0xaa, 0x02, 0x1e, 0x47, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x57, 0x65, 0x6c, 0x6c,
	0x4b, 0x6e, 0x6f, 0x77, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_google_protobuf_struct_proto_rawDescOnce sync.Once
	file_google_protobuf_struct_proto_rawDescData = file_google_protobuf_struct_proto_rawDesc
)

func file_google_protobuf_struct_proto_rawDescGZIP() []byte {
	file_google_protobuf_struct_proto_rawDescOnce.Do(func() {
		file_google_protobuf_struct_proto_rawDescData = protoimpl.X.CompressGZIP(file_google_protobuf_struct_proto_rawDescData)
	})
	return file_google_protobuf_struct_proto_rawDescData
}

var file_google_protobuf_struct_proto_enumTypes = make([]prototype.Enum, 1)
var file_google_protobuf_struct_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_google_protobuf_struct_proto_goTypes = []interface{}{
	(NullValue)(0),    // 0: google.protobuf.NullValue
	(*Struct)(nil),    // 1: google.protobuf.Struct
	(*Value)(nil),     // 2: google.protobuf.Value
	(*ListValue)(nil), // 3: google.protobuf.ListValue
	nil,               // 4: google.protobuf.Struct.FieldsEntry
}
var file_google_protobuf_struct_proto_depIdxs = []int32{
	4, // google.protobuf.Struct.fields:type_name -> google.protobuf.Struct.FieldsEntry
	0, // google.protobuf.Value.null_value:type_name -> google.protobuf.NullValue
	1, // google.protobuf.Value.struct_value:type_name -> google.protobuf.Struct
	3, // google.protobuf.Value.list_value:type_name -> google.protobuf.ListValue
	2, // google.protobuf.ListValue.values:type_name -> google.protobuf.Value
	2, // google.protobuf.Struct.FieldsEntry.value:type_name -> google.protobuf.Value
	6, // starting offset of method output_type sub-list
	6, // starting offset of method input_type sub-list
	6, // starting offset of extension type_name sub-list
	6, // starting offset of extension extendee sub-list
	0, // starting offset of field type_name sub-list
}

func init() { file_google_protobuf_struct_proto_init() }
func file_google_protobuf_struct_proto_init() {
	if File_google_protobuf_struct_proto != nil {
		return
	}
	file_google_protobuf_struct_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*Value_NullValue)(nil),
		(*Value_NumberValue)(nil),
		(*Value_StringValue)(nil),
		(*Value_BoolValue)(nil),
		(*Value_StructValue)(nil),
		(*Value_ListValue)(nil),
	}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			RawDescriptor: file_google_protobuf_struct_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_google_protobuf_struct_proto_goTypes,
		DependencyIndexes: file_google_protobuf_struct_proto_depIdxs,
		MessageInfos:      file_google_protobuf_struct_proto_msgTypes,
	}.Build()
	File_google_protobuf_struct_proto = out.File
	file_google_protobuf_struct_proto_enumTypes = out.Enums
	file_google_protobuf_struct_proto_rawDesc = nil
	file_google_protobuf_struct_proto_goTypes = nil
	file_google_protobuf_struct_proto_depIdxs = nil
}
