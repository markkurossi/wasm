//
// binary_decode.go
//
// Copyright (c) 2022 Markku Rossi
//
// All rights reserved.
//

package module

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
)

// DecodeBinary decodes module from binary encoding.
func DecodeBinary(data io.Reader) error {
	in := &Decoder{}
	in.PushInput(bufio.NewReader(data))
	magic, err := in.MSBu32()
	if err != nil {
		return err
	}
	fmt.Printf("magic: %08x\n", magic)

	version, err := in.MSBu32()
	if err != nil {
		return err
	}
	fmt.Printf("version: %08x\n", version)

	// Read sections.
	for {
		b, err := in.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		id := SectionID(b)

		size, err := in.LEB128u32()
		if err != nil {
			return err
		}

		in.PushInput(io.LimitReader(in.in, int64(size)))

		section, err := in.decodeSection(id)
		if err != nil {
			return err
		}

		in.PopInput()

		if section.ID == SectionCustom {
			fmt.Printf("section '%v':\n%s", section.ID, hex.Dump(section.Data))
		} else {
			fmt.Printf("section '%v': size=%v\n", id, size)
		}
	}

	return nil
}

// Decoder implements primite functions for WASM binary format
// decoding.
type Decoder struct {
	inputStack []io.Reader
	in         io.Reader
	ofs        int64
}

// PushInput pushes a new input to the decoder.
func (d *Decoder) PushInput(in io.Reader) {
	d.inputStack = append(d.inputStack, in)
	d.in = in
}

// PopInput removes the topmost input from the decoder.
func (d *Decoder) PopInput() {
	d.inputStack = d.inputStack[0 : len(d.inputStack)-1]
	d.in = d.inputStack[len(d.inputStack)-1]
}

// ReadByte reads one byte from the input.
func (d *Decoder) ReadByte() (byte, error) {
	var buf [1]byte
	_, err := d.in.Read(buf[:])
	if err != nil {
		return 0, err
	}
	d.ofs++
	return buf[0], nil
}

// MSBu32 decodes a most-significant bit (MSB) encoded uint32 number.
func (d *Decoder) MSBu32() (uint32, error) {
	var v uint32

	for i := 0; i < 4; i++ {
		b, err := d.ReadByte()
		if err != nil {
			return 0, err
		}
		v <<= 8
		v |= uint32(b)
	}
	return v, nil
}

// LEB128u32 decodes a Little Endian Base 128 encoded uint32 number.
func (d *Decoder) LEB128u32() (uint32, error) {
	var v uint32

	for i := 0; i < 5; i++ {
		b, err := d.ReadByte()
		if err != nil {
			return 0, err
		}
		bits := b & 0b01111111
		if i == 4 && (bits&0b01110000) != 0 {
			return 0, fmt.Errorf("extra bits in u32 at offset %v", d.ofs)
		}
		v |= uint32(bits) << (i * 7)

		if (b & 0b10000000) == 0 {
			return v, nil
		}
	}

	return 0, fmt.Errorf("malformed u32 at offset %v", d.ofs-5)
}

func (d *Decoder) decodeSection(id SectionID) (*Section, error) {
	data, err := io.ReadAll(d.in)
	if err != nil {
		return nil, err
	}
	return &Section{
		ID:   id,
		Data: data,
	}, nil
}
