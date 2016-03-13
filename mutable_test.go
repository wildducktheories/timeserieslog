package tsl

import (
	"reflect"
	"testing"
)

func Test_MutableRange_Empty(t *testing.T) {
	d := NewElements([]int{})
	r := &mutableRange{}
	r.Add(d)
	cursor := r.Freeze().Open()
	got := Elements([]Element{})
	for next := cursor.Next(); next != nil; next = cursor.Next() {
		got = append(got, next)
	}
	expected := NewElements([]int{})
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("sort failed got: %+v, expected: %+v", got, expected)
	}
}

func Test_MutableRange_Single(t *testing.T) {
	d := NewElements([]int{0})
	r := &mutableRange{}
	r.Add(d)
	cursor := r.Freeze().Open()
	got := Elements([]Element{})
	for next := cursor.Next(); next != nil; next = cursor.Next() {
		got = append(got, next)
	}
	expected := NewElements([]int{0})
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("sort failed got: %+v, expected: %+v", got, expected)
	}
}

func Test_MutableRange_Duplicated(t *testing.T) {
	d := NewElements([]int{0, 0})
	r := &mutableRange{}
	r.Add(d)
	cursor := r.Freeze().Open()
	got := Elements([]Element{})
	for next := cursor.Next(); next != nil; next = cursor.Next() {
		got = append(got, next)
	}
	expected := NewElements([]int{0})
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("sort failed got: %+v, expected: %+v", got, expected)
	}
}

func Test_MutableRange_Two(t *testing.T) {
	d := NewElements([]int{0, 1})
	r := &mutableRange{}
	r.Add(d)
	cursor := r.Freeze().Open()
	got := Elements([]Element{})
	for next := cursor.Next(); next != nil; next = cursor.Next() {
		got = append(got, next)
	}
	expected := NewElements([]int{0, 1})
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("sort failed got: %+v, expected: %+v", got, expected)
	}
}

func Test_MutableRange_Reversed(t *testing.T) {
	d := NewElements([]int{1, 0})
	r := &mutableRange{}
	r.Add(d)
	cursor := r.Freeze().Open()
	got := Elements([]Element{})
	for next := cursor.Next(); next != nil; next = cursor.Next() {
		got = append(got, next)
	}
	expected := NewElements([]int{0, 1})
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("sort failed got: %+v, expected: %+v", got, expected)
	}
}

func Test_MutableRange_Triple(t *testing.T) {
	d := NewElements([]int{0, 1, 2})
	r := &mutableRange{}
	r.Add(d)
	cursor := r.Freeze().Open()
	got := Elements([]Element{})
	for next := cursor.Next(); next != nil; next = cursor.Next() {
		got = append(got, next)
	}
	expected := NewElements([]int{0, 1, 2})
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("sort failed got: %+v, expected: %+v", got, expected)
	}
}

func Test_MutableRange_TripleReversed(t *testing.T) {
	d := NewElements([]int{2, 1, 0})
	r := &mutableRange{}
	r.Add(d)
	cursor := r.Freeze().Open()
	got := Elements([]Element{})
	for next := cursor.Next(); next != nil; next = cursor.Next() {
		got = append(got, next)
	}
	expected := NewElements([]int{0, 1, 2})
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("sort failed got: %+v, expected: %+v", got, expected)
	}
}

func Test_MutableRange_TriplePermuted(t *testing.T) {
	d := NewElements([]int{2, 0, 1})
	r := &mutableRange{}
	r.Add(d)
	cursor := r.Freeze().Open()
	got := Elements([]Element{})
	for next := cursor.Next(); next != nil; next = cursor.Next() {
		got = append(got, next)
	}
	expected := NewElements([]int{0, 1, 2})
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("sort failed got: %+v, expected: %+v", got, expected)
	}
}

func Test_MutableRange(t *testing.T) {
	d := NewElements([]int{0, 2, 3, 4, 6, 6, 3, 2, 1, 5, 7})
	r := &mutableRange{}
	r.Add(d)
	frozen := r.Freeze()
	cursor := frozen.Open()
	got := Elements([]Element{})
	for next := cursor.Next(); next != nil; next = cursor.Next() {
		got = append(got, next)
	}
	expected := NewElements([]int{0, 1, 2, 3, 4, 5, 6, 7})
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("sort failed got: %+v, expected: %+v", got, expected)
	}

	cursor = frozen.Open()
	got = Elements([]Element{})
	for next := cursor.Next(); next != nil; next = cursor.Next() {
		got = append(got, next)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("sort failed got: %+v, expected: %+v", got, expected)
	}
}

