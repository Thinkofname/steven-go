package nibble

// Array is an array of 4 bit values.
type Array []byte

// New creates a new nibble array that can store at least
// the requested number of values.
func New(size int) []byte {
	return make(Array, (size+1)>>1)
}

// Get returns the value at the index.
func (a Array) Get(idx int) byte {
	val := a[idx>>1]
	if idx&1 == 0 {
		return val & 0xF
	}
	return val >> 4
}

// Set sets the value a the index.
func (a Array) Set(idx int, val byte) {
	i := idx >> 1
	o := a[i]
	if idx&1 == 0 {
		a[i] = (o & 0xF0) | (val & 0xF)
	} else {
		a[i] = (o & 0x0F) | ((val & 0xF) << 4)
	}
}
