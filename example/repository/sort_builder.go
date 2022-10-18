/*
Provide struct and methods to build sort query
Code generatated with repogen
Do not Edit
*/
package repository

import (
	"github.com/0B1t322/RepoGen/pkg/sortorder"
	"github.com/samber/mo"
)

type sortBuilder struct {
	sort []SortFields
}

func SortBuilder() sortBuilder {
	return sortBuilder{}
}

func (s sortBuilder) Build() []SortFields {
	return s.sort
}

func (s sortBuilder) CreatedAt(order sortorder.SortOrder) sortBuilder {
	s.sort = append(s.sort, SortFields{
		CreatedAt: mo.Some(order),
	})
	return s
}

func (s sortBuilder) Name(order sortorder.SortOrder) sortBuilder {
	s.sort = append(s.sort, SortFields{
		Name: mo.Some(order),
	})
	return s
}
