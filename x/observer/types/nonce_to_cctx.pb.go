// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: observer/nonce_to_cctx.proto

package types

import (
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"

	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/zeta-chain/zetacore/common"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// store key is tss+chainid+nonce
type NonceToCctx struct {
	ChainId   int64  `protobuf:"varint,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	Nonce     int64  `protobuf:"varint,2,opt,name=nonce,proto3" json:"nonce,omitempty"`
	CctxIndex string `protobuf:"bytes,3,opt,name=cctxIndex,proto3" json:"cctxIndex,omitempty"`
	Tss       string `protobuf:"bytes,4,opt,name=tss,proto3" json:"tss,omitempty"`
}

func (m *NonceToCctx) Reset()         { *m = NonceToCctx{} }
func (m *NonceToCctx) String() string { return proto.CompactTextString(m) }
func (*NonceToCctx) ProtoMessage()    {}
func (*NonceToCctx) Descriptor() ([]byte, []int) {
	return fileDescriptor_6f9bbe8a689fa6e4, []int{0}
}
func (m *NonceToCctx) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NonceToCctx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NonceToCctx.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *NonceToCctx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NonceToCctx.Merge(m, src)
}
func (m *NonceToCctx) XXX_Size() int {
	return m.Size()
}
func (m *NonceToCctx) XXX_DiscardUnknown() {
	xxx_messageInfo_NonceToCctx.DiscardUnknown(m)
}

var xxx_messageInfo_NonceToCctx proto.InternalMessageInfo

func (m *NonceToCctx) GetChainId() int64 {
	if m != nil {
		return m.ChainId
	}
	return 0
}

func (m *NonceToCctx) GetNonce() int64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *NonceToCctx) GetCctxIndex() string {
	if m != nil {
		return m.CctxIndex
	}
	return ""
}

func (m *NonceToCctx) GetTss() string {
	if m != nil {
		return m.Tss
	}
	return ""
}

func init() {
	proto.RegisterType((*NonceToCctx)(nil), "zetachain.zetacore.observer.NonceToCctx")
}

func init() { proto.RegisterFile("observer/nonce_to_cctx.proto", fileDescriptor_6f9bbe8a689fa6e4) }

var fileDescriptor_6f9bbe8a689fa6e4 = []byte{
	// 242 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xc9, 0x4f, 0x2a, 0x4e,
	0x2d, 0x2a, 0x4b, 0x2d, 0xd2, 0xcf, 0xcb, 0xcf, 0x4b, 0x4e, 0x8d, 0x2f, 0xc9, 0x8f, 0x4f, 0x4e,
	0x2e, 0xa9, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x92, 0xae, 0x4a, 0x2d, 0x49, 0x4c, 0xce,
	0x48, 0xcc, 0xcc, 0xd3, 0x03, 0xb3, 0xf2, 0x8b, 0x52, 0xf5, 0x60, 0x1a, 0xa4, 0x84, 0x93, 0xf3,
	0x73, 0x73, 0xf3, 0xf3, 0xf4, 0x21, 0x14, 0x44, 0x87, 0x94, 0x48, 0x7a, 0x7e, 0x7a, 0x3e, 0x98,
	0xa9, 0x0f, 0x62, 0x41, 0x44, 0x95, 0xf2, 0xb8, 0xb8, 0xfd, 0x40, 0xc6, 0x87, 0xe4, 0x3b, 0x27,
	0x97, 0x54, 0x08, 0x49, 0x72, 0x71, 0x80, 0x0d, 0x8d, 0xcf, 0x4c, 0x91, 0x60, 0x54, 0x60, 0xd4,
	0x60, 0x0e, 0x62, 0x07, 0xf3, 0x3d, 0x53, 0x84, 0x44, 0xb8, 0x58, 0xc1, 0x0e, 0x91, 0x60, 0x02,
	0x8b, 0x43, 0x38, 0x42, 0x32, 0x5c, 0x9c, 0x20, 0x57, 0x79, 0xe6, 0xa5, 0xa4, 0x56, 0x48, 0x30,
	0x2b, 0x30, 0x6a, 0x70, 0x06, 0x21, 0x04, 0x84, 0x04, 0xb8, 0x98, 0x4b, 0x8a, 0x8b, 0x25, 0x58,
	0xc0, 0xe2, 0x20, 0xa6, 0x93, 0xe7, 0x89, 0x47, 0x72, 0x8c, 0x17, 0x1e, 0xc9, 0x31, 0x3e, 0x78,
	0x24, 0xc7, 0x38, 0xe1, 0xb1, 0x1c, 0xc3, 0x85, 0xc7, 0x72, 0x0c, 0x37, 0x1e, 0xcb, 0x31, 0x44,
	0xe9, 0xa7, 0x67, 0x96, 0x64, 0x94, 0x26, 0xe9, 0x25, 0xe7, 0xe7, 0xea, 0x83, 0xbc, 0xa4, 0x0b,
	0xb6, 0x58, 0x1f, 0xe6, 0x3b, 0xfd, 0x0a, 0x7d, 0x78, 0x80, 0x94, 0x54, 0x16, 0xa4, 0x16, 0x27,
	0xb1, 0x81, 0x7d, 0x60, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x18, 0xf1, 0x0b, 0xf9, 0x29, 0x01,
	0x00, 0x00,
}

func (m *NonceToCctx) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NonceToCctx) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *NonceToCctx) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Tss) > 0 {
		i -= len(m.Tss)
		copy(dAtA[i:], m.Tss)
		i = encodeVarintNonceToCctx(dAtA, i, uint64(len(m.Tss)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.CctxIndex) > 0 {
		i -= len(m.CctxIndex)
		copy(dAtA[i:], m.CctxIndex)
		i = encodeVarintNonceToCctx(dAtA, i, uint64(len(m.CctxIndex)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Nonce != 0 {
		i = encodeVarintNonceToCctx(dAtA, i, uint64(m.Nonce))
		i--
		dAtA[i] = 0x10
	}
	if m.ChainId != 0 {
		i = encodeVarintNonceToCctx(dAtA, i, uint64(m.ChainId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintNonceToCctx(dAtA []byte, offset int, v uint64) int {
	offset -= sovNonceToCctx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *NonceToCctx) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ChainId != 0 {
		n += 1 + sovNonceToCctx(uint64(m.ChainId))
	}
	if m.Nonce != 0 {
		n += 1 + sovNonceToCctx(uint64(m.Nonce))
	}
	l = len(m.CctxIndex)
	if l > 0 {
		n += 1 + l + sovNonceToCctx(uint64(l))
	}
	l = len(m.Tss)
	if l > 0 {
		n += 1 + l + sovNonceToCctx(uint64(l))
	}
	return n
}

func sovNonceToCctx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozNonceToCctx(x uint64) (n int) {
	return sovNonceToCctx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *NonceToCctx) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNonceToCctx
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
			return fmt.Errorf("proto: NonceToCctx: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: NonceToCctx: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainId", wireType)
			}
			m.ChainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNonceToCctx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ChainId |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nonce", wireType)
			}
			m.Nonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNonceToCctx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Nonce |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CctxIndex", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNonceToCctx
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
				return ErrInvalidLengthNonceToCctx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthNonceToCctx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.CctxIndex = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tss", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNonceToCctx
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
				return ErrInvalidLengthNonceToCctx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthNonceToCctx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Tss = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipNonceToCctx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthNonceToCctx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipNonceToCctx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowNonceToCctx
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
					return 0, ErrIntOverflowNonceToCctx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowNonceToCctx
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
				return 0, ErrInvalidLengthNonceToCctx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupNonceToCctx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthNonceToCctx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthNonceToCctx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowNonceToCctx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupNonceToCctx = fmt.Errorf("proto: unexpected end of group")
)
