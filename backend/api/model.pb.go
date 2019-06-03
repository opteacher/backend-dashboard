// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: api/model.proto

// package 命名使用 {appid}.{version} 的方式, version 形如 v1, v2 ..

package api

import (
	context "context"
	fmt "fmt"
	_ "github.com/bilibili/kratos/tool/protobuf/pkg/extensions/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	io "io"
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
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type GrpcReqs struct {
	Version              string   `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	Method               string   `protobuf:"bytes,2,opt,name=method,proto3" json:"method,omitempty"`
	Params               [][]byte `protobuf:"bytes,3,rep,name=params,proto3" json:"params,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GrpcReqs) Reset()         { *m = GrpcReqs{} }
func (m *GrpcReqs) String() string { return proto.CompactTextString(m) }
func (*GrpcReqs) ProtoMessage()    {}
func (*GrpcReqs) Descriptor() ([]byte, []int) {
	return fileDescriptor_43c98abc783f4709, []int{0}
}
func (m *GrpcReqs) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GrpcReqs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GrpcReqs.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GrpcReqs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GrpcReqs.Merge(m, src)
}
func (m *GrpcReqs) XXX_Size() int {
	return m.Size()
}
func (m *GrpcReqs) XXX_DiscardUnknown() {
	xxx_messageInfo_GrpcReqs.DiscardUnknown(m)
}

var xxx_messageInfo_GrpcReqs proto.InternalMessageInfo

type GrpcResp struct {
	Status               int32    `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GrpcResp) Reset()         { *m = GrpcResp{} }
func (m *GrpcResp) String() string { return proto.CompactTextString(m) }
func (*GrpcResp) ProtoMessage()    {}
func (*GrpcResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_43c98abc783f4709, []int{1}
}
func (m *GrpcResp) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GrpcResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GrpcResp.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GrpcResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GrpcResp.Merge(m, src)
}
func (m *GrpcResp) XXX_Size() int {
	return m.Size()
}
func (m *GrpcResp) XXX_DiscardUnknown() {
	xxx_messageInfo_GrpcResp.DiscardUnknown(m)
}

var xxx_messageInfo_GrpcResp proto.InternalMessageInfo

func init() {
	proto.RegisterType((*GrpcReqs)(nil), "model.service.v1.GrpcReqs")
	proto.RegisterType((*GrpcResp)(nil), "model.service.v1.GrpcResp")
}

func init() { proto.RegisterFile("api/model.proto", fileDescriptor_43c98abc783f4709) }

var fileDescriptor_43c98abc783f4709 = []byte{
	// 283 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0xc1, 0x4a, 0x85, 0x40,
	0x14, 0x86, 0x9b, 0x6b, 0xdd, 0x7b, 0x9b, 0x84, 0x62, 0xa8, 0x10, 0x09, 0x11, 0x57, 0xd2, 0x42,
	0xa9, 0x76, 0x2d, 0xdb, 0xd4, 0xa2, 0x20, 0xa6, 0x5a, 0xd4, 0x6e, 0xd2, 0xc1, 0x04, 0xf5, 0x4c,
	0x9e, 0xc9, 0x07, 0xe8, 0x15, 0x7a, 0xa9, 0xbb, 0x0c, 0x7a, 0x81, 0x92, 0x1e, 0x24, 0x66, 0xd4,
	0x4d, 0x70, 0x77, 0xe7, 0x9b, 0xff, 0xcc, 0x7f, 0xf8, 0x7f, 0xba, 0x2b, 0x54, 0x99, 0xd6, 0x90,
	0xcb, 0x2a, 0x51, 0x2d, 0x68, 0x60, 0x7b, 0x03, 0xa0, 0x6c, 0xbb, 0x32, 0x93, 0x49, 0x77, 0xe2,
	0xef, 0x17, 0x50, 0x80, 0x15, 0x53, 0x33, 0x0d, 0x7b, 0xfe, 0x51, 0x01, 0x50, 0x54, 0x32, 0x35,
	0xff, 0x45, 0xd3, 0x80, 0x16, 0xba, 0x84, 0x06, 0x07, 0x35, 0xba, 0xa7, 0xcb, 0xcb, 0x56, 0x65,
	0x5c, 0xbe, 0x22, 0xf3, 0xe8, 0xa2, 0x93, 0x2d, 0x96, 0xd0, 0x78, 0x24, 0x24, 0xf1, 0x36, 0x9f,
	0x90, 0x1d, 0xd2, 0x79, 0x2d, 0xf5, 0x0b, 0xe4, 0xde, 0xcc, 0x0a, 0x23, 0x99, 0x77, 0x25, 0x5a,
	0x51, 0xa3, 0xe7, 0x84, 0x4e, 0xec, 0xf2, 0x91, 0xa2, 0xdb, 0xc9, 0x15, 0x95, 0xd9, 0x41, 0x2d,
	0xf4, 0x1b, 0x5a, 0xd3, 0x2d, 0x3e, 0x92, 0xb9, 0x56, 0x4b, 0x44, 0x51, 0xc8, 0xd1, 0x74, 0x42,
	0xc6, 0xe8, 0x66, 0x2e, 0xb4, 0xf0, 0x9c, 0x90, 0xc4, 0x2e, 0xb7, 0xf3, 0xa9, 0xa4, 0xcb, 0x1b,
	0x93, 0xf7, 0xae, 0xcb, 0xd8, 0x23, 0xdd, 0xb9, 0xca, 0x2b, 0x8b, 0x0f, 0xfc, 0x9a, 0xf9, 0xc9,
	0xff, 0x26, 0x92, 0x29, 0x92, 0xbf, 0x56, 0x43, 0x15, 0xb1, 0xf7, 0xaf, 0xdf, 0x8f, 0x99, 0x1b,
	0x2d, 0x86, 0x5a, 0xf1, 0x9c, 0x1c, 0x5f, 0x1c, 0xac, 0x7e, 0x82, 0x8d, 0x55, 0x1f, 0x90, 0xcf,
	0x3e, 0x20, 0xdf, 0x7d, 0x40, 0x9e, 0x1c, 0xa1, 0xca, 0xe7, 0xb9, 0x2d, 0xeb, 0xec, 0x2f, 0x00,
	0x00, 0xff, 0xff, 0xe9, 0x30, 0xe5, 0xe8, 0x85, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ModelSvcClient is the client API for ModelSvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ModelSvcClient interface {
	HdlModelURL(ctx context.Context, in *GrpcReqs, opts ...grpc.CallOption) (*GrpcResp, error)
}

type modelSvcClient struct {
	cc *grpc.ClientConn
}

func NewModelSvcClient(cc *grpc.ClientConn) ModelSvcClient {
	return &modelSvcClient{cc}
}

func (c *modelSvcClient) HdlModelURL(ctx context.Context, in *GrpcReqs, opts ...grpc.CallOption) (*GrpcResp, error) {
	out := new(GrpcResp)
	err := c.cc.Invoke(ctx, "/model.service.v1.ModelSvc/HdlModelURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ModelSvcServer is the server API for ModelSvc service.
type ModelSvcServer interface {
	HdlModelURL(context.Context, *GrpcReqs) (*GrpcResp, error)
}

func RegisterModelSvcServer(s *grpc.Server, srv ModelSvcServer) {
	s.RegisterService(&_ModelSvc_serviceDesc, srv)
}

func _ModelSvc_HdlModelURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GrpcReqs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ModelSvcServer).HdlModelURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/model.service.v1.ModelSvc/HdlModelURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ModelSvcServer).HdlModelURL(ctx, req.(*GrpcReqs))
	}
	return interceptor(ctx, in, info, handler)
}

var _ModelSvc_serviceDesc = grpc.ServiceDesc{
	ServiceName: "model.service.v1.ModelSvc",
	HandlerType: (*ModelSvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HdlModelURL",
			Handler:    _ModelSvc_HdlModelURL_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/model.proto",
}

func (m *GrpcReqs) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GrpcReqs) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Version) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintModel(dAtA, i, uint64(len(m.Version)))
		i += copy(dAtA[i:], m.Version)
	}
	if len(m.Method) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintModel(dAtA, i, uint64(len(m.Method)))
		i += copy(dAtA[i:], m.Method)
	}
	if len(m.Params) > 0 {
		for _, b := range m.Params {
			dAtA[i] = 0x1a
			i++
			i = encodeVarintModel(dAtA, i, uint64(len(b)))
			i += copy(dAtA[i:], b)
		}
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func (m *GrpcResp) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GrpcResp) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Status != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintModel(dAtA, i, uint64(m.Status))
	}
	if len(m.Message) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintModel(dAtA, i, uint64(len(m.Message)))
		i += copy(dAtA[i:], m.Message)
	}
	if len(m.Data) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintModel(dAtA, i, uint64(len(m.Data)))
		i += copy(dAtA[i:], m.Data)
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintModel(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *GrpcReqs) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Version)
	if l > 0 {
		n += 1 + l + sovModel(uint64(l))
	}
	l = len(m.Method)
	if l > 0 {
		n += 1 + l + sovModel(uint64(l))
	}
	if len(m.Params) > 0 {
		for _, b := range m.Params {
			l = len(b)
			n += 1 + l + sovModel(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *GrpcResp) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Status != 0 {
		n += 1 + sovModel(uint64(m.Status))
	}
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovModel(uint64(l))
	}
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovModel(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovModel(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozModel(x uint64) (n int) {
	return sovModel(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GrpcReqs) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowModel
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GrpcReqs: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GrpcReqs: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowModel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthModel
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthModel
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Version = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Method", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowModel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthModel
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthModel
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Method = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowModel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthModel
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthModel
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Params = append(m.Params, make([]byte, postIndex-iNdEx))
			copy(m.Params[len(m.Params)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipModel(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthModel
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthModel
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *GrpcResp) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowModel
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GrpcResp: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GrpcResp: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			m.Status = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowModel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Status |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowModel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthModel
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthModel
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowModel
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthModel
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthModel
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipModel(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthModel
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthModel
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipModel(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowModel
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowModel
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowModel
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthModel
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthModel
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowModel
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipModel(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthModel
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthModel = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowModel   = fmt.Errorf("proto: integer overflow")
)
