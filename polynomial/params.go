package polynomial

import (
	"github.com/dedis/student_18_lattices/bigint"
)

type NttParams struct {
	n, nReverse uint32
	q bigint.Int
	PsiReverse []bigint.Int
	PsiReverseMontgomery []bigint.Int
	PsiInvReverse []bigint.Int
	PsiInvReverseMontgomery []bigint.Int
	bitLen uint32  // param of montgomery reduction
	qInv bigint.Int // (2^bitLen * (inverse(2^bitLen mod q)) - 1) / q, param of montgomery reduction
}

// generateNTTParameters generates the parameters for NTT and inverse NTT transformations.
func generateNTTParameters(N uint32, Q bigint.Int) (*NttParams, error) {
	newNttParams := new(NttParams)
	// set n
	newNttParams.n = N
	// set nReverse
	var temp bigint.Int
	temp.Inv(bigint.NewInt(int64(N)), &Q)
	newNttParams.nReverse = temp.Uint32()
	// set q
	newNttParams.q.SetBigInt(&Q)

	// In the following, we calculate PsiReverse, PsiReverseMontgomery, PsiInvReverse, PsiInvReverseMontgomery.
	// 1. First, set primitive root g = 2, and fi = q-1
	g := primitiveRoot(&Q)
	fi := new(bigint.Int)
	fi.Sub(&Q, bigint.NewInt(1))  // fi = q - 1

	// 2. Second, calculate 2N-th root of unity and its inverse, i.e. psi = g^(fi/2N) mod q, and psi^-1 mod q.
	_2n := bigint.NewInt(2)
	_2n.Mul(_2n, bigint.NewInt(int64(N)))
	power := new(bigint.Int)
	power.SetBigInt(fi)
	power.Div(fi, _2n)
	psi := new(bigint.Int)
	psi.Exp(g, power, &Q)

	powerInv := new(bigint.Int)
	powerInv.Sub(fi, power)
	psiInv := new(bigint.Int)
	psiInv.Exp(g, powerInv, &Q)

	// 3. Third, calculate powers of psi and psiInv in bit-reversed order
	newNttParams.PsiReverse = make([]bigint.Int, N)
	newNttParams.PsiInvReverse = make([]bigint.Int, N)

	// computing the bit length of N
	var bitLenofN uint32
	for i := 32-1; i >= 0; i-- {
		if N & (1 << uint(i)) != 0 {
			bitLenofN = uint32(i)
			break
		}
	}

	// computing PsiReverse and PsiInvReverse
	var indexReverse uint32
	var _i bigint.Int
	for i := uint32(0); i < N; i++ {
		_i.SetInt(int64(i))
		indexReverse = bitReverse(i, bitLenofN)
		newNttParams.PsiReverse[indexReverse].Exp(psi, &_i, &Q)
		newNttParams.PsiInvReverse[indexReverse].Exp(psiInv, &_i, &Q)
	}

	// computing PsiReverseMontgomery and PsiInvReverseMontgomery
	newNttParams.PsiReverseMontgomery = make([]bigint.Int, N)
	newNttParams.PsiInvReverseMontgomery = make([]bigint.Int, N)

	qBitLen := Q.Value.BitLen() + 5
	r := new(bigint.Int).Exp(bigint.NewInt(2), bigint.NewInt(int64(qBitLen)), &Q)
	for i := uint32(0); i < N; i++ {
		newNttParams.PsiReverseMontgomery[i].Mul(r, &newNttParams.PsiReverse[i])
		newNttParams.PsiReverseMontgomery[i].Mod(&newNttParams.PsiReverseMontgomery[i], &Q)
		newNttParams.PsiInvReverseMontgomery[i].Mul(r, &newNttParams.PsiInvReverse[i])
		newNttParams.PsiInvReverseMontgomery[i].Mod(&newNttParams.PsiInvReverseMontgomery[i], &Q)
	}

	// set bitLen
	newNttParams.bitLen = uint32(qBitLen)
	// set qInv
	r = r.Lsh(bigint.NewInt(1), uint32(qBitLen))
	rInv := new(bigint.Int).Inv(r, &Q)
	newNttParams.qInv.Mul(r, rInv)
	newNttParams.qInv.Sub(&newNttParams.qInv, bigint.NewInt(1))
	newNttParams.qInv.Div(&newNttParams.qInv, &Q)

	return newNttParams, nil
}

