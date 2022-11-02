//
// module.go
//
// Copyright (c) 2022 Markku Rossi
//
// All rights reserved.
//

// Package module implements WebAssembly modules and provides data
// structures to represent the decoded module. It also provides
// functions for encoding and decoding modules in binary and text
// formats.
package module

import (
	"fmt"
)

// Section implements a module section.
type Section struct {
	ID   SectionID
	Data []byte
}

// SectionID defines the ID of module sections.
type SectionID byte

// Known section IDs.
const (
	SectionCustom SectionID = iota
	SectionType
	SectionImport
	SectionFunction
	SectionTable
	SectionMemory
	SectionGlobal
	SectionExport
	SectionStart
	SectionElement
	SectionCode
	SectionData
	SectionDataCount
)

var sectionIDs = map[SectionID]string{
	SectionCustom:    "custom",
	SectionType:      "type",
	SectionImport:    "import",
	SectionFunction:  "function",
	SectionTable:     "table",
	SectionMemory:    "memory",
	SectionGlobal:    "global",
	SectionExport:    "export",
	SectionStart:     "start",
	SectionElement:   "element",
	SectionCode:      "code",
	SectionData:      "data",
	SectionDataCount: "data count",
}

func (id SectionID) String() string {
	name, ok := sectionIDs[id]
	if ok {
		return name
	}
	return fmt.Sprintf("{section %v}", byte(id))
}
