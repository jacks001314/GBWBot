// Code generated by protoc-gen-go.
// source: RefreshAuthorizationPolicyProtocol.proto
// DO NOT EDIT!

package hadoop_common

import proto "github.com/golang/protobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

// *
//  Refresh service acl request.
type RefreshServiceAclRequestProto struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *RefreshServiceAclRequestProto) Reset()         { *m = RefreshServiceAclRequestProto{} }
func (m *RefreshServiceAclRequestProto) String() string { return proto.CompactTextString(m) }
func (*RefreshServiceAclRequestProto) ProtoMessage()    {}

// *
// void response
type RefreshServiceAclResponseProto struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *RefreshServiceAclResponseProto) Reset()         { *m = RefreshServiceAclResponseProto{} }
func (m *RefreshServiceAclResponseProto) String() string { return proto.CompactTextString(m) }
func (*RefreshServiceAclResponseProto) ProtoMessage()    {}

func init() {
}