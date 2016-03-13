package tsl

var emptyRange *immutableRange

func init() {
	emptyRange = &immutableRange{
		basicRange: basicRange{
			first:    nil,
			last:     nil,
			elements: nil,
		},
	}
	EmptyRange = emptyRange
}

// immutableRange is a SortedRange whose Elements never change.
type immutableRange struct {
	basicRange
}

// newImmutableRange accepts a sorted slice of element and answers
// an immutableRange for that slice. It is the caller's responsibility
// to guarantee that the supplied slice is both immutable and
// sorted.
func newImmutableRange(sorted []Element) *immutableRange {
	if len(sorted) > 0 {
		return &immutableRange{
			basicRange: basicRange{
				first:    sorted[0],
				last:     sorted[len(sorted)-1],
				elements: sorted,
			},
		}
	} else {
		return emptyRange
	}
}

// Open a cursor that iterates over the immutable range.
func (r *immutableRange) Open() Cursor {
	return &basicCursor{
		next:     0,
		elements: r.elements,
	}
}

// Partition use a binary search to partition the receiver into a pair disjoint
// SortedRanges such that o(i, e) is true for each element i of the first member of the returned pair
// and o(i, e) is false for each element i of the second member of the returned pair.
func (r *immutableRange) Partition(e Element, o Order) (SortedRange, SortedRange) {
	if len(r.elements) == 0 {
		return emptyRange, emptyRange
	}
	found := func() int {
		lower := 0
		upper := len(r.elements)
		for lower < upper {
			mid := (lower + upper) / 2
			sortsLeft := o(r.elements[mid], e)
			sortsRight := o(e, r.elements[mid])
			if sortsRight && !sortsLeft {
				upper = mid - 1
			} else if sortsLeft && !sortsRight {
				lower = mid + 1
			} else {
				return mid
			}
		}
		return upper
	}()
	if found >= 0 && found < len(r.elements) && o(r.elements[found], e) {
		found = found + 1
	}
	if found <= 0 {
		return emptyRange, &immutableRange{
			basicRange: basicRange{
				first:    r.first,
				last:     r.last,
				elements: r.elements,
			},
		}
	} else if found >= len(r.elements) {
		return &immutableRange{
			basicRange: basicRange{
				first:    r.first,
				last:     r.last,
				elements: r.elements,
			},
		}, emptyRange
	} else {
		return &immutableRange{
				basicRange: basicRange{
					first:    r.elements[0],
					last:     r.elements[found-1],
					elements: r.elements[0:found],
				},
			}, &immutableRange{
				basicRange: basicRange{
					first:    r.elements[found],
					last:     r.elements[len(r.elements)-1],
					elements: r.elements[found:],
				},
			}
	}
}

func (r *immutableRange) String() string {
	return "immutableRange" + r.basicRange.String()
}
