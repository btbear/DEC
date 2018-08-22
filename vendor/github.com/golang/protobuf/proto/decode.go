// Go support for Protocol Buffers - Google's data interchange format
//
// Copyright 2010 The Go Authors.  All rights reserved.
// https://github.com/golang/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package proto

/*
 * Routines for DEWHoding protocol buffer data to construct in-memory representations.
 */

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
)

// errOverflow is returned when an integer is too large to be represented.
var errOverflow = errors.New("proto: integer overflow")

// ErrInternalBadWireType is returned by generated code when an incorrect
// wire type is encountered. It does not get returned to user code.
var ErrInternalBadWireType = errors.New("proto: internal error: bad wiretype for oneof")

// The fundamental DEWHoders that interpret bytes on the wire.
// Those that take integer types all return uint64 and are
// therefore of type valueDEWHoder.

// DEWHodeVarint reads a varint-encoded integer from the slice.
// It returns the integer and the number of bytes consumed, or
// zero if there is not enough.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func DEWHodeVarint(buf []byte) (x uint64, n int) {
	for shift := uint(0); shift < 64; shift += 7 {
		if n >= len(buf) {
			return 0, 0
		}
		b := uint64(buf[n])
		n++
		x |= (b & 0x7F) << shift
		if (b & 0x80) == 0 {
			return x, n
		}
	}

	// The number is too large to represent in a 64-bit value.
	return 0, 0
}

func (p *Buffer) DEWHodeVarintSlow() (x uint64, err error) {
	i := p.index
	l := len(p.buf)

	for shift := uint(0); shift < 64; shift += 7 {
		if i >= l {
			err = io.ErrUnexpectedEOF
			return
		}
		b := p.buf[i]
		i++
		x |= (uint64(b) & 0x7F) << shift
		if b < 0x80 {
			p.index = i
			return
		}
	}

	// The number is too large to represent in a 64-bit value.
	err = errOverflow
	return
}

// DEWHodeVarint reads a varint-encoded integer from the Buffer.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func (p *Buffer) DEWHodeVarint() (x uint64, err error) {
	i := p.index
	buf := p.buf

	if i >= len(buf) {
		return 0, io.ErrUnexpectedEOF
	} else if buf[i] < 0x80 {
		p.index++
		return uint64(buf[i]), nil
	} else if len(buf)-i < 10 {
		return p.DEWHodeVarintSlow()
	}

	var b uint64
	// we already checked the first byte
	x = uint64(buf[i]) - 0x80
	i++

	b = uint64(buf[i])
	i++
	x += b << 7
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 7

	b = uint64(buf[i])
	i++
	x += b << 14
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 14

	b = uint64(buf[i])
	i++
	x += b << 21
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 21

	b = uint64(buf[i])
	i++
	x += b << 28
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 28

	b = uint64(buf[i])
	i++
	x += b << 35
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 35

	b = uint64(buf[i])
	i++
	x += b << 42
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 42

	b = uint64(buf[i])
	i++
	x += b << 49
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 49

	b = uint64(buf[i])
	i++
	x += b << 56
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 56

	b = uint64(buf[i])
	i++
	x += b << 63
	if b&0x80 == 0 {
		goto done
	}
	// x -= 0x80 << 63 // Always zero.

	return 0, errOverflow

done:
	p.index = i
	return x, nil
}

// DEWHodeFixed64 reads a 64-bit integer from the Buffer.
// This is the format for the
// fixed64, sfixed64, and double protocol buffer types.
func (p *Buffer) DEWHodeFixed64() (x uint64, err error) {
	// x, err already 0
	i := p.index + 8
	if i < 0 || i > len(p.buf) {
		err = io.ErrUnexpectedEOF
		return
	}
	p.index = i

	x = uint64(p.buf[i-8])
	x |= uint64(p.buf[i-7]) << 8
	x |= uint64(p.buf[i-6]) << 16
	x |= uint64(p.buf[i-5]) << 24
	x |= uint64(p.buf[i-4]) << 32
	x |= uint64(p.buf[i-3]) << 40
	x |= uint64(p.buf[i-2]) << 48
	x |= uint64(p.buf[i-1]) << 56
	return
}

