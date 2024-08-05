package random

import "testing"

func BenchmarkNewRandomURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewRandomURL()
	}
}
