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
	Producers      []Producer
	TargetFeatures []Feature
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
