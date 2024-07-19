/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package proxyprotocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"strconv"
)

// Protocol Headers
var (
	SIGV1 = []byte{'\x50', '\x52', '\x4F', '\x58', '\x59'}
	SIGV2 = []byte{'\x0D', '\x0A', '\x0D', '\x0A', '\x00', '\x0D', '\x0A', '\x51', '\x55', '\x49', '\x54', '\x0A'}
)

// Errors
var (
	ErrCantReadVersion1Header               = errors.New("proxyproto: can't read version 1 header")
	ErrVersion1HeaderTooLong                = errors.New("proxyproto: version 1 header must be 107 bytes or less")
	ErrLineMustEndWithCrlf                  = errors.New("proxyproto: version 1 header is invalid, must end with \\r\\n")
	ErrCantReadProtocolVersionAndCommand    = errors.New("proxyproto: can't read proxy protocol version and command")
	ErrCantReadAddressFamilyAndProtocol     = errors.New("proxyproto: can't read address family or protocol")
	ErrCantReadLength                       = errors.New("proxyproto: can't read length")
	ErrCantResolveSourceUnixAddress         = errors.New("proxyproto: can't resolve source Unix address")
	ErrCantResolveDestinationUnixAddress    = errors.New("proxyproto: can't resolve destination Unix address")
	ErrNoProxyProtocol                      = errors.New("proxyproto: proxy protocol signature not present")
	ErrUnknownProxyProtocolVersion          = errors.New("proxyproto: unknown proxy protocol version")
	ErrUnsupportedProtocolVersionAndCommand = errors.New("proxyproto: unsupported proxy protocol version and command")
	ErrUnsupportedAddressFamilyAndProtocol  = errors.New("proxyproto: unsupported address family and protocol")
	ErrInvalidLength                        = errors.New("proxyproto: invalid length")
	ErrInvalidAddress                       = errors.New("proxyproto: invalid address")
	ErrInvalidPortNumber                    = errors.New("proxyproto: invalid port number")
	ErrSuperfluousProxyHeader               = errors.New("proxyproto: upstream connection sent PROXY header but isn't allowed to send one")
)

// Header is the placeholder for proxy protocol header.
type Header struct {
	Version           byte
	Command           ProtocolVersionAndCommand
	TransportProtocol AddressFamilyAndProtocol
	SourceAddr        net.Addr
	DestinationAddr   net.Addr
	rawTLVs           []byte
}

// TCPAddrs returns the tcp addresses for the proxy protocol header.
func (header *Header) TCPAddrs() (sourceAddr, destAddr *net.TCPAddr, ok bool) {
	if !header.TransportProtocol.IsStream() {
		return nil, nil, false
	}
	sourceAddr, sourceOK := header.SourceAddr.(*net.TCPAddr)
	destAddr, destOK := header.DestinationAddr.(*net.TCPAddr)
	return sourceAddr, destAddr, sourceOK && destOK
}

// UDPAddrs returns the udp addresses for the proxy protocol header.
func (header *Header) UDPAddrs() (sourceAddr, destAddr *net.UDPAddr, ok bool) {
	if !header.TransportProtocol.IsDatagram() {
		return nil, nil, false
	}
	sourceAddr, sourceOK := header.SourceAddr.(*net.UDPAddr)
	destAddr, destOK := header.DestinationAddr.(*net.UDPAddr)
	return sourceAddr, destAddr, sourceOK && destOK
}

// UnixAddrs returns the uds addresses for the proxy protocol header.
func (header *Header) UnixAddrs() (sourceAddr, destAddr *net.UnixAddr, ok bool) {
	if !header.TransportProtocol.IsUnix() {
		return nil, nil, false
	}
	sourceAddr, sourceOK := header.SourceAddr.(*net.UnixAddr)
	destAddr, destOK := header.DestinationAddr.(*net.UnixAddr)
	return sourceAddr, destAddr, sourceOK && destOK
}

// IPs returns the ip addresses for the proxy protocol header.
func (header *Header) IPs() (sourceIP, destIP net.IP, ok bool) {
	if sourceAddr, destAddr, ok := header.TCPAddrs(); ok {
		return sourceAddr.IP, destAddr.IP, true
	} else if sourceAddr, destAddr, ok := header.UDPAddrs(); ok {
		return sourceAddr.IP, destAddr.IP, true
	} else {
		return nil, nil, false
	}
}

