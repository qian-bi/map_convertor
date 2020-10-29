package main

import (
	"testing"
)

func BenchmarkConvertor(b *testing.B) {
	for n := 0; n < b.N; n++ {
		main()
	}
}
