package tsl

import (
	"fmt"
	"sync"
)

// mergeableRange represents a possibly incomplete merge between an
// immutableRange and an unsorted range. A mergeable range is closed
// to further additions but needs to be locked until the merge is complete.
// Once the merge is complete, the range becomes immutable.
//
// The unsorted range is sorted with an O(n.log(n)) sort
// when the first cursor is opened or the range is first partitioned.
type mergeableRange struct {
	immutableRange
	mu       sync.RWMutex
	left     SortedRange
	right    SortedRange
	unsorted *unsortedRange
	lx       *mergeCursor // a peekable cursor into the left arm of the merge
	rx       *mergeCursor // a peekable cursor into the right arm of the merge
	mx       int          // number of copied elements
	nx       int          // number of deduplicated elements
}

// mergeOne advances nx so that it represents the length of the merged, deduplicated slice and advances
// mx to point to the location of the next element to be written into the merged slice.
func (r *mergeableRange) mergeOne() {
	if r.left == nil {
		return
	}
	if r.rx == nil {
		panic("illegal state: r.rx is nil")
	}
	notexhausted := func() bool {
		return r.lx.peek() != nil || r.rx.peek() != nil
	}
	r.nx = r.mx
	for notexhausted() && r.nx == r.mx {
		leftPeek := r.lx.peek()
		rightPeek := r.rx.peek()
		if rightPeek == nil || (leftPeek != nil && leftPeek.Less(rightPeek)) {
			r.elements[r.mx] = r.lx.next()
		} else if leftPeek == nil || rightPeek.Less(leftPeek) {
			r.elements[r.mx] = r.rx.next()
		} else {
			r.elements[r.mx] = r.rx.next()
			r.lx.next()
		}

		if r.mx > 0 {
			if r.elements[r.mx-1].Less(r.elements[r.mx]) {
				r.mx++
			} else {
				r.elements[r.mx-1] = r.elements[r.mx]
				r.elements[r.mx] = nil
			}
		} else {
			r.mx++
			r.nx = r.mx
		}
	}
	if !notexhausted() {
		r.left = nil
		r.lx = nil
		r.right = nil
		r.rx = nil
		r.unsorted = nil
		r.nx = r.mx
		r.elements = r.elements[0:r.mx]
	}
}

func newMergeableRange(first Element, last Element, left SortedRange, right SortedRange, unsorted *unsortedRange) *mergeableRange {

	if last == nil {
		if first != nil || left.Limit() != 0 || right.Limit() != 0 {
			panic("inconsistent")
		}
		return &mergeableRange{
			immutableRange: immutableRange{
				basicRange: basicRange{
					first:    nil,
					last:     nil,
					elements: nil,
				},
			},
		}
	}

	if last.Less(first) {
		first, last = last, first
	}

	limit := left.Limit()
	var rx *mergeCursor

	if right == nil {
		limit += unsorted.Limit()
	} else {
		limit += right.Limit()
		rx = &mergeCursor{
			underlying: right.Open(),
		}
	}

	return &mergeableRange{
		immutableRange: immutableRange{
			basicRange: basicRange{
				first:    first,
				last:     last,
				elements: make([]Element, limit),
			},
		},
		left:     left,
		right:    right,
		unsorted: unsorted,
		lx: &mergeCursor{
			underlying: left.Open(),
		},
		rx: rx,
	}
}

// Open opens a cursor over a mergeable range. If the range
// has already been merged then the cursor is opened over the
// result of the merge. Otherwise, a cursor that progresses
// the merge on each iteration is used.
//
// The first time a cursor is opened, it incurs on O(n.log(n))
// cost to sort the unsorted part of the range.
func (r *mergeableRange) Open() Cursor {

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.left != nil && r.right == nil {
		r.right = r.unsorted.freeze()
		r.rx = &mergeCursor{underlying: r.right.Open()}
	}

	if r.left == nil {
		return r.immutableRange.Open()
	} else {
		if r.rx == nil {
			panic(fmt.Errorf("r.rx is nil. r.left is %v. r.right is %v", r.left, r.right))
		}
		return &mergeableCursor{
			basicCursor: basicCursor{
				next:     0,
				elements: r.elements,
			},
			advance: r.mergeOne,
			mu:      &r.mu,
			nx:      &r.nx,
		}
	}
}