// Ports returns the ports for the proxy protocol header.
func (header *Header) Ports() (sourcePort, destPort int, ok bool) {
	if sourceAddr, destAddr, ok := header.TCPAddrs(); ok {
		return sourceAddr.Port, destAddr.Port, true
	} else if sourceAddr, destAddr, ok := header.UDPAddrs(); ok {
		return sourceAddr.Port, destAddr.Port, true
	} else {
		return 0, 0, false
	}
}

// EqualTo returns true if headers are equivalent, false otherwise.
// Deprecated: use EqualsTo instead. This method will eventually be removed.
func (header *Header) EqualTo(otherHeader *Header) bool {
	return header.EqualsTo(otherHeader)
}

// EqualsTo returns true if headers are equivalent, false otherwise.
func (header *Header) EqualsTo(otherHeader *Header) bool {
	if otherHeader == nil {
		return false
	}
	// TLVs only exist for version 2
	if header.Version == 2 && !bytes.Equal(header.rawTLVs, otherHeader.rawTLVs) {
		return false
	}
	if header.Version != otherHeader.Version || header.Command != otherHeader.Command || header.TransportProtocol != otherHeader.TransportProtocol {
		return false
	}
	// Return early for header with ProtocolVersionAndCommandLocal command, which contains no address information
	if header.Command == ProtocolVersionAndCommandLocal {
		return true
	}
	return header.SourceAddr.String() == otherHeader.SourceAddr.String() &&
		header.DestinationAddr.String() == otherHeader.DestinationAddr.String()
}

// WriteTo renders a proxy protocol header in a format and writes it to an io.Writer.
func (header *Header) WriteTo(w io.Writer) (int64, error) {
	buf, err := header.Format()
	if err != nil {
		return 0, err
	}
	return bytes.NewBuffer(buf).WriteTo(w)
}

// Format renders a proxy protocol header in a format to write over the wire.
func (header *Header) Format() ([]byte, error) {
	switch header.Version {
	case 1:
		return header.formatVersion1()
	case 2:
		return header.formatVersion2()
	default:
		return nil, ErrUnknownProxyProtocolVersion
	}
}

// TLVs returns the TLVs stored into this header, if they exist.  TLVs are optional for v2 of the protocol.
func (header *Header) TLVs() ([]TLV, error) {
	return SplitTLVs(header.rawTLVs)
}

// SetTLVs sets the TLVs stored in this header. This method replaces any
// previous TLV.
func (header *Header) SetTLVs(tlvs []TLV) error {
	raw, err := JoinTLVs(tlvs)
	if err != nil {
		return err
	}
	header.rawTLVs = raw
	return nil
}

const (
	crlf      = "\r\n"
	separator = " "
)

func (header *Header) formatVersion1() ([]byte, error) {
	// As of version 1, only "TCP4" ( \x54 \x43 \x50 \x34 ) for TCP over IPv4,
	// and "TCP6" ( \x54 \x43 \x50 \x36 ) for TCP over IPv6 are allowed.
	var proto string
	switch header.TransportProtocol {
	case AddressFamilyAndProtocolTCPv4:
		proto = "TCP4"
	case AddressFamilyAndProtocolTCPv6:
		proto = "TCP6"
	default:
		// Unknown connection (short form)
		return []byte("PROXY UNKNOWN" + crlf), nil
	}

	sourceAddr, sourceOK := header.SourceAddr.(*net.TCPAddr)
	destAddr, destOK := header.DestinationAddr.(*net.TCPAddr)
	if !sourceOK || !destOK {
		return nil, ErrInvalidAddress
	}

	sourceIP, destIP := sourceAddr.IP, destAddr.IP
	switch header.TransportProtocol {
	case AddressFamilyAndProtocolTCPv4:
		sourceIP = sourceIP.To4()
		destIP = destIP.To4()
	case AddressFamilyAndProtocolTCPv6:
		sourceIP = sourceIP.To16()
		destIP = destIP.To16()
	}
	if sourceIP == nil || destIP == nil {
		return nil, ErrInvalidAddress
	}

	buf := bytes.NewBuffer(make([]byte, 0, 108))
	buf.Write(SIGV1)
	buf.WriteString(separator)
	buf.WriteString(proto)
	buf.WriteString(separator)
	buf.WriteString(sourceIP.String())
	buf.WriteString(separator)
	buf.WriteString(destIP.String())
	buf.WriteString(separator)
	buf.WriteString(strconv.Itoa(sourceAddr.Port))
	buf.WriteString(separator)
	buf.WriteString(strconv.Itoa(destAddr.Port))
	buf.WriteString(crlf)

	return buf.Bytes(), nil
}