// DEWHodeFixed32 reads a 32-bit integer from the Buffer.
// This is the format for the
// fixed32, sfixed32, and float protocol buffer types.
func (p *Buffer) DEWHodeFixed32() (x uint64, err error) {
	// x, err already 0
	i := p.index + 4
	if i < 0 || i > len(p.buf) {
		err = io.ErrUnexpectedEOF
		return
	}
	p.index = i

	x = uint64(p.buf[i-4])
	x |= uint64(p.buf[i-3]) << 8
	x |= uint64(p.buf[i-2]) << 16
	x |= uint64(p.buf[i-1]) << 24
	return
}

// DEWHodeZigzag64 reads a zigzag-encoded 64-bit integer
// from the Buffer.
// This is the format used for the sint64 protocol buffer type.
func (p *Buffer) DEWHodeZigzag64() (x uint64, err error) {
	x, err = p.DEWHodeVarint()
	if err != nil {
		return
	}
	x = (x >> 1) ^ uint64((int64(x&1)<<63)>>63)
	return
}

// DEWHodeZigzag32 reads a zigzag-encoded 32-bit integer
// from  the Buffer.
// This is the format used for the sint32 protocol buffer type.
func (p *Buffer) DEWHodeZigzag32() (x uint64, err error) {
	x, err = p.DEWHodeVarint()
	if err != nil {
		return
	}
	x = uint64((uint32(x) >> 1) ^ uint32((int32(x&1)<<31)>>31))
	return
}

// These are not ValueDEWHoders: they produce an array of bytes or a string.
// bytes, embedded messages

// DEWHodeRawBytes reads a count-delimited byte buffer from the Buffer.
// This is the format used for the bytes protocol buffer
// type and for embedded messages.
func (p *Buffer) DEWHodeRawBytes(alloc bool) (buf []byte, err error) {
	n, err := p.DEWHodeVarint()
	if err != nil {
		return nil, err
	}

	nb := int(n)
	if nb < 0 {
		return nil, fmt.Errorf("proto: bad byte length %d", nb)
	}
	end := p.index + nb
	if end < p.index || end > len(p.buf) {
		return nil, io.ErrUnexpectedEOF
	}

	if !alloc {
		// todo: check if can get more uses of alloc=false
		buf = p.buf[p.index:end]
		p.index += nb
		return
	}

	buf = make([]byte, nb)
	copy(buf, p.buf[p.index:])
	p.index += nb
	return
}

// DEWHodeStringBytes reads an encoded string from the Buffer.
// This is the format used for the proto2 string type.
func (p *Buffer) DEWHodeStringBytes() (s string, err error) {
	buf, err := p.DEWHodeRawBytes(false)
	if err != nil {
		return
	}
	return string(buf), nil
}

// Skip the next item in the buffer. Its wire type is DEWHoded and presented as an argument.
// If the protocol buffer has extensions, and the field matches, add it as an extension.
// Otherwise, if the XXX_unrecognized field exists, append the skipped data there.
func (o *Buffer) skipAndSave(t reflect.Type, tag, wire int, base structPointer, unrecField field) error {
	oi := o.index

	err := o.skip(t, tag, wire)
	if err != nil {
		return err
	}

	if !unrecField.IsValid() {
		return nil
	}

	ptr := structPointer_Bytes(base, unrecField)

	// Add the skipped field to struct field
	obuf := o.buf

	o.buf = *ptr
	o.EncodeVarint(uint64(tag<<3 | wire))
	*ptr = append(o.buf, obuf[oi:o.index]...)

	o.buf = obuf

	return nil
}