func Test_MutableRange_Quintuple(t *testing.T) {
	d := NewElements([]int{0, 2, 1, 4, 3})
	r := &mutableRange{}
	r.Add(d)

	cursor := r.Freeze().Open()
	got := Elements([]Element{})
	for next := cursor.Next(); next != nil; next = cursor.Next() {
		got = append(got, next)
	}
	expected := NewElements([]int{0, 1, 2, 3, 4})
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("sort failed got: %+v, expected: %+v", got, expected)
	}
}

func Test_MutableRange_Partition_Less(t *testing.T) {
	d := NewElements([]int{1, 0})
	r := &mutableRange{}
	r.Add(d)
	p1, p2 := r.Freeze().Partition(intElement{1}, LessOrder)

	got := p1
	expected := NewElements([]int{0})
	if !reflect.DeepEqual(Elements(AsSlice(got)), expected) {
		t.Fatalf("partition failed. got %v, expected: %v", got, expected)
	}
	if err := checkSortedRangeInvariants(got); err != nil {
		t.Fatalf("got: %v. %v", got, err)
	}

	got = p2
	expected = NewElements([]int{1})
	if !reflect.DeepEqual(Elements(AsSlice(got)), expected) {
		t.Fatalf("partition failed. got %v, expected: %v", got, expected)
	}
	if err := checkSortedRangeInvariants(got); err != nil {
		t.Fatalf("got: %v. %v", got, err)
	}
}

func Test_MutableRange_Partition_LessOrEqual(t *testing.T) {
	d := NewElements([]int{1, 0})
	r := &mutableRange{}
	r.Add(d)
	p1, p2 := r.Freeze().Partition(intElement{1}, LessOrEqualOrder)

	got := p1
	expected := NewElements([]int{0, 1})
	if !reflect.DeepEqual(Elements(AsSlice(got)), expected) {
		t.Fatalf("partition failed p1. got %v, expected: %v", got, expected)
	}
	if err := checkSortedRangeInvariants(got); err != nil {
		t.Fatalf("got: %v. %v", got, err)
	}

	got = p2
	if !reflect.DeepEqual(p2, EmptyRange) {
		t.Fatalf("partition failed p2. got %v, expected: %v", got, EmptyRange)
	}
	if err := checkSortedRangeInvariants(got); err != nil {
		t.Fatalf("got: %v. %v", got, err)
	}
}

func Test_MutableRange_1_0_3_2_Partition_2_LessOrEqual(t *testing.T) {
	d := NewElements([]int{1, 0, 3, 2})
	r := &mutableRange{}
	r.Add(d)
	p1, p2 := r.Freeze().Partition(intElement{2}, LessOrEqualOrder)

	got := p1
	expected := NewElements([]int{0, 1, 2})
	if !reflect.DeepEqual(Elements(AsSlice(got)), expected) {
		t.Fatalf("partition failed p1. got %v, expected: %v", got, expected)
	}
	if err := checkSortedRangeInvariants(got); err != nil {
		t.Fatalf("got: %v. %v", got, err)
	}

	got = p2
	expected = NewElements([]int{3})
	if !reflect.DeepEqual(Elements(AsSlice(got)), expected) {
		t.Fatalf("partition failed p1. got %v, expected: %v", got, expected)
	}
	if err := checkSortedRangeInvariants(got); err != nil {
		t.Fatalf("got: %v. %v", got, err)
	}
}

func Test_MutableRange_1_0_3_2_Partition_2_Less(t *testing.T) {
	d := NewElements([]int{1, 0, 3, 2})
	r := &mutableRange{}
	r.Add(d)
	p1, p2 := r.Freeze().Partition(intElement{2}, LessOrder)

	got := p1
	expected := NewElements([]int{0, 1})
	if !reflect.DeepEqual(Elements(AsSlice(got)), expected) {
		t.Fatalf("partition failed p1. got %v, expected: %v", got, expected)
	}
	if err := checkSortedRangeInvariants(got); err != nil {
		t.Fatalf("got: %v. %v", got, err)
	}

	got = p2
	expected = NewElements([]int{2, 3})
	if !reflect.DeepEqual(Elements(AsSlice(got)), expected) {
		t.Fatalf("partition failed p1. got %v, expected: %v", got, expected)
	}
	if err := checkSortedRangeInvariants(got); err != nil {
		t.Fatalf("got: %v. %v", got, err)
	}
}
