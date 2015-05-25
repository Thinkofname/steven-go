package bit

import "testing"

func TestSet(t *testing.T) {
	s := NewSet(200)
	for i := 0; i < 200; i++ {
		if i%3 == 0 {
			s.Set(i, true)
		}
	}
	for i := 0; i < 200; i++ {
		if s.Get(i) != (i%3 == 0) {
			t.Fail()
		}
	}
}
