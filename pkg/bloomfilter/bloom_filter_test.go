package bloomfilter

import (
	"testing"
)

type PairI64I64 struct {
	first, second int64
}

func TestNew(t *testing.T) {
	// input - expected
	caps := []*PairI64I64{
		{10000, 190000},
		{20000, 380000},
		{5400, 102600},
	}

	for _, p := range caps {
		cap, expectedSize := p.first, p.second
		bf := New(cap)

		// verify filter size
		if bf.Size() != expectedSize {
			t.Fatalf(`bloomfilter.New(): expected bf.Size() is %d, but actual is %d`,
				expectedSize, bf.Size())
		}
	}
}

func TestAdd(t *testing.T) {
	bf := New(int64(10000))
	value := "This is Trieu"

	result, err := bf.Add(value)
	if !result || err != nil {
		t.Fatalf(`bloomfilter.Add(): expected result is %v,%v, but actual is %v,%v`,
			true, nil, result, err)
	}
}

func TestMightContain(t *testing.T) {
	bf := New(int64(10000))
	hitValues := []string{
		"Trieu", "Million", "Ti", "Billion",
	}
	missValues := []string{
		"TRiEU", "mILLioN", "ti", "stranger",
	}

	// add to bloom filter
	for _, val := range hitValues {
		bf.Add(val)
	}

	// verify positive
	for _, val := range hitValues {
		result := bf.MightContain(val)
		if !result {
			t.Fatalf(`bloomfilter.MightContain(%v): expected result is %v, but actual is %v`,
				val, true, result)
		}
	}

	// verify negative
	for _, val := range missValues {
		result := bf.MightContain(val)
		if result {
			t.Fatalf(`bloomfilter.MightContain(%v): expected result is %v, but actual is %v`,
				val, false, result)
		}
	}
}
