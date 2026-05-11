package path

import (
	"example/settings/data/tree"
	"strconv"
	"strings"
)

type Path interface {
	Current() string // full, example tg=>body=>MVT=>departure
	Segment() string // short, example MVT
}

type Link interface {
	ContextPayload() any
	ReasonPayload() any
	Path
}

type Visitor interface {
	Path() Path
	tree.Visitor
}

type Builder struct {
	builder   *strings.Builder
	visitor   Visitor
	separator string
}

func NewBuilder(visitor Visitor) *Builder {
	return &Builder{
		builder:   &strings.Builder{},
		visitor:   visitor,
		separator: "@",
	}
}

func (b *Builder) BuildPath(idx int, node tree.Node) Path {
	b.builder.WriteString(strconv.Itoa(idx))
	b.builder.WriteString(") ")
	if b.builder.Len() != 0 {
		b.builder.WriteString(b.separator)
	}

	b.visitor.Visit(node)
	path := b.visitor.Path()
	b.builder.WriteString(path.Current())

	return path
}

func (b *Builder) Builder() *strings.Builder {
	return b.builder
}
