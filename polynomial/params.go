package polynomial

import (
	"github.com/dedis/student_18_lattices/bigint"
)

type params struct {
	n uint32
	q bigint.Int
	PsiReverse []bigint.Int
	PsiInvReverse []bigint.Int
}

// generateNTTParameters generates the parameters for NTT and inverse NTT transformations.
func GenerateNTTParameters(N uint32, Q bigint.Int) (*params, error) {
	// TODO Check if Q is a prime number
	// TODO Check if Q/(2N) is an integer
	newNttParams := new(params)
	newNttParams.n = N
	newNttParams.q.SetBigInt(&Q)

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
	// computing psiReverse and psiInvReverse
	var indexReverse uint32
	var _i bigint.Int
	for i := uint32(0); i < N; i++ {
		_i.SetInt(int64(i))
		indexReverse = bitReverse(i, bitLenofN)
		newNttParams.PsiReverse[indexReverse].Exp(psi, &_i, &Q)
		newNttParams.PsiInvReverse[indexReverse].Exp(psiInv, &_i, &Q)
	}

	return newNttParams, nil
}

func (p *params) GetPsiReverseUint32() []uint32 {
	psi := make([]uint32, p.n)
	for i := uint32(0); i < p.n; i++ {
		psi[i] = uint32(p.PsiReverse[i].Int64())
	}
	return psi
}

func (p *params) GetPsiInvReverseUint32() []uint32 {
	psiInv := make([]uint32, p.n)
	for i := uint32(0); i < p.n; i++ {
		psiInv[i] = uint32(p.PsiInvReverse[i].Int64())
	}
	return psiInv
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
	one := bigint.NewInt(1)
	ten := bigint.NewInt(10)

	// c is to change the polynomial used in Pollard's Rho algorithm,
	// Every time the algorithm fails to get a factor, increasing c to retry,
	// because Pollard's Rho algorithm sometimes will miss some small prime factors.
	for c = bigint.NewInt(1); !c.EqualTo(ten); c.Add(c, one){
		x, y, d = bigint.NewInt(2), bigint.NewInt(2), bigint.NewInt(1)
		for d.EqualTo(one) {
			x = polynomialPollardsRho(x, m, c)
			y = polynomialPollardsRho(polynomialPollardsRho(y, m, c), m, c)
			sub := new(bigint.Int).Sub(x, y)
			d.Value.GCD(nil, nil, sub.Value.Abs(&sub.Value), &m.Value)
		}
	}

	if d.EqualTo(m) {
		return one
	}
	return d
}

// getFactors returns all the prime factors of m
func getFactors(n *bigint.Int) []bigint.Int {
	var factor *bigint.Int
	var factors []bigint.Int
	var subFactors []bigint.Int
	var m, tmp bigint.Int
	m.SetBigInt(n)
	zero := bigint.NewInt(0)
	two := bigint.NewInt(2)

	// first, turn m into odd, and add 2 as a factor
	for tmp.Mod(&m, two).EqualTo(zero) {
		m.Div(&m, two)
	}
	if !m.EqualTo(n) {
		factors = append(factors, *two)
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
		factors = append(factors, subFactors...)
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
