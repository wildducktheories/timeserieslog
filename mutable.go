package tsl

import (
	"sync"
)

// mutableRange is a basicRange that can be extended with additional
// elements. In addition to a slice of sorted elements, it
// also has an unsortedRange for elements that
// arrive out of order, a mutex to guard access to the mutable
// state and frozen reference which is initialized with the
// SortedRange that may used to access to the contents of the
// range after it is has been frozen.
type mutableRange struct {
	basicRange
	mu       sync.RWMutex
	unsorted unsortedRange
	frozen   SortedRange
}

func (r *mutableRange) First() Element {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.first
}

func (r *mutableRange) Last() Element {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.last
}

func (r *mutableRange) Limit() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.elements) + r.unsorted.Limit()
}

func (r *mutableRange) addOne(e Element) {
	if len(r.elements) > 0 {
		if !r.last.Less(e) {
			r.unsorted.add(e)
		} else {
			r.elements = append(r.elements, e)
			r.last = e
		}
		if e.Less(r.first) {
			r.first = e
		}
	} else {
		r.elements = []Element{e}
		r.first = e
		r.last = e
	}
}

func (r *mutableRange) Add(elements []Element) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.frozen != nil {
		return ErrAlreadyFrozen
	}

	for _, e := range elements {
		r.addOne(e)
	}
	return nil
}

func (r *mutableRange) Freeze() SortedRange {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.frozen == nil {
		if r.unsorted.Limit() > 0 {
			r.frozen = newMergeableRange(
				r.first,
				r.last,
				&immutableRange{
					basicRange: basicRange{
						first:    r.elements[0],
						last:     r.elements[len(r.elements)-1],
						elements: r.elements,
					},
				},
				nil,
				&r.unsorted)
		} else {
			r.frozen = &immutableRange{
				basicRange: basicRange{
					first:    r.first,
					last:     r.last,
					elements: r.elements,
				},
			}
		}
	}

	return r.frozen
}
