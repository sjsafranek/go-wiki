package utils

import (
	"testing"
)

func TestNewJobId(t *testing.T) {
	job_id := newJobId(10)
	if 10 != len(job_id) {
		t.Errorf("job_id is too long: %v\n", job_id)
	}
}

func BenchmarkNewJobId(b *testing.B) {
	test_slice := []int{6, 7, 8, 9, 10}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newJobId(test_slice[i%5])
	}
}
