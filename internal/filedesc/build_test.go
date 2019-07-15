// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filedesc_test

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"

	testpb "google.golang.org/protobuf/internal/testprotos/test"
	_ "google.golang.org/protobuf/internal/testprotos/test/weak1"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestInit(t *testing.T) {
	// Compare the FileDescriptorProto for the same test file from two different sources:
	//
	// 1. The result of passing the filedesc-produced FileDescriptor through protodesc.
	// 2. The protoc-generated wire-encoded message.
	//
	// This serves as a test of both filedesc and protodesc.
	got := protodesc.ToFileDescriptorProto(testpb.File_test_test_proto)

	want := &descriptorpb.FileDescriptorProto{}
	zb, _ := (&testpb.TestAllTypes{}).Descriptor()
	r, _ := gzip.NewReader(bytes.NewBuffer(zb))
	b, _ := ioutil.ReadAll(r)
	if err := proto.Unmarshal(b, want); err != nil {
		t.Fatal(err)
	}

	if !proto.Equal(got, want) {
		t.Errorf("protodesc.ToFileDescriptorProto(testpb.Test_protoFile) is not equal to the protoc-generated FileDescriptorProto for internal/testprotos/test/test.proto")
	}

	// Verify that the test proto file provides exhaustive coverage of all descriptor fields.
	seen := make(map[protoreflect.FullName]bool)
	visitFields(want.ProtoReflect(), func(field protoreflect.FieldDescriptor) {
		seen[field.FullName()] = true
	})
	ignore := map[protoreflect.FullName]bool{
		// The protoreflect descriptors don't include source info.
		"google.protobuf.FileDescriptorProto.source_code_info": true,
		"google.protobuf.FileDescriptorProto.syntax":           true,

		// TODO: Test oneof and extension options. Testing these requires extending the
		// options messages (because they contain no user-settable fields), but importing
		// decriptor.proto from test.proto currently causes an import cycle. Add test
		// cases when that import cycle has been fixed.
		"google.protobuf.OneofDescriptorProto.options": true,
	}
	for _, messageName := range []protoreflect.Name{
		"FileDescriptorProto",
		"DescriptorProto",
		"FieldDescriptorProto",
		"OneofDescriptorProto",
		"EnumDescriptorProto",
		"EnumValueDescriptorProto",
		"ServiceDescriptorProto",
		"MethodDescriptorProto",
	} {
		message := descriptorpb.File_google_protobuf_descriptor_proto.Messages().ByName(messageName)
		for i, fields := 0, message.Fields(); i < fields.Len(); i++ {
			if name := fields.Get(i).FullName(); !seen[name] && !ignore[name] {
				t.Errorf("No test for descriptor field: %v", name)
			}
		}
	}

	// Verify that message descriptors for map entries have no Go type info.
	mapEntryName := protoreflect.FullName("goproto.proto.test.TestAllTypes.MapInt32Int32Entry")
	d := testpb.File_test_test_proto.Messages().ByName("TestAllTypes").Fields().ByName("map_int32_int32").Message()
	if gotName, wantName := d.FullName(), mapEntryName; gotName != wantName {
		t.Fatalf("looked up wrong descriptor: got %v, want %v", gotName, wantName)
	}
	if _, ok := d.(protoreflect.MessageType); ok {
		t.Errorf("message descriptor for %v must not implement protoreflect.MessageType", mapEntryName)
	}
}

// visitFields calls f for every field set in m and its children.
func visitFields(m protoreflect.Message, f func(protoreflect.FieldDescriptor)) {
	m.Range(func(fd protoreflect.FieldDescriptor, value protoreflect.Value) bool {
		f(fd)
		switch fd.Kind() {
		case protoreflect.MessageKind, protoreflect.GroupKind:
			if fd.IsList() {
				for i, list := 0, value.List(); i < list.Len(); i++ {
					visitFields(list.Get(i).Message(), f)
				}
			} else {
				visitFields(value.Message(), f)
			}
		}
		return true
	})
}

func TestWeakInit(t *testing.T) {
	file := testpb.File_test_test_proto

	// We do not expect to get a placeholder since weak1 is imported.
	fd1 := file.Messages().ByName("TestWeak").Fields().ByName("weak_message1")
	if got, want := fd1.IsWeak(), true; got != want {
		t.Errorf("field %v: IsWeak() = %v, want %v", fd1.FullName(), got, want)
	}
	if got, want := fd1.Message().IsPlaceholder(), false; got != want {
		t.Errorf("field %v: Message.IsPlaceholder() = %v, want %v", fd1.FullName(), got, want)
	}
	if got, want := fd1.Message().Fields().Len(), 1; got != want {
		t.Errorf("field %v: Message().Fields().Len() == %d, want %d", fd1.FullName(), got, want)
	}

	// We do expect to get a placeholder since weak2 is not imported.
	fd2 := file.Messages().ByName("TestWeak").Fields().ByName("weak_message2")
	if got, want := fd2.IsWeak(), true; got != want {
		t.Errorf("field %v: IsWeak() = %v, want %v", fd2.FullName(), got, want)
	}
	if got, want := fd2.Message().IsPlaceholder(), true; got != want {
		t.Errorf("field %v: Message.IsPlaceholder() = %v, want %v", fd2.FullName(), got, want)
	}
	if got, want := fd2.Message().Fields().Len(), 0; got != want {
		t.Errorf("field %v: Message().Fields().Len() == %d, want %d", fd2.FullName(), got, want)
	}
}