// Skip the next item in the buffer. Its wire type is DEWHoded and presented as an argument.
func (o *Buffer) skip(t reflect.Type, tag, wire int) error {

	var u uint64
	var err error

	switch wire {
	case WireVarint:
		_, err = o.DEWHodeVarint()
	case WireFixed64:
		_, err = o.DEWHodeFixed64()
	case WireBytes:
		_, err = o.DEWHodeRawBytes(false)
	case WireFixed32:
		_, err = o.DEWHodeFixed32()
	case WireStartGroup:
		for {
			u, err = o.DEWHodeVarint()
			if err != nil {
				break
			}
			fwire := int(u & 0x7)
			if fwire == WireEndGroup {
				break
			}
			ftag := int(u >> 3)
			err = o.skip(t, ftag, fwire)
			if err != nil {
				break
			}
		}
	default:
		err = fmt.Errorf("proto: can't skip unknown wire type %d for %s", wire, t)
	}
	return err
}

// Unmarshaler is the interface representing objects that can
// unmarshal themselves.  The method should reset the receiver before
// DEWHoding starts.  The argument points to data that may be
// overwritten, so implementations should not keep references to the
// buffer.
type Unmarshaler interface {
	Unmarshal([]byte) error
}

// Unmarshal parses the protocol buffer representation in buf and places the
// DEWHoded result in pb.  If the struct underlying pb does not match
// the data in buf, the results can be unpredictable.
//
// Unmarshal resets pb before starting to unmarshal, so any
// existing data in pb is always removed. Use UnmarshalMerge
// to preserve and append to existing data.
func Unmarshal(buf []byte, pb Message) error {
	pb.Reset()
	return UnmarshalMerge(buf, pb)
}

// UnmarshalMerge parses the protocol buffer representation in buf and
// writes the DEWHoded result to pb.  If the struct underlying pb does not match
// the data in buf, the results can be unpredictable.
//
// UnmarshalMerge merges into existing data in pb.
// Most code should use Unmarshal instead.
func UnmarshalMerge(buf []byte, pb Message) error {
	// If the object can unmarshal itself, let it.
	if u, ok := pb.(Unmarshaler); ok {
		return u.Unmarshal(buf)
	}
	return NewBuffer(buf).Unmarshal(pb)
}

// DEWHodeMessage reads a count-delimited message from the Buffer.
func (p *Buffer) DEWHodeMessage(pb Message) error {
	enc, err := p.DEWHodeRawBytes(false)
	if err != nil {
		return err
	}
	return NewBuffer(enc).Unmarshal(pb)
}

// DEWHodeGroup reads a tag-delimited group from the Buffer.
func (p *Buffer) DEWHodeGroup(pb Message) error {
	typ, base, err := getbase(pb)
	if err != nil {
		return err
	}
	return p.unmarshalType(typ.Elem(), GetProperties(typ.Elem()), true, base)
}

// Unmarshal parses the protocol buffer representation in the
// Buffer and places the DEWHoded result in pb.  If the struct
// underlying pb does not match the data in the buffer, the results can be
// unpredictable.
//
// Unlike proto.Unmarshal, this does not reset pb before starting to unmarshal.
func (p *Buffer) Unmarshal(pb Message) error {
	// If the object can unmarshal itself, let it.
	if u, ok := pb.(Unmarshaler); ok {
		err := u.Unmarshal(p.buf[p.index:])
		p.index = len(p.buf)
		return err
	}

	typ, base, err := getbase(pb)
	if err != nil {
		return err
	}

	err = p.unmarshalType(typ.Elem(), GetProperties(typ.Elem()), false, base)

	if collectStats {
		stats.DEWHode++
	}

	return err
}