var (
	lengthUnspec      = uint16(0)
	lengthV4          = uint16(12)
	lengthV6          = uint16(36)
	lengthUnix        = uint16(216)
	lengthUnspecBytes = func() []byte {
		a := make([]byte, 2)
		binary.BigEndian.PutUint16(a, lengthUnspec)
		return a
	}()
	lengthV4Bytes = func() []byte {
		a := make([]byte, 2)
		binary.BigEndian.PutUint16(a, lengthV4)
		return a
	}()
	lengthV6Bytes = func() []byte {
		a := make([]byte, 2)
		binary.BigEndian.PutUint16(a, lengthV6)
		return a
	}()
	lengthUnixBytes = func() []byte {
		a := make([]byte, 2)
		binary.BigEndian.PutUint16(a, lengthUnix)
		return a
	}()
	errUint16Overflow = errors.New("proxyproto: uint16 overflow")
)

func (header *Header) formatVersion2() ([]byte, error) {
	var buf bytes.Buffer
	buf.Write(SIGV2)
	buf.WriteByte(header.Command.toByte())
	buf.WriteByte(header.TransportProtocol.toByte())
	if header.TransportProtocol.IsUnspec() {
		// For UNSPEC, write no addresses and ports but only TLVs if they are present
		hdrLen, err := addTLVLen(lengthUnspecBytes, len(header.rawTLVs))
		if err != nil {
			return nil, err
		}
		buf.Write(hdrLen)
	} else {
		var addrSrc, addrDst []byte
		if header.TransportProtocol.IsIPv4() {
			hdrLen, err := addTLVLen(lengthV4Bytes, len(header.rawTLVs))
			if err != nil {
				return nil, err
			}
			buf.Write(hdrLen)
			sourceIP, destIP, _ := header.IPs()
			addrSrc = sourceIP.To4()
			addrDst = destIP.To4()
		} else if header.TransportProtocol.IsIPv6() {
			hdrLen, err := addTLVLen(lengthV6Bytes, len(header.rawTLVs))
			if err != nil {
				return nil, err
			}
			buf.Write(hdrLen)
			sourceIP, destIP, _ := header.IPs()
			addrSrc = sourceIP.To16()
			addrDst = destIP.To16()
		} else if header.TransportProtocol.IsUnix() {
			buf.Write(lengthUnixBytes)
			sourceAddr, destAddr, ok := header.UnixAddrs()
			if !ok {
				return nil, ErrInvalidAddress
			}
			addrSrc = formatUnixName(sourceAddr.Name)
			addrDst = formatUnixName(destAddr.Name)
		}

		if addrSrc == nil || addrDst == nil {
			return nil, ErrInvalidAddress
		}
		buf.Write(addrSrc)
		buf.Write(addrDst)

		if sourcePort, destPort, ok := header.Ports(); ok {
			portBytes := make([]byte, 2)

			binary.BigEndian.PutUint16(portBytes, uint16(sourcePort))
			buf.Write(portBytes)

			binary.BigEndian.PutUint16(portBytes, uint16(destPort))
			buf.Write(portBytes)
		}
	}

	if len(header.rawTLVs) > 0 {
		buf.Write(header.rawTLVs)
	}

	return buf.Bytes(), nil
}

// ProtocolVersionAndCommand represents the command in proxy protocol v2.
// Command doesn't exist in v1 but it should be set since other parts of
// this library may rely on it for determining connection details.
type ProtocolVersionAndCommand byte

const (
	// ProtocolVersionAndCommandLocal represents the ProtocolVersionAndCommandLocal command in v2 or UNKNOWN transport in v1,
	// in which case no address information is expected.
	ProtocolVersionAndCommandLocal ProtocolVersionAndCommand = '\x20'
	// ProtocolVersionAndCommandProxy represents the PROXY command in v2 or transport is not UNKNOWN in v1,
	// in which case valid local/remote address and port information is expected.
	ProtocolVersionAndCommandProxy ProtocolVersionAndCommand = '\x21'
)

var supportedCommand = map[ProtocolVersionAndCommand]bool{
	ProtocolVersionAndCommandLocal: true,
	ProtocolVersionAndCommandProxy: true,
}

