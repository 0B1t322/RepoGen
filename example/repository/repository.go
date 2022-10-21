package repository

import (
	"github.com/0B1t322/RepoGen/pkg/filter"
	"github.com/0B1t322/RepoGen/pkg/queryexpression"
	"github.com/0B1t322/RepoGen/pkg/sortorder"
	"github.com/samber/mo"
)

//go:generate go run -mod=mod github.com/0B1t322/RepoGen

//repogen:filter
type (
	FilterQuery = queryexpression.QueryExpression[FilterFields]

	FilterFields struct {
		IDs   mo.Option[filter.FilterField[[]string]]
		Name  mo.Option[filter.FilterField[string]]
		Names mo.Option[filter.FilterField[[]string]]

		SimpleFilter mo.Option[string]
	}
)

type GetObjectQuery struct {
	Filter FilterQuery
}

//repogen:sort
type SortFields struct {
	CreatedAt mo.Option[sortorder.SortOrder]
	Name      mo.Option[sortorder.SortOrder]
}

//repogen:edges
type (
	Edges struct {
		NestedObjectFirst mo.Option[NestedObjectFirstEdge]
	}

	NestedObjectFirstEdge struct {
		HasName mo.Option[bool]
		Offset  mo.Option[int]
		Limit   mo.Option[int]
		Edges   NestedObjectFirstEdges
	}

	NestedObjectFirstEdges struct {
		NestedObjectSecond mo.Option[NestedObjectSecondEdge]
	}

	NestedObjectSecondEdge struct {
		Offset mo.Option[int]
		Limit  mo.Option[int]
		Edges  NestedObjectSecondEdges
	}

	NestedObjectSecondEdges struct {
	}
)
