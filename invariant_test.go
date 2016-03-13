package tsl

import (
	"fmt"
)

func checkRangeInvariants(r Range) error {
	if r.Limit() == 0 {
		if r.First() != nil {
			return fmt.Errorf("first should be nil when limit is zero. got %v, expected: nil", r.First())
		}
		if r.Last() != nil {
			return fmt.Errorf("last should be nil when limit is zero. got %v, expected: nil", r.Last())
		}
	} else {
		if r.First() == nil {
			return fmt.Errorf("first should not be nil when limit is not zero. got %v, expected: not nil", r.First())
		}
		if r.Last() == nil {
			return fmt.Errorf("last should not be nil when limit is not zero. got %v, expected: not nil", r.Last())
		}
		firstLessThanLast := r.First().Less(r.Last())
		lastLessThanFirst := r.Last().Less(r.First())

		if !firstLessThanLast {
			if lastLessThanFirst {
				return fmt.Errorf("if first is not less than last, then last must not be less than first. got: %v, expected: false", lastLessThanFirst)
			}
		}
	}
	return nil
}

func checkSortedRangeInvariants(r SortedRange) error {
	if err := checkRangeInvariants(r); err != nil {
		return err
	} else {
		initialLimit := r.Limit()
		c := r.Open()
		if c.Next() == nil {
			if initialLimit != 0 {
				return fmt.Errorf("if c.Next() on fresh cursor is nil, then r.Limit() must be zero. got: %d, expected: 0", initialLimit)
			}
		} else {
			if initialLimit == 0 {
				return fmt.Errorf("if c.Next() is on fresh cursor is not nil, then r.Limit() must be non zero. got: %d, expected: >0", initialLimit)
			}
		}
		slice := AsSlice(r)
		finalLimit := r.Limit()
		if finalLimit > initialLimit {
			return fmt.Errorf("final limit must always be less than or equal to initial limit. got %d, expected: %d", finalLimit, initialLimit)
		}
		if finalLimit != len(slice) {
			return fmt.Errorf("final limit must be identical to actual size of SortedRange. got: %d, expected: %d", finalLimit, len(slice))
		}
		for i := 1; i < len(slice); i++ {
			if !LessOrder(slice[i-1], slice[i]) {
				fmt.Errorf("adjacent elements must always satisfy LessOrder. got: false. expected: true. i: %d", i)
			}
		}
		if len(slice) > 0 {
			if r.First() != slice[0] {
				return fmt.Errorf("first element of slice does not equal c.First(). got: %v, expected: %v.", slice[0], r.First())
			}
			if r.Last() != slice[len(slice)-1] {
				return fmt.Errorf("last element of slice does not equal c.Last(). got: %v, expected: %v.", slice[len(slice)-1], r.Last())
			}
		}
		return nil
	}
}
