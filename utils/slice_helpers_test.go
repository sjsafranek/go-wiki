package utils

import (
	"testing"
)

func TestStringInSlice(t *testing.T) {
	test_slice := []string{"sam", "i", "am"}
	if StringInSlice("fail", test_slice) {
		t.Errorf("string is not in slice: %v %v\n", "fail", test_slice)
	}

	if !StringInSlice("sam", test_slice) {
		t.Errorf("string is in slice: %v %v\n", "sam", test_slice)
	}
}

func BenchmarkStringInSlice(b *testing.B) {
	test_slice := []string{"sam", "i", "am"}
	test_values := []string{"sam", "i", "fail"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StringInSlice(test_values[i%3], test_slice)
	}
}

func TestIntInSlice(t *testing.T) {
	test_slice := []int{1, 2, 3}
	if intInSlice(4, test_slice) {
		t.Errorf("string is not in slice: %v %v\n", 4, test_slice)
	}

	if !intInSlice(1, test_slice) {
		t.Errorf("string is in slice: %v %v\n", 1, test_slice)
	}
}

func BenchmarkIntInSlice(b *testing.B) {
	test_slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	test_values := []int{1, 10}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		intInSlice(test_values[i%2], test_slice)
	}
}
