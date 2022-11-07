//
// binary.go
//
// Copyright (c) 2022 Markku Rossi
//
// All rights reserved.
//

// Package binary implements the WebAsembly binary encoding and
// decoding.
package binary

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

// NameID defines name section IDs.
type NameID byte

// Known name IDs.
const (
	ModuleName NameID = iota
	FunctionNames
	LocalNames
	LabelNames
	TypeNames
	TableNames
	MemoryNames
	GlobalNames
	ElemSegmentNames
	DataSegmentNames
)

var nameNames = map[NameID]string{
	ModuleName:       "ModuleName",
	FunctionNames:    "FunctionNames",
	LocalNames:       "LocalNames",
	LabelNames:       "LabelNames",
	TypeNames:        "TypeNames",
	TableNames:       "TableNames",
	MemoryNames:      "MemoryNames",
	GlobalNames:      "GlobalNames",
	ElemSegmentNames: "ElemSegmentNames",
	DataSegmentNames: "DataSegmentNames",
}

func (id NameID) String() string {
	name, ok := nameNames[id]
	if ok {
		return name
	}
	return fmt.Sprintf("{NameID %v}", byte(id))
}
