package tsl

import (
	"fmt"
)

// disjointRanges represents a slice of SortedRanges which do not
// overlap. Sorted iteration across disjoint ranges is quicker than
// iteration across overlapping ranges since it avoids the need
// for a comparison in order to select the cursor to read from next.
type disjointRanges struct {
	first    Element
	last     Element
	segments []SortedRange
}

func (d *disjointRanges) First() Element {
	return d.first
}

func (d *disjointRanges) Last() Element {
	return d.last
}

func (d *disjointRanges) Open() Cursor {
	return &disjointCursor{
		next:     0,
		cursor:   d.segments[0].Open(),
		segments: d.segments,
	}
}

func (d *disjointRanges) Limit() int {
	limit := 0
	for _, r := range d.segments {
		limit += r.Limit()
	}
	return limit
}

func (d *disjointRanges) Partition(e Element, o Order) (SortedRange, SortedRange) {
	if !o(d.first, e) {
		return EmptyRange, d
	}

	for i, r := range d.segments {
		if r.First() == nil {
			panic("r.First() is nil!")
		}
		if r.Last() == nil {
			panic("r.Last() is nil!")
		}
		if o(r.First(), e) && o(e, r.Last()) {
			p1, p2 := r.Partition(e, o)
			var r1, r2 SortedRange
			if p1.Limit() == 0 {
				if p2.First() == nil {
					panic("p2.First() is nil!")
				}
				r1, r2 = &disjointRanges{
					first:    d.first,
					last:     d.segments[i-1].Last(),
					segments: d.segments[0:i],
				}, &disjointRanges{
					first:    p2.First(),
					last:     d.last,
					segments: append([]SortedRange{p2}, d.segments[i+1:]...),
				}
				return EmptyRange, p2
			} else if p2.Limit() == 0 {
				if p1.Last() == nil {
					panic("p1.Last() is nil!")
				}
				r1, r2 = &disjointRanges{
					first:    d.first,
					last:     p1.Last(),
					segments: append(d.segments[0:i], p1),
				}, &disjointRanges{
					first:    d.segments[i+1].First(),
					last:     d.last,
					segments: d.segments[i+1:],
				}
			} else {
				if p1.Last() == nil {
					panic("p1.Last() is nil!")
				}
				if p2.First() == nil {
					panic("p2.First() is nil!")
				}
				r1, r2 = &disjointRanges{
					first:    d.first,
					last:     p1.Last(),
					segments: append(d.segments[0:i], p1),
				}, &disjointRanges{
					first:    p2.First(),
					last:     d.last,
					segments: append([]SortedRange{p2}, d.segments[i+1:]...),
				}
			}
			return r1, r2
		}
	}

	return d, EmptyRange
}

type disjointCursor struct {
	next     int
	cursor   Cursor
	segments []SortedRange
}

func (c *disjointCursor) nextCursor() Cursor {
	c.next++
	if c.next < len(c.segments) {
		return c.segments[c.next].Open()
	} else {
		return nil
	}
}

func (c *disjointCursor) Next() Element {
	var next Element
	for c.cursor != nil {
		next = c.cursor.Next()
		if next == nil {
			c.cursor = c.nextCursor()
		} else {
			break
		}
	}
	return next
}

func (c *disjointCursor) Fill(buffer []Element) int {
	max := len(buffer)
	next := 0
	for next < max && c.cursor != nil {
		filled := c.cursor.Fill(buffer[next:max])
		next += filled
		if next < max {
			c.cursor = c.nextCursor()
		}
	}
	return next
}

func merge(a SortedRange, b SortedRange) SortedRange {
	if a.Limit() == 0 {
		return b
	} else if b.Limit() == 0 {
		return a
	} else if a.Last().Less(b.First()) {
		var last Element
		if !b.Last().Less(a.Last()) {
			last = b.Last()
		} else {
			last = a.Last()
		}
		return &disjointRanges{
			first:    a.First(),
			last:     last,
			segments: []SortedRange{a, b},
		}
	} else {
		var first, last Element
		if a.First().Less(b.First()) {
			first = a.First()
		} else {
			first = b.First()
		}

		if !b.Last().Less(a.Last()) {
			last = b.Last()
		} else {
			last = a.Last()
		}

		p1, p2 := a.Partition(b.First(), LessOrder)
		p3, p4 := b.Partition(a.Last(), LessOrEqualOrder)
		m23 := useEmptyRangeIfEmpty(newMergeableRange(p2.First(), p3.Last(), p2, p3, nil))
		if m23 == EmptyRange {
			if p1.Limit() == 0 {
				return p4
			} else if p4.Limit() == 0 {
				return p1
			}
			return &disjointRanges{
				first: first,
				last:  last,
				segments: []SortedRange{
					p1,
					p4,
				},
			}
		} else {
			if p1.Limit() == 0 {
				if p4.Limit() == 0 {
					return m23
				}
				return &disjointRanges{
					first: first,
					last:  last,
					segments: []SortedRange{
						m23,
						p4,
					},
				}
			} else if p4.Limit() == 0 {
				if p1.Limit() == 0 {
					return m23
				}
				return &disjointRanges{
					first: first,
					last:  last,
					segments: []SortedRange{
						p1,
						m23,
					},
				}
			}
			return &disjointRanges{
				first: first,
				last:  last,
				segments: []SortedRange{
					p1,
					m23,
					p4,
				},
			}
		}

	}
}

func (d *disjointRanges) String() string {
	buf := fmt.Sprintf("disjointRanges{first: %v, last: %v, segments: [", d.first, d.last)
	for i, s := range d.segments {
		if i > 0 {
			buf = buf + ","
		}
		buf = buf + fmt.Sprintf("%v", s)
	}
	buf = buf + "]}"
	return buf
}
