package storage

import (
	"strconv"
	"testing"
)

func BenchmarkInMemory_Add(b *testing.B) {
	memStore := NewMemoryURLStore()

	for i := 0; i < b.N; i++ {
		memStore.Add("1"+strconv.Itoa(i), "http://www.google.com")
	}
}

func BenchmarkInMemory_Get(b *testing.B) {
	memStore := NewMemoryURLStore()

	for i := 0; i < b.N; i++ {
		memStore.Get("1" + strconv.Itoa(i))
	}
}