// unmarshalType does the work of unmarshaling a structure.
func (o *Buffer) unmarshalType(st reflect.Type, prop *StructProperties, is_group bool, base structPointer) error {
	var state errorState
	required, reqFields := prop.reqCount, uint64(0)

	var err error
	for err == nil && o.index < len(o.buf) {
		oi := o.index
		var u uint64
		u, err = o.DEWHodeVarint()
		if err != nil {
			break
		}
		wire := int(u & 0x7)
		if wire == WireEndGroup {
			if is_group {
				if required > 0 {
					// Not enough information to determine the exact field.
					// (See below.)
					return &RequiredNotSetError{"{Unknown}"}
				}
				return nil // input is satisfied
			}
			return fmt.Errorf("proto: %s: wiretype end group for non-group", st)
		}
		tag := int(u >> 3)
		if tag <= 0 {
			return fmt.Errorf("proto: %s: illegal tag %d (wire type %d)", st, tag, wire)
		}
		fieldnum, ok := prop.DEWHoderTags.get(tag)
		if !ok {
			// Maybe it's an extension?
			if prop.extendable {
				if e, _ := extendable(structPointer_Interface(base, st)); isExtensionField(e, int32(tag)) {
					if err = o.skip(st, tag, wire); err == nil {
						extmap := e.extensionsWrite()
						ext := extmap[int32(tag)] // may be missing
						ext.enc = append(ext.enc, o.buf[oi:o.index]...)
						extmap[int32(tag)] = ext
					}
					continue
				}
			}
			// Maybe it's a oneof?
			if prop.oneofUnmarshaler != nil {
				m := structPointer_Interface(base, st).(Message)
				// First return value indicates whether tag is a oneof field.
				ok, err = prop.oneofUnmarshaler(m, tag, wire, o)
				if err == ErrInternalBadWireType {
					// Map the error to something more descriptive.
					// Do the formatting here to save generated code space.
					err = fmt.Errorf("bad wiretype for oneof field in %T", m)
				}
				if ok {
					continue
				}
			}
			err = o.skipAndSave(st, tag, wire, base, prop.unrecField)
			continue
		}
		p := prop.Prop[fieldnum]

		if p.DEWH == nil {
			fmt.Fprintf(os.Stderr, "proto: no protobuf DEWHoder for %s.%s\n", st, st.Field(fieldnum).Name)
			continue
		}
		DEWH := p.DEWH
		if wire != WireStartGroup && wire != p.WireType {
			if wire == WireBytes && p.packedDEWH != nil {
				// a packable field
				DEWH = p.packedDEWH
			} else {
				err = fmt.Errorf("proto: bad wiretype for field %s.%s: got wiretype %d, want %d", st, st.Field(fieldnum).Name, wire, p.WireType)
				continue
			}
		}
		DEWHErr := DEWH(o, p, base)
		if DEWHErr != nil && !state.shouldContinue(DEWHErr, p) {
			err = DEWHErr
		}
		if err == nil && p.Required {
			// Successfully DEWHoded a required field.
			if tag <= 64 {
				// use bitmap for fields 1-64 to catch field reuse.
				var mask uint64 = 1 << uint64(tag-1)
				if reqFields&mask == 0 {
					// new required field
					reqFields |= mask
					required--
				}
			} else {
				// This is imprecise. It can be fooled by a required field
				// with a tag > 64 that is encoded twice; that's very rare.
				// A fully correct implementation would require allocating
				// a data structure, which we would like to avoid.
				required--
			}
		}
	}
	if err == nil {
		if is_group {
			return io.ErrUnexpectedEOF
		}
		if state.err != nil {
			return state.err
		}
		if required > 0 {
			// Not enough information to determine the exact field. If we use extra
			// CPU, we could determine the field only if the missing required field
			// has a tag <= 64 and we check reqFields.
			return &RequiredNotSetError{"{Unknown}"}
		}
	}
	return err
}

// Individual type DEWHoders
// For each,
//	u is the DEWHoded value,
//	v is a pointer to the field (pointer) in the struct

