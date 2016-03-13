package tsl

import (
	"sort"
	"sync"
)

// unsortedRange is a basicRange that can be extended with
// additional unsorted elements. Freezing an unsorted range
// incurs an O(n.log(n)) sort.
type unsortedRange struct {
	basicRange
	mu     sync.RWMutex
	frozen *immutableRange
}

// add adds a single element to the unsorted range, updating
// the first and last members as appropriate.
func (r *unsortedRange) add(e Element) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.frozen != nil {
		panic("attempt to add to frozen range")
	}

	if len(r.elements) == 0 {
		r.first = e
		r.last = e
	} else {
		if e.Less(r.first) {
			r.first = e
		}
		if !e.Less(r.last) {
			r.last = e
		}
	}
	r.elements = append(r.elements, e)
}

// freeze locks the receiver to prevent further updates
// and creates an immutableRange from the sorted, deduplicated
// elements.
func (r *unsortedRange) freeze() *immutableRange {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.frozen == nil {
		sort.Stable(Elements(r.elements))
		r.deduplicate()
		r.frozen = &immutableRange{
			basicRange: basicRange{
				first:    r.first,
				last:     r.last,
				elements: r.elements,
			},
		}
	}

	return r.frozen
}

// deduplicate ensures on the the last of equal Elements
// is kept.
func (r *unsortedRange) deduplicate() {
	if len(r.elements) < 2 {
		return
	}
	j := 0
	for _, e := range r.elements {
		if j > 0 {
			if !r.elements[j-1].Less(e) {
				j--
			}
		}
		r.elements[j] = e
		j++
	}
	r.elements = r.elements[0:j]
}

func (r *unsortedRange) String() string {
	return "unsortedRange" + r.basicRange.String()
}
