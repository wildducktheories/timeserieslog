package tsl

import (
	"reflect"
	"testing"
)

func Test_Merge_Empty(t *testing.T) {
	got := Merge(EmptyRange, EmptyRange)
	expected := EmptyRange
	if got != expected {
		t.Fatalf("Merge empty + empty -> yield empty. got: %v, expected :%v", got, expected)
	}
	if err := checkSortedRangeInvariants(got); err != nil {
		t.Fatalf("got: %v. %v", got, err)
	}
}

func Test_Merge_Empty_Single(t *testing.T) {
	d := newImmutableRange(NewElements([]int{0}))
	got := Merge(EmptyRange, d)
	expected := d
	if got != expected {
		t.Fatalf("Merge empty + single -> yield single. got: %v, expected :%v", got, expected)
	}
	if err := checkSortedRangeInvariants(got); err != nil {
		t.Fatalf("got: %v. %v", got, err)
	}
}

func Test_Merge_Single_Empty(t *testing.T) {
	d := newImmutableRange(NewElements([]int{0}))
	got := Merge(d, EmptyRange)
	expected := d
	if got != expected {
		t.Fatalf("Merge single + empty -> yield single. got: %v, expected :%v", got, expected)
	}
	if err := checkSortedRangeInvariants(got); err != nil {
		t.Fatalf("got: %v. %v", got, err)
	}
}

func Test_Merge_Single_Single_Disjoint_InOrder(t *testing.T) {
	a := newImmutableRange(NewElements([]int{1}))
	b := newImmutableRange(NewElements([]int{2}))
	c := NewElements([]int{1, 2})
	got := Merge(a, b)
	expected := c
	if !reflect.DeepEqual(Elements(AsSlice(got)), expected) {
		t.Fatalf("Merge single + single -> yield pair. got: %v, expected :%v", AsSlice(got), expected)
	}
	if err := checkSortedRangeInvariants(got); err != nil {
		t.Fatalf("got: %v. %v", got, err)
	}
}

func Test_Merge_Single_Single_Disjoint_OutOfOrder(t *testing.T) {
	a := newImmutableRange(NewElements([]int{2}))
	b := newImmutableRange(NewElements([]int{1}))
	c := NewElements([]int{1, 2})
	got := Merge(a, b)
	expected := c
	if !reflect.DeepEqual(Elements(AsSlice(got)), expected) {
		t.Fatalf("Merge single + single -> yield pair. got: %v, expected :%v", AsSlice(got), expected)
	}
	if err := checkSortedRangeInvariants(got); err != nil {
		t.Fatalf("got: %v. %v", got, err)
	}
}

func Test_Merge_Pair_Pair_Disjoint_InOrder(t *testing.T) {
	a := newImmutableRange(NewElements([]int{1, 2}))
	b := newImmutableRange(NewElements([]int{3, 4}))
	c := NewElements([]int{1, 2, 3, 4})
	got := Merge(a, b)
	expected := c
	if !reflect.DeepEqual(Elements(AsSlice(got)), expected) {
		t.Fatalf("Merge pair + pair -> yield quadruple. got: %v, expected :%v", AsSlice(got), expected)
	}
	if err := checkSortedRangeInvariants(got); err != nil {
		t.Fatalf("got: %v. %v", got, err)
	}
}

func Test_Merge_Pair_Pair_Overlapping_InOrder(t *testing.T) {
	a := newImmutableRange(NewElements([]int{1, 3}))
	b := newImmutableRange(NewElements([]int{2, 4}))
	c := NewElements([]int{1, 2, 3, 4})
	got := Merge(a, b)
	expected := c
	if !reflect.DeepEqual(Elements(AsSlice(got)), expected) {
		t.Fatalf("Merge pair + pair -> yield quadruple. got: %v, expected :%v", AsSlice(got), expected)
	}
	if err := checkSortedRangeInvariants(got); err != nil {
		t.Fatalf("got: %v. %v", got, err)
	}
}

func Test_Merge_Pair_Pair_Overlapping_OutOfOrder(t *testing.T) {
	a := newImmutableRange(NewElements([]int{2, 4}))
	b := newImmutableRange(NewElements([]int{1, 3}))
	c := NewElements([]int{1, 2, 3, 4})
	got := Merge(a, b)
	expected := c
	if !reflect.DeepEqual(Elements(AsSlice(got)), expected) {
		t.Fatalf("Merge pair + pair -> yield quadruple. got: %v, expected :%v", AsSlice(got), expected)
	}
	if err := checkSortedRangeInvariants(got); err != nil {
		t.Fatalf("got: %v. %v", got, err)
	}
}

func Test_Merge_Pair_Pair_Overlapping_Duplicate(t *testing.T) {
	a := newImmutableRange(NewElements([]int{1, 2}))
	b := newImmutableRange(NewElements([]int{2, 3}))
	c := NewElements([]int{1, 2, 3})
	got := Merge(a, b)
	expected := c
	if !reflect.DeepEqual(Elements(AsSlice(got)), expected) {
		t.Fatalf("Merge pair + pair -> yield triple. got: %v, expected :%v", AsSlice(got), expected)
	}
	if err := checkSortedRangeInvariants(got); err != nil {
		t.Fatalf("got: %v. %v", got, err)
	}
}