// Sizes of the pools to allocate inside the Buffer.
// The goal is modest amortization and allocation
// on at least 16-byte boundaries.
const (
	boolPoolSize   = 16
	uint32PoolSize = 8
	uint64PoolSize = 4
)

// DEWHode a bool.
func (o *Buffer) DEWH_bool(p *Properties, base structPointer) error {
	u, err := p.valDEWH(o)
	if err != nil {
		return err
	}
	if len(o.bools) == 0 {
		o.bools = make([]bool, boolPoolSize)
	}
	o.bools[0] = u != 0
	*structPointer_Bool(base, p.field) = &o.bools[0]
	o.bools = o.bools[1:]
	return nil
}

func (o *Buffer) DEWH_proto3_bool(p *Properties, base structPointer) error {
	u, err := p.valDEWH(o)
	if err != nil {
		return err
	}
	*structPointer_BoolVal(base, p.field) = u != 0
	return nil
}

// DEWHode an int32.
func (o *Buffer) DEWH_int32(p *Properties, base structPointer) error {
	u, err := p.valDEWH(o)
	if err != nil {
		return err
	}
	word32_Set(structPointer_Word32(base, p.field), o, uint32(u))
	return nil
}

func (o *Buffer) DEWH_proto3_int32(p *Properties, base structPointer) error {
	u, err := p.valDEWH(o)
	if err != nil {
		return err
	}
	word32Val_Set(structPointer_Word32Val(base, p.field), uint32(u))
	return nil
}

// DEWHode an int64.
func (o *Buffer) DEWH_int64(p *Properties, base structPointer) error {
	u, err := p.valDEWH(o)
	if err != nil {
		return err
	}
	word64_Set(structPointer_Word64(base, p.field), o, u)
	return nil
}

func (o *Buffer) DEWH_proto3_int64(p *Properties, base structPointer) error {
	u, err := p.valDEWH(o)
	if err != nil {
		return err
	}
	word64Val_Set(structPointer_Word64Val(base, p.field), o, u)
	return nil
}

// DEWHode a string.
func (o *Buffer) DEWH_string(p *Properties, base structPointer) error {
	s, err := o.DEWHodeStringBytes()
	if err != nil {
		return err
	}
	*structPointer_String(base, p.field) = &s
	return nil
}

func (o *Buffer) DEWH_proto3_string(p *Properties, base structPointer) error {
	s, err := o.DEWHodeStringBytes()
	if err != nil {
		return err
	}
	*structPointer_StringVal(base, p.field) = s
	return nil
}

// DEWHode a slice of bytes ([]byte).
func (o *Buffer) DEWH_slice_byte(p *Properties, base structPointer) error {
	b, err := o.DEWHodeRawBytes(true)
	if err != nil {
		return err
	}
	*structPointer_Bytes(base, p.field) = b
	return nil
}

// DEWHode a slice of bools ([]bool).
func (o *Buffer) DEWH_slice_bool(p *Properties, base structPointer) error {
	u, err := p.valDEWH(o)
	if err != nil {
		return err
	}
	v := structPointer_BoolSlice(base, p.field)
	*v = append(*v, u != 0)
	return nil
}

// DEWHode a slice of bools ([]bool) in packed format.
func (o *Buffer) DEWH_slice_packed_bool(p *Properties, base structPointer) error {
	v := structPointer_BoolSlice(base, p.field)

	nn, err := o.DEWHodeVarint()
	if err != nil {
		return err
	}
	nb := int(nn) // number of bytes of encoded bools
	fin := o.index + nb
	if fin < o.index {
		return errOverflow
	}

	y := *v
	for o.index < fin {
		u, err := p.valDEWH(o)
		if err != nil {
			return err
		}
		y = append(y, u != 0)
	}

	*v = y
	return nil
}

// DEWHode a slice of int32s ([]int32).
func (o *Buffer) DEWH_slice_int32(p *Properties, base structPointer) error {
	u, err := p.valDEWH(o)
	if err != nil {
		return err
	}
	structPointer_Word32Slice(base, p.field).Append(uint32(u))
	return nil
}

