package checks

import (
	"fmt"

	"github.com/mesosphere/bun/v2/bundle"
)

// SearchCheckBuilder builds a check which searches for the specified
// string in the the specified files. If the pattern
// is found, the check is considered problematic.
// The number of the found line and its content appear in the Check.Problems of the check.
// The check searches only for the first appearance of the line.
type SearchCheckBuilder struct {
	Name         string `yaml:"name"`         // Required
	Description  string `yaml:"description"`  // Optional
	FileTypeName string `yaml:"fileTypeName"` // Required
	SearchString string `yaml:"searchString"` // Required
}

// Build creates a bun.Check.
func (b SearchCheckBuilder) Build() Check {
	if b.FileTypeName == "" {
		panic("FileTypeName should be specified.")
	}
	if b.SearchString == "" {
		panic("SearchString should be set.")
	}
	builder := CheckBuilder{
		Name:        b.Name,
		Description: b.Description,
		Aggregate:   DefaultAggregate,
	}
	t := bundle.GetFileType(b.FileTypeName)
	for _, dirType := range t.DirTypes {
		switch dirType {
		case bundle.DTMaster:
			builder.CollectFromMasters = b.collect
		case bundle.DTAgent:
			builder.CollectFromAgents = b.collect
		case bundle.DTPublicAgent:
			builder.CollectFromPublicAgents = b.collect
		}
	}
	return builder.Build()
}

func (b SearchCheckBuilder) collect(host bundle.Host) (ok bool, details interface{},
	err error) {
	n, line, err := host.FindLine(b.FileTypeName, b.SearchString)
	if err != nil {
		return
	}
	if n != 0 {
		details = fmt.Sprintf("%v: %v", n, line)
		return
	}
	ok = true
	return
}
