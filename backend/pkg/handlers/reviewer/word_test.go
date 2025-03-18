package reviewer_test

import (
	"math"
	"testing"

	"math/rand"
)


func log_builtin(n int) int {
	return int(math.Ceil(math.Log2(float64(n))))
}

func log_shift(n int) int {
	n--
	if n <= 0 {
		return 0
	}
	count := 1
	for n > 1 {
		n >>= 1
		count++
	}
	return count
}
func BenchmarkLogBuiltin(b *testing.B) {
	input := generateInput(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log_builtin(input[i])
	}
}

func BenchmarkLogShift(b *testing.B) {
	input := generateInput(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log_shift(input[i])
	}
}



func TestLog(t *testing.T) {
	input := generateInput(100000)
	for i := 0; i < len(input); i++ {
		builtin := log_builtin(input[i])
		shift := log_shift(input[i])
		if builtin != shift {
			t.Errorf("Mismatch at index %d: builtin %d, shift %d", i, builtin, shift)
		}
	}
}

func TestLogShift2(t *testing.T) {
	t.Run("log 1", func(t *testing.T) {
		if log_shift(1) != 0 {
			t.Errorf("Expected 0, got %d", log_shift(1))
		}
	})
	t.Run("log 2", func(t *testing.T) {
		if log_shift(2) != 1 {
			t.Errorf("Expected 1, got %d", log_shift(2))
		}
	})
	t.Run("log 3", func(t *testing.T) {
		if log_shift(3) != 2 {
			t.Errorf("Expected 2, got %d", log_shift(3))
		}
	})
	t.Run("log 4", func(t *testing.T) {
		if log_shift(4) != 2 {
			t.Errorf("Expected 2, got %d", log_shift(4))
		}
	})
}

func generateInput(n int) []int {
	r := rand.New(rand.NewSource(123)) // Seed the random number generator
	input := make([]int, n)
	for i := 0; i < n; i++ {
		input[i] = r.Intn(1000000) // Generate random integers
	}
	return input
}