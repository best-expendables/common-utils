package filter

import "strings"

type BasicFilter struct {
	where  Where
	joins  Joins
	keys   map[string]bool
	groups Groups
}

func NewBasicFilter() *BasicFilter {
	return &BasicFilter{
		where:  Where{},
		joins:  Joins{},
		keys:   Keys{},
		groups: Groups{},
	}
}

// implement repository.Filter interface
func (f *BasicFilter) GetLimit() int {
	return IgnoreLimit
}

// implement repository.Filter interface
func (f *BasicFilter) GetOffset() int {
	return IgnoreOffset
}

// implement repository.Filter interface
func (f *BasicFilter) GetWhere() Where {
	return f.where
}

// GetJoins implement repository.Filter interface
func (f *BasicFilter) GetJoins() Joins {
	return f.joins
}

func (f *BasicFilter) AddWhere(key string, query string, values ...interface{}) *BasicFilter {
	f.where[query] = values
	f.keys[key] = true
	return f
}

func (f *BasicFilter) AddJoin(join string, values ...interface{}) *BasicFilter {
	f.joins[join] = values
	return f
}

func (f *BasicFilter) AddKey(key string) *BasicFilter {
	f.keys[key] = true
	return f
}

func (f *BasicFilter) GetKeys() Keys {
	return f.keys
}

func (f *BasicFilter) AddGroup(query string) *BasicFilter {
	f.groups[query] = true
	return f
}
func (f *BasicFilter) GetGroups() string {
	var queries []string
	for query := range f.groups {
		queries = append(queries, query)
	}
	return strings.Join(queries, ",")
}
