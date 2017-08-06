// Package bloom - implements a bloom filter
package bloom

import (
	"math"

	"github.com/dchest/siphash"
)

const (
	// Siphash input
	k0 = 17697571051839533707
	k1 = 15128385881502100741
	// 64 bit word
	word = 64
	// Modulus -  64 bit word
	mod = word - 1
	// Number of right shifts necessary to divide by 64
	div = 6
)

var ()

// Filter - bloom filter
type Filter struct {
	bits    []uint64
	nbits   uint64
	nhashes uint64
	one     uint64
}

// New - creates a new bloom filter for n items with a false positive probability of p.
func New(n int, p float64) *Filter {

	if n < 0 {
		panic("bloom.New: incorrect number of elements")
	}

	// Optimal number of bits, m = -(n log(p) / pow(log(2), 2))
	// pow(log(2),2)
	const lnsq = 0.480453013918201424667102526326664971730552951594545586866864133623665382259834472199948263443926990932715597661358897481255128413358268503177555294880844290839184664798896404335252423673643658092881230886029639112807153031
	nfloat := float64(n)
	m := -math.Ceil(nfloat * math.Log(p) / lnsq)
	nbits := ((uint64(m) + mod) >> div) * word
	nhashes := uint64(-int(math.Ceil(math.Log(p) / math.Ln2)))
	return &Filter{
		bits:    make([]uint64, nbits>>div),
		nbits:   nbits,
		nhashes: nhashes,
	}

}

// Add - adds an element (string) to the bloom filter
func (f *Filter) Add(key string) { f.AddBytes([]byte(key)) }

// AddBytes - adds an element (bytes) to the bloom filter
func (f *Filter) AddBytes(key []byte) {
	a, b := siphash.Hash128(k0, k1, key)
	for h := uint64(0); h <= f.nhashes; h++ {
		i := (a + h*b) % f.nbits
		mask := uint64(1) << (i & mod)
		// if bit is not set in a specific word (i>>div)
		if f.bits[i>>div]&mask == 0 {
			f.bits[i>>div] |= mask
			f.one++
		}
	}
}

// Has - verifies if the element (string) exists in the bloom filter
func (f *Filter) Has(key string) bool { return f.HasBytes([]byte(key)) }

// HasBytes - verifies if the element (bytes) exists in the bloom filter
func (f *Filter) HasBytes(key []byte) bool {
	a, b := siphash.Hash128(k0, k1, key)
	for h := uint64(0); h <= f.nhashes; h++ {
		i := (a + h*b) % f.nbits
		mask := uint64(1) << (i & mod)
		if f.bits[i>>div]&mask == 0 {
			return false
		}
	}
	return true
}

// Info - returns basic filter information
func (f *Filter) Info() (nhashes, nbits, one uint64) {
	return f.nhashes, f.nbits, f.one
}

// Clear - resets the bloom filter
func (f *Filter) Clear() {
	for i := range f.bits {
		f.bits[i] = 0
	}
}