func Test_Merge_Pair_Pair_Partition(t *testing.T) {
	a := newImmutableRange(NewElements([]int{1, 3}))
	b := newImmutableRange(NewElements([]int{2, 4}))
	c := NewElements([]int{1, 2, 3, 4})
	got := Merge(a, b)
	if err := checkSortedRangeInvariants(got); err != nil {
		t.Fatalf("got: %v", err)
	}
	expected := c
	if !reflect.DeepEqual(Elements(AsSlice(got)), expected) {
		t.Fatalf("Merge pair + pair -> yield quadruple. got: %v, expected :%v", AsSlice(got), expected)
	}

	p1, p2 := got.Partition(intElement{3}, LessOrEqualOrder)
	if err := checkSortedRangeInvariants(p1); err != nil {
		t.Fatalf("p1: %v", err)
	}
	if err := checkSortedRangeInvariants(p2); err != nil {
		t.Fatalf("p2: %v, got: %v", err, p2)
	}
	expected = NewElements([]int{1, 2, 3})
	got = p1
	if !reflect.DeepEqual(Elements(AsSlice(got)), expected) {
		t.Fatalf("partition p1 got: %v, expected :%v", AsSlice(got), expected)
	}

	expected = NewElements([]int{4})
	got = p2
	if !reflect.DeepEqual(Elements(AsSlice(got)), expected) {
		t.Fatalf("partition p2 got: %v, expected :%v", AsSlice(got), expected)
	}
}

func Test_Merge_Pair1_Pair3_Partition_0(t *testing.T) {
	a := newImmutableRange(NewElements([]int{0, 1}))
	b := newImmutableRange(NewElements([]int{4, 5}))
	aa, bb := Merge(a, b).Partition(intElement{-1}, LessOrder)

	got := Elements(AsSlice(aa))
	expected := NewElements([]int{})
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("partition got: %v, expected :%v", got, expected)
	}

	got = Elements(AsSlice(bb))
	expected = NewElements([]int{0, 1, 4, 5})
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("partition got: %v, expected :%v", got, expected)
	}
}

func Test_Merge_Pair1_Pair3_Partition_2(t *testing.T) {
	a := newImmutableRange(NewElements([]int{0, 1}))
	b := newImmutableRange(NewElements([]int{4, 5}))
	aa, bb := Merge(a, b).Partition(intElement{2}, LessOrder)

	got := Elements(AsSlice(aa))
	expected := NewElements([]int{0, 1})
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("partition got: %v, expected :%v", got, expected)
	}

	got = Elements(AsSlice(bb))
	expected = NewElements([]int{4, 5})
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("partition got: %v, expected :%v", got, expected)
	}
}

func Test_Merge_Pair1_Pair3_Partition_4(t *testing.T) {
	a := newImmutableRange(NewElements([]int{0, 1}))
	b := newImmutableRange(NewElements([]int{4, 5}))
	aa, bb := Merge(a, b).Partition(intElement{4}, LessOrEqualOrder)

	got := Elements(AsSlice(aa))
	expected := NewElements([]int{0, 1, 4})
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("partition got: %v, expected :%v", got, expected)
	}

	got = Elements(AsSlice(bb))
	expected = NewElements([]int{5})
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("partition got: %v, expected :%v", got, expected)
	}
}

func Test_Merge_Pair1_Pair3_Partition_5(t *testing.T) {
	a := newImmutableRange(NewElements([]int{0, 1}))
	b := newImmutableRange(NewElements([]int{4, 5}))
	aa, bb := Merge(a, b).Partition(intElement{5}, LessOrEqualOrder)

	got := Elements(AsSlice(aa))
	expected := NewElements([]int{0, 1, 4, 5})
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("partition got: %v, expected :%v", got, expected)
	}

	got = Elements(AsSlice(bb))
	expected = NewElements([]int{})
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("partition got: %v, expected :%v", got, expected)
	}
}

func Test_Merge_Pair1_Pair3_Pair2(t *testing.T) {
	a := newImmutableRange(NewElements([]int{0, 1}))
	b := newImmutableRange(NewElements([]int{4, 5}))
	c := newImmutableRange(NewElements([]int{2, 3}))
	got := Merge(Merge(a, b), c)

	expected := NewElements([]int{0, 1, 2, 3, 4, 5})
	if !reflect.DeepEqual(Elements(AsSlice(got)), expected) {
		t.Fatalf("partition got: %v, expected :%v", AsSlice(got), expected)
	}
}

func Test_Merge_Pair1_Pair2_Pair3(t *testing.T) {
	a := newImmutableRange(NewElements([]int{0, 1}))
	b := newImmutableRange(NewElements([]int{2, 3}))
	c := newImmutableRange(NewElements([]int{4, 5}))
	got := Merge(Merge(a, b), c)

	expected := NewElements([]int{0, 1, 2, 3, 4, 5})
	if !reflect.DeepEqual(Elements(AsSlice(got)), expected) {
		t.Fatalf("partition got: %v, expected :%v", AsSlice(got), expected)
	}
}
