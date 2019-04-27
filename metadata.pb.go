// Code generated by protoc-gen-go. DO NOT EDIT.
// source: metadata.proto

package teflon

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type PersistentMeta struct {
	Proto                string            `protobuf:"bytes,1,opt,name=Proto,proto3" json:"Proto,omitempty"`
	Instances            []string          `protobuf:"bytes,3,rep,name=Instances,proto3" json:"Instances,omitempty"`
	UserData             map[string]string `protobuf:"bytes,2,rep,name=UserData,proto3" json:"UserData,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ImgInfo              *ImgInfo          `protobuf:"bytes,4,opt,name=ImgInfo,proto3" json:"ImgInfo,omitempty"`
	Seq                  *Seq              `protobuf:"bytes,5,opt,name=Seq,proto3" json:"Seq,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *PersistentMeta) Reset()         { *m = PersistentMeta{} }
func (m *PersistentMeta) String() string { return proto.CompactTextString(m) }
func (*PersistentMeta) ProtoMessage()    {}
func (*PersistentMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_56d9f74966f40d04, []int{0}
}

func (m *PersistentMeta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PersistentMeta.Unmarshal(m, b)
}
func (m *PersistentMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PersistentMeta.Marshal(b, m, deterministic)
}
func (m *PersistentMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PersistentMeta.Merge(m, src)
}
func (m *PersistentMeta) XXX_Size() int {
	return xxx_messageInfo_PersistentMeta.Size(m)
}
func (m *PersistentMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_PersistentMeta.DiscardUnknown(m)
}

var xxx_messageInfo_PersistentMeta proto.InternalMessageInfo

func (m *PersistentMeta) GetProto() string {
	if m != nil {
		return m.Proto
	}
	return ""
}

func (m *PersistentMeta) GetInstances() []string {
	if m != nil {
		return m.Instances
	}
	return nil
}

func (m *PersistentMeta) GetUserData() map[string]string {
	if m != nil {
		return m.UserData
	}
	return nil
}

func (m *PersistentMeta) GetImgInfo() *ImgInfo {
	if m != nil {
		return m.ImgInfo
	}
	return nil
}

func (m *PersistentMeta) GetSeq() *Seq {
	if m != nil {
		return m.Seq
	}
	return nil
}

type ImgInfo struct {
	Width                int32    `protobuf:"varint,1,opt,name=Width,proto3" json:"Width,omitempty"`
	Height               int32    `protobuf:"varint,2,opt,name=Height,proto3" json:"Height,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ImgInfo) Reset()         { *m = ImgInfo{} }
func (m *ImgInfo) String() string { return proto.CompactTextString(m) }
func (*ImgInfo) ProtoMessage()    {}
func (*ImgInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_56d9f74966f40d04, []int{1}
}

func (m *ImgInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ImgInfo.Unmarshal(m, b)
}
func (m *ImgInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ImgInfo.Marshal(b, m, deterministic)
}
func (m *ImgInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImgInfo.Merge(m, src)
}
func (m *ImgInfo) XXX_Size() int {
	return xxx_messageInfo_ImgInfo.Size(m)
}
func (m *ImgInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ImgInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ImgInfo proto.InternalMessageInfo

func (m *ImgInfo) GetWidth() int32 {
	if m != nil {
		return m.Width
	}
	return 0
}

func (m *ImgInfo) GetHeight() int32 {
	if m != nil {
		return m.Height
	}
	return 0
}

type Seq struct {
	BaseName             string   `protobuf:"bytes,1,opt,name=BaseName,proto3" json:"BaseName,omitempty"`
	First                int32    `protobuf:"varint,2,opt,name=First,proto3" json:"First,omitempty"`
	Last                 int32    `protobuf:"varint,3,opt,name=Last,proto3" json:"Last,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Seq) Reset()         { *m = Seq{} }
func (m *Seq) String() string { return proto.CompactTextString(m) }
func (*Seq) ProtoMessage()    {}
func (*Seq) Descriptor() ([]byte, []int) {
	return fileDescriptor_56d9f74966f40d04, []int{2}
}

func (m *Seq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Seq.Unmarshal(m, b)
}
func (m *Seq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Seq.Marshal(b, m, deterministic)
}
func (m *Seq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Seq.Merge(m, src)
}
func (m *Seq) XXX_Size() int {
	return xxx_messageInfo_Seq.Size(m)
}
func (m *Seq) XXX_DiscardUnknown() {
	xxx_messageInfo_Seq.DiscardUnknown(m)
}

var xxx_messageInfo_Seq proto.InternalMessageInfo

func (m *Seq) GetBaseName() string {
	if m != nil {
		return m.BaseName
	}
	return ""
}

func (m *Seq) GetFirst() int32 {
	if m != nil {
		return m.First
	}
	return 0
}

func (m *Seq) GetLast() int32 {
	if m != nil {
		return m.Last
	}
	return 0
}

func init() {
	proto.RegisterType((*PersistentMeta)(nil), "teflon.PersistentMeta")
	proto.RegisterMapType((map[string]string)(nil), "teflon.PersistentMeta.UserDataEntry")
	proto.RegisterType((*ImgInfo)(nil), "teflon.ImgInfo")
	proto.RegisterType((*Seq)(nil), "teflon.Seq")
}

func init() { proto.RegisterFile("metadata.proto", fileDescriptor_56d9f74966f40d04) }

var fileDescriptor_56d9f74966f40d04 = []byte{
	// 293 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x91, 0x4b, 0x4b, 0xc3, 0x40,
	0x14, 0x85, 0x49, 0xa6, 0xa9, 0xed, 0x2d, 0x56, 0x19, 0x44, 0x86, 0xa2, 0x10, 0x82, 0x8b, 0xb8,
	0xc9, 0xa2, 0x2e, 0x14, 0xdd, 0x88, 0xa8, 0x18, 0x7c, 0x50, 0xa6, 0x88, 0xeb, 0xd1, 0xde, 0xb6,
	0xc1, 0x66, 0x62, 0x33, 0x57, 0xa1, 0x3f, 0xc2, 0xff, 0x2c, 0x33, 0x79, 0x48, 0x77, 0xf7, 0x3b,
	0x73, 0xce, 0x7d, 0x30, 0x30, 0xcc, 0x91, 0xd4, 0x4c, 0x91, 0x4a, 0xbe, 0xca, 0x82, 0x0a, 0xde,
	0x25, 0x9c, 0xaf, 0x0a, 0x1d, 0xfd, 0xfa, 0x30, 0x9c, 0x60, 0x69, 0x32, 0x43, 0xa8, 0xe9, 0x19,
	0x49, 0xf1, 0x03, 0x08, 0x26, 0xd6, 0x23, 0xbc, 0xd0, 0x8b, 0xfb, 0xb2, 0x02, 0x7e, 0x04, 0xfd,
	0x54, 0x1b, 0x52, 0xfa, 0x03, 0x8d, 0x60, 0x21, 0x8b, 0xfb, 0xf2, 0x5f, 0xe0, 0xd7, 0xd0, 0x7b,
	0x35, 0x58, 0xde, 0x2a, 0x52, 0xc2, 0x0f, 0x59, 0x3c, 0x18, 0x9f, 0x24, 0xd5, 0x84, 0x64, 0xbb,
	0x7b, 0xd2, 0xd8, 0xee, 0x34, 0x95, 0x1b, 0xd9, 0xa6, 0xf8, 0x29, 0xec, 0xa4, 0xf9, 0x22, 0xd5,
	0xf3, 0x42, 0x74, 0x42, 0x2f, 0x1e, 0x8c, 0xf7, 0x9a, 0x06, 0xb5, 0x2c, 0x9b, 0x77, 0x7e, 0x0c,
	0x6c, 0x8a, 0x6b, 0x11, 0x38, 0xdb, 0xa0, 0xb1, 0x4d, 0x71, 0x2d, 0xad, 0x3e, 0xba, 0x82, 0xdd,
	0xad, 0x21, 0x7c, 0x1f, 0xd8, 0x27, 0x6e, 0xea, 0x73, 0x6c, 0x69, 0x4f, 0xfc, 0x51, 0xab, 0x6f,
	0x14, 0x7e, 0x75, 0xa2, 0x83, 0x4b, 0xff, 0xc2, 0x8b, 0xce, 0xdb, 0x35, 0xac, 0xe9, 0x2d, 0x9b,
	0xd1, 0xd2, 0x05, 0x03, 0x59, 0x01, 0x3f, 0x84, 0xee, 0x03, 0x66, 0x8b, 0x25, 0xb9, 0x6c, 0x20,
	0x6b, 0x8a, 0x1e, 0xdd, 0x52, 0x7c, 0x04, 0xbd, 0x1b, 0x65, 0xf0, 0x45, 0xe5, 0x58, 0x0f, 0x6c,
	0xd9, 0x36, 0xbc, 0xcf, 0x4a, 0xd3, 0x24, 0x2b, 0xe0, 0x1c, 0x3a, 0x4f, 0xca, 0x90, 0x60, 0x4e,
	0x74, 0xf5, 0x7b, 0xd7, 0x7d, 0xd2, 0xd9, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf2, 0x7b, 0x4d,
	0x8c, 0xb6, 0x01, 0x00, 0x00,
}