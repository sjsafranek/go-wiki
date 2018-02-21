package utils

import (
	"testing"
)

func TestStrIsInt(t *testing.T) {
	if strIsInt("fail") {
		t.Error("string is not integer: \"fail\"")
	}

	if !strIsInt("100") {
		t.Error("string is integer: 100")
	}

	if strIsInt("100.9") {
		t.Error("string is float: 100.95")
	}
}

func BenchmarkStrIsInt(b *testing.B) {
	test_slice := []string{"fail", "100", "100.95"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strIsInt(test_slice[i%3])
	}
}

func TestStrIsFloat(t *testing.T) {
	if strIsFloat("fail") {
		t.Error("string is not float: \"fail\"")
	}

	if strIsFloat("100") {
		t.Error("string is integer: 100")
	}

	if !strIsFloat("100.9") {
		t.Error("string is float: 100.9")
	}
}

func BenchmarkStrIsFloat(b *testing.B) {
	test_slice := []string{"fail", "100", "100.95"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strIsFloat(test_slice[i%3])
	}
}
