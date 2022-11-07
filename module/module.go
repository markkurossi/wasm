//
// module.go
//
// Copyright (c) 2022 Markku Rossi
//
// All rights reserved.
//

// Package module implements WebAssembly modules.
package module

import (
	"fmt"
)

// New creates a new module.
func New() *Module {
	return new(Module)
}

// Module implements a WebAssembly module.
type Module struct {
	Names struct {
		Module       string
		Functions    NameMap
		Globals      NameMap
		DataSegments NameMap
	}
	Producers      []Producer
	TargetFeatures []Feature
}

// NameMap defines a mappings from names to indices in the given index
// space.
type NameMap []NameAssoc

func (nmap NameMap) String() string {
	result := "{\n"
	for _, assoc := range nmap {
		result += fmt.Sprintf("  %v\t%v\n", assoc.Idx, assoc.Name)
	}
	return result + "}"
}

// NameAssoc defines a name index mapping.
type NameAssoc struct {
	Idx  uint32
	Name string
}

// Producer defines a tool that produced this module.
type Producer struct {
	Name   string
	Values []VersionedName
}

func (p Producer) String() string {
	str := fmt.Sprintf("%s=", p.Name)
	for i, v := range p.Values {
		if i > 0 {
			str += ","
		}
		str += fmt.Sprintf("%s:%s", v.Name, v.Version)
	}
	return str
}

// VersionedName defines a name with version information.
type VersionedName struct {
	Name    string
	Version string
}

// Feature defines WebAssembly features
type Feature struct {
	Prefix byte
	Name   string
}

func (f Feature) String() string {
	return fmt.Sprintf("%c%v", f.Prefix, f.Name)
}

// Known features.
const (
	Atomics            = "atomics"
	BulkMemory         = "bulk-memory"
	ExceptionHandling  = "exception-handling"
	Multivalue         = "multivalue"
	MutableGlobals     = "mutable-globals"
	NontrappingFptoint = "nontrapping-fptoint"
	SignExt            = "sign-ext"
	SIMD128            = "simd128"
	TailCall           = "tail-call"
)
