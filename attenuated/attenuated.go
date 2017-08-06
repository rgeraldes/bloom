// Package attenuated - implements an attenuated bloom filter
package attenuated

import (
	"github.com/rgeraldes/bloom/standard"
)

// Filter - attenuated bloom filter
type Filter struct {
	depth  uint64
	filter []*standard.Filter
}

// New - creates a new attenuaded bloom filter with a specific depth
func New(depth uint64) *Filter {
	return &Filter{
		depth:  depth,
		filter: make([]*standard.Filter, depth),
	}
}
