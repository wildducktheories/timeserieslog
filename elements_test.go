package tsl

type intElement struct {
	value int
}

func (i intElement) Less(e Element) bool {
	return i.value < e.(intElement).value
}

func NewElements(values []int) Elements {
	result := make(Elements, len(values))
	for i, e := range values {
		result[i] = intElement{e}
	}
	return result
}
