// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prototext_test

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/internal/detrand"
	"google.golang.org/protobuf/internal/encoding/pack"
	pimpl "google.golang.org/protobuf/internal/impl"
	"google.golang.org/protobuf/internal/scalar"
	"google.golang.org/protobuf/proto"
	preg "google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/runtime/protoiface"

	"google.golang.org/protobuf/encoding/testprotos/pb2"
	"google.golang.org/protobuf/encoding/testprotos/pb3"
	"google.golang.org/protobuf/types/known/anypb"
)

func init() {
	// Disable detrand to enable direct comparisons on outputs.
	detrand.Disable()
}

// TODO: Use proto.SetExtension when available.
func setExtension(m proto.Message, xd *protoiface.ExtensionDescV1, val interface{}) {
	m.ProtoReflect().Set(xd.Type, xd.Type.ValueOf(val))
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		desc    string
		mo      prototext.MarshalOptions
		input   proto.Message
		want    string
		wantErr bool // TODO: Verify error message content.
	}{{
		desc:  "proto2 optional scalars not set",
		input: &pb2.Scalars{},
		want:  "\n",
	}, {
		desc:  "proto3 scalars not set",
		input: &pb3.Scalars{},
		want:  "\n",
	}, {
		desc: "proto2 optional scalars set to zero values",
		input: &pb2.Scalars{
			OptBool:     scalar.Bool(false),
			OptInt32:    scalar.Int32(0),
			OptInt64:    scalar.Int64(0),
			OptUint32:   scalar.Uint32(0),
			OptUint64:   scalar.Uint64(0),
			OptSint32:   scalar.Int32(0),
			OptSint64:   scalar.Int64(0),
			OptFixed32:  scalar.Uint32(0),
			OptFixed64:  scalar.Uint64(0),
			OptSfixed32: scalar.Int32(0),
			OptSfixed64: scalar.Int64(0),
			OptFloat:    scalar.Float32(0),
			OptDouble:   scalar.Float64(0),
			OptBytes:    []byte{},
			OptString:   scalar.String(""),
		},
		want: `opt_bool: false
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
	}, {
		desc: "proto3 scalars set to zero values",
		input: &pb3.Scalars{
			SBool:     false,
			SInt32:    0,
			SInt64:    0,
			SUint32:   0,
			SUint64:   0,
			SSint32:   0,
			SSint64:   0,
			SFixed32:  0,
			SFixed64:  0,
			SSfixed32: 0,
			SSfixed64: 0,
			SFloat:    0,
			SDouble:   0,
			SBytes:    []byte{},
			SString:   "",
		},
		want: "\n",
	}, {
		desc: "proto2 optional scalars set to some values",
		input: &pb2.Scalars{
			OptBool:     scalar.Bool(true),
			OptInt32:    scalar.Int32(0xff),
			OptInt64:    scalar.Int64(0xdeadbeef),
			OptUint32:   scalar.Uint32(47),
			OptUint64:   scalar.Uint64(0xdeadbeef),
			OptSint32:   scalar.Int32(-1001),
			OptSint64:   scalar.Int64(-0xffff),
			OptFixed64:  scalar.Uint64(64),
			OptSfixed32: scalar.Int32(-32),
			OptFloat:    scalar.Float32(1.02),
			OptDouble:   scalar.Float64(1.0199999809265137),
			OptBytes:    []byte("\xe8\xb0\xb7\xe6\xad\x8c"),
			OptString:   scalar.String("谷歌"),
		},
		want: `opt_bool: true
opt_int32: 255
opt_int64: 3735928559
opt_uint32: 47
opt_uint64: 3735928559
opt_sint32: -1001
opt_sint64: -65535
opt_fixed64: 64
opt_sfixed32: -32
opt_float: 1.02
opt_double: 1.0199999809265137
opt_bytes: "谷歌"
opt_string: "谷歌"
`,
	}, {
		desc: "string with invalid UTF-8",
		input: &pb3.Scalars{
			SString: "abc\xff",
		},
		wantErr: true,
	}, {
		desc: "float nan",
		input: &pb3.Scalars{
			SFloat: float32(math.NaN()),
		},
		want: "s_float: nan\n",
	}, {
		desc: "float positive infinity",
		input: &pb3.Scalars{
			SFloat: float32(math.Inf(1)),
		},
		want: "s_float: inf\n",
	}, {
		desc: "float negative infinity",
		input: &pb3.Scalars{
			SFloat: float32(math.Inf(-1)),
		},
		want: "s_float: -inf\n",
	}, {
		desc: "double nan",
		input: &pb3.Scalars{
			SDouble: math.NaN(),
		},
		want: "s_double: nan\n",
	}, {
		desc: "double positive infinity",
		input: &pb3.Scalars{
			SDouble: math.Inf(1),
		},
		want: "s_double: inf\n",
	}, {
		desc: "double negative infinity",
		input: &pb3.Scalars{
			SDouble: math.Inf(-1),
		},
		want: "s_double: -inf\n",
	}, {
		desc:  "proto2 enum not set",
		input: &pb2.Enums{},
		want:  "\n",
	}, {
		desc: "proto2 enum set to zero value",
		input: &pb2.Enums{
			OptEnum:       pb2.Enum(0).Enum(),
			OptNestedEnum: pb2.Enums_NestedEnum(0).Enum(),
		},
		want: `opt_enum: 0
opt_nested_enum: 0
`,
	}, {
		desc: "proto2 enum",
		input: &pb2.Enums{
			OptEnum:       pb2.Enum_ONE.Enum(),
			OptNestedEnum: pb2.Enums_UNO.Enum(),
		},
		want: `opt_enum: ONE
opt_nested_enum: UNO
`,
	}, {
		desc: "proto2 enum set to numeric values",
		input: &pb2.Enums{
			OptEnum:       pb2.Enum(2).Enum(),
			OptNestedEnum: pb2.Enums_NestedEnum(2).Enum(),
		},
		want: `opt_enum: TWO
opt_nested_enum: DOS
`,
	}, {
		desc: "proto2 enum set to unnamed numeric values",
		input: &pb2.Enums{
			OptEnum:       pb2.Enum(101).Enum(),
			OptNestedEnum: pb2.Enums_NestedEnum(-101).Enum(),
		},
		want: `opt_enum: 101
opt_nested_enum: -101
`,
	}, {
		desc:  "proto3 enum not set",
		input: &pb3.Enums{},
		want:  "\n",
	}, {
		desc: "proto3 enum set to zero value",
		input: &pb3.Enums{
			SEnum:       pb3.Enum_ZERO,
			SNestedEnum: pb3.Enums_CERO,
		},
		want: "\n",
	}, {
		desc: "proto3 enum",
		input: &pb3.Enums{
			SEnum:       pb3.Enum_ONE,
			SNestedEnum: pb3.Enums_UNO,
		},
		want: `s_enum: ONE
s_nested_enum: UNO
`,
	}, {
		desc: "proto3 enum set to numeric values",
		input: &pb3.Enums{
			SEnum:       2,
			SNestedEnum: 2,
		},
		want: `s_enum: TWO
s_nested_enum: DOS
`,
	}, {
		desc: "proto3 enum set to unnamed numeric values",
		input: &pb3.Enums{
			SEnum:       -47,
			SNestedEnum: 47,
		},
		want: `s_enum: -47
s_nested_enum: 47
`,
	}, {
		desc:  "proto2 nested message not set",
		input: &pb2.Nests{},
		want:  "\n",
	}, {
		desc: "proto2 nested message set to empty",
		input: &pb2.Nests{
			OptNested: &pb2.Nested{},
			Optgroup:  &pb2.Nests_OptGroup{},
		},
		want: `opt_nested: {}
OptGroup: {}
`,
	}, {
		desc: "proto2 nested messages",
		input: &pb2.Nests{
			OptNested: &pb2.Nested{
				OptString: scalar.String("nested message"),
				OptNested: &pb2.Nested{
					OptString: scalar.String("another nested message"),
				},
			},
		},
		want: `opt_nested: {
  opt_string: "nested message"
  opt_nested: {
    opt_string: "another nested message"
  }
}
`,
	}, {
		desc: "proto2 groups",
		input: &pb2.Nests{
			Optgroup: &pb2.Nests_OptGroup{
				OptString: scalar.String("inside a group"),
				OptNested: &pb2.Nested{
					OptString: scalar.String("nested message inside a group"),
				},
				Optnestedgroup: &pb2.Nests_OptGroup_OptNestedGroup{
					OptFixed32: scalar.Uint32(47),
				},
			},
		},
		want: `OptGroup: {
  opt_string: "inside a group"
  opt_nested: {
    opt_string: "nested message inside a group"
  }
  OptNestedGroup: {
    opt_fixed32: 47
  }
}
`,
	}, {
		desc:  "proto3 nested message not set",
		input: &pb3.Nests{},
		want:  "\n",
	}, {
		desc: "proto3 nested message set to empty",
		input: &pb3.Nests{
			SNested: &pb3.Nested{},
		},
		want: "s_nested: {}\n",
	}, {
		desc: "proto3 nested message",
		input: &pb3.Nests{
			SNested: &pb3.Nested{
				SString: "nested message",
				SNested: &pb3.Nested{
					SString: "another nested message",
				},
			},
		},
		want: `s_nested: {
  s_string: "nested message"
  s_nested: {
    s_string: "another nested message"
  }
}
`,
	}, {
		desc: "proto3 nested message contains invalid UTF-8",
		input: &pb3.Nests{
			SNested: &pb3.Nested{
				SString: "abc\xff",
			},
		},
		wantErr: true,
	}, {
		desc:  "oneof not set",
		input: &pb3.Oneofs{},
		want:  "\n",
	}, {
		desc: "oneof set to empty string",
		input: &pb3.Oneofs{
			Union: &pb3.Oneofs_OneofString{},
		},
		want: `oneof_string: ""
`,
	}, {
		desc: "oneof set to string",
		input: &pb3.Oneofs{
			Union: &pb3.Oneofs_OneofString{
				OneofString: "hello",
			},
		},
		want: `oneof_string: "hello"
`,
	}, {
		desc: "oneof set to enum",
		input: &pb3.Oneofs{
			Union: &pb3.Oneofs_OneofEnum{
				OneofEnum: pb3.Enum_ZERO,
			},
		},
		want: `oneof_enum: ZERO
`,
	}, {
		desc: "oneof set to empty message",
		input: &pb3.Oneofs{
			Union: &pb3.Oneofs_OneofNested{
				OneofNested: &pb3.Nested{},
			},
		},
		want: "oneof_nested: {}\n",
	}, {
		desc: "oneof set to message",
		input: &pb3.Oneofs{
			Union: &pb3.Oneofs_OneofNested{
				OneofNested: &pb3.Nested{
					SString: "nested message",
				},
			},
		},
		want: `oneof_nested: {
  s_string: "nested message"
}
`,
	}, {
		desc:  "repeated fields not set",
		input: &pb2.Repeats{},
		want:  "\n",
	}, {
		desc: "repeated fields set to empty slices",
		input: &pb2.Repeats{
			RptBool:   []bool{},
			RptInt32:  []int32{},
			RptInt64:  []int64{},
			RptUint32: []uint32{},
			RptUint64: []uint64{},
			RptFloat:  []float32{},
			RptDouble: []float64{},
			RptBytes:  [][]byte{},
		},
		want: "\n",
	}, {
		desc: "repeated fields set to some values",
		input: &pb2.Repeats{
			RptBool:   []bool{true, false, true, true},
			RptInt32:  []int32{1, 6, 0, 0},
			RptInt64:  []int64{-64, 47},
			RptUint32: []uint32{0xff, 0xffff},
			RptUint64: []uint64{0xdeadbeef},
			RptFloat:  []float32{float32(math.NaN()), float32(math.Inf(1)), float32(math.Inf(-1)), 1.034},
			RptDouble: []float64{math.NaN(), math.Inf(1), math.Inf(-1), 1.23e-308},
			RptString: []string{"hello", "世界"},
			RptBytes: [][]byte{
				[]byte("hello"),
				[]byte("\xe4\xb8\x96\xe7\x95\x8c"),
			},
		},
		want: `rpt_bool: true
rpt_bool: false
rpt_bool: true
rpt_bool: true
rpt_int32: 1
rpt_int32: 6
rpt_int32: 0
rpt_int32: 0
rpt_int64: -64
rpt_int64: 47
rpt_uint32: 255
rpt_uint32: 65535
rpt_uint64: 3735928559
rpt_float: nan
rpt_float: inf
rpt_float: -inf
rpt_float: 1.034
rpt_double: nan
rpt_double: inf
rpt_double: -inf
rpt_double: 1.23e-308
rpt_string: "hello"
rpt_string: "世界"
rpt_bytes: "hello"
rpt_bytes: "世界"
`,
	}, {
		desc: "repeated contains invalid UTF-8",
		input: &pb2.Repeats{
			RptString: []string{"abc\xff"},
		},
		wantErr: true,
	}, {
		desc: "repeated enums",
		input: &pb2.Enums{
			RptEnum:       []pb2.Enum{pb2.Enum_ONE, 2, pb2.Enum_TEN, 42},
			RptNestedEnum: []pb2.Enums_NestedEnum{2, 47, 10},
		},
		want: `rpt_enum: ONE
rpt_enum: TWO
rpt_enum: TEN
rpt_enum: 42
rpt_nested_enum: DOS
rpt_nested_enum: 47
rpt_nested_enum: DIEZ
`,
	}, {
		desc: "repeated messages set to empty",
		input: &pb2.Nests{
			RptNested: []*pb2.Nested{},
			Rptgroup:  []*pb2.Nests_RptGroup{},
		},
		want: "\n",
	}, {
		desc: "repeated messages",
		input: &pb2.Nests{
			RptNested: []*pb2.Nested{
				{
					OptString: scalar.String("repeat nested one"),
				},
				{
					OptString: scalar.String("repeat nested two"),
					OptNested: &pb2.Nested{
						OptString: scalar.String("inside repeat nested two"),
					},
				},
				{},
			},
		},
		want: `rpt_nested: {
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
	}, {
		desc: "repeated messages contains nil value",
		input: &pb2.Nests{
			RptNested: []*pb2.Nested{nil, {}},
		},
		want: `rpt_nested: {}
rpt_nested: {}
`,
	}, {
		desc: "repeated groups",
		input: &pb2.Nests{
			Rptgroup: []*pb2.Nests_RptGroup{
				{
					RptString: []string{"hello", "world"},
				},
				{},
				nil,
			},
		},
		want: `RptGroup: {
  rpt_string: "hello"
  rpt_string: "world"
}
RptGroup: {}
RptGroup: {}
`,
	}, {
		desc:  "map fields not set",
		input: &pb3.Maps{},
		want:  "\n",
	}, {
		desc: "map fields set to empty",
		input: &pb3.Maps{
			Int32ToStr:   map[int32]string{},
			BoolToUint32: map[bool]uint32{},
			Uint64ToEnum: map[uint64]pb3.Enum{},
			StrToNested:  map[string]*pb3.Nested{},
			StrToOneofs:  map[string]*pb3.Oneofs{},
		},
		want: "\n",
	}, {
		desc: "map fields 1",
		input: &pb3.Maps{
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
		want: `int32_to_str: {
  key: -101
  value: "-101"
}
int32_to_str: {
  key: 0
  value: "zero"
}
int32_to_str: {
  key: 255
  value: "0xff"
}
bool_to_uint32: {
  key: false
  value: 101
}
bool_to_uint32: {
  key: true
  value: 42
}
`,
	}, {
		desc: "map fields 2",
		input: &pb3.Maps{
			Uint64ToEnum: map[uint64]pb3.Enum{
				1:  pb3.Enum_ONE,
				2:  pb3.Enum_TWO,
				10: pb3.Enum_TEN,
				47: 47,
			},
		},
		want: `uint64_to_enum: {
  key: 1
  value: ONE
}
uint64_to_enum: {
  key: 2
  value: TWO
}
uint64_to_enum: {
  key: 10
  value: TEN
}
uint64_to_enum: {
  key: 47
  value: 47
}
`,
	}, {
		desc: "map fields 3",
		input: &pb3.Maps{
			StrToNested: map[string]*pb3.Nested{
				"nested": &pb3.Nested{
					SString: "nested in a map",
				},
			},
		},
		want: `str_to_nested: {
  key: "nested"
  value: {
    s_string: "nested in a map"
  }
}
`,
	}, {
		desc: "map fields 4",
		input: &pb3.Maps{
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
		want: `str_to_oneofs: {
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
	}, {
		desc: "map field value contains invalid UTF-8",
		input: &pb3.Maps{
			Int32ToStr: map[int32]string{
				101: "abc\xff",
			},
		},
		wantErr: true,
	}, {
		desc: "map field key contains invalid UTF-8",
		input: &pb3.Maps{
			StrToNested: map[string]*pb3.Nested{
				"abc\xff": {},
			},
		},
		wantErr: true,
	}, {
		desc: "map field contains nil value",
		input: &pb3.Maps{
			StrToNested: map[string]*pb3.Nested{
				"nil": nil,
			},
		},
		want: `str_to_nested: {
  key: "nil"
  value: {}
}
`,
	}, {
		desc:    "required fields not set",
		input:   &pb2.Requireds{},
		want:    "\n",
		wantErr: true,
	}, {
		desc: "required fields partially set",
		input: &pb2.Requireds{
			ReqBool:     scalar.Bool(false),
			ReqSfixed64: scalar.Int64(0xbeefcafe),
			ReqDouble:   scalar.Float64(math.NaN()),
			ReqString:   scalar.String("hello"),
			ReqEnum:     pb2.Enum_ONE.Enum(),
		},
		want: `req_bool: false
req_sfixed64: 3203386110
req_double: nan
req_string: "hello"
req_enum: ONE
`,
		wantErr: true,
	}, {
		desc: "required fields not set with AllowPartial",
		mo:   prototext.MarshalOptions{AllowPartial: true},
		input: &pb2.Requireds{
			ReqBool:     scalar.Bool(false),
			ReqSfixed64: scalar.Int64(0xbeefcafe),
			ReqDouble:   scalar.Float64(math.NaN()),
			ReqString:   scalar.String("hello"),
			ReqEnum:     pb2.Enum_ONE.Enum(),
		},
		want: `req_bool: false
req_sfixed64: 3203386110
req_double: nan
req_string: "hello"
req_enum: ONE
`,
	}, {
		desc: "required fields all set",
		input: &pb2.Requireds{
			ReqBool:     scalar.Bool(false),
			ReqSfixed64: scalar.Int64(0),
			ReqDouble:   scalar.Float64(1.23),
			ReqString:   scalar.String(""),
			ReqEnum:     pb2.Enum_ONE.Enum(),
			ReqNested:   &pb2.Nested{},
		},
		want: `req_bool: false
req_sfixed64: 0
req_double: 1.23
req_string: ""
req_enum: ONE
req_nested: {}
`,
	}, {
		desc: "indirect required field",
		input: &pb2.IndirectRequired{
			OptNested: &pb2.NestedWithRequired{},
		},
		want:    "opt_nested: {}\n",
		wantErr: true,
	}, {
		desc: "indirect required field with AllowPartial",
		mo:   prototext.MarshalOptions{AllowPartial: true},
		input: &pb2.IndirectRequired{
			OptNested: &pb2.NestedWithRequired{},
		},
		want: "opt_nested: {}\n",
	}, {
		desc: "indirect required field in empty repeated",
		input: &pb2.IndirectRequired{
			RptNested: []*pb2.NestedWithRequired{},
		},
		want: "\n",
	}, {
		desc: "indirect required field in repeated",
		input: &pb2.IndirectRequired{
			RptNested: []*pb2.NestedWithRequired{
				&pb2.NestedWithRequired{},
			},
		},
		want:    "rpt_nested: {}\n",
		wantErr: true,
	}, {
		desc: "indirect required field in repeated with AllowPartial",
		mo:   prototext.MarshalOptions{AllowPartial: true},
		input: &pb2.IndirectRequired{
			RptNested: []*pb2.NestedWithRequired{
				&pb2.NestedWithRequired{},
			},
		},
		want: "rpt_nested: {}\n",
	}, {
		desc: "indirect required field in empty map",
		input: &pb2.IndirectRequired{
			StrToNested: map[string]*pb2.NestedWithRequired{},
		},
		want: "\n",
	}, {
		desc: "indirect required field in map",
		input: &pb2.IndirectRequired{
			StrToNested: map[string]*pb2.NestedWithRequired{
				"fail": &pb2.NestedWithRequired{},
			},
		},
		want: `str_to_nested: {
  key: "fail"
  value: {}
}
`,
		wantErr: true,
	}, {
		desc: "indirect required field in map with AllowPartial",
		mo:   prototext.MarshalOptions{AllowPartial: true},
		input: &pb2.IndirectRequired{
			StrToNested: map[string]*pb2.NestedWithRequired{
				"fail": &pb2.NestedWithRequired{},
			},
		},
		want: `str_to_nested: {
  key: "fail"
  value: {}
}
`,
	}, {
		desc: "indirect required field in oneof",
		input: &pb2.IndirectRequired{
			Union: &pb2.IndirectRequired_OneofNested{
				OneofNested: &pb2.NestedWithRequired{},
			},
		},
		want:    "oneof_nested: {}\n",
		wantErr: true,
	}, {
		desc: "indirect required field in oneof with AllowPartial",
		mo:   prototext.MarshalOptions{AllowPartial: true},
		input: &pb2.IndirectRequired{
			Union: &pb2.IndirectRequired_OneofNested{
				OneofNested: &pb2.NestedWithRequired{},
			},
		},
		want: "oneof_nested: {}\n",
	}, {
		desc: "unknown varint and fixed types",
		input: func() proto.Message {
			m := &pb2.Scalars{
				OptString: scalar.String("this message contains unknown fields"),
			}
			m.ProtoReflect().SetUnknown(pack.Message{
				pack.Tag{101, pack.VarintType}, pack.Bool(true),
				pack.Tag{102, pack.VarintType}, pack.Varint(0xff),
				pack.Tag{103, pack.Fixed32Type}, pack.Uint32(47),
				pack.Tag{104, pack.Fixed64Type}, pack.Int64(0xdeadbeef),
			}.Marshal())
			return m
		}(),
		want: `opt_string: "this message contains unknown fields"
101: 1
102: 255
103: 47
104: 3735928559
`,
	}, {
		desc: "unknown length-delimited",
		input: func() proto.Message {
			m := new(pb2.Scalars)
			m.ProtoReflect().SetUnknown(pack.Message{
				pack.Tag{101, pack.BytesType}, pack.LengthPrefix{pack.Bool(true), pack.Bool(false)},
				pack.Tag{102, pack.BytesType}, pack.String("hello world"),
				pack.Tag{103, pack.BytesType}, pack.Bytes("\xe4\xb8\x96\xe7\x95\x8c"),
			}.Marshal())
			return m
		}(),
		want: `101: "\x01\x00"
102: "hello world"
103: "世界"
`,
	}, {
		desc: "unknown group type",
		input: func() proto.Message {
			m := new(pb2.Scalars)
			m.ProtoReflect().SetUnknown(pack.Message{
				pack.Tag{101, pack.StartGroupType}, pack.Tag{101, pack.EndGroupType},
				pack.Tag{102, pack.StartGroupType},
				pack.Tag{101, pack.VarintType}, pack.Bool(false),
				pack.Tag{102, pack.BytesType}, pack.String("inside a group"),
				pack.Tag{102, pack.EndGroupType},
			}.Marshal())
			return m
		}(),
		want: `101: {}
102: {
  101: 0
  102: "inside a group"
}
`,
	}, {
		desc: "unknown unpack repeated field",
		input: func() proto.Message {
			m := new(pb2.Scalars)
			m.ProtoReflect().SetUnknown(pack.Message{
				pack.Tag{101, pack.BytesType}, pack.LengthPrefix{pack.Bool(true), pack.Bool(false), pack.Bool(true)},
				pack.Tag{102, pack.BytesType}, pack.String("hello"),
				pack.Tag{101, pack.VarintType}, pack.Bool(true),
				pack.Tag{102, pack.BytesType}, pack.String("世界"),
			}.Marshal())
			return m
		}(),
		want: `101: "\x01\x00\x01"
102: "hello"
101: 1
102: "世界"
`,
	}, {
		desc: "extensions of non-repeated fields",
		input: func() proto.Message {
			m := &pb2.Extensions{
				OptString: scalar.String("non-extension field"),
				OptBool:   scalar.Bool(true),
				OptInt32:  scalar.Int32(42),
			}
			setExtension(m, pb2.E_OptExtBool, true)
			setExtension(m, pb2.E_OptExtString, "extension field")
			setExtension(m, pb2.E_OptExtEnum, pb2.Enum_TEN)
			setExtension(m, pb2.E_OptExtNested, &pb2.Nested{
				OptString: scalar.String("nested in an extension"),
				OptNested: &pb2.Nested{
					OptString: scalar.String("another nested in an extension"),
				},
			})
			return m
		}(),
		want: `opt_string: "non-extension field"
opt_bool: true
opt_int32: 42
[pb2.opt_ext_bool]: true
[pb2.opt_ext_enum]: TEN
[pb2.opt_ext_nested]: {
  opt_string: "nested in an extension"
  opt_nested: {
    opt_string: "another nested in an extension"
  }
}
[pb2.opt_ext_string]: "extension field"
`,
	}, {
		desc: "extension field contains invalid UTF-8",
		input: func() proto.Message {
			m := &pb2.Extensions{}
			setExtension(m, pb2.E_OptExtString, "abc\xff")
			return m
		}(),
		wantErr: true,
	}, {
		desc: "extension partial returns error",
		input: func() proto.Message {
			m := &pb2.Extensions{}
			setExtension(m, pb2.E_OptExtPartial, &pb2.PartialRequired{
				OptString: scalar.String("partial1"),
			})
			setExtension(m, pb2.E_ExtensionsContainer_OptExtPartial, &pb2.PartialRequired{
				OptString: scalar.String("partial2"),
			})
			return m
		}(),
		want: `[pb2.ExtensionsContainer.opt_ext_partial]: {
  opt_string: "partial2"
}
[pb2.opt_ext_partial]: {
  opt_string: "partial1"
}
`,
		wantErr: true,
	}, {
		desc: "extension partial with AllowPartial",
		mo:   prototext.MarshalOptions{AllowPartial: true},
		input: func() proto.Message {
			m := &pb2.Extensions{}
			setExtension(m, pb2.E_OptExtPartial, &pb2.PartialRequired{
				OptString: scalar.String("partial1"),
			})
			return m
		}(),
		want: `[pb2.opt_ext_partial]: {
  opt_string: "partial1"
}
`,
	}, {
		desc: "extensions of repeated fields",
		input: func() proto.Message {
			m := &pb2.Extensions{}
			setExtension(m, pb2.E_RptExtEnum, &[]pb2.Enum{pb2.Enum_TEN, 101, pb2.Enum_ONE})
			setExtension(m, pb2.E_RptExtFixed32, &[]uint32{42, 47})
			setExtension(m, pb2.E_RptExtNested, &[]*pb2.Nested{
				&pb2.Nested{OptString: scalar.String("one")},
				&pb2.Nested{OptString: scalar.String("two")},
				&pb2.Nested{OptString: scalar.String("three")},
			})
			return m
		}(),
		want: `[pb2.rpt_ext_enum]: TEN
[pb2.rpt_ext_enum]: 101
[pb2.rpt_ext_enum]: ONE
[pb2.rpt_ext_fixed32]: 42
[pb2.rpt_ext_fixed32]: 47
[pb2.rpt_ext_nested]: {
  opt_string: "one"
}
[pb2.rpt_ext_nested]: {
  opt_string: "two"
}
[pb2.rpt_ext_nested]: {
  opt_string: "three"
}
`,
	}, {
		desc: "extensions of non-repeated fields in another message",
		input: func() proto.Message {
			m := &pb2.Extensions{}
			setExtension(m, pb2.E_ExtensionsContainer_OptExtBool, true)
			setExtension(m, pb2.E_ExtensionsContainer_OptExtString, "extension field")
			setExtension(m, pb2.E_ExtensionsContainer_OptExtEnum, pb2.Enum_TEN)
			setExtension(m, pb2.E_ExtensionsContainer_OptExtNested, &pb2.Nested{
				OptString: scalar.String("nested in an extension"),
				OptNested: &pb2.Nested{
					OptString: scalar.String("another nested in an extension"),
				},
			})
			return m
		}(),
		want: `[pb2.ExtensionsContainer.opt_ext_bool]: true
[pb2.ExtensionsContainer.opt_ext_enum]: TEN
[pb2.ExtensionsContainer.opt_ext_nested]: {
  opt_string: "nested in an extension"
  opt_nested: {
    opt_string: "another nested in an extension"
  }
}
[pb2.ExtensionsContainer.opt_ext_string]: "extension field"
`,
	}, {
		desc: "extensions of repeated fields in another message",
		input: func() proto.Message {
			m := &pb2.Extensions{
				OptString: scalar.String("non-extension field"),
				OptBool:   scalar.Bool(true),
				OptInt32:  scalar.Int32(42),
			}
			setExtension(m, pb2.E_ExtensionsContainer_RptExtEnum, &[]pb2.Enum{pb2.Enum_TEN, 101, pb2.Enum_ONE})
			setExtension(m, pb2.E_ExtensionsContainer_RptExtString, &[]string{"hello", "world"})
			setExtension(m, pb2.E_ExtensionsContainer_RptExtNested, &[]*pb2.Nested{
				&pb2.Nested{OptString: scalar.String("one")},
				&pb2.Nested{OptString: scalar.String("two")},
				&pb2.Nested{OptString: scalar.String("three")},
			})
			return m
		}(),
		want: `opt_string: "non-extension field"
opt_bool: true
opt_int32: 42
[pb2.ExtensionsContainer.rpt_ext_enum]: TEN
[pb2.ExtensionsContainer.rpt_ext_enum]: 101
[pb2.ExtensionsContainer.rpt_ext_enum]: ONE
[pb2.ExtensionsContainer.rpt_ext_nested]: {
  opt_string: "one"
}
[pb2.ExtensionsContainer.rpt_ext_nested]: {
  opt_string: "two"
}
[pb2.ExtensionsContainer.rpt_ext_nested]: {
  opt_string: "three"
}
[pb2.ExtensionsContainer.rpt_ext_string]: "hello"
[pb2.ExtensionsContainer.rpt_ext_string]: "world"
`,
	}, {
		desc: "MessageSet",
		input: func() proto.Message {
			m := &pb2.MessageSet{}
			setExtension(m, pb2.E_MessageSetExtension_MessageSetExtension, &pb2.MessageSetExtension{
				OptString: scalar.String("a messageset extension"),
			})
			setExtension(m, pb2.E_MessageSetExtension_NotMessageSetExtension, &pb2.MessageSetExtension{
				OptString: scalar.String("not a messageset extension"),
			})
			setExtension(m, pb2.E_MessageSetExtension_ExtNested, &pb2.Nested{
				OptString: scalar.String("just a regular extension"),
			})
			return m
		}(),
		want: `[pb2.MessageSetExtension]: {
  opt_string: "a messageset extension"
}
[pb2.MessageSetExtension.ext_nested]: {
  opt_string: "just a regular extension"
}
[pb2.MessageSetExtension.not_message_set_extension]: {
  opt_string: "not a messageset extension"
}
`,
	}, {
		desc: "not real MessageSet 1",
		input: func() proto.Message {
			m := &pb2.FakeMessageSet{}
			setExtension(m, pb2.E_FakeMessageSetExtension_MessageSetExtension, &pb2.FakeMessageSetExtension{
				OptString: scalar.String("not a messageset extension"),
			})
			return m
		}(),
		want: `[pb2.FakeMessageSetExtension.message_set_extension]: {
  opt_string: "not a messageset extension"
}
`,
	}, {
		desc: "not real MessageSet 2",
		input: func() proto.Message {
			m := &pb2.MessageSet{}
			setExtension(m, pb2.E_MessageSetExtension, &pb2.FakeMessageSetExtension{
				OptString: scalar.String("another not a messageset extension"),
			})
			return m
		}(),
		want: `[pb2.message_set_extension]: {
  opt_string: "another not a messageset extension"
}
`,
	}, {
		desc: "Any not expanded",
		mo: prototext.MarshalOptions{
			Resolver: preg.NewTypes(),
		},
		input: func() proto.Message {
			m := &pb2.Nested{
				OptString: scalar.String("embedded inside Any"),
				OptNested: &pb2.Nested{
					OptString: scalar.String("inception"),
				},
			}
			b, err := proto.MarshalOptions{Deterministic: true}.Marshal(m)
			if err != nil {
				t.Fatalf("error in binary marshaling message for Any.value: %v", err)
			}
			return &anypb.Any{
				TypeUrl: "pb2.Nested",
				Value:   b,
			}
		}(),
		want: `type_url: "pb2.Nested"
value: "\n\x13embedded inside Any\x12\x0b\n\tinception"
`,
	}, {
		desc: "Any expanded",
		mo: prototext.MarshalOptions{
			Resolver: preg.NewTypes(pimpl.Export{}.MessageTypeOf(&pb2.Nested{})),
		},
		input: func() proto.Message {
			m := &pb2.Nested{
				OptString: scalar.String("embedded inside Any"),
				OptNested: &pb2.Nested{
					OptString: scalar.String("inception"),
				},
			}
			b, err := proto.MarshalOptions{Deterministic: true}.Marshal(m)
			if err != nil {
				t.Fatalf("error in binary marshaling message for Any.value: %v", err)
			}
			return &anypb.Any{
				TypeUrl: "foo/pb2.Nested",
				Value:   b,
			}
		}(),
		want: `[foo/pb2.Nested]: {
  opt_string: "embedded inside Any"
  opt_nested: {
    opt_string: "inception"
  }
}
`,
	}, {
		desc: "Any expanded with missing required",
		mo: prototext.MarshalOptions{
			Resolver: preg.NewTypes(pimpl.Export{}.MessageTypeOf(&pb2.PartialRequired{})),
		},
		input: func() proto.Message {
			m := &pb2.PartialRequired{
				OptString: scalar.String("embedded inside Any"),
			}
			b, err := proto.MarshalOptions{
				AllowPartial:  true,
				Deterministic: true,
			}.Marshal(m)
			if err != nil {
				t.Fatalf("error in binary marshaling message for Any.value: %v", err)
			}
			return &anypb.Any{
				TypeUrl: string(m.ProtoReflect().Descriptor().FullName()),
				Value:   b,
			}
		}(),
		want: `[pb2.PartialRequired]: {
  opt_string: "embedded inside Any"
}
`,
	}, {
		desc: "Any with invalid value",
		mo: prototext.MarshalOptions{
			Resolver: preg.NewTypes(pimpl.Export{}.MessageTypeOf(&pb2.Nested{})),
		},
		input: &anypb.Any{
			TypeUrl: "foo/pb2.Nested",
			Value:   []byte("\x80"),
		},
		want: `type_url: "foo/pb2.Nested"
value: "\x80"
`,
	}}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			// Use 2-space indentation on all MarshalOptions.
			tt.mo.Indent = "  "
			b, err := tt.mo.Marshal(tt.input)
			if err != nil && !tt.wantErr {
				t.Errorf("Marshal() returned error: %v\n", err)
			}
			if err == nil && tt.wantErr {
				t.Error("Marshal() got nil error, want error\n")
			}
			got := string(b)
			if tt.want != "" && got != tt.want {
				t.Errorf("Marshal()\n<got>\n%v\n<want>\n%v\n", got, tt.want)
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("Marshal() diff -want +got\n%v\n", diff)
				}
			}
		})
	}
}