func useEmptyRangeIfEmpty(s SortedRange) SortedRange {
	if s.Limit() == 0 {
		return EmptyRange
	} else {
		return s
	}
}

// Partition partitions the immutableRange if the merge is already done or
// the underlying ranges otherwise. This operation can be O(n.log(n)) in the
// size of the unsorted range, if there is one but is O(log(n)) otherwise.
func (r *mergeableRange) Partition(e Element, o Order) (SortedRange, SortedRange) {
	r.mu.Lock()
	if r.left != nil && r.right == nil {
		r.right = r.unsorted.freeze()
	}
	if r.left == nil {
		return r.immutableRange.Partition(e, o)
	}
	r.mu.Unlock()

	r.mu.RLock()
	defer r.mu.RUnlock()

	p1, p2 := r.left.Partition(e, o)
	p3, p4 := r.right.Partition(e, o)

	selectFirst := func(p, q SortedRange) Element {
		if p.Limit() == 0 {
			return q.First()
		} else if q.Limit() == 0 {
			return p.First()
		}
		if o(p.First(), q.First()) {
			return p.First()
		} else {
			return q.First()
		}
	}

	selectLast := func(p, q SortedRange) Element {
		if p.Limit() == 0 {
			return q.Last()
		} else if q.Limit() == 0 {
			return p.Last()
		}
		if o(q.Last(), p.Last()) {
			return p.Last()
		} else {
			return q.Last()
		}
	}

	return useEmptyRangeIfEmpty(newMergeableRange(selectFirst(p1, p3), selectLast(p1, p3), p1, p3, nil)),
		useEmptyRangeIfEmpty(newMergeableRange(selectFirst(p2, p4), selectLast(p2, p4), p2, p4, nil))
}

func (r *mergeableRange) String() string {
	return "mergableRange" + r.basicRange.String()
}

// mergeableCursor knows how to iterate across the merged part of a
// mergeableRange and then advance merging operation by one element
// so as to advance the extent of the merged part of the mergeable
// range.
type mergeableCursor struct {
	basicCursor
	mu      *sync.RWMutex
	nx      *int
	advance func()
}

// Next() checks to see if the cursor has reached the end of the
// merged region. If so, it acquires a write lock and advances
// the cursor by one.
func (c *mergeableCursor) Next() Element {
	c.mu.RLock()

	if c.next == *c.nx {
		c.mu.RUnlock()

		c.mu.Lock()
		c.advance()
		c.mu.Unlock()

		c.mu.RLock()
		if c.next == *c.nx {
			c.mu.RUnlock()
			return nil
		}
	}

	defer func() {
		c.next++
		c.mu.RUnlock()
	}()
	return c.elements[c.next]
}

// Fill copies the merged part of the range with a slice copy, then iterates
// over the unmerged part.
func (c *mergeableCursor) Fill(buffer []Element) int {
	max := len(buffer)
	next := 0

	// efficiently copy the merged part of the range
	c.mu.RLock()
	limit := *c.nx - c.next
	if limit > max {
		limit = max
	}
	if limit > 0 {
		copy(buffer, c.elements[c.next:c.next+limit])
	}
	c.mu.RUnlock()

	next = limit
	c.next += limit

	// iterate over the remainder of the range, upto the
	// the buffer limit
	for next < max {
		e := c.Next()
		if e != nil {
			buffer[next] = e
			next++
		} else {
			return next
		}
	}
	return max
}

// mergeCursor is a cursor that iterates over the sorted, deduplicated elements
// of an underlying cursor with a lookahead of 1. It is used to iterate
// over one of the arms of a merge.
type mergeCursor struct {
	underlying Cursor
	peeked     Element
	done       bool
}

// fill ensures that mc.peeked is not empty or mc.done is true.
func (mc *mergeCursor) fill() {
	if !mc.done {
		mc.peeked = mc.underlying.Next()
		mc.done = (mc.peeked == nil)
	}
}

// next returns the next element in the sequence, clearing the
// peeked element if there was one.
func (mc *mergeCursor) next() Element {
	if mc.peeked == nil {
		mc.fill()
	}
	peek := mc.peeked
	mc.peeked = nil
	return peek
}

// peek returns the next element in the sequence without
// consuming it.
func (mc *mergeCursor) peek() Element {
	if mc == nil {
		panic("illegal state: mc is nil")
	}
	if mc.peeked == nil {
		mc.fill()
	}
	return mc.peeked
}
