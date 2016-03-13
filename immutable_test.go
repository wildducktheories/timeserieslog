package tsl

import (
	"reflect"
	"testing"
)

func Test_ImmutableRange_Empty(t *testing.T) {
	d := newImmutableRange(NewElements([]int{}))

	pivot := 2
	left, right := d.Partition(intElement{pivot}, LessOrder)
	expectedLeft := emptyRange
	expectedRight := emptyRange
	expected := []Range{expectedLeft, expectedRight}
	got := []Range{left, right}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("partition of %d at %d failed. got: %+v, expected: %+v", d, pivot, got, expected)
	}
	if err := checkSortedRangeInvariants(left); err != nil {
		t.Fatalf("got: %v. %v", left, err)
	}
	if err := checkSortedRangeInvariants(right); err != nil {
		t.Fatalf("got: %v. %v", right, err)
	}
}

func Test_ImmutableRange_Pivot_2(t *testing.T) {
	d := newImmutableRange(NewElements([]int{0, 1, 3, 4}))

	pivot := 2
	left, right := d.Partition(intElement{pivot}, LessOrder)
	expectedLeft := newImmutableRange(NewElements([]int{0, 1}))
	expectedRight := newImmutableRange(NewElements([]int{3, 4}))
	expected := []Range{expectedLeft, expectedRight}
	got := []Range{left, right}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("partition of %d at %d failed. got: %+v, expected: %+v", d, pivot, got, expected)
	}
	if err := checkSortedRangeInvariants(left); err != nil {
		t.Fatalf("got: %v. %v", left, err)
	}
	if err := checkSortedRangeInvariants(right); err != nil {
		t.Fatalf("got: %v. %v", right, err)
	}
}

func Test_ImmutableRange_Pivot_1(t *testing.T) {
	d := newImmutableRange(NewElements([]int{0, 1, 3, 4}))

	pivot := 1
	left, right := d.Partition(intElement{pivot}, LessOrder)
	expectedLeft := newImmutableRange(NewElements([]int{0}))
	expectedRight := newImmutableRange(NewElements([]int{1, 3, 4}))
	expected := []Range{expectedLeft, expectedRight}
	got := []Range{left, right}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("partition of %d at %d failed. got: %+v, expected: %+v", d, pivot, got, expected)
	}
	if err := checkSortedRangeInvariants(left); err != nil {
		t.Fatalf("got: %v. %v", left, err)
	}
	if err := checkSortedRangeInvariants(right); err != nil {
		t.Fatalf("got: %v. %v", right, err)
	}
}

func Test_ImmutableRange_Pivot_0(t *testing.T) {
	d := newImmutableRange(NewElements([]int{0, 1, 3, 4}))

	pivot := 0
	left, right := d.Partition(intElement{pivot}, LessOrder)
	got := []Range{left, right}
	expectedLeft := emptyRange
	expectedRight := newImmutableRange(NewElements([]int{0, 1, 3, 4}))
	expected := []Range{expectedLeft, expectedRight}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("partition of %d at %d failed. got: %+v, expected: %+v", d, pivot, got, expected)
	}
	if err := checkSortedRangeInvariants(left); err != nil {
		t.Fatalf("got: %v. %v", left, err)
	}
	if err := checkSortedRangeInvariants(right); err != nil {
		t.Fatalf("got: %v. %v", right, err)
	}
}

func Test_ImmutableRange_Pivot_Minus1(t *testing.T) {
	d := newImmutableRange(NewElements([]int{0, 1, 3, 4}))

	pivot := -1
	left, right := d.Partition(intElement{pivot}, LessOrder)
	got := []Range{left, right}
	expectedLeft := emptyRange
	expectedRight := newImmutableRange(NewElements([]int{0, 1, 3, 4}))
	expected := []Range{expectedLeft, expectedRight}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("partition of %d at %d failed. got: %+v, expected: %+v", d, pivot, got, expected)
	}
	if err := checkSortedRangeInvariants(left); err != nil {
		t.Fatalf("got: %v. %v", left, err)
	}
	if err := checkSortedRangeInvariants(right); err != nil {
		t.Fatalf("got: %v. %v", right, err)
	}
}

func Test_ImmutableRange_Pivot_5(t *testing.T) {
	d := newImmutableRange(NewElements([]int{0, 1, 3, 4}))

	pivot := 5
	left, right := d.Partition(intElement{pivot}, LessOrder)
	got := []Range{left, right}
	expectedLeft := newImmutableRange(NewElements([]int{0, 1, 3, 4}))
	expectedRight := emptyRange
	expected := []Range{expectedLeft, expectedRight}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("partition of %d at %d failed. got: %+v, expected: %+v", d, pivot, got, expected)
	}
	if err := checkSortedRangeInvariants(left); err != nil {
		t.Fatalf("got: %v. %v", left, err)
	}
	if err := checkSortedRangeInvariants(right); err != nil {
		t.Fatalf("got: %v. %v", right, err)
	}
}

func Test_ImmutableRange_Pivot_3_Odd(t *testing.T) {
	d := newImmutableRange(NewElements([]int{0, 1, 3, 4, 5}))

	pivot := 3
	left, right := d.Partition(intElement{pivot}, LessOrder)
	got := []Range{left, right}
	expectedLeft := newImmutableRange(NewElements([]int{0, 1}))
	expectedRight := newImmutableRange(NewElements([]int{3, 4, 5}))
	expected := []Range{expectedLeft, expectedRight}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("partition of %d at %d failed. got: %+v, expected: %+v", d, pivot, got, expected)
	}
	if err := checkSortedRangeInvariants(left); err != nil {
		t.Fatalf("got: %v. %v", left, err)
	}
	if err := checkSortedRangeInvariants(right); err != nil {
		t.Fatalf("got: %v. %v", right, err)
	}
}

