package bundle

import (
	"fmt"
	"sync"
)

// ContentType defines type of the content in the bundle file.
type ContentType string

const (
	// CTJson represents CTJson files.
	CTJson ContentType = "JSON"
	// CTJournal represents CTJournal files.
	CTJournal = "journal"
	// CTDmesg represents dmesg files.
	CTDmesg = "dmesg"
	// CTOutput is a output of a command.
	CTOutput = "output"
	//CTOther file types
	CTOther = "other"
)

type FileTypeName string

// FileType Describes a kind of files in the bundle (e.g. dcos-marathon.service).
type FileType struct {
	Name        FileTypeName `yaml:"name"`
	ContentType ContentType  `yaml:"contentType"`
	Paths       []string     `yaml:"paths"`
	Description string       `yaml:"description"`
	// DirTypes defines on which host types this file can be found.
	// For example, dcos-marathon.service file can be found only on the masters.
	DirTypes []DirType `yaml:"dirTypes"`
}

var (
	fileTypes   = make(map[FileTypeName]FileType)
	fileTypesMu sync.RWMutex
)

// RegisterFileType adds the file type to the file type registry. It panics
// if the file type with the same name is already registered.
func RegisterFileType(f FileType) {
	fileTypesMu.Lock()
	defer fileTypesMu.Unlock()
	if _, dup := fileTypes[f.Name]; dup {
		panic(fmt.Sprintf("bun.RegisterFileType: called twice for file type %v", f.Name))
	}
	dirTypes := make(map[DirType]struct{})
	for _, t := range f.DirTypes {
		if _, ok := dirTypes[t]; ok {
			panic(fmt.Sprintf("bun.RegisterFileType: duplicate DirType: %v in file type %v", t, f.Name))
		}
		dirTypes[t] = struct{}{}
	}
	fileTypes[f.Name] = f
}

// GetFileType returns a file type by its name. It panics if the file type
// is not in the registry.
func GetFileType(typeName FileTypeName) FileType {
	fileTypesMu.RLock()
	defer fileTypesMu.RUnlock()
	fileType, ok := fileTypes[typeName]
	if !ok {
		panic(fmt.Sprintf("bun.RegisterFileType: No such fileType: %v", typeName))
	}
	return fileType
}
