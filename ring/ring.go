package ring

import (
	"github.com/dedis/student_18_lattices/bigint"
	"github.com/dedis/student_18_lattices/polynomial"
	"errors"
	"github.com/LoCCS/bliss/sampler"
	"github.com/LoCCS/bliss/poly"
)

type Ring struct {
	N uint32
	Q bigint.Int
	Poly *polynomial.Poly
}

func NewRing(n uint32, q bigint.Int) *Ring {
	r := new(Ring)
	r.N = n
	r.Q = q
	r.Poly, _ = polynomial.NewPolynomial(n, q)
	return r
}

func NewGaussPoly(n uint32, q bigint.Int, m bigint.Int) *Ring {
	r := new(Ring)
	r.N = n
	r.Q = q
	r.Poly, _ = polynomial.NewPolynomial(n, q)
	coeffs := make([]bigint.Int, n)
	var coeff bigint.Int
	negM := new(bigint.Int).Neg(&m, &q)
	for i := range coeffs {
		for {
			coeff.SetInt(int64(GaussSampling(118, 10)))
			coeff.Mod(&coeff, &q)
			if coeff.Compare(&m) == 1 {
				if coeff.Compare(negM) == -1 {
					continue
				}
			}
			coeffs[i].SetBigInt(&coeff)
			break
		}
	}
	r.Poly.SetCoefficients(coeffs)
	return r
}

func NewGaussPolyFromBLISS(n uint32, q bigint.Int) *Ring {
	seed := make([]uint8, sampler.SHA_512_DIGEST_LENGTH)
	for i := 0; i < len(seed); i++ {
		seed[i] = uint8(i % 8)
	}
	entropy, _ := sampler.NewEntropy(seed)
	mySampler, _ := sampler.New(0, entropy)
	gaussPoly := poly.GaussPoly(0, mySampler)

	r := new(Ring)
	r.N = n
	r.Q = q
	r.Poly, _ = polynomial.NewPolynomial(n, q)
	coeffs := make([]bigint.Int, n)
	_coeffs := gaussPoly.GetData()
	for i := range coeffs {
		coeffs[i].SetInt(int64(_coeffs[i]))
	}
	r.Poly.SetCoefficients(coeffs)
	return r
}

func NewUniformPoly(n uint32, q bigint.Int, m bigint.Int) *Ring {
	r := new(Ring)
	r.N = n
	r.Q = q
	r.Poly, _ = polynomial.NewPolynomial(n, q)
	coeffs := make([]bigint.Int, n)
	for i := range coeffs {
		coeffs[i].SetInt(int64(randUniform(m.Uint32())))
	}
	r.Poly.SetCoefficients(coeffs)
	return r
}

func (r *Ring) GetCoefficients() []bigint.Int{
	return r.Poly.GetCoefficients()
}

func (r *Ring) GetCoefficientsInt64() []int64{
	return r.Poly.GetCoefficientsInt64()
}

func (r *Ring) Add(r1, r2 *Ring) (*Ring, error) {
	if r.N != r1.N || !r.Q.EqualTo(&r1.Q) ||
		r.N != r2.N || !r.Q.EqualTo(&r2.Q) ||
		r1.N != r2.N || !r1.Q.EqualTo(&r2.Q) {
		return nil, errors.New("unmatched degree or module")
	}
	r.Poly.AddMod(r1.Poly, r2.Poly)
	return r, nil
}

func (r *Ring) Sub(r1, r2 *Ring) (*Ring, error) {
	if r.N != r1.N || !r.Q.EqualTo(&r1.Q) ||
		r.N != r2.N || !r.Q.EqualTo(&r2.Q) ||
		r1.N != r2.N || !r1.Q.EqualTo(&r2.Q) {
		return nil, errors.New("unmatched degree or module")
	}
	r.Poly.SubMod(r1.Poly, r2.Poly)
	return r, nil
}

func (r *Ring) Neg(r1 *Ring) (*Ring, error) {
	if r.N != r1.N || !r.Q.EqualTo(&r1.Q) {
		return nil, errors.New("unmatched degree or module")
	}
	r.Poly.Neg(r1.Poly)
	return r, nil
}

func (r *Ring) MulPoly(r1, r2 *Ring) (*Ring, error) {
	if r1.N != r2.N || !r1.Q.EqualTo(&r2.Q) {
		return nil, errors.New("unmatched degree or module")
	}
	r.Poly.MulPoly(r1.Poly, r2.Poly)
	return r, nil
}

func (r *Ring) DebugMulPoly(r1, r2 *Ring) (*Ring, error) {
	if r1.N != r2.N || !r1.Q.EqualTo(&r2.Q) {
		return nil, errors.New("unmatched degree or module")
	}
	r.Poly.DebugMulPoly(r1.Poly, r2.Poly)
	return r, nil
}

func (r *Ring) MulCoeffs(r1, r2 *Ring) (*Ring, error) {
	if r.N != r1.N || !r.Q.EqualTo(&r1.Q) ||
		r.N != r2.N || !r.Q.EqualTo(&r2.Q) ||
		r1.N != r2.N || !r1.Q.EqualTo(&r2.Q) {
		return nil, errors.New("unmatched degree or module")
	}
	r.Poly.MulCoeffs(r1.Poly, r2.Poly)
	return r, nil
}

func (r *Ring) MulScalar(r1 *Ring, scalar bigint.Int) (*Ring, error) {
	if r.N != r1.N || !r.Q.EqualTo(&r1.Q) {
		return nil, errors.New("unmatched degree or module")
	}
	r.Poly.MulScalar(r1.Poly, scalar)
	return r, nil
}

func (r *Ring) Div(r1 *Ring, scalar bigint.Int) (*Ring, error) {
	if r.N != r1.N || !r.Q.EqualTo(&r1.Q) {
		return nil, errors.New("unmatched degree or module")
	}
	if scalar.EqualTo(bigint.NewInt(int64(0))) {
		return nil, errors.New("divisor cannot be zero")
	}
	r.Poly.Div(r1.Poly, scalar)
	return r, nil
}

func (r *Ring) DivRound(r1 *Ring, scalar bigint.Int) (*Ring, error) {
	if r.N != r1.N || !r.Q.EqualTo(&r1.Q) {
		return nil, errors.New("unmatched degree or module")
	}
	if scalar.EqualTo(bigint.NewInt(int64(0))) {
		return nil, errors.New("divisor cannot be zero")
	}
	r.Poly.DivRound(r1.Poly, scalar)
	return r, nil
}

func (r *Ring) Mod(r1 *Ring, m bigint.Int) (*Ring, error) {
	if r.N != r1.N || !r.Q.EqualTo(&r1.Q) {
		return nil, errors.New("unmatched degree or module")
	}
	r.Poly.Mod(r1.Poly, m)
	return r, nil
}

