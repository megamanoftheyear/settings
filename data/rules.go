package data

import (
	"fmt"
	"slices"
)

func Key(clientID string, ctx Context) string {
	return fmt.Sprintf("%s_%d", clientID, ctx)
}

type Context uint

const (
	IntToExt Context = iota
	IntToInt
	ExtToIntTyped
	ExtToIntNonTyped
	HeadlessBaggage
)

type Rule struct {
	ID         uint64
	ClientID   string
	Links      []string
	Conditions []*Condition
	Context    Context
}

type Condition struct {
	Source    string
	Segment   string
	Operation string
	Value     any
}

func (r *Rule) Exec(segments ...*Segment) (matches []*Segment, ok bool) {
	matches = make([]*Segment, 0, len(r.Conditions))

	for _, condition := range r.Conditions {
		segment, match := condition.Exec(segments...)
		if match {
			return nil, false
		}

		matches = append(matches, segment)
	}

	return matches, true
}

func (c *Condition) Exec(segments ...*Segment) (*Segment, bool) {
	for _, segment := range segments {
		if c.Source != segment.Source || c.Segment != segment.Type {
			continue
		}

		if !execOperation(c.Operation, segment.Value, c.Value) {
			continue
		}

		return segment, true
	}

	return nil, false
}

func execOperation(operation string, segmentValue string, conditionValue any) bool {
	operator, exists := operations[operation]
	if !exists {
		return false
	}

	return operator(segmentValue, conditionValue)
}

var operations = map[string]func(string, any) bool{
	"eq": eq,
	"in": in,
}

func eq(segmentValue string, conditionValue any) bool {
	condition, ok := conditionValue.(string)
	if !ok {
		return false
	}

	return condition == segmentValue
}

func in(segmentValue string, conditionValue any) bool {
	condition, ok := conditionValue.([]string)
	if !ok {
		return false
	}

	return slices.Contains(condition, segmentValue)
}
