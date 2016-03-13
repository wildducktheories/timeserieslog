package tsl

import (
	"fmt"
)

// basicRange is a type that encapsulates the essential element of all Range
// implementations. That is, it has a first element, a last element and a limit
// which specifies an upper bound for the maximum number of elements in the range.
type basicRange struct {
	first    Element
	last     Element
	elements []Element
}

func (r *basicRange) First() Element {
	return r.first
}

func (r *basicRange) Last() Element {
	return r.last
}

func (r *basicRange) Limit() int {
	return len(r.elements)
}

func (r *basicRange) String() string {
	out := fmt.Sprintf("{first: %v, last: %v, elements: [", r.first, r.last)
	first := true
	for _, e := range r.elements {
		if !first {
			out = out + ", "
		} else {
			first = false
		}
		out = out + fmt.Sprintf("%v", e)
	}
	out = out + "]}"
	return out
}

// basicCursor represents a cursor over an immutable, sorted, deduplicated slice
// of elements.
type basicCursor struct {
	next     int
	elements []Element
}

func (c *basicCursor) Next() Element {
	if c.next < len(c.elements) {
		tmp := c.elements[c.next]
		c.next++
		return tmp
	} else {
		return nil
	}
}

func (c *basicCursor) Fill(buffer []Element) int {
	max := len(buffer)
	if max > (len(c.elements) - c.next) {
		max = len(c.elements) - c.next
	}
	copy(buffer, c.elements[c.next:c.next+max])
	c.next += max
	return max
}