// IsLocal returns true if the command in v2 is ProtocolVersionAndCommandLocal or the transport in v1 is UNKNOWN,
// i.e. when no address information is expected, false otherwise.
func (pvc ProtocolVersionAndCommand) IsLocal() bool {
	return ProtocolVersionAndCommandLocal == pvc
}

// IsProxy returns true if the command in v2 is PROXY or the transport in v1 is not UNKNOWN,
// i.e. when valid local/remote address and port information is expected, false otherwise.
func (pvc ProtocolVersionAndCommand) IsProxy() bool {
	return ProtocolVersionAndCommandProxy == pvc
}

// IsUnspec returns true if the command is unspecified, false otherwise.
func (pvc ProtocolVersionAndCommand) IsUnspec() bool {
	return !(pvc.IsLocal() || pvc.IsProxy())
}

func (pvc ProtocolVersionAndCommand) toByte() byte {
	if pvc.IsLocal() {
		return byte(ProtocolVersionAndCommandLocal)
	} else if pvc.IsProxy() {
		return byte(ProtocolVersionAndCommandProxy)
	}

	return byte(ProtocolVersionAndCommandLocal)
}

// AddressFamilyAndProtocol represents address family and transport protocol.
type AddressFamilyAndProtocol byte

// Address family and protocol constants
const (
	AddressFamilyAndProtocolUnknown      AddressFamilyAndProtocol = '\x00'
	AddressFamilyAndProtocolTCPv4        AddressFamilyAndProtocol = '\x11'
	AddressFamilyAndProtocolUDPv4        AddressFamilyAndProtocol = '\x12'
	AddressFamilyAndProtocolTCPv6        AddressFamilyAndProtocol = '\x21'
	AddressFamilyAndProtocolUDPv6        AddressFamilyAndProtocol = '\x22'
	AddressFamilyAndProtocolUnixStream   AddressFamilyAndProtocol = '\x31'
	AddressFamilyAndProtocolUnixDatagram AddressFamilyAndProtocol = '\x32'
)

// IsIPv4 returns true if the address family is IPv4 (AF_INET4), false otherwise.
func (ap AddressFamilyAndProtocol) IsIPv4() bool {
	return 0x10 == ap&0xF0
}

// IsIPv6 returns true if the address family is IPv6 (AF_INET6), false otherwise.
func (ap AddressFamilyAndProtocol) IsIPv6() bool {
	return 0x20 == ap&0xF0
}

// IsUnix returns true if the address family is UNIX (AF_UNIX), false otherwise.
func (ap AddressFamilyAndProtocol) IsUnix() bool {
	return 0x30 == ap&0xF0
}

// IsStream returns true if the transport protocol is TCP or STREAM (SOCK_STREAM), false otherwise.
func (ap AddressFamilyAndProtocol) IsStream() bool {
	return 0x01 == ap&0x0F
}

// IsDatagram returns true if the transport protocol is UDP or DGRAM (SOCK_DGRAM), false otherwise.
func (ap AddressFamilyAndProtocol) IsDatagram() bool {
	return 0x02 == ap&0x0F
}

// IsUnspec returns true if the transport protocol or address family is unspecified, false otherwise.
func (ap AddressFamilyAndProtocol) IsUnspec() bool {
	return (0x00 == ap&0xF0) || (0x00 == ap&0x0F)
}

func (ap AddressFamilyAndProtocol) toByte() byte {
	if ap.IsIPv4() && ap.IsStream() {
		return byte(AddressFamilyAndProtocolTCPv4)
	} else if ap.IsIPv4() && ap.IsDatagram() {
		return byte(AddressFamilyAndProtocolUDPv4)
	} else if ap.IsIPv6() && ap.IsStream() {
		return byte(AddressFamilyAndProtocolTCPv6)
	} else if ap.IsIPv6() && ap.IsDatagram() {
		return byte(AddressFamilyAndProtocolUDPv6)
	} else if ap.IsUnix() && ap.IsStream() {
		return byte(AddressFamilyAndProtocolUnixStream)
	} else if ap.IsUnix() && ap.IsDatagram() {
		return byte(AddressFamilyAndProtocolUnixDatagram)
	}

	return byte(AddressFamilyAndProtocolUnknown)
}

