// Package tsl provides an abstraction of a timeseries log.
//
// The purpose of a timeseries log is to allow writers to extend
// the log without being blocked by readers who need to efficiently
// access a deduplicated, sorted view of the log.
//
// Implementations of this abstraction might be useful to
// implementors of timeseries databases. Indeed, the original
// inspiration for this idea is the multi-snapshot based
// log design that was expressed in influxdb prior to 017c24c98
// when it was removed in favour of a simpler single-snapshot design.
//
// The intent of this project is to develop a standalone abstraction
// that allows me to freely experiment with the ideas without being
// restricted by the constraints of supporting already released code.
// Should the ideas work out, I am happy to workout how to fold the
// ideas into influx should it make sense to do so.
//
// Principles:
//
// The log has three types of clients:
//
// - writers who wish to write extensions to the log with a minimum of delay. Writes should
// never block because of reading or archiving activities; they may either block or fail
// in the case that the write rate exceeds the archiving rate so as to prevent
// out of memory conditions.
//
// - readers who wish to read a sorted, deduplicated view of the log containing
// everything written into the log prior to the commencement of the read
//
// - archivers who wish copy sections of the log to persistent storage and then truncate
// the in memory portion of the log so as to recover memory consumed by those portions.
//
// In addition, we want to minimise the amount of re-sorting that is performed
// during both normal usage and in failure scenarios, such as restart. The primary
// means to achieve this will be to use merge sorts, where possible, to take advantage
// sorting that has already occurred.
//
//
package tsl

import (
	"errors"
)

var (
	// ErrAlreadyFrozen is returned by UnsortedRange.Add if the range has been frozen
	ErrAlreadyFrozen = errors.New("error attempting to add elements to a frozen range.")
)

// An Element is any type which can be compared to another Element that has
// the same type. Two Elements, a and b, are equal if a.Less(b) and b.Less(a) are
// both false.
type Element interface {
	Less(other Element) bool
}

// Elements are slices of Element
type Elements []Element

// A Range knows the first and last Elements of its range and
// knows the maximum number of Elements that it may contain.
//
// A Range is empty if, and only if, Limit() == 0 && First() == nil && Last() == nil.
type Range interface {
	// Limit is the maximum number of Elements in the Range. The actual number may be less.
	Limit() int
	// First answers the first Element in the Range.
	First() Element
	// Last answers the last Element in the range.
	Last() Element
}

// Order is type of a function compares two elements a, b according to a total order and returns true
// if a is less than b in the sense implied by the total order embodied in the function.
type Order func(a Element, b Element) bool

// LessOrder is an Order function that returns true iff a is strictly less than b, in particular
// iff a.Less(b) is true.
func LessOrder(a Element, b Element) bool {
	return a.Less(b)
}

// LessOrEqualOrder is an Order function that returns true iff a is less than or equal to b, in
// particular iff b.Less(a) is false.
func LessOrEqualOrder(a Element, b Element) bool {
	return !b.Less(a)
}

// A Cursor is used to obtain the next Element from
// a SortedRange. There may be multiple Cursors over a single SortedRange.
type Cursor interface {
	// Next provides access to the next sorted slice of []Element from the underlying
	// range or cursors.
	Next() Element
	// Fill buffer with at most len(buffer) elements, returning the number of elements
	// actually filled.
	Fill(buffer []Element) int
}

// A SortedRange is a Range that can provide a Cursor that performs a sorted, deduplicated
// iteration over the contents of a SortedRange. A SortedRange can also be partitioned
// into a pair of (possibly empty) sub-ranges which are also sorted, the concatenation of
// which has the same elements as the originally partitioned range.
type SortedRange interface {
	Range
	// Open a cursor that iterates over the deduplicated elements of the receiver in sorted order.
	Open() Cursor
	// Partition the receiver into two SortedRanges A and B such that for each element a in A
	// o(a,e) is true and for each element b in B, o(b, e) is false. It is assumed that for each pair
	// of consecutive elements, (p,q), in the receiver that o(p,q) is true.
	Partition(e Element, o Order) (SortedRange, SortedRange)
}

// AsSlice converts a SortedRange into a slice.
func AsSlice(r SortedRange) []Element {
	result := make([]Element, r.Limit(), r.Limit())
	c := r.Open()
	n := c.Fill(result)
	result = result[0:n]
	return result
}

// Merge merges two SortedRange to produce a third SortedRange which represents, the
// merged, deduplicated merge of the two original ranges. Where two elements from a and
// b are equal, the resulting SortedRange contains the element from b.
func Merge(a SortedRange, b SortedRange) SortedRange {
	return merge(a, b)
}

// EmptyRange is a SortedRange that has no elements.
var EmptyRange SortedRange

// A UnsortedRange can have slices of elements added to it, but it cannot be
// read directly. To read elements from an UnsortedRange, call Freeze() to obtain
// a SortedRange. It is illegal to call Add() on an UnsortedRange which is already frozen.
type UnsortedRange interface {
	Range
	// Adds the specified elements to the receiver. Returns an error if the
	// receiver has already been frozen. If the receiver has already been frozen
	// then ErrAlreadyFrozen is returned.
	Add([]Element) error
	// Freezes the UnsortedRange, returning a SortedRange for the contained elements.
	// This method is idempotent.
	Freeze() SortedRange
}

// NewUnsortedRange returns an UnsortedRange that can be extended by calling the Add method.
func NewUnsortedRange() UnsortedRange {
	return &mutableRange{}
}
