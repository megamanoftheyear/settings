package tree

import (
	"example/settings/data/tree/path"
	"time"
)

type Container[N Node] []N

type Segment interface {
	Links() []*Link
	Value() any
	Node
}

type Node interface {
	TrackID() string
	Version() int
	Path() string
	Timestamp() time.Time
	StartText() int
	EndText() int
	Meta() Meta
}

func (cnt Container[N]) TrackID() string      { return cnt.first().TrackID() }
func (cnt Container[N]) Version() int         { return cnt.first().Version() }
func (cnt Container[N]) Timestamp() time.Time { return cnt.first().Timestamp() }
func (cnt Container[N]) StartText() int       { return cnt.first().StartText() }
func (cnt Container[N]) EndText() int         { return cnt.last().EndText() }
func (cnt Container[N]) Meta() Meta           { return Meta{} }

func (cnt Container[N]) Path() path.Path {
	builder := path.NewBuilder(nil)

	for idx, node := range cnt {
		builder.BuildPath(idx, node)
	}

	return
}

func (cnt Container[N]) first() N {
	if len(cnt) == 0 {
		return *new(N)
	}

	return cnt[0]
}

func (cnt Container[N]) last() N {
	if len(cnt) == 0 {
		return *new(N)
	}

	return cnt[len(cnt)-1]
}

type Line struct {
	segments Container[Segment]
	Node
}

func (l *Line) Segments() []Segment { return l.segments }
