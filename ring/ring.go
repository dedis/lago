package ring

import (
	"github.com/dedis/lago/bigint"
	"github.com/dedis/lago/polynomial"
)

type Ring struct {
	N uint32
	Q bigint.Int
	Poly *polynomial.Poly
}

// NewRing creates a new polynomial ring with given parameters.
func NewRing(n uint32, q bigint.Int, nttParams *polynomial.NttParams) (*Ring, error) {
	r := new(Ring)
	err := *new(error)
	r.N = n
	r.Q = q
	r.Poly, err = polynomial.NewPolynomial(n, q, nttParams)
	return r, err
}

// CopyRing copies a polynomial ring.
func CopyRing(r1 *Ring) (*Ring, error) {
	r := new(Ring)
	err := *new(error)
	r.N = r1.N
	r.Q = r1.Q
	r.Poly, err = polynomial.NewPolynomial(r1.N, r1.Q, r1.Poly.GetNTTParams())
	return r, err
}

// NewGaussPoly creates a new polynomial ring,
// the parameters of which obey discrete gaussian distribution with derivation sigma
func NewGaussPoly(n uint32, q bigint.Int, nttParams *polynomial.NttParams, sigma float64) (*Ring, error) {
	r := new(Ring)
	err := *new(error)
	r.N = n
	r.Q = q
	r.Poly, err = polynomial.NewPolynomial(n, q, nttParams)
	coeffs := make([]bigint.Int, n)
	var coeff bigint.Int

	boundMultiplier := float64(6)  // this parameter is the same as the SEAL library
	positiveBound := bigint.NewInt(int64(boundMultiplier * sigma))  // the suggested sigma from SEAL is 3.19
	negativeBound := bigint.NewInt(-int64(boundMultiplier * sigma))
	for i := range coeffs {
		for {
			coeff.SetInt(int64(GaussSampling(sigma)))
			coeff.Mod(&coeff, &q)
			if coeff.Compare(positiveBound) == 1 {
				if coeff.Compare(negativeBound) == -1 {
					continue
				}
			}
			coeffs[i].SetBigInt(&coeff)
			break
		}
	}

	r.Poly.SetCoefficients(coeffs)
	return r, err
}

// NewUniformPoly creates a new polynomial ring,
// the parameters of which obey uniform distribution [0, v)
func NewUniformPoly(n uint32, q bigint.Int, nttParams *polynomial.NttParams, v bigint.Int) (*Ring, error) {
	r := new(Ring)
	err := *new(error)
	r.N = n
	r.Q = q
	r.Poly, err = polynomial.NewPolynomial(n, q, nttParams)
	coeffs := make([]bigint.Int, n)
	for i := range coeffs {
		coeffs[i].SetInt(int64(randUniform(v.Uint32())))
	}

	r.Poly.SetCoefficients(coeffs)
	return r, err
}

func (r *Ring) GetCoefficients() []bigint.Int{
	return r.Poly.GetCoefficients()
}

func (r *Ring) GetCoefficientsInt64() []int64{
	return r.Poly.GetCoefficientsInt64()
}

func (r *Ring) Add(r1, r2 *Ring) (*Ring, error) {
	_, err := r.Poly.AddMod(r1.Poly, r2.Poly)
	return r, err
}

func (r *Ring) Sub(r1, r2 *Ring) (*Ring, error) {
	_, err := r.Poly.SubMod(r1.Poly, r2.Poly)
	return r, err
}

func (r *Ring) Neg(r1 *Ring) (*Ring, error) {
	_, err := r.Poly.Neg(r1.Poly)
	return r, err
}

func (r *Ring) MulPoly(r1, r2 *Ring) (*Ring, error) {
	_, err := r.Poly.MulPoly(r1.Poly, r2.Poly)
	return r, err
}

func (r *Ring) MulCoeffs(r1, r2 *Ring) (*Ring, error) {
	_, err := r.Poly.MulCoeffs(r1.Poly, r2.Poly)
	return r, err
}

func (r *Ring) MulScalar(r1 *Ring, scalar bigint.Int) (*Ring, error) {
	_, err := r.Poly.MulScalar(r1.Poly, scalar)
	return r, err
}

func (r *Ring) Div(r1 *Ring, scalar bigint.Int) (*Ring, error) {
	_, err := r.Poly.Div(r1.Poly, scalar)
	return r, err
}

func (r *Ring) DivRound(r1 *Ring, scalar bigint.Int) (*Ring, error) {
	_, err := r.Poly.DivRound(r1.Poly, scalar)
	return r, err
}

func (r *Ring) Mod(r1 *Ring, m bigint.Int) (*Ring, error) {
	_, err := r.Poly.Mod(r1.Poly, m)
	return r, err
}

func (r *Ring) And(r1 *Ring, m bigint.Int) (*Ring, error) {
	_, err := r.Poly.And(r1.Poly, m)
	return r, err
}

func (r *Ring) Lsh(r1 *Ring, m uint32) (*Ring, error) {
	_, err := r.Poly.Lsh(r1.Poly, m)
	return r, err
}

func (r *Ring) Rsh(r1 *Ring, m uint32) (*Ring, error) {
	_, err := r.Poly.Rsh(r1.Poly, m)
	return r, err
}
