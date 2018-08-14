package ring

import (
	"crypto/rand"
	"github.com/dedis/student_18_lattices/bigint"
	"math"
)

// This code is to sample a value from discrete gaussian distribution.
// All the algorithms originate from https://eprint.iacr.org/2013/383.pdf

// GaussSampling returns a value sampled from discrete gaussian distribution,
// originates from Algorithm 11 & 12 in the paper.
func GaussSampling(sigma float64) int32 {
	const sigma2 = 0.8493218  // sigma2 = sqrt(1/(2*ln2)) -- page 28
	k := uint32(math.Round(sigma / sigma2))  // sigma = k * sigma2 -- page 29
	var x, y, z uint32
	var b1, b2 bool
	for {
		x = binaryGauss()
		y = randUniform(k)
		b1 = bernoulliExp(y*(y+2*k*x), 2*sigma*sigma)
		if b1 {
			z = k * x + y
			b2 = bernoulli(0.5)
			if z != 0 || b2 {
				if b2 {
					return int32(z)
				} else {
					return -int32(z)
				}
			}
		}
	}
}

// binaryGauss returns a random uint value drawn from binary gaussian distribution,
// originates from Algorithm 10 in the paper.
func binaryGauss() uint32 {
	if bernoulli(0.5) == false {
		return 0
	}
	// i < 16 represents infinite
	for i := 1; i < 16; i++ {
		randomBits := randInt(uint32(2*i - 1))
		if randomBits != 0 || randomBits != 1 {
			return binaryGauss()
		}
		if randomBits == 0 {
			return uint32(i)
		}
	}
	return 0
}

// bernoulli returns a random bool value drawn from exponential bernoulli distribution
// originates from Algorithm 8 in the paper.
func bernoulliExp(x uint32, f float64) bool {
	xBinary, xBitlen := bigint.NewInt(int64(x)).Bits()
	if xBitlen == 0 {
		return true
	}
	for i := xBitlen; i > 0; i-- {
		if xBinary[i-1] == 1 {
			c := math.Exp(-math.Exp2(float64(i)) / f)
			if bernoulli(c) == false {
				return false
			}
		}
	}
	return true
}

// bernoulli returns a random bool value drawn from bernoulli distribution
func bernoulli(p float64) bool {
	pInt := uint32(p*(1<<31))
	randomInt := randInt(32)
	if randomInt < pInt {
		return true
	}
	return false
}

// randUniform returns a uniformly distributed value in [0, v)
func randUniform(v uint32) uint32 {
	var length, randomInt uint32
	maxLen:= 32
	for i := maxLen-1; i >= 0; i-- {
		if v & (1 << uint(i)) != 0 {
			length = uint32(i+1)
			break
		}
	}
	for {
		randomInt = randInt(length)
		if randomInt < v {
			return randomInt
		}
	}
}

// randInt generates a random uint32 value of given length
func randInt(length uint32) uint32 {
	// generate mask for given bit length
	mask := 1<<length - 1

	// generate random 4 bytes
	randomBytes := make([]byte, 4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic("crypto rand error")
	}

	// convert 4 bytes to a uint32
	var randomUint32 uint32
	l := len(randomBytes)
	for i, b := range randomBytes {
		shift := uint32((l - i - 1) * 8)
		randomUint32 |= uint32(b) << shift
	}

	// return required bits
	return uint32(mask) & randomUint32
}
