/*
Provide struct and methods to build edges query
Code generatated with repogen
Do not Edit
*/
package repository

import (
	"github.com/samber/mo"
)

type EdgesBuilder struct {
	Edges
}

func NewEdgesBuilder() EdgesBuilder {
	return EdgesBuilder{}
}

func (b EdgesBuilder) WithNestedObjectFirst(with NestedObjectFirstEdgeBuilder) EdgesBuilder {
	b.Edges.NestedObjectFirst = mo.Some(with.Build())
	return b
}

func (b EdgesBuilder) Build() Edges {
	return b.Edges
}

type NestedObjectFirstEdgeBuilder struct {
	NestedObjectFirstEdge
}

func NewNestedObjectFirstEdgeBuilder() NestedObjectFirstEdgeBuilder {
	return NestedObjectFirstEdgeBuilder{}
}

func (b NestedObjectFirstEdgeBuilder) WithNestedObjectSecond(with NestedObjectSecondEdgeBuilder) NestedObjectFirstEdgeBuilder {
	b.NestedObjectFirstEdge.Edges.NestedObjectSecond = mo.Some(with.Build())
	return b
}

func (b NestedObjectFirstEdgeBuilder) SetHasName(hasName bool) NestedObjectFirstEdgeBuilder {
	b.NestedObjectFirstEdge.HasName = mo.Some(hasName)
	return b
}

func (b NestedObjectFirstEdgeBuilder) SetOffset(offset int) NestedObjectFirstEdgeBuilder {
	b.NestedObjectFirstEdge.Offset = mo.Some(offset)
	return b
}

func (b NestedObjectFirstEdgeBuilder) SetLimit(limit int) NestedObjectFirstEdgeBuilder {
	b.NestedObjectFirstEdge.Limit = mo.Some(limit)
	return b
}

func (b NestedObjectFirstEdgeBuilder) Build() NestedObjectFirstEdge {
	return b.NestedObjectFirstEdge
}

type NestedObjectFirstEdgesBuilder struct {
	NestedObjectFirstEdges
}

func NewNestedObjectFirstEdgesBuilder() NestedObjectFirstEdgesBuilder {
	return NestedObjectFirstEdgesBuilder{}
}

func (b NestedObjectFirstEdgesBuilder) WithNestedObjectSecond(with NestedObjectSecondEdgeBuilder) NestedObjectFirstEdgesBuilder {
	b.NestedObjectFirstEdges.NestedObjectSecond = mo.Some(with.Build())
	return b
}

func (b NestedObjectFirstEdgesBuilder) Build() NestedObjectFirstEdges {
	return b.NestedObjectFirstEdges
}

type NestedObjectSecondEdgeBuilder struct {
	NestedObjectSecondEdge
}

func NewNestedObjectSecondEdgeBuilder() NestedObjectSecondEdgeBuilder {
	return NestedObjectSecondEdgeBuilder{}
}

func (b NestedObjectSecondEdgeBuilder) SetOffset(offset int) NestedObjectSecondEdgeBuilder {
	b.NestedObjectSecondEdge.Offset = mo.Some(offset)
	return b
}

func (b NestedObjectSecondEdgeBuilder) SetLimit(limit int) NestedObjectSecondEdgeBuilder {
	b.NestedObjectSecondEdge.Limit = mo.Some(limit)
	return b
}

func (b NestedObjectSecondEdgeBuilder) Build() NestedObjectSecondEdge {
	return b.NestedObjectSecondEdge
}

type NestedObjectSecondEdgesBuilder struct {
	NestedObjectSecondEdges
}

func NewNestedObjectSecondEdgesBuilder() NestedObjectSecondEdgesBuilder {
	return NestedObjectSecondEdgesBuilder{}
}

func (b NestedObjectSecondEdgesBuilder) Build() NestedObjectSecondEdges {
	return b.NestedObjectSecondEdges
}