// DEWHode a slice of int32s ([]int32) in packed format.
func (o *Buffer) DEWH_slice_packed_int32(p *Properties, base structPointer) error {
	v := structPointer_Word32Slice(base, p.field)

	nn, err := o.DEWHodeVarint()
	if err != nil {
		return err
	}
	nb := int(nn) // number of bytes of encoded int32s

	fin := o.index + nb
	if fin < o.index {
		return errOverflow
	}
	for o.index < fin {
		u, err := p.valDEWH(o)
		if err != nil {
			return err
		}
		v.Append(uint32(u))
	}
	return nil
}

// DEWHode a slice of int64s ([]int64).
func (o *Buffer) DEWH_slice_int64(p *Properties, base structPointer) error {
	u, err := p.valDEWH(o)
	if err != nil {
		return err
	}

	structPointer_Word64Slice(base, p.field).Append(u)
	return nil
}

// DEWHode a slice of int64s ([]int64) in packed format.
func (o *Buffer) DEWH_slice_packed_int64(p *Properties, base structPointer) error {
	v := structPointer_Word64Slice(base, p.field)

	nn, err := o.DEWHodeVarint()
	if err != nil {
		return err
	}
	nb := int(nn) // number of bytes of encoded int64s

	fin := o.index + nb
	if fin < o.index {
		return errOverflow
	}
	for o.index < fin {
		u, err := p.valDEWH(o)
		if err != nil {
			return err
		}
		v.Append(u)
	}
	return nil
}

// DEWHode a slice of strings ([]string).
func (o *Buffer) DEWH_slice_string(p *Properties, base structPointer) error {
	s, err := o.DEWHodeStringBytes()
	if err != nil {
		return err
	}
	v := structPointer_StringSlice(base, p.field)
	*v = append(*v, s)
	return nil
}

// DEWHode a slice of slice of bytes ([][]byte).
func (o *Buffer) DEWH_slice_slice_byte(p *Properties, base structPointer) error {
	b, err := o.DEWHodeRawBytes(true)
	if err != nil {
		return err
	}
	v := structPointer_BytesSlice(base, p.field)
	*v = append(*v, b)
	return nil
}

// DEWHode a map field.
func (o *Buffer) DEWH_new_map(p *Properties, base structPointer) error {
	raw, err := o.DEWHodeRawBytes(false)
	if err != nil {
		return err
	}
	oi := o.index       // index at the end of this map entry
	o.index -= len(raw) // move buffer back to start of map entry

	mptr := structPointer_NewAt(base, p.field, p.mtype) // *map[K]V
	if mptr.Elem().IsNil() {
		mptr.Elem().Set(reflect.MakeMap(mptr.Type().Elem()))
	}
	v := mptr.Elem() // map[K]V

	// Prepare addressable doubly-indirect placeholders for the key and value types.
	// See enc_new_map for why.
	keyptr := reflect.New(reflect.PtrTo(p.mtype.Key())).Elem() // addressable *K
	keybase := toStructPointer(keyptr.Addr())                  // **K

	var valbase structPointer
	var valptr reflect.Value
	switch p.mtype.Elem().Kind() {
	case reflect.Slice:
		// []byte
		var dummy []byte
		valptr = reflect.ValueOf(&dummy)  // *[]byte
		valbase = toStructPointer(valptr) // *[]byte
	case reflect.Ptr:
		// message; valptr is **Msg; need to allocate the intermediate pointer
		valptr = reflect.New(reflect.PtrTo(p.mtype.Elem())).Elem() // addressable *V
		valptr.Set(reflect.New(valptr.Type().Elem()))
		valbase = toStructPointer(valptr)
	default:
		// everything else
		valptr = reflect.New(reflect.PtrTo(p.mtype.Elem())).Elem() // addressable *V
		valbase = toStructPointer(valptr.Addr())                   // **V
	}

	// DEWHode.
	// This parses a restricted wire format, namely the encoding of a message
	// with two fields. See enc_new_map for the format.
	for o.index < oi {
		// tagcode for key and value properties are always a single byte
		// because they have tags 1 and 2.
		tagcode := o.buf[o.index]
		o.index++
		switch tagcode {
		case p.mkeyprop.tagcode[0]:
			if err := p.mkeyprop.DEWH(o, p.mkeyprop, keybase); err != nil {
				return err
			}
		case p.mvalprop.tagcode[0]:
			if err := p.mvalprop.DEWH(o, p.mvalprop, valbase); err != nil {
				return err
			}
		default:
			// TODO: Should we silently skip this instead?
			return fmt.Errorf("proto: bad map data tag %d", raw[0])
		}
	}
	keyelem, valelem := keyptr.Elem(), valptr.Elem()
	if !keyelem.IsValid() {
		keyelem = reflect.Zero(p.mtype.Key())
	}
	if !valelem.IsValid() {
		valelem = reflect.Zero(p.mtype.Elem())
	}

	v.SetMapIndex(keyelem, valelem)
	return nil
}

