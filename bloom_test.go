package bloom

import (
	"fmt"
	"testing"

	"github.com/dchest/siphash"
	"github.com/stretchr/testify/assert"
)

const (
	// Number of elements
	n = 65000
	// False positive probability
	p = 0.1
)

func TestNbits(t *testing.T) {
	//Optimal number of bits, m = -(n log(p) / pow(log(2), 2))
	// p = 0.1 ; n = 65000
	// m ~= 311514 > 311552 (multiple of 64 - 4868 64bit-words)
	// Given a filter with n = 65000 and p = 0.1
	filter := New(n, p)
	// Then, the number of bits must be equal to 311552
	assert.Equal(t, uint64(311552), filter.nbits)
}

func TestNHashes(t *testing.T) {
	// uint64(-int(math.Ceil(math.Log(p) / math.Ln2)))
	// Given a filter with n = 65000 and p = 0.1
	filter := New(n, p)
	// Then, the number of bits must be equal to 3 (math.Ceil(-3.321928094887362) == 3)
	assert.Equal(t, filter.nhashes, uint64(3))
}

func TestAdd(t *testing.T) {
	// Given an empty filter
	value := "random"
	filter := New(n, p)
	// If we add a value
	filter.Add(value)
	// Then, the value must exist in the filter
	a, b := siphash.Hash128(k0, k1, []byte(value))
	for h := uint64(0); h <= filter.nhashes; h++ {
		i := (a + h*b) % filter.nbits
		mask := uint64(1) << (i & mod)
		assert.NotZero(t, filter.bits[i>>div]&mask)
	}
}

func TestHas(t *testing.T) {
	// Given a filter that has a value x
	value := "random"
	filter := New(n, p)
	a, b := siphash.Hash128(k0, k1, []byte(value))
	for h := uint64(0); h <= filter.nhashes; h++ {
		i := (a + h*b) % filter.nbits
		mask := uint64(1) << (i & mod)
		if filter.bits[i>>div]&mask == 0 {
			filter.bits[i>>div] |= mask
		}
	}
	// If the has operation is executed to verify x presence
	result := filter.Has(value)
	fmt.Println(result)
	// Then, the result must be true
	assert.True(t, result)
}

func TestClear(t *testing.T) {
	// Given a filter with values
	filter := New(n, p)
	filter.bits = []uint64{300, 55}
	// If the clear option is executed
	filter.Clear()
	// Then, all the words must be equal to 0
	for i := range filter.bits {
		assert.Zero(t, filter.bits[i])
	}
}