// bitReverse calculates the bit-reverse index.
// for example, given index=6 (110) and its bit-length bitLen=3, the indexReverse would be 3 (011)
func bitReverse(index, bitLen uint32)  uint32{
	indexReverse := uint32(0)
	for i := uint32(0); i < bitLen; i++ {
		if (index >> i) & 1 != 0 {
			indexReverse |= 1 << (bitLen - 1 - i)
		}
	}
	return indexReverse
}

// polynomialPollardsRho calculates x1^2 + c mod x2, and is used in factorizationPollardsRho
func polynomialPollardsRho(x1, x2, c *bigint.Int) *bigint.Int{
	two := bigint.NewInt(2)
	z := new(bigint.Int).Exp(x1, two, x2) // x1^2 mod x2
	z.Add(z, c)                   // (x1^2 mod x2) + 1
	z.Mod(z, x2)                     // (x1^2 + 1) mod x2
	return z
}

// factorizationPollardsRho realizes Pollard's Rho algorithm for fast prime factorization,
// but this function only returns one factor a time
func factorizationPollardsRho (m *bigint.Int) *bigint.Int {
	var x, y, d, c *bigint.Int
	zero := bigint.NewInt(0)
	one := bigint.NewInt(1)
	ten := bigint.NewInt(10)

	// c is to change the polynomial used in Pollard's Rho algorithm,
	// Every time the algorithm fails to get a factor, increasing c to retry,
	// because Pollard's Rho algorithm sometimes will miss some small prime factors.
	for c = bigint.NewInt(1); !c.EqualTo(ten); c.Add(c, one){
		x, y, d = bigint.NewInt(2), bigint.NewInt(2), bigint.NewInt(1)
		for !d.EqualTo(zero) {
			x = polynomialPollardsRho(x, m, c)
			y = polynomialPollardsRho(polynomialPollardsRho(y, m, c), m, c)
			sub := new(bigint.Int).Sub(x, y)
			d.Value.GCD(nil, nil, sub.Value.Abs(&sub.Value), &m.Value)
			if d.Compare(one) == 1.0 {
				return d
			}
		}
	}
	return d
}

// getFactors returns all the prime factors of m
func getFactors(n *bigint.Int) []bigint.Int {
	var factor *bigint.Int
	var factors []bigint.Int
	var m, tmp bigint.Int
	m.SetBigInt(n)
	zero := bigint.NewInt(0)
	one := bigint.NewInt(1)

	// first, append small prime factors
	for i := range smallPrimes {
		smallPrime := bigint.NewInt(smallPrimes[i])
		addFactor := false
		for tmp.Mod(&m, smallPrime).EqualTo(zero) {
			m.Div(&m, smallPrime)
			addFactor = true
		}
		if addFactor {
			factors = append(factors, *smallPrime)
		}
	}

	if m.EqualTo(one) {
		return factors
	}

	// second, find other prime factors
	for {
		factor = factorizationPollardsRho(&m)
		if factor.EqualTo(zero) {
			factors = append(factors, m)
			break
		}
		m.Div(&m, factor)
		if len(factors) > 0 && factor.EqualTo(&factors[len(factors)-1]) {
			continue
		}
		factors = append(factors, *factor)
	}
	return factors
}

// primitiveRoot calculates one primitive root of prime q
func primitiveRoot(q *bigint.Int) *bigint.Int {
	tmp := new(bigint.Int)
	notFoundPrimitiveRoot := true
	qMinusOne := new(bigint.Int).Sub(q, bigint.NewInt(1))
	factors := getFactors(qMinusOne)
	g := bigint.NewInt(2)
	one := bigint.NewInt(1)
	for notFoundPrimitiveRoot {
		g.Add(g, one)
		for _, factor := range factors {
			tmp.Div(qMinusOne, &factor)
			// once exist g^(q-1)/factor = 1 mod q, g is not a primitive root
			if tmp.Exp(g, tmp, q).EqualTo(one) {
				notFoundPrimitiveRoot = true
				break
			}
			notFoundPrimitiveRoot = false
		}
	}
	return g
}

var smallPrimes = []int64 {
	2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 73, 79, 83, 89, 97,
	101, 103, 107, 109, 113, 127, 131, 139, 149, 151, 163, 167, 173, 179, 181, 191, 193, 197, 199,
}
