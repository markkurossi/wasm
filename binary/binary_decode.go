//
// binary_decode.go
//
// Copyright (c) 2022 Markku Rossi
//
// All rights reserved.
//

package binary

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"math"

	"github.com/markkurossi/wasm/module"
)

// DecodeBinary decodes module from binary encoding.
func DecodeBinary(data io.Reader) error {
	in := &Decoder{}
	in.PushInput(bufio.NewReader(data), math.MaxInt64)
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

	m := module.New()

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

		err = in.decodeSection(m, id, int64(size))
		if err != nil {
			return err
		}
	}

	return nil
}

// Decoder implements primite functions for WASM binary format
// decoding.
type Decoder struct {
	inputStack []*input
	in         *input
	ofs        int64
}

type input struct {
	in    io.Reader
	start int64
	size  int64
}

// PushInput pushes a new input to the decoder.
func (d *Decoder) PushInput(in io.Reader, size int64) {
	d.in = &input{
		in:    in,
		start: d.ofs,
		size:  size,
	}
	d.inputStack = append(d.inputStack, d.in)
}

// PopInput removes the topmost input from the decoder.
func (d *Decoder) PopInput() {
	d.inputStack = d.inputStack[0 : len(d.inputStack)-1]
	d.in = d.inputStack[len(d.inputStack)-1]
}

// Avail returns the number of bytes available in the current decoder
// input.
func (d *Decoder) Avail() int64 {
	return d.in.size - (d.ofs - d.in.start)
}

// ReadByte reads one byte from the input.
func (d *Decoder) ReadByte() (byte, error) {
	var buf [1]byte
	_, err := d.in.in.Read(buf[:])
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

// Name decodes an UTF-8 name.
func (d *Decoder) Name() (string, error) {
	l, err := d.LEB128u32()
	if err != nil {
		return "", err
	}
	if int64(l) > d.Avail() {
		return "", fmt.Errorf("truncated name: %v > %v", l, d.Avail())
	}
	buf := make([]byte, l)
	_, err = io.ReadFull(d.in.in, buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (d *Decoder) decodeSection(m *module.Module, id SectionID,
	size int64) error {

	if size > d.Avail() {
		return fmt.Errorf("malformed section, need %v bytes, got %v",
			size, d.Avail())
	}

	d.PushInput(io.LimitReader(d.in.in, size), size)
	defer d.PopInput()

	switch id {
	case SectionCustom:
		name, err := d.Name()
		if err != nil {
			return err
		}
		switch name {
		case "name":
			return d.decodeCustomSectionName(m)

		case "producers":
			return d.decodeCustomSectionProducers(m)

		case "target_features":
			return d.decodeCustomSectionTargetFeatures(m)

		default:
			return fmt.Errorf("unsupported custom section %v", name)
		}

	case SectionType, SectionImport, SectionFunction, SectionTable,
		SectionMemory, SectionGlobal, SectionExport, SectionStart,
		SectionElement, SectionCode, SectionData, SectionDataCount:
		data, err := io.ReadAll(d.in.in)
		if err != nil {
			return err
		}
		if false {
			fmt.Printf("section %s:\n%s", id, hex.Dump(data))
		}
		return nil

	default:
		_, err := io.Copy(io.Discard, d.in.in)
		if err != nil {
			return err
		}
		return fmt.Errorf("unknown section: %v", id)
	}

}

func (d *Decoder) decodeCustomSectionName(m *module.Module) error {
	for {
		b, err := d.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		id := NameID(b)
		size, err := d.LEB128u32()
		if err != nil {
			return err
		}
		switch id {
		case FunctionNames:
			m.Names.Functions, err = d.parseNameMap()
			if err != nil {
				return err
			}

		case GlobalNames:
			m.Names.Globals, err = d.parseNameMap()
			if err != nil {
				return err
			}

		case DataSegmentNames:
			m.Names.DataSegments, err = d.parseNameMap()
			if err != nil {
				return err
			}

		default:
			fmt.Printf("section custom %v: len=%v\n", id, size)
			for i := 0; i < int(size); i++ {
				_, err = d.ReadByte()
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (d *Decoder) parseNameMap() (module.NameMap, error) {
	count, err := d.LEB128u32()
	if err != nil {
		return nil, err
	}

	var result module.NameMap
	for i := 0; i < int(count); i++ {
		idx, err := d.LEB128u32()
		if err != nil {
			return nil, err
		}
		name, err := d.Name()
		if err != nil {
			return nil, err
		}
		result = append(result, module.NameAssoc{
			Idx:  idx,
			Name: name,
		})
	}
	return result, nil
}

func (d *Decoder) decodeCustomSectionProducers(m *module.Module) error {
	count, err := d.LEB128u32()
	if err != nil {
		return err
	}
	for i := 0; i < int(count); i++ {
		name, err := d.Name()
		if err != nil {
			return err
		}

		producer := module.Producer{
			Name: name,
		}

		nValues, err := d.LEB128u32()
		if err != nil {
			return err
		}
		for j := 0; j < int(nValues); j++ {
			value, err := d.Name()
			if err != nil {
				return err
			}
			version, err := d.Name()
			if err != nil {
				return err
			}

			producer.Values = append(producer.Values, module.VersionedName{
				Name:    value,
				Version: version,
			})
		}
		m.Producers = append(m.Producers, producer)
	}
	fmt.Printf("Producers: %v\n", m.Producers)
	return nil
}

func (d *Decoder) decodeCustomSectionTargetFeatures(m *module.Module) error {
	count, err := d.LEB128u32()
	if err != nil {
		return err
	}
	for i := 0; i < int(count); i++ {
		prefix, err := d.ReadByte()
		if err != nil {
			return err
		}
		name, err := d.Name()
		if err != nil {
			return err
		}
		m.TargetFeatures = append(m.TargetFeatures, module.Feature{
			Prefix: prefix,
			Name:   name,
		})
	}
	fmt.Printf("Target Fetures: %v\n", m.TargetFeatures)
	return nil
}
