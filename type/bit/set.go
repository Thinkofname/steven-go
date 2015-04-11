package bit

// Set is a collection of booleans stored as bits
type Set []uint64

// NewSet allocates a new bit set that can store up to the
// passed number of bits.
func NewSet(size int) Set {
	return make(Set, (size+63)>>6)
}

// Set changes the value of the bit at the location.
func (s Set) Set(i int, v bool) {
	if v {
		s[i>>6] |= 1 << uint(i&0x3F)
	} else {
		s[i>>6] &= ^(1 << uint(i&0x3F))
	}
}

// Get returns the value of the bit at the location
func (s Set) Get(i int) bool {
	v := s[i>>6] & (1 << uint(i&0x3F))
	return v != 0
}
