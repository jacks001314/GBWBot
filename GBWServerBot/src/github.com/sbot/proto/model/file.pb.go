// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.12.3
// source: proto/model/file.proto

package model

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DownloadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// the node information that sends download request
	Node *Node `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
	//the file name that wants to dowload
	Fname string `protobuf:"bytes,2,opt,name=fname,proto3" json:"fname,omitempty"`
}

func (x *DownloadRequest) Reset() {
	*x = DownloadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_model_file_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownloadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadRequest) ProtoMessage() {}

func (x *DownloadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_model_file_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadRequest.ProtoReflect.Descriptor instead.
func (*DownloadRequest) Descriptor() ([]byte, []int) {
	return file_proto_model_file_proto_rawDescGZIP(), []int{0}
}

func (x *DownloadRequest) GetNode() *Node {
	if x != nil {
		return x.Node
	}
	return nil
}

func (x *DownloadRequest) GetFname() string {
	if x != nil {
		return x.Fname
	}
	return ""
}

type FilePart struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Fpath string `protobuf:"bytes,1,opt,name=fpath,proto3" json:"fpath,omitempty"`
	//the total size of file that been downloaded/uploaded
	Tbytes int64 `protobuf:"varint,2,opt,name=tbytes,proto3" json:"tbytes,omitempty"`
	//the content size of this chunk
	Bytes int64 `protobuf:"varint,3,opt,name=bytes,proto3" json:"bytes,omitempty"`
	//is last parts of file
	IsLastParts bool `protobuf:"varint,4,opt,name=isLastParts,proto3" json:"isLastParts,omitempty"`
	// the md5 that been downloaded/uploaded
	Md5 string `protobuf:"bytes,5,opt,name=md5,proto3" json:"md5,omitempty"`
	// the buffer used to store file contents
	Contents []byte `protobuf:"bytes,6,opt,name=contents,proto3" json:"contents,omitempty"`
}

func (x *FilePart) Reset() {
	*x = FilePart{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_model_file_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilePart) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilePart) ProtoMessage() {}

func (x *FilePart) ProtoReflect() protoreflect.Message {
	mi := &file_proto_model_file_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilePart.ProtoReflect.Descriptor instead.
func (*FilePart) Descriptor() ([]byte, []int) {
	return file_proto_model_file_proto_rawDescGZIP(), []int{1}
}

func (x *FilePart) GetFpath() string {
	if x != nil {
		return x.Fpath
	}
	return ""
}

func (x *FilePart) GetTbytes() int64 {
	if x != nil {
		return x.Tbytes
	}
	return 0
}

func (x *FilePart) GetBytes() int64 {
	if x != nil {
		return x.Bytes
	}
	return 0
}

func (x *FilePart) GetIsLastParts() bool {
	if x != nil {
		return x.IsLastParts
	}
	return false
}

func (x *FilePart) GetMd5() string {
	if x != nil {
		return x.Md5
	}
	return ""
}

func (x *FilePart) GetContents() []byte {
	if x != nil {
		return x.Contents
	}
	return nil
}

type UPloadStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//0 ---ok,-1 failed
	Status int64  `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Fpath  string `protobuf:"bytes,2,opt,name=fpath,proto3" json:"fpath,omitempty"`
}

func (x *UPloadStatus) Reset() {
	*x = UPloadStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_model_file_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UPloadStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UPloadStatus) ProtoMessage() {}

func (x *UPloadStatus) ProtoReflect() protoreflect.Message {
	mi := &file_proto_model_file_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UPloadStatus.ProtoReflect.Descriptor instead.
func (*UPloadStatus) Descriptor() ([]byte, []int) {
	return file_proto_model_file_proto_rawDescGZIP(), []int{2}
}

func (x *UPloadStatus) GetStatus() int64 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *UPloadStatus) GetFpath() string {
	if x != nil {
		return x.Fpath
	}
	return ""
}

var File_proto_model_file_proto protoreflect.FileDescriptor

var file_proto_model_file_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x66, 0x69,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x73, 0x62, 0x6f, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x1a, 0x16, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x53, 0x0a, 0x0f, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x73, 0x62, 0x6f, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x6e, 0x6f, 0x64,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x66, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x9e, 0x01, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65,
	0x50, 0x61, 0x72, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x66, 0x70, 0x61, 0x74, 0x68, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x62,
	0x79, 0x74, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x74, 0x62, 0x79, 0x74,
	0x65, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x05, 0x62, 0x79, 0x74, 0x65, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x69, 0x73, 0x4c, 0x61,
	0x73, 0x74, 0x50, 0x61, 0x72, 0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x69,
	0x73, 0x4c, 0x61, 0x73, 0x74, 0x50, 0x61, 0x72, 0x74, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x64,
	0x35, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x64, 0x35, 0x12, 0x1a, 0x0a, 0x08,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x3c, 0x0a, 0x0c, 0x55, 0x50, 0x6c, 0x6f,
	0x61, 0x64, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x14, 0x0a, 0x05, 0x66, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x66, 0x70, 0x61, 0x74, 0x68, 0x42, 0x6d, 0x0a, 0x1c, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x62,
	0x77, 0x33, 0x62, 0x61, 0x6f, 0x2e, 0x73, 0x62, 0x6f, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x42, 0x08, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64,
	0x50, 0x01, 0x5a, 0x1b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73,
	0x62, 0x6f, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0xaa,
	0x02, 0x10, 0x53, 0x42, 0x6f, 0x74, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x6f, 0x64,
	0x65, 0x6c, 0xca, 0x02, 0x10, 0x53, 0x42, 0x6f, 0x74, 0x5c, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x5c,
	0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_model_file_proto_rawDescOnce sync.Once
	file_proto_model_file_proto_rawDescData = file_proto_model_file_proto_rawDesc
)

func file_proto_model_file_proto_rawDescGZIP() []byte {
	file_proto_model_file_proto_rawDescOnce.Do(func() {
		file_proto_model_file_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_model_file_proto_rawDescData)
	})
	return file_proto_model_file_proto_rawDescData
}

var file_proto_model_file_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_model_file_proto_goTypes = []interface{}{
	(*DownloadRequest)(nil), // 0: sbot.proto.model.DownloadRequest
	(*FilePart)(nil),        // 1: sbot.proto.model.FilePart
	(*UPloadStatus)(nil),    // 2: sbot.proto.model.UPloadStatus
	(*Node)(nil),            // 3: sbot.proto.model.Node
}
var file_proto_model_file_proto_depIdxs = []int32{
	3, // 0: sbot.proto.model.DownloadRequest.node:type_name -> sbot.proto.model.Node
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_model_file_proto_init() }
func file_proto_model_file_proto_init() {
	if File_proto_model_file_proto != nil {
		return
	}
	file_proto_model_node_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_proto_model_file_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DownloadRequest); i {
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
		file_proto_model_file_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilePart); i {
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
		file_proto_model_file_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UPloadStatus); i {
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
			RawDescriptor: file_proto_model_file_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_model_file_proto_goTypes,
		DependencyIndexes: file_proto_model_file_proto_depIdxs,
		MessageInfos:      file_proto_model_file_proto_msgTypes,
	}.Build()
	File_proto_model_file_proto = out.File
	file_proto_model_file_proto_rawDesc = nil
	file_proto_model_file_proto_goTypes = nil
	file_proto_model_file_proto_depIdxs = nil
}
