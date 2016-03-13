package tsl

// Less returns true iff s[i].Less(s[j]) is true.
func (s Elements) Less(i, j int) bool {
	return s[i].Less(s[j])
}

// Swap swaps s[i] and s[j]
func (s Elements) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Len returns len(s)
func (s Elements) Len() int {
	return len(s)
}
