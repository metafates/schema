// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: example.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Person_PhoneType int32

const (
	Person_MOBILE Person_PhoneType = 0
	Person_HOME   Person_PhoneType = 1
	Person_WORK   Person_PhoneType = 2
)

// Enum value maps for Person_PhoneType.
var (
	Person_PhoneType_name = map[int32]string{
		0: "MOBILE",
		1: "HOME",
		2: "WORK",
	}
	Person_PhoneType_value = map[string]int32{
		"MOBILE": 0,
		"HOME":   1,
		"WORK":   2,
	}
)

func (x Person_PhoneType) Enum() *Person_PhoneType {
	p := new(Person_PhoneType)
	*p = x
	return p
}

func (x Person_PhoneType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Person_PhoneType) Descriptor() protoreflect.EnumDescriptor {
	return file_example_proto_enumTypes[0].Descriptor()
}

func (Person_PhoneType) Type() protoreflect.EnumType {
	return &file_example_proto_enumTypes[0]
}

func (x Person_PhoneType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Person_PhoneType.Descriptor instead.
func (Person_PhoneType) EnumDescriptor() ([]byte, []int) {
	return file_example_proto_rawDescGZIP(), []int{0, 0}
}

type Person struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Id            int32                  `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"` // Unique ID number for this person.
	Email         string                 `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Phones        []*Person_PhoneNumber  `protobuf:"bytes,4,rep,name=phones,proto3" json:"phones,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Person) Reset() {
	*x = Person{}
	mi := &file_example_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Person) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Person) ProtoMessage() {}

func (x *Person) ProtoReflect() protoreflect.Message {
	mi := &file_example_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Person.ProtoReflect.Descriptor instead.
func (*Person) Descriptor() ([]byte, []int) {
	return file_example_proto_rawDescGZIP(), []int{0}
}

func (x *Person) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Person) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Person) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *Person) GetPhones() []*Person_PhoneNumber {
	if x != nil {
		return x.Phones
	}
	return nil
}

// Our address book file is just one of these.
type AddressBook struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	People        []*Person              `protobuf:"bytes,1,rep,name=people,proto3" json:"people,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddressBook) Reset() {
	*x = AddressBook{}
	mi := &file_example_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddressBook) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddressBook) ProtoMessage() {}

func (x *AddressBook) ProtoReflect() protoreflect.Message {
	mi := &file_example_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddressBook.ProtoReflect.Descriptor instead.
func (*AddressBook) Descriptor() ([]byte, []int) {
	return file_example_proto_rawDescGZIP(), []int{1}
}

func (x *AddressBook) GetPeople() []*Person {
	if x != nil {
		return x.People
	}
	return nil
}

type Person_PhoneNumber struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Number        string                 `protobuf:"bytes,1,opt,name=number,proto3" json:"number,omitempty"`
	Type          Person_PhoneType       `protobuf:"varint,2,opt,name=type,proto3,enum=pb.Person_PhoneType" json:"type,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Person_PhoneNumber) Reset() {
	*x = Person_PhoneNumber{}
	mi := &file_example_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Person_PhoneNumber) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Person_PhoneNumber) ProtoMessage() {}

func (x *Person_PhoneNumber) ProtoReflect() protoreflect.Message {
	mi := &file_example_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Person_PhoneNumber.ProtoReflect.Descriptor instead.
func (*Person_PhoneNumber) Descriptor() ([]byte, []int) {
	return file_example_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Person_PhoneNumber) GetNumber() string {
	if x != nil {
		return x.Number
	}
	return ""
}

func (x *Person_PhoneNumber) GetType() Person_PhoneType {
	if x != nil {
		return x.Type
	}
	return Person_MOBILE
}

var File_example_proto protoreflect.FileDescriptor

const file_example_proto_rawDesc = "" +
	"\n" +
	"\rexample.proto\x12\x02pb\"\xf0\x01\n" +
	"\x06Person\x12\x12\n" +
	"\x04name\x18\x01 \x01(\tR\x04name\x12\x0e\n" +
	"\x02id\x18\x02 \x01(\x05R\x02id\x12\x14\n" +
	"\x05email\x18\x03 \x01(\tR\x05email\x12.\n" +
	"\x06phones\x18\x04 \x03(\v2\x16.pb.Person.PhoneNumberR\x06phones\x1aO\n" +
	"\vPhoneNumber\x12\x16\n" +
	"\x06number\x18\x01 \x01(\tR\x06number\x12(\n" +
	"\x04type\x18\x02 \x01(\x0e2\x14.pb.Person.PhoneTypeR\x04type\"+\n" +
	"\tPhoneType\x12\n" +
	"\n" +
	"\x06MOBILE\x10\x00\x12\b\n" +
	"\x04HOME\x10\x01\x12\b\n" +
	"\x04WORK\x10\x02\"1\n" +
	"\vAddressBook\x12\"\n" +
	"\x06people\x18\x01 \x03(\v2\n" +
	".pb.PersonR\x06peopleB4Z2github.com/metafates/schame/examples/parse-grpc/pbb\x06proto3"

var (
	file_example_proto_rawDescOnce sync.Once
	file_example_proto_rawDescData []byte
)

func file_example_proto_rawDescGZIP() []byte {
	file_example_proto_rawDescOnce.Do(func() {
		file_example_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_example_proto_rawDesc), len(file_example_proto_rawDesc)))
	})
	return file_example_proto_rawDescData
}

var file_example_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_example_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_example_proto_goTypes = []any{
	(Person_PhoneType)(0),      // 0: pb.Person.PhoneType
	(*Person)(nil),             // 1: pb.Person
	(*AddressBook)(nil),        // 2: pb.AddressBook
	(*Person_PhoneNumber)(nil), // 3: pb.Person.PhoneNumber
}
var file_example_proto_depIdxs = []int32{
	3, // 0: pb.Person.phones:type_name -> pb.Person.PhoneNumber
	1, // 1: pb.AddressBook.people:type_name -> pb.Person
	0, // 2: pb.Person.PhoneNumber.type:type_name -> pb.Person.PhoneType
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_example_proto_init() }
func file_example_proto_init() {
	if File_example_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_example_proto_rawDesc), len(file_example_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_example_proto_goTypes,
		DependencyIndexes: file_example_proto_depIdxs,
		EnumInfos:         file_example_proto_enumTypes,
		MessageInfos:      file_example_proto_msgTypes,
	}.Build()
	File_example_proto = out.File
	file_example_proto_goTypes = nil
	file_example_proto_depIdxs = nil
}
