package lazy_fd

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

const LARGE_BUFFER_SIZE = 1024 * 50 // 50k
const SMALL_BUFFER_SIZE = 5         // 5 bytes

func ReadSum(f io.Reader) int64 {
	reader := json.NewDecoder(f)

	var sum int64 = 0
	for reader.More() {
		var value int64
		err := reader.Decode(&value)
		if err != nil {
			panic(err)
		}
		sum += value
	}

	return sum
}

func Test_Simple(t *testing.T) {
	filename := "test_files/golden_file_1.txt"
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	correct_sum := ReadSum(f)

	lf := NewLazyFileReaderSimple(filename)
	test_sum := ReadSum(lf)

	if test_sum != correct_sum {
		t.Fail()
	}
}

func Test_BufferLarge(t *testing.T) {
	filename := "test_files/golden_file_1.txt"
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	correct_sum := ReadSum(f)

	lf := NewLazyFileReaderBuffer(filename, LARGE_BUFFER_SIZE)
	test_sum := ReadSum(lf)

	if test_sum != correct_sum {
		t.Fail()
	}
}

func Test_BufferSmall(t *testing.T) {
	filename := "test_files/golden_file_1.txt"
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	correct_sum := ReadSum(f)

	lf := NewLazyFileReaderBuffer(filename, SMALL_BUFFER_SIZE)
	test_sum := ReadSum(lf)

	if test_sum != correct_sum {
		t.Fail()
	}
}

func Benchmark_Baseline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		filename := "test_files/golden_file_1.txt"
		f, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		ReadSum(f)
	}
}

func Benchmark_Simple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		filename := "test_files/golden_file_1.txt"
		lf := NewLazyFileReaderSimple(filename)
		ReadSum(lf)
	}
}

func Benchmark_Buffer_Large(b *testing.B) {
	for i := 0; i < b.N; i++ {
		filename := "test_files/golden_file_1.txt"
		lf := NewLazyFileReaderBuffer(filename, LARGE_BUFFER_SIZE)
		ReadSum(lf)
	}
}

func Benchmark_Buffer_Small(b *testing.B) {
	for i := 0; i < b.N; i++ {
		filename := "test_files/golden_file_1.txt"
		lf := NewLazyFileReaderBuffer(filename, SMALL_BUFFER_SIZE)
		ReadSum(lf)
	}
}
