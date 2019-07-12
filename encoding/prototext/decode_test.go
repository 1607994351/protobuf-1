// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prototext_test

import (
	"math"
	"testing"

	"google.golang.org/protobuf/encoding/prototext"
	pimpl "google.golang.org/protobuf/internal/impl"
	"google.golang.org/protobuf/proto"
	preg "google.golang.org/protobuf/reflect/protoregistry"

	"google.golang.org/protobuf/encoding/testprotos/pb2"
	"google.golang.org/protobuf/encoding/testprotos/pb3"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		desc         string
		umo          prototext.UnmarshalOptions
		inputMessage proto.Message
		inputText    string
		wantMessage  proto.Message
		wantErr      bool // TODO: Verify error message content.
	}{{
		desc:         "proto2 empty message",
		inputMessage: &pb2.Scalars{},
		wantMessage:  &pb2.Scalars{},
	}, {
		desc:         "proto2 optional scalars set to zero values",
		inputMessage: &pb2.Scalars{},
		inputText: `opt_bool: false
opt_int32: 0
opt_int64: 0
opt_uint32: 0
opt_uint64: 0
opt_sint32: 0
opt_sint64: 0
opt_fixed32: 0
opt_fixed64: 0
opt_sfixed32: 0
opt_sfixed64: 0
opt_float: 0
opt_double: 0
opt_bytes: ""
opt_string: ""
`,
		wantMessage: &pb2.Scalars{
			OptBool:     proto.Bool(false),
			OptInt32:    proto.Int32(0),
			OptInt64:    proto.Int64(0),
			OptUint32:   proto.Uint32(0),
			OptUint64:   proto.Uint64(0),
			OptSint32:   proto.Int32(0),
			OptSint64:   proto.Int64(0),
			OptFixed32:  proto.Uint32(0),
			OptFixed64:  proto.Uint64(0),
			OptSfixed32: proto.Int32(0),
			OptSfixed64: proto.Int64(0),
			OptFloat:    proto.Float32(0),
			OptDouble:   proto.Float64(0),
			OptBytes:    []byte{},
			OptString:   proto.String(""),
		},
	}, {
		desc:         "proto3 scalars set to zero values",
		inputMessage: &pb3.Scalars{},
		inputText: `s_bool: false
s_int32: 0
s_int64: 0
s_uint32: 0
s_uint64: 0
s_sint32: 0
s_sint64: 0
s_fixed32: 0
s_fixed64: 0
s_sfixed32: 0
s_sfixed64: 0
s_float: 0
s_double: 0
s_bytes: ""
s_string: ""
`,
		wantMessage: &pb3.Scalars{},
	}, {
		desc:         "proto2 optional scalars",
		inputMessage: &pb2.Scalars{},
		inputText: `opt_bool: true
opt_int32: 255
opt_int64: 3735928559
opt_uint32: 0xff
opt_uint64: 0xdeadbeef
opt_sint32: -1001
opt_sint64: -0xffff
opt_fixed64: 64
opt_sfixed32: -32
opt_float: 1.234
opt_double: 1.23e+100
opt_bytes: "\xe8\xb0\xb7\xe6\xad\x8c"
opt_string: "谷歌"
`,
		wantMessage: &pb2.Scalars{
			OptBool:     proto.Bool(true),
			OptInt32:    proto.Int32(0xff),
			OptInt64:    proto.Int64(0xdeadbeef),
			OptUint32:   proto.Uint32(0xff),
			OptUint64:   proto.Uint64(0xdeadbeef),
			OptSint32:   proto.Int32(-1001),
			OptSint64:   proto.Int64(-0xffff),
			OptFixed64:  proto.Uint64(64),
			OptSfixed32: proto.Int32(-32),
			OptFloat:    proto.Float32(1.234),
			OptDouble:   proto.Float64(1.23e100),
			OptBytes:    []byte("\xe8\xb0\xb7\xe6\xad\x8c"),
			OptString:   proto.String("谷歌"),
		},
	}, {
		desc:         "case sensitive",
		inputMessage: &pb3.Scalars{},
		inputText:    `S_BOOL: true`,
		wantErr:      true,
	}, {
		desc:         "proto3 scalars",
		inputMessage: &pb3.Scalars{},
		inputText: `s_bool: true
s_int32: 255
s_int64: 3735928559
s_uint32: 0xff
s_uint64: 0xdeadbeef
s_sint32: -1001
s_sint64: -0xffff
s_fixed64: 64
s_sfixed32: -32
s_float: 1.234
s_double: 1.23e+100
s_bytes: "\xe8\xb0\xb7\xe6\xad\x8c"
s_string: "谷歌"
`,
		wantMessage: &pb3.Scalars{
			SBool:     true,
			SInt32:    0xff,
			SInt64:    0xdeadbeef,
			SUint32:   0xff,
			SUint64:   0xdeadbeef,
			SSint32:   -1001,
			SSint64:   -0xffff,
			SFixed64:  64,
			SSfixed32: -32,
			SFloat:    1.234,
			SDouble:   1.23e100,
			SBytes:    []byte("\xe8\xb0\xb7\xe6\xad\x8c"),
			SString:   "谷歌",
		},
	}, {
		desc:         "string with invalid UTF-8",
		inputMessage: &pb3.Scalars{},
		inputText:    `s_string: "abc\xff"`,
		wantErr:      true,
	}, {
		desc:         "proto2 message contains unknown field",
		inputMessage: &pb2.Scalars{},
		inputText:    "unknown_field: 123",
		wantErr:      true,
	}, {
		desc:         "proto3 message contains unknown field",
		inputMessage: &pb3.Scalars{},
		inputText:    "unknown_field: 456",
		wantErr:      true,
	}, {
		desc:         "proto2 numeric key field",
		inputMessage: &pb2.Scalars{},
		inputText:    "1: true",
		wantErr:      true,
	}, {
		desc:         "proto3 numeric key field",
		inputMessage: &pb3.Scalars{},
		inputText:    "1: true",
		wantErr:      true,
	}, {
		desc:         "invalid bool value",
		inputMessage: &pb3.Scalars{},
		inputText:    "s_bool: 123",
		wantErr:      true,
	}, {
		desc:         "invalid int32 value",
		inputMessage: &pb3.Scalars{},
		inputText:    "s_int32: not_a_num",
		wantErr:      true,
	}, {
		desc:         "invalid int64 value",
		inputMessage: &pb3.Scalars{},
		inputText:    "s_int64: 'not a num either'",
		wantErr:      true,
	}, {
		desc:         "invalid uint32 value",
		inputMessage: &pb3.Scalars{},
		inputText:    "s_fixed32: -42",
		wantErr:      true,
	}, {
		desc:         "invalid uint64 value",
		inputMessage: &pb3.Scalars{},
		inputText:    "s_uint64: -47",
		wantErr:      true,
	}, {
		desc:         "invalid sint32 value",
		inputMessage: &pb3.Scalars{},
		inputText:    "s_sint32: '42'",
		wantErr:      true,
	}, {
		desc:         "invalid sint64 value",
		inputMessage: &pb3.Scalars{},
		inputText:    "s_sint64: '-47'",
		wantErr:      true,
	}, {
		desc:         "invalid fixed32 value",
		inputMessage: &pb3.Scalars{},
		inputText:    "s_fixed32: -42",
		wantErr:      true,
	}, {
		desc:         "invalid fixed64 value",
		inputMessage: &pb3.Scalars{},
		inputText:    "s_fixed64: -42",
		wantErr:      true,
	}, {
		desc:         "invalid sfixed32 value",
		inputMessage: &pb3.Scalars{},
		inputText:    "s_sfixed32: 'not valid'",
		wantErr:      true,
	}, {
		desc:         "invalid sfixed64 value",
		inputMessage: &pb3.Scalars{},
		inputText:    "s_sfixed64: bad",
		wantErr:      true,
	}, {
		desc:         "float positive infinity",
		inputMessage: &pb3.Scalars{},
		inputText:    "s_float: inf",
		wantMessage: &pb3.Scalars{
			SFloat: float32(math.Inf(1)),
		},
	}, {
		desc:         "float negative infinity",
		inputMessage: &pb3.Scalars{},
		inputText:    "s_float: -inf",
		wantMessage: &pb3.Scalars{
			SFloat: float32(math.Inf(-1)),
		},
	}, {
		desc:         "double positive infinity",
		inputMessage: &pb3.Scalars{},
		inputText:    "s_double: inf",
		wantMessage: &pb3.Scalars{
			SDouble: math.Inf(1),
		},
	}, {
		desc:         "double negative infinity",
		inputMessage: &pb3.Scalars{},
		inputText:    "s_double: -inf",
		wantMessage: &pb3.Scalars{
			SDouble: math.Inf(-1),
		},
	}, {
		desc:         "invalid string value",
		inputMessage: &pb3.Scalars{},
		inputText:    "s_string: invalid_string",
		wantErr:      true,
	}, {
		desc:         "proto2 bytes set to empty string",
		inputMessage: &pb2.Scalars{},
		inputText:    "opt_bytes: ''",
		wantMessage: &pb2.Scalars{
			OptBytes: []byte(""),
		},
	}, {
		desc:         "proto3 bytes set to empty string",
		inputMessage: &pb3.Scalars{},
		inputText:    "s_bytes: ''",
		wantMessage:  &pb3.Scalars{},
	}, {
		desc:         "proto2 duplicate singular field",
		inputMessage: &pb2.Scalars{},
		inputText: `
opt_bool: true
opt_bool: false
`,
		wantErr: true,
	}, {
		desc:         "proto2 more duplicate singular field",
		inputMessage: &pb2.Scalars{},
		inputText: `
opt_bool: true
opt_string: "hello"
opt_bool: false
`,
		wantErr: true,
	}, {
		desc:         "proto2 invalid singular field",
		inputMessage: &pb2.Scalars{},
		inputText: `
opt_bool: [true, false]
`,
		wantErr: true,
	}, {
		desc:         "proto3 duplicate singular field",
		inputMessage: &pb3.Scalars{},
		inputText: `
s_bool: false
s_bool: true
`,
		wantErr: true,
	}, {
		desc:         "proto3 more duplicate singular field",
		inputMessage: &pb3.Scalars{},
		inputText: `
s_bool: false
s_string: ""
s_bool: true
`,
		wantErr: true,
	}, {
		desc:         "proto2 enum",
		inputMessage: &pb2.Enums{},
		inputText: `
opt_enum: ONE
opt_nested_enum: UNO
`,
		wantMessage: &pb2.Enums{
			OptEnum:       pb2.Enum_ONE.Enum(),
			OptNestedEnum: pb2.Enums_UNO.Enum(),
		},
	}, {
		desc:         "proto2 enum set to numeric values",
		inputMessage: &pb2.Enums{},
		inputText: `
opt_enum: 2
opt_nested_enum: 2
`,
		wantMessage: &pb2.Enums{
			OptEnum:       pb2.Enum_TWO.Enum(),
			OptNestedEnum: pb2.Enums_DOS.Enum(),
		},
	}, {
		desc:         "proto2 enum set to unnamed numeric values",
		inputMessage: &pb2.Enums{},
		inputText: `
opt_enum: 101
opt_nested_enum: -101
`,
		wantMessage: &pb2.Enums{
			OptEnum:       pb2.Enum(101).Enum(),
			OptNestedEnum: pb2.Enums_NestedEnum(-101).Enum(),
		},
	}, {
		desc:         "proto2 enum set to invalid named",
		inputMessage: &pb2.Enums{},
		inputText: `
opt_enum: UNNAMED
opt_nested_enum: UNNAMED_TOO
`,
		wantErr: true,
	}, {
		desc:         "proto3 enum name value",
		inputMessage: &pb3.Enums{},
		inputText: `
s_enum: ONE
s_nested_enum: DIEZ
`,
		wantMessage: &pb3.Enums{
			SEnum:       pb3.Enum_ONE,
			SNestedEnum: pb3.Enums_DIEZ,
		},
	}, {
		desc:         "proto3 enum numeric value",
		inputMessage: &pb3.Enums{},
		inputText: `
s_enum: 2
s_nested_enum: 2
`,
		wantMessage: &pb3.Enums{
			SEnum:       pb3.Enum_TWO,
			SNestedEnum: pb3.Enums_DOS,
		},
	}, {
		desc:         "proto3 enum unnamed numeric value",
		inputMessage: &pb3.Enums{},
		inputText: `
s_enum: 0x7fffffff
s_nested_enum: -0x80000000
`,
		wantMessage: &pb3.Enums{
			SEnum:       0x7fffffff,
			SNestedEnum: -0x80000000,
		},
	}, {
		desc:         "proto2 nested empty messages",
		inputMessage: &pb2.Nests{},
		inputText: `
opt_nested: {}
OptGroup: {}
`,
		wantMessage: &pb2.Nests{
			OptNested: &pb2.Nested{},
			Optgroup:  &pb2.Nests_OptGroup{},
		},
	}, {
		desc:         "group field name",
		inputMessage: &pb2.Nests{},
		inputText:    `optgroup: {}`,
		wantErr:      true,
	}, {
		desc:         "proto2 nested messages",
		inputMessage: &pb2.Nests{},
		inputText: `
opt_nested: {
  opt_string: "nested message"
  opt_nested: {
    opt_string: "another nested message"
  }
}
`,
		wantMessage: &pb2.Nests{
			OptNested: &pb2.Nested{
				OptString: proto.String("nested message"),
				OptNested: &pb2.Nested{
					OptString: proto.String("another nested message"),
				},
			},
		},
	}, {
		desc:         "proto3 nested empty message",
		inputMessage: &pb3.Nests{},
		inputText:    "s_nested: {}",
		wantMessage: &pb3.Nests{
			SNested: &pb3.Nested{},
		},
	}, {
		desc:         "proto3 nested message",
		inputMessage: &pb3.Nests{},
		inputText: `
s_nested: {
  s_string: "nested message"
  s_nested: {
    s_string: "another nested message"
  }
}
`,
		wantMessage: &pb3.Nests{
			SNested: &pb3.Nested{
				SString: "nested message",
				SNested: &pb3.Nested{
					SString: "another nested message",
				},
			},
		},
	}, {
		desc:         "proto3 nested message contains invalid UTF-8",
		inputMessage: &pb3.Nests{},
		inputText: `s_nested: {
  s_string: "abc\xff"
}
`,
		wantErr: true,
	}, {
		desc:         "oneof set to empty string",
		inputMessage: &pb3.Oneofs{},
		inputText:    "oneof_string: ''",
		wantMessage: &pb3.Oneofs{
			Union: &pb3.Oneofs_OneofString{},
		},
	}, {
		desc:         "oneof set to string",
		inputMessage: &pb3.Oneofs{},
		inputText:    "oneof_string: 'hello'",
		wantMessage: &pb3.Oneofs{
			Union: &pb3.Oneofs_OneofString{
				OneofString: "hello",
			},
		},
	}, {
		desc:         "oneof set to enum",
		inputMessage: &pb3.Oneofs{},
		inputText:    "oneof_enum: TEN",
		wantMessage: &pb3.Oneofs{
			Union: &pb3.Oneofs_OneofEnum{
				OneofEnum: pb3.Enum_TEN,
			},
		},
	}, {
		desc:         "oneof set to empty message",
		inputMessage: &pb3.Oneofs{},
		inputText:    "oneof_nested: {}",
		wantMessage: &pb3.Oneofs{
			Union: &pb3.Oneofs_OneofNested{
				OneofNested: &pb3.Nested{},
			},
		},
	}, {
		desc:         "oneof set to message",
		inputMessage: &pb3.Oneofs{},
		inputText: `
oneof_nested: {
  s_string: "nested message"
}
`,
		wantMessage: &pb3.Oneofs{
			Union: &pb3.Oneofs_OneofNested{
				OneofNested: &pb3.Nested{
					SString: "nested message",
				},
			},
		},
	}, {
		desc:         "oneof set to more than one field",
		inputMessage: &pb3.Oneofs{},
		inputText: `
oneof_enum: ZERO
oneof_string: "hello"
`,
		wantErr: true,
	}, {
		desc:         "repeated scalar using same field name",
		inputMessage: &pb2.Repeats{},
		inputText: `
rpt_string: "a"
rpt_string: "b"
rpt_int32: 0xff
rpt_float: 1.23
rpt_bytes: "bytes"
`,
		wantMessage: &pb2.Repeats{
			RptString: []string{"a", "b"},
			RptInt32:  []int32{0xff},
			RptFloat:  []float32{1.23},
			RptBytes:  [][]byte{[]byte("bytes")},
		},
	}, {
		desc:         "repeated using mix of [] and repeated field name",
		inputMessage: &pb2.Repeats{},
		inputText: `
rpt_string: "a"
rpt_bool: true
rpt_string: ["x", "y"]
rpt_bool: [ false, true ]
rpt_string: "b"
`,
		wantMessage: &pb2.Repeats{
			RptString: []string{"a", "x", "y", "b"},
			RptBool:   []bool{true, false, true},
		},
	}, {
		desc:         "repeated contains invalid UTF-8",
		inputMessage: &pb2.Repeats{},
		inputText:    `rpt_string: "abc\xff"`,
		wantErr:      true,
	}, {
		desc:         "repeated enums",
		inputMessage: &pb2.Enums{},
		inputText: `
rpt_enum: TEN
rpt_enum: 1
rpt_nested_enum: [DOS, 2]
rpt_enum: 42
rpt_nested_enum: -47
`,
		wantMessage: &pb2.Enums{
			RptEnum:       []pb2.Enum{pb2.Enum_TEN, pb2.Enum_ONE, 42},
			RptNestedEnum: []pb2.Enums_NestedEnum{pb2.Enums_DOS, pb2.Enums_DOS, -47},
		},
	}, {
		desc:         "repeated nested messages",
		inputMessage: &pb2.Nests{},
		inputText: `
rpt_nested: {
  opt_string: "repeat nested one"
}
rpt_nested: {
  opt_string: "repeat nested two"
  opt_nested: {
    opt_string: "inside repeat nested two"
  }
}
rpt_nested: {}
`,
		wantMessage: &pb2.Nests{
			RptNested: []*pb2.Nested{
				{
					OptString: proto.String("repeat nested one"),
				},
				{
					OptString: proto.String("repeat nested two"),
					OptNested: &pb2.Nested{
						OptString: proto.String("inside repeat nested two"),
					},
				},
				{},
			},
		},
	}, {
		desc:         "repeated group fields",
		inputMessage: &pb2.Nests{},
		inputText: `
RptGroup: {
  rpt_string: "hello"
  rpt_string: "world"
}
RptGroup: {}
`,
		wantMessage: &pb2.Nests{
			Rptgroup: []*pb2.Nests_RptGroup{
				{
					RptString: []string{"hello", "world"},
				},
				{},
			},
		},
	}, {
		desc:         "map fields 1",
		inputMessage: &pb3.Maps{},
		inputText: `
int32_to_str: {
  key: -101
  value: "-101"
}
int32_to_str: {
  key: 0
  value: "zero"
}
bool_to_uint32: {
  key: false
  value: 101
}
int32_to_str: {
  key: 255
  value: "0xff"
}
bool_to_uint32: {
  key: true
  value: 42
}
`,
		wantMessage: &pb3.Maps{
			Int32ToStr: map[int32]string{
				-101: "-101",
				0xff: "0xff",
				0:    "zero",
			},
			BoolToUint32: map[bool]uint32{
				true:  42,
				false: 101,
			},
		},
	}, {
		desc:         "map fields 2",
		inputMessage: &pb3.Maps{},
		inputText: `
uint64_to_enum: {
  key: 1
  value: ONE
}
uint64_to_enum: {
  key: 2
  value: 2
}
uint64_to_enum: {
  key: 10
  value: 101
}
`,
		wantMessage: &pb3.Maps{
			Uint64ToEnum: map[uint64]pb3.Enum{
				1:  pb3.Enum_ONE,
				2:  pb3.Enum_TWO,
				10: 101,
			},
		},
	}, {
		desc:         "map fields 3",
		inputMessage: &pb3.Maps{},
		inputText: `
str_to_nested: {
  key: "nested_one"
  value: {
    s_string: "nested in a map"
  }
}
`,
		wantMessage: &pb3.Maps{
			StrToNested: map[string]*pb3.Nested{
				"nested_one": &pb3.Nested{
					SString: "nested in a map",
				},
			},
		},
	}, {
		desc:         "map fields 4",
		inputMessage: &pb3.Maps{},
		inputText: `
str_to_oneofs: {
  key: "nested"
  value: {
    oneof_nested: {
      s_string: "nested oneof in map field value"
    }
  }
}
str_to_oneofs: {
  key: "string"
  value: {
    oneof_string: "hello"
  }
}
`,
		wantMessage: &pb3.Maps{
			StrToOneofs: map[string]*pb3.Oneofs{
				"string": &pb3.Oneofs{
					Union: &pb3.Oneofs_OneofString{
						OneofString: "hello",
					},
				},
				"nested": &pb3.Oneofs{
					Union: &pb3.Oneofs_OneofNested{
						OneofNested: &pb3.Nested{
							SString: "nested oneof in map field value",
						},
					},
				},
			},
		},
	}, {
		desc:         "map contains duplicate keys",
		inputMessage: &pb3.Maps{},
		inputText: `
int32_to_str: {
  key: 0
  value: "cero"
}
int32_to_str: {
  key: 0
  value: "zero"
}
`,
		wantMessage: &pb3.Maps{
			Int32ToStr: map[int32]string{
				0: "zero",
			},
		},
	}, {
		desc:         "map contains duplicate key fields",
		inputMessage: &pb3.Maps{},
		inputText: `
int32_to_str: {
  key: 0
  key: 1
  value: "cero"
}
`,
		wantErr: true,
	}, {
		desc:         "map contains duplicate value fields",
		inputMessage: &pb3.Maps{},
		inputText: `
int32_to_str: {
  key: 1
  value: "cero"
  value: "uno"
}
`,
		wantErr: true,
	}, {
		desc:         "map contains missing key",
		inputMessage: &pb3.Maps{},
		inputText: `
int32_to_str: {
  value: "zero"
}
bool_to_uint32: {
  value: 47
}
str_to_nested: {
  value: {}
}
`,
		wantMessage: &pb3.Maps{
			Int32ToStr: map[int32]string{
				0: "zero",
			},
			BoolToUint32: map[bool]uint32{
				false: 47,
			},
			StrToNested: map[string]*pb3.Nested{
				"": {},
			},
		},
	}, {
		desc:         "map contains missing value",
		inputMessage: &pb3.Maps{},
		inputText: `
int32_to_str: {
  key: 100
}
bool_to_uint32: {
  key: true
}
uint64_to_enum: {
  key: 101
}
str_to_nested: {
  key: "hello"
}
`,
		wantMessage: &pb3.Maps{
			Int32ToStr: map[int32]string{
				100: "",
			},
			BoolToUint32: map[bool]uint32{
				true: 0,
			},
			Uint64ToEnum: map[uint64]pb3.Enum{
				101: pb3.Enum_ZERO,
			},
			StrToNested: map[string]*pb3.Nested{
				"hello": {},
			},
		},
	}, {
		desc:         "map contains missing key and value",
		inputMessage: &pb3.Maps{},
		inputText: `
int32_to_str: {}
bool_to_uint32: {}
uint64_to_enum: {}
str_to_nested: {}
`,
		wantMessage: &pb3.Maps{
			Int32ToStr: map[int32]string{
				0: "",
			},
			BoolToUint32: map[bool]uint32{
				false: 0,
			},
			Uint64ToEnum: map[uint64]pb3.Enum{
				0: pb3.Enum_ZERO,
			},
			StrToNested: map[string]*pb3.Nested{
				"": {},
			},
		},
	}, {
		desc:         "map contains overriding entries",
		inputMessage: &pb3.Maps{},
		inputText: `
int32_to_str: {
  key: 0
}
int32_to_str: {
  value: "empty"
}
int32_to_str: {}
`,
		wantMessage: &pb3.Maps{
			Int32ToStr: map[int32]string{
				0: "",
			},
		},
	}, {
		desc:         "map field value contains invalid UTF-8",
		inputMessage: &pb3.Maps{},
		inputText: `int32_to_str: {
  key: 101
  value: "abc\xff"
}
`,
		wantErr: true,
	}, {
		desc:         "map field key contains invalid UTF-8",
		inputMessage: &pb3.Maps{},
		inputText: `str_to_nested: {
  key: "abc\xff"
  value: {}
}
`,
		wantErr: true,
	}, {
		desc:         "map contains unknown field",
		inputMessage: &pb3.Maps{},
		inputText: `
int32_to_str: {
  key: 0
  value: "cero"
  unknown: "bad"
}
`,
		wantErr: true,
	}, {
		desc:         "map contains extension-like key field",
		inputMessage: &pb3.Maps{},
		inputText: `
int32_to_str: {
  [key]: 10
  value: "ten"
}
`,
		wantErr: true,
	}, {
		desc:         "map contains invalid key",
		inputMessage: &pb3.Maps{},
		inputText: `
int32_to_str: {
  key: "invalid"
  value: "cero"
}
`,
		wantErr: true,
	}, {
		desc:         "map contains invalid value",
		inputMessage: &pb3.Maps{},
		inputText: `
int32_to_str: {
  key: 100
  value: 101
}
`,
		wantErr: true,
	}, {
		desc:         "map using mix of [] and repeated",
		inputMessage: &pb3.Maps{},
		inputText: `
int32_to_str: {
  key: 1
  value: "one"
}
int32_to_str: [
  {
    key: 2
    value: "not this"
  },
  {
  },
  {
    key: 3
    value: "three"
  }
]
int32_to_str: {
  key: 2
  value: "two"
}
`,
		wantMessage: &pb3.Maps{
			Int32ToStr: map[int32]string{
				0: "",
				1: "one",
				2: "two",
				3: "three",
			},
		},
	}, {
		desc:         "required fields not set",
		inputMessage: &pb2.Requireds{},
		wantErr:      true,
	}, {
		desc:         "required field set",
		inputMessage: &pb2.PartialRequired{},
		inputText:    "req_string: 'this is required'",
		wantMessage: &pb2.PartialRequired{
			ReqString: proto.String("this is required"),
		},
	}, {
		desc:         "required fields partially set",
		inputMessage: &pb2.Requireds{},
		inputText: `
req_bool: false
req_sfixed64: 3203386110
req_string: "hello"
req_enum: ONE
`,
		wantMessage: &pb2.Requireds{
			ReqBool:     proto.Bool(false),
			ReqSfixed64: proto.Int64(0xbeefcafe),
			ReqString:   proto.String("hello"),
			ReqEnum:     pb2.Enum_ONE.Enum(),
		},
		wantErr: true,
	}, {
		desc:         "required fields partially set with AllowPartial",
		umo:          prototext.UnmarshalOptions{AllowPartial: true},
		inputMessage: &pb2.Requireds{},
		inputText: `
req_bool: false
req_sfixed64: 3203386110
req_string: "hello"
req_enum: ONE
`,
		wantMessage: &pb2.Requireds{
			ReqBool:     proto.Bool(false),
			ReqSfixed64: proto.Int64(0xbeefcafe),
			ReqString:   proto.String("hello"),
			ReqEnum:     pb2.Enum_ONE.Enum(),
		},
	}, {
		desc:         "required fields all set",
		inputMessage: &pb2.Requireds{},
		inputText: `
req_bool: false
req_sfixed64: 0
req_double: 0
req_string: ""
req_enum: ONE
req_nested: {}
`,
		wantMessage: &pb2.Requireds{
			ReqBool:     proto.Bool(false),
			ReqSfixed64: proto.Int64(0),
			ReqDouble:   proto.Float64(0),
			ReqString:   proto.String(""),
			ReqEnum:     pb2.Enum_ONE.Enum(),
			ReqNested:   &pb2.Nested{},
		},
	}, {
		desc:         "indirect required field",
		inputMessage: &pb2.IndirectRequired{},
		inputText:    "opt_nested: {}",
		wantMessage: &pb2.IndirectRequired{
			OptNested: &pb2.NestedWithRequired{},
		},
		wantErr: true,
	}, {
		desc:         "indirect required field with AllowPartial",
		umo:          prototext.UnmarshalOptions{AllowPartial: true},
		inputMessage: &pb2.IndirectRequired{},
		inputText:    "opt_nested: {}",
		wantMessage: &pb2.IndirectRequired{
			OptNested: &pb2.NestedWithRequired{},
		},
	}, {
		desc:         "indirect required field in repeated",
		inputMessage: &pb2.IndirectRequired{},
		inputText: `
rpt_nested: {
  req_string: "one"
}
rpt_nested: {}
`,
		wantMessage: &pb2.IndirectRequired{
			RptNested: []*pb2.NestedWithRequired{
				{
					ReqString: proto.String("one"),
				},
				{},
			},
		},
		wantErr: true,
	}, {
		desc:         "indirect required field in repeated with AllowPartial",
		umo:          prototext.UnmarshalOptions{AllowPartial: true},
		inputMessage: &pb2.IndirectRequired{},
		inputText: `
rpt_nested: {
  req_string: "one"
}
rpt_nested: {}
`,
		wantMessage: &pb2.IndirectRequired{
			RptNested: []*pb2.NestedWithRequired{
				{
					ReqString: proto.String("one"),
				},
				{},
			},
		},
	}, {
		desc:         "indirect required field in map",
		inputMessage: &pb2.IndirectRequired{},
		inputText: `
str_to_nested: {
  key: "missing"
}
str_to_nested: {
  key: "contains"
  value: {
    req_string: "here"
  }
}
`,
		wantMessage: &pb2.IndirectRequired{
			StrToNested: map[string]*pb2.NestedWithRequired{
				"missing": &pb2.NestedWithRequired{},
				"contains": &pb2.NestedWithRequired{
					ReqString: proto.String("here"),
				},
			},
		},
		wantErr: true,
	}, {
		desc:         "indirect required field in map with AllowPartial",
		umo:          prototext.UnmarshalOptions{AllowPartial: true},
		inputMessage: &pb2.IndirectRequired{},
		inputText: `
str_to_nested: {
  key: "missing"
}
str_to_nested: {
  key: "contains"
  value: {
    req_string: "here"
  }
}
`,
		wantMessage: &pb2.IndirectRequired{
			StrToNested: map[string]*pb2.NestedWithRequired{
				"missing": &pb2.NestedWithRequired{},
				"contains": &pb2.NestedWithRequired{
					ReqString: proto.String("here"),
				},
			},
		},
	}, {
		desc:         "indirect required field in oneof",
		inputMessage: &pb2.IndirectRequired{},
		inputText: `oneof_nested: {}
`,
		wantMessage: &pb2.IndirectRequired{
			Union: &pb2.IndirectRequired_OneofNested{
				OneofNested: &pb2.NestedWithRequired{},
			},
		},
		wantErr: true,
	}, {
		desc:         "indirect required field in oneof with AllowPartial",
		umo:          prototext.UnmarshalOptions{AllowPartial: true},
		inputMessage: &pb2.IndirectRequired{},
		inputText: `oneof_nested: {}
`,
		wantMessage: &pb2.IndirectRequired{
			Union: &pb2.IndirectRequired_OneofNested{
				OneofNested: &pb2.NestedWithRequired{},
			},
		},
	}, {
		desc:         "ignore reserved field",
		inputMessage: &pb2.Nests{},
		inputText:    "reserved_field: 'ignore this'",
		wantMessage:  &pb2.Nests{},
	}, {
		desc:         "extensions of non-repeated fields",
		inputMessage: &pb2.Extensions{},
		inputText: `opt_string: "non-extension field"
[pb2.opt_ext_bool]: true
opt_bool: true
[pb2.opt_ext_nested]: {
  opt_string: "nested in an extension"
  opt_nested: {
    opt_string: "another nested in an extension"
  }
}
[pb2.opt_ext_string]: "extension field"
opt_int32: 42
[pb2.opt_ext_enum]: TEN
`,
		wantMessage: func() proto.Message {
			m := &pb2.Extensions{
				OptString: proto.String("non-extension field"),
				OptBool:   proto.Bool(true),
				OptInt32:  proto.Int32(42),
			}
			setExtension(m, pb2.E_OptExtBool, true)
			setExtension(m, pb2.E_OptExtString, "extension field")
			setExtension(m, pb2.E_OptExtEnum, pb2.Enum_TEN)
			setExtension(m, pb2.E_OptExtNested, &pb2.Nested{
				OptString: proto.String("nested in an extension"),
				OptNested: &pb2.Nested{
					OptString: proto.String("another nested in an extension"),
				},
			})
			return m
		}(),
	}, {
		desc:         "extension field contains invalid UTF-8",
		inputMessage: &pb2.Extensions{},
		inputText:    `[pb2.opt_ext_string]: "abc\xff"`,
		wantErr:      true,
	}, {
		desc:         "extensions of repeated fields",
		inputMessage: &pb2.Extensions{},
		inputText: `[pb2.rpt_ext_enum]: TEN
[pb2.rpt_ext_enum]: 101
[pb2.rpt_ext_fixed32]: 42
[pb2.rpt_ext_enum]: ONE
[pb2.rpt_ext_nested]: {
  opt_string: "one"
}
[pb2.rpt_ext_nested]: {
  opt_string: "two"
}
[pb2.rpt_ext_fixed32]: 47
[pb2.rpt_ext_nested]: {
  opt_string: "three"
}
`,
		wantMessage: func() proto.Message {
			m := &pb2.Extensions{}
			setExtension(m, pb2.E_RptExtEnum, &[]pb2.Enum{pb2.Enum_TEN, 101, pb2.Enum_ONE})
			setExtension(m, pb2.E_RptExtFixed32, &[]uint32{42, 47})
			setExtension(m, pb2.E_RptExtNested, &[]*pb2.Nested{
				&pb2.Nested{OptString: proto.String("one")},
				&pb2.Nested{OptString: proto.String("two")},
				&pb2.Nested{OptString: proto.String("three")},
			})
			return m
		}(),
	}, {
		desc:         "extensions of non-repeated fields in another message",
		inputMessage: &pb2.Extensions{},
		inputText: `[pb2.ExtensionsContainer.opt_ext_bool]: true
[pb2.ExtensionsContainer.opt_ext_enum]: TEN
[pb2.ExtensionsContainer.opt_ext_nested]: {
  opt_string: "nested in an extension"
  opt_nested: {
    opt_string: "another nested in an extension"
  }
}
[pb2.ExtensionsContainer.opt_ext_string]: "extension field"
`,
		wantMessage: func() proto.Message {
			m := &pb2.Extensions{}
			setExtension(m, pb2.E_ExtensionsContainer_OptExtBool, true)
			setExtension(m, pb2.E_ExtensionsContainer_OptExtString, "extension field")
			setExtension(m, pb2.E_ExtensionsContainer_OptExtEnum, pb2.Enum_TEN)
			setExtension(m, pb2.E_ExtensionsContainer_OptExtNested, &pb2.Nested{
				OptString: proto.String("nested in an extension"),
				OptNested: &pb2.Nested{
					OptString: proto.String("another nested in an extension"),
				},
			})
			return m
		}(),
	}, {
		desc:         "extensions of repeated fields in another message",
		inputMessage: &pb2.Extensions{},
		inputText: `opt_string: "non-extension field"
opt_bool: true
opt_int32: 42
[pb2.ExtensionsContainer.rpt_ext_nested]: {
  opt_string: "one"
}
[pb2.ExtensionsContainer.rpt_ext_enum]: TEN
[pb2.ExtensionsContainer.rpt_ext_nested]: {
  opt_string: "two"
}
[pb2.ExtensionsContainer.rpt_ext_enum]: 101
[pb2.ExtensionsContainer.rpt_ext_string]: "hello"
[pb2.ExtensionsContainer.rpt_ext_enum]: ONE
[pb2.ExtensionsContainer.rpt_ext_nested]: {
  opt_string: "three"
}
[pb2.ExtensionsContainer.rpt_ext_string]: "world"
`,
		wantMessage: func() proto.Message {
			m := &pb2.Extensions{
				OptString: proto.String("non-extension field"),
				OptBool:   proto.Bool(true),
				OptInt32:  proto.Int32(42),
			}
			setExtension(m, pb2.E_ExtensionsContainer_RptExtEnum, &[]pb2.Enum{pb2.Enum_TEN, 101, pb2.Enum_ONE})
			setExtension(m, pb2.E_ExtensionsContainer_RptExtString, &[]string{"hello", "world"})
			setExtension(m, pb2.E_ExtensionsContainer_RptExtNested, &[]*pb2.Nested{
				&pb2.Nested{OptString: proto.String("one")},
				&pb2.Nested{OptString: proto.String("two")},
				&pb2.Nested{OptString: proto.String("three")},
			})
			return m
		}(),
	}, {
		desc:         "invalid extension field name",
		inputMessage: &pb2.Extensions{},
		inputText:    "[pb2.invalid_message_field]: true",
		wantErr:      true,
	}, {
		desc:         "MessageSet",
		inputMessage: &pb2.MessageSet{},
		inputText: `
[pb2.MessageSetExtension]: {
  opt_string: "a messageset extension"
}
[pb2.MessageSetExtension.ext_nested]: {
  opt_string: "just a regular extension"
}
[pb2.MessageSetExtension.not_message_set_extension]: {
  opt_string: "not a messageset extension"
}
`,
		wantMessage: func() proto.Message {
			m := &pb2.MessageSet{}
			setExtension(m, pb2.E_MessageSetExtension_MessageSetExtension, &pb2.MessageSetExtension{
				OptString: proto.String("a messageset extension"),
			})
			setExtension(m, pb2.E_MessageSetExtension_NotMessageSetExtension, &pb2.MessageSetExtension{
				OptString: proto.String("not a messageset extension"),
			})
			setExtension(m, pb2.E_MessageSetExtension_ExtNested, &pb2.Nested{
				OptString: proto.String("just a regular extension"),
			})
			return m
		}(),
	}, {
		desc:         "not real MessageSet 1",
		inputMessage: &pb2.FakeMessageSet{},
		inputText: `
[pb2.FakeMessageSetExtension.message_set_extension]: {
  opt_string: "not a messageset extension"
}
`,
		wantMessage: func() proto.Message {
			m := &pb2.FakeMessageSet{}
			setExtension(m, pb2.E_FakeMessageSetExtension_MessageSetExtension, &pb2.FakeMessageSetExtension{
				OptString: proto.String("not a messageset extension"),
			})
			return m
		}(),
	}, {
		desc:         "not real MessageSet 2",
		inputMessage: &pb2.FakeMessageSet{},
		inputText: `
[pb2.FakeMessageSetExtension]: {
  opt_string: "not a messageset extension"
}
`,
		wantErr: true,
	}, {
		desc:         "not real MessageSet 3",
		inputMessage: &pb2.MessageSet{},
		inputText: `
[pb2.message_set_extension]: {
  opt_string: "another not a messageset extension"
}`,
		wantMessage: func() proto.Message {
			m := &pb2.MessageSet{}
			setExtension(m, pb2.E_MessageSetExtension, &pb2.FakeMessageSetExtension{
				OptString: proto.String("another not a messageset extension"),
			})
			return m
		}(),
	}, {
		desc:         "Any not expanded",
		inputMessage: &anypb.Any{},
		inputText: `
type_url: "pb2.Nested"
value: "some bytes"
`,
		wantMessage: &anypb.Any{
			TypeUrl: "pb2.Nested",
			Value:   []byte("some bytes"),
		},
	}, {
		desc:         "Any not expanded missing value",
		inputMessage: &anypb.Any{},
		inputText:    `type_url: "pb2.Nested"`,
		wantMessage: &anypb.Any{
			TypeUrl: "pb2.Nested",
		},
	}, {
		desc:         "Any not expanded missing type_url",
		inputMessage: &anypb.Any{},
		inputText:    `value: "some bytes"`,
		wantMessage: &anypb.Any{
			Value: []byte("some bytes"),
		},
	}, {
		desc: "Any expanded",
		umo: prototext.UnmarshalOptions{
			Resolver: preg.NewTypes(pimpl.Export{}.MessageTypeOf(&pb2.Nested{})),
		},
		inputMessage: &anypb.Any{},
		inputText: `
[foobar/pb2.Nested]: {
  opt_string: "embedded inside Any"
  opt_nested: {
    opt_string: "inception"
  }
}
`,
		wantMessage: func() proto.Message {
			m := &pb2.Nested{
				OptString: proto.String("embedded inside Any"),
				OptNested: &pb2.Nested{
					OptString: proto.String("inception"),
				},
			}
			b, err := proto.MarshalOptions{Deterministic: true}.Marshal(m)
			if err != nil {
				t.Fatalf("error in binary marshaling message for Any.value: %v", err)
			}
			return &anypb.Any{
				TypeUrl: "foobar/pb2.Nested",
				Value:   b,
			}
		}(),
	}, {
		desc: "Any expanded with empty value",
		umo: prototext.UnmarshalOptions{
			Resolver: preg.NewTypes(pimpl.Export{}.MessageTypeOf(&pb2.Nested{})),
		},
		inputMessage: &anypb.Any{},
		inputText:    `[foo.com/pb2.Nested]: {}`,
		wantMessage: &anypb.Any{
			TypeUrl: "foo.com/pb2.Nested",
		},
	}, {
		desc: "Any expanded with missing required",
		umo: prototext.UnmarshalOptions{
			Resolver: preg.NewTypes(pimpl.Export{}.MessageTypeOf(&pb2.PartialRequired{})),
		},
		inputMessage: &anypb.Any{},
		inputText: `
[pb2.PartialRequired]: {
  opt_string: "embedded inside Any"
}
`,
		wantMessage: func() proto.Message {
			m := &pb2.PartialRequired{
				OptString: proto.String("embedded inside Any"),
			}
			b, err := proto.MarshalOptions{
				AllowPartial:  true,
				Deterministic: true,
			}.Marshal(m)
			if err != nil {
				t.Fatalf("error in binary marshaling message for Any.value: %v", err)
			}
			return &anypb.Any{
				TypeUrl: "pb2.PartialRequired",
				Value:   b,
			}
		}(),
	}, {
		desc: "Any with invalid UTF-8",
		umo: prototext.UnmarshalOptions{
			Resolver: preg.NewTypes(pimpl.Export{}.MessageTypeOf(&pb3.Nested{})),
		},
		inputMessage: &anypb.Any{},
		inputText: `
[pb3.Nested]: {
  s_string: "abc\xff"
}
`,
		wantErr: true,
	}, {
		desc:         "Any expanded with unregistered type",
		umo:          prototext.UnmarshalOptions{Resolver: preg.NewTypes()},
		inputMessage: &anypb.Any{},
		inputText:    `[SomeMessage]: {}`,
		wantErr:      true,
	}, {
		desc: "Any expanded with invalid value",
		umo: prototext.UnmarshalOptions{
			Resolver: preg.NewTypes(pimpl.Export{}.MessageTypeOf(&pb2.Nested{})),
		},
		inputMessage: &anypb.Any{},
		inputText:    `[pb2.Nested]: 123`,
		wantErr:      true,
	}, {
		desc: "Any expanded with unknown fields",
		umo: prototext.UnmarshalOptions{
			Resolver: preg.NewTypes(pimpl.Export{}.MessageTypeOf(&pb2.Nested{})),
		},
		inputMessage: &anypb.Any{},
		inputText: `
[pb2.Nested]: {}
unknown: ""
`,
		wantErr: true,
	}, {
		desc: "Any contains expanded and unexpanded fields",
		umo: prototext.UnmarshalOptions{
			Resolver: preg.NewTypes(pimpl.Export{}.MessageTypeOf(&pb2.Nested{})),
		},
		inputMessage: &anypb.Any{},
		inputText: `
[pb2.Nested]: {}
type_url: "pb2.Nested"
`,
		wantErr: true,
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			err := tt.umo.Unmarshal([]byte(tt.inputText), tt.inputMessage)
			if err != nil && !tt.wantErr {
				t.Errorf("Unmarshal() returned error: %v\n\n", err)
			}
			if err == nil && tt.wantErr {
				t.Error("Unmarshal() got nil error, want error\n\n")
			}
			if tt.wantMessage != nil && !proto.Equal(tt.inputMessage, tt.wantMessage) {
				t.Errorf("Unmarshal()\n<got>\n%v\n<want>\n%v\n", tt.inputMessage, tt.wantMessage)
			}
		})
	}
}
