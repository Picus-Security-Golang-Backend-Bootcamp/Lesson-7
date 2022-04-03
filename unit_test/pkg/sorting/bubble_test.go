package sorting

import (
	"testing"
)

func TestBubbleSort(t *testing.T) {
	// Init
	elements := []int{3, 4, 7, 10, 9, 1, 6, 5, 2, 8}

	// Execute
	BubbleSort(elements)

	// Validate
	want := 10
	got := elements[0]

	if got != want {
		t.Error("First element should be 10")
	}

	if elements[len(elements)-1] != 1 {
		t.Error("First element should be 1")
	}
}
func BenchmarkBubbleSort(b *testing.B) {
	b.ReportAllocs()
	// Execute
	for i := 0; i < b.N; i++ {
		// Init
		elements := []int{3, 4, 7, 10, 9, 1, 6, 5, 2, 8}
		BubbleSort(elements)

	}
}