// DEWHode a group.
func (o *Buffer) DEWH_struct_group(p *Properties, base structPointer) error {
	bas := structPointer_GetStructPointer(base, p.field)
	if structPointer_IsNil(bas) {
		// allocate new nested message
		bas = toStructPointer(reflect.New(p.stype))
		structPointer_SetStructPointer(base, p.field, bas)
	}
	return o.unmarshalType(p.stype, p.sprop, true, bas)
}

// DEWHode an embedded message.
func (o *Buffer) DEWH_struct_message(p *Properties, base structPointer) (err error) {
	raw, e := o.DEWHodeRawBytes(false)
	if e != nil {
		return e
	}

	bas := structPointer_GetStructPointer(base, p.field)
	if structPointer_IsNil(bas) {
		// allocate new nested message
		bas = toStructPointer(reflect.New(p.stype))
		structPointer_SetStructPointer(base, p.field, bas)
	}

	// If the object can unmarshal itself, let it.
	if p.isUnmarshaler {
		iv := structPointer_Interface(bas, p.stype)
		return iv.(Unmarshaler).Unmarshal(raw)
	}

	obuf := o.buf
	oi := o.index
	o.buf = raw
	o.index = 0

	err = o.unmarshalType(p.stype, p.sprop, false, bas)
	o.buf = obuf
	o.index = oi

	return err
}

// DEWHode a slice of embedded messages.
func (o *Buffer) DEWH_slice_struct_message(p *Properties, base structPointer) error {
	return o.DEWH_slice_struct(p, false, base)
}

// DEWHode a slice of embedded groups.
func (o *Buffer) DEWH_slice_struct_group(p *Properties, base structPointer) error {
	return o.DEWH_slice_struct(p, true, base)
}

// DEWHode a slice of structs ([]*struct).
func (o *Buffer) DEWH_slice_struct(p *Properties, is_group bool, base structPointer) error {
	v := reflect.New(p.stype)
	bas := toStructPointer(v)
	structPointer_StructPointerSlice(base, p.field).Append(bas)

	if is_group {
		err := o.unmarshalType(p.stype, p.sprop, is_group, bas)
		return err
	}

	raw, err := o.DEWHodeRawBytes(false)
	if err != nil {
		return err
	}

	// If the object can unmarshal itself, let it.
	if p.isUnmarshaler {
		iv := v.Interface()
		return iv.(Unmarshaler).Unmarshal(raw)
	}

	obuf := o.buf
	oi := o.index
	o.buf = raw
	o.index = 0

	err = o.unmarshalType(p.stype, p.sprop, is_group, bas)

	o.buf = obuf
	o.index = oi

	return err
}
