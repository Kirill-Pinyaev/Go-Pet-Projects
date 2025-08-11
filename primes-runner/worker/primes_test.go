package worker

import (
	"fmt"
	"runtime"
	"testing"
)

// TestIsPrime проверяет корректность детектора простых чисел в табличном формате.
func TestIsPrime(t *testing.T) {
	cases := []struct {
		n    int
		want bool
	}{
		{0, false},
		{1, false},
		{2, true},
		{3, true},
		{4, false},
		{5, true},
		{6, false},
		{7, true},
		{17, true},
		{18, false},
		{19, true},
		{97, true},
		{99, false},
		{7919, true}, // большое простое
		{7920, false},
	}

	for _, c := range cases {
		if got := isPrime(c.n); got != c.want {
			t.Errorf("isPrime(%d) = %v, want %v", c.n, got, c.want)
		}
	}
}

// benchPrimeRunner замеряет время и аллокации для одного набора параметров.
func benchPrimeRunner(b *testing.B, max int, workers int, gomax int) {
	prev := runtime.GOMAXPROCS(gomax)
	defer runtime.GOMAXPROCS(prev)

	b.ReportAllocs()
	b.ResetTimer()

	var primes int
	for i := 0; i < b.N; i++ {
		primes = Dispatcher(max, workers)
	}

	// пользовательская метрика: сколько простых найдено за прогон
	b.ReportMetric(float64(primes), "primes")
}

// BenchmarkPrimeRunner перебирает workers ∈ {1,4,8} и GOMAXPROCS = 1..NumCPU (степени двойки).
func BenchmarkPrimeRunner(b *testing.B) {
	max := 1_000_000

	for _, workers := range []int{1, 4, 8} {
		for gomax := 1; gomax <= runtime.NumCPU(); gomax *= 2 {
			name := fmt.Sprintf("max=%d/W%d/G%d", max, workers, gomax)
			b.Run(name, func(b *testing.B) {
				benchPrimeRunner(b, max, workers, gomax)
			})
		}
	}
}