// addTLVLen adds the length of the TLV to the header length or errors on uint16 overflow.
func addTLVLen(cur []byte, tlvLen int) ([]byte, error) {
	if tlvLen == 0 {
		return cur, nil
	}
	curLen := binary.BigEndian.Uint16(cur)
	newLen := int(curLen) + tlvLen
	if newLen >= 1<<16 {
		return nil, errUint16Overflow
	}
	a := make([]byte, 2)
	binary.BigEndian.PutUint16(a, uint16(newLen))
	return a, nil
}

/*
const (
	// Section 2.2
	PP2_TYPE_ALPN           PP2Type = 0x01
	PP2_TYPE_AUTHORITY      PP2Type = 0x02
	PP2_TYPE_CRC32C         PP2Type = 0x03
	PP2_TYPE_NOOP           PP2Type = 0x04
	PP2_TYPE_UNIQUE_ID      PP2Type = 0x05
	PP2_TYPE_SSL            PP2Type = 0x20
	PP2_SUBTYPE_SSL_VERSION PP2Type = 0x21
	PP2_SUBTYPE_SSL_CN      PP2Type = 0x22
	PP2_SUBTYPE_SSL_CIPHER  PP2Type = 0x23
	PP2_SUBTYPE_SSL_SIG_ALG PP2Type = 0x24
	PP2_SUBTYPE_SSL_KEY_ALG PP2Type = 0x25
	PP2_TYPE_NETNS          PP2Type = 0x30

	// Section 2.2.7, reserved types
	PP2_TYPE_MIN_CUSTOM     PP2Type = 0xE0
	PP2_TYPE_MAX_CUSTOM     PP2Type = 0xEF
	PP2_TYPE_MIN_EXPERIMENT PP2Type = 0xF0
	PP2_TYPE_MAX_EXPERIMENT PP2Type = 0xF7
	PP2_TYPE_MIN_FUTURE     PP2Type = 0xF8
	PP2_TYPE_MAX_FUTURE     PP2Type = 0xFF
)
*/

// Proxy Protocol Type 2 constants
const (
	PP2TypeNoop      PP2Type = 0x04
	PP2TypeAuthority PP2Type = 0x02
)

// Error constants
var (
	ErrTruncatedTLV    = errors.New("proxyproto: truncated TLV")
	ErrMalformedTLV    = errors.New("proxyproto: malformed TLV Value")
	ErrIncompatibleTLV = errors.New("proxyproto: incompatible TLV type")
)

// PP2Type is the proxy protocol v2 type
type PP2Type byte

// TLV is a uninterpreted Type-Length-Value for V2 protocol, see section 2.2
type TLV struct {
	Type  PP2Type
	Value []byte
}

// SplitTLVs splits the Type-Length-Value vector, returns the vector or an error.
func SplitTLVs(raw []byte) ([]TLV, error) {
	var tlvs []TLV
	for i := 0; i < len(raw); {
		tlv := TLV{
			Type: PP2Type(raw[i]),
		}
		if len(raw)-i <= 2 {
			return nil, ErrTruncatedTLV
		}
		tlvLen := int(binary.BigEndian.Uint16(raw[i+1 : i+3])) // Max length = 65K
		i += 3
		if i+tlvLen > len(raw) {
			return nil, ErrTruncatedTLV
		}
		// Ignore no-op padding
		if tlv.Type != PP2TypeNoop {
			tlv.Value = make([]byte, tlvLen)
			copy(tlv.Value, raw[i:i+tlvLen])
		}
		i += tlvLen
		tlvs = append(tlvs, tlv)
	}
	return tlvs, nil
}

// JoinTLVs joins multiple Type-Length-Value records.
func JoinTLVs(tlvs []TLV) ([]byte, error) {
	var raw []byte
	for _, tlv := range tlvs {
		if len(tlv.Value) > math.MaxUint16 {
			return nil, fmt.Errorf("proxyproto: cannot format TLV %v with length %d", tlv.Type, len(tlv.Value))
		}
		var length [2]byte
		binary.BigEndian.PutUint16(length[:], uint16(len(tlv.Value)))
		raw = append(raw, byte(tlv.Type))
		raw = append(raw, length[:]...)
		raw = append(raw, tlv.Value...)
	}
	return raw, nil
}

func formatUnixName(name string) []byte {
	n := int(lengthUnix) / 2
	if len(name) >= n {
		return []byte(name[:n])
	}
	pad := make([]byte, n-len(name))
	return append([]byte(name), pad...)
}