func Test_ImmutableRange_Pivot_2_Odd(t *testing.T) {
	d := newImmutableRange(NewElements([]int{0, 1, 3, 4, 5}))

	pivot := 2
	left, right := d.Partition(intElement{pivot}, LessOrder)
	got := []Range{left, right}
	expectedLeft := newImmutableRange(NewElements([]int{0, 1}))
	expectedRight := newImmutableRange(NewElements([]int{3, 4, 5}))
	expected := []Range{expectedLeft, expectedRight}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("partition of %d at %d failed. got: %+v, expected: %+v", d, pivot, got, expected)
	}
	if err := checkSortedRangeInvariants(left); err != nil {
		t.Fatalf("got: %v. %v", left, err)
	}
	if err := checkSortedRangeInvariants(right); err != nil {
		t.Fatalf("got: %v. %v", right, err)
	}
}

func Test_ImmutableRange_Pivot_1_Odd(t *testing.T) {
	d := newImmutableRange(NewElements([]int{0, 1, 3, 4, 5}))

	pivot := 1
	left, right := d.Partition(intElement{pivot}, LessOrder)
	got := []Range{left, right}
	expectedLeft := newImmutableRange(NewElements([]int{0}))
	expectedRight := newImmutableRange(NewElements([]int{1, 3, 4, 5}))
	expected := []Range{expectedLeft, expectedRight}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("partition of %d at %d failed. got: %+v, expected: %+v", d, pivot, got, expected)
	}
	if err := checkSortedRangeInvariants(left); err != nil {
		t.Fatalf("got: %v. %v", left, err)
	}
	if err := checkSortedRangeInvariants(right); err != nil {
		t.Fatalf("got: %v. %v", right, err)
	}
}

func Test_ImmutableRange_Pivot_1_Odd_LessOrEqual(t *testing.T) {
	d := newImmutableRange(NewElements([]int{0, 1, 3, 4, 5}))

	pivot := 1
	left, right := d.Partition(intElement{pivot}, LessOrEqualOrder)
	got := []Range{left, right}
	expectedLeft := newImmutableRange(NewElements([]int{0, 1}))
	expectedRight := newImmutableRange(NewElements([]int{3, 4, 5}))
	expected := []Range{expectedLeft, expectedRight}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("partition of %d at %d failed. got: %+v, expected: %+v", d, pivot, got, expected)
	}
	if err := checkSortedRangeInvariants(left); err != nil {
		t.Fatalf("got: %v. %v", left, err)
	}
	if err := checkSortedRangeInvariants(right); err != nil {
		t.Fatalf("got: %v. %v", right, err)
	}
}

func Test_ImmutableRange_Pivot_5_LessOrEqual(t *testing.T) {
	d := newImmutableRange(NewElements([]int{0, 1, 3, 4}))

	pivot := 5
	left, right := d.Partition(intElement{pivot}, LessOrEqualOrder)
	got := []Range{left, right}
	expectedLeft := newImmutableRange(NewElements([]int{0, 1, 3, 4}))
	expectedRight := emptyRange
	expected := []Range{expectedLeft, expectedRight}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("partition of %d at %d failed. got: %+v, expected: %+v", d, pivot, got, expected)
	}
	if err := checkSortedRangeInvariants(left); err != nil {
		t.Fatalf("got: %v. %v", left, err)
	}
	if err := checkSortedRangeInvariants(right); err != nil {
		t.Fatalf("got: %v. %v", right, err)
	}
}

func Test_ImmutableRange_Pivot_4_LessOrEqual(t *testing.T) {
	d := newImmutableRange(NewElements([]int{0, 1, 3, 4}))

	pivot := 4
	left, right := d.Partition(intElement{pivot}, LessOrEqualOrder)
	got := []Range{left, right}
	expectedLeft := newImmutableRange(NewElements([]int{0, 1, 3, 4}))
	expectedRight := emptyRange
	expected := []Range{expectedLeft, expectedRight}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("partition of %d at %d failed. got: %+v, expected: %+v", d, pivot, got, expected)
	}
	if err := checkSortedRangeInvariants(left); err != nil {
		t.Fatalf("got: %v. %v", left, err)
	}
	if err := checkSortedRangeInvariants(right); err != nil {
		t.Fatalf("got: %v. %v", right, err)
	}
}

func Test_ImmutableRange_Pivot_Minus1_Odd_LessOrEqual(t *testing.T) {
	d := newImmutableRange(NewElements([]int{0, 1, 3, 4, 5}))

	pivot := -1
	left, right := d.Partition(intElement{pivot}, LessOrEqualOrder)
	got := []Range{left, right}
	expectedLeft := emptyRange
	expectedRight := newImmutableRange(NewElements([]int{0, 1, 3, 4, 5}))
	expected := []Range{expectedLeft, expectedRight}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("partition of %d at %d failed. got: %+v, expected: %+v", d, pivot, got, expected)
	}
	if err := checkSortedRangeInvariants(left); err != nil {
		t.Fatalf("got: %v. %v", left, err)
	}
	if err := checkSortedRangeInvariants(right); err != nil {
		t.Fatalf("got: %v. %v", right, err)
	}
}

func Test_ImmutableRange_Fill(t *testing.T) {
	expected := NewElements([]int{0, 1, 3, 4, 5})
	d := newImmutableRange(expected)
	got := Elements(make([]Element, d.Limit(), d.Limit()))
	d.Open().Fill(got)
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("fill failed. got: %v, expected: %v", got, expected)
	}
}
