package polynomial

import (
	"errors"
	"github.com/dedis/student_18_lattices/bigint"
)

type Poly struct {
	coeffs []bigint.Int
	n      uint32
	q      bigint.Int
	nttParams *NttParams
}

// NewPolynomial creates a new polynomial with a given degree N and module Q
func NewPolynomial(N uint32, Q bigint.Int, NttParams *NttParams) (*Poly, error) {
	p := &Poly{make([]bigint.Int, N), N, Q, NttParams}
	return p, nil
}

// GenerateNTTParams generates the ntt params of polynomial p
func GenerateNTTParams(N uint32, Q bigint.Int) *NttParams {
	//TODO : Primality Test
	if (N & (N - 1)) != 0 { // if N is power of 2
		panic("polynomial degree N has to be power of 2")
	}
	if !new(bigint.Int).Mod( // if Q mod 2N = 1
		&Q, new(bigint.Int).Mul(bigint.NewInt(2), bigint.NewInt(int64(N)))).EqualTo(bigint.NewInt(1)) {
			panic("polynomial modulus Q mod 2*N should be 1")
	}
	nttparams, err := generateNTTParameters(N, Q)
	if err != nil {
		panic(err)
	}
	return nttparams
}

// SetNTTParams sets the nttParams of polynomial p to the given nttparams
func (p *Poly) SetNTTParams(nttparams *NttParams) error {
	if nttparams == nil {
		return errors.New("invalid ntt params")
	}
	p.n = nttparams.n
	p.q = nttparams.q
	p.nttParams = nttparams
	return nil
}

// GetNTTParams returns the nttParams of polynomial p
func (p *Poly) GetNTTParams() *NttParams {
	return p.nttParams
}

// SetCoefficients sets the coefficient of target polynomial p to coeffs
func (p *Poly) SetCoefficients(coeffs []bigint.Int) error {
	if uint32(len(coeffs)) != p.n {
		return errors.New("provided coeffs has different length with target polynomial")
	}
	for i, c := range coeffs {
		p.coeffs[i].SetBigInt(&c)
	}
	return nil
}

// GetCoefficients returns the coefficients of target polynomial p
func (p *Poly) GetCoefficients() []bigint.Int {
	return p.coeffs
}

// GetCoefficientsInt64 returns the low 64 bits of coefficients of target polynomial p as int64
func (p *Poly) GetCoefficientsInt64() []int64 {
	coeffs := make([]int64, p.n)
	for i := range p.coeffs {
		coeffs[i] = p.coeffs[i].Int64()
	}
	return coeffs
}

// AddMod adds then mod the coefficients of p1 and p2
func (p *Poly) AddMod(p1, p2 *Poly) (*Poly, error) {
	if p.n != p1.n || !p.q.EqualTo(&p1.q) ||
		p.n != p2.n || !p.q.EqualTo(&p2.q) ||
		p1.n != p2.n || !p1.q.EqualTo(&p2.q) {
		return nil, errors.New("unmatched degree or module")
	}
	for i := range p.coeffs {
		p.coeffs[i].Add(&p1.coeffs[i], &p2.coeffs[i])
		p.coeffs[i].Mod(&p.coeffs[i], &p.q)
	}
	return p, nil
}

// SubMod subtracts then mod the coefficients of p1 and p2
func (p *Poly) SubMod(p1, p2 *Poly) (*Poly, error) {
	if p.n != p1.n || !p.q.EqualTo(&p1.q) ||
		p.n != p2.n || !p.q.EqualTo(&p2.q) ||
		p1.n != p2.n || !p1.q.EqualTo(&p2.q) {
		return nil, errors.New("unmatched degree or module")
	}
	for i := range p.coeffs {
		p.coeffs[i].Sub(&p1.coeffs[i], &p2.coeffs[i])
		p.coeffs[i].Mod(&p.coeffs[i], &p.q)
	}
	return p, nil
}

// Neg sets the coefficients of polynomial p to the negative of p1'coefficients
func (p *Poly) Neg(p1 *Poly) (*Poly, error) {
	if p.n != p1.n || !p.q.EqualTo(&p1.q) {
		return nil, errors.New("unmatched degree or module")
	}
	for i := range p.coeffs {
		p.coeffs[i].Neg(&p1.coeffs[i], &p.q)
	}
	return p, nil
}

// InnerProduct multiplies polynomials p1 and p2 in coefficient-wise
func (p *Poly) MulCoeffs(p1, p2 *Poly) (*Poly, error) {
	if p.n != p1.n || !p.q.EqualTo(&p1.q) ||
		p.n != p2.n || !p.q.EqualTo(&p2.q) ||
		p1.n != p2.n || !p1.q.EqualTo(&p2.q) {
		return nil, errors.New("unmatched degree or module")
	}
	for i := range p.coeffs {
		p.coeffs[i].Mul(&p1.coeffs[i], &p2.coeffs[i])
	}
	return p, nil
}

// MulScalar multiplies each coefficients of p with scalar
func (p *Poly) MulScalar(p1 *Poly, scalar bigint.Int) (*Poly, error) {
	if p.n != p1.n || !p.q.EqualTo(&p1.q) {
		return nil, errors.New("unmatched degree or module")
	}
	for i := range p.coeffs {
		p.coeffs[i].Mul(&p1.coeffs[i], &scalar)
	}
	return p, nil
}

// MulPoly multiplies p1 and p2 in polynomial style
func (p *Poly) MulPoly(p1, p2 *Poly) (*Poly, error) {
	if p.n != p1.n || !p.q.EqualTo(&p1.q) {
		return nil, errors.New("unmatched degree or module")
	}
	p1.NTT()
	p2.NTT()
	p.MulCoeffs(p1, p2)
	p.Mod(p, p.q)
	p.InverseNTT()
	if p != p1 {
		p1.InverseNTT()
	}
	if p != p2 {
		p2.InverseNTT()
	}
	return p, nil
}

// NaiveMultPoly implements a very basic polynomial multiplication,
// its results are the same with MulPoly.
func (p *Poly) NaiveMultPoly(p1, p2 *Poly) (*Poly, error) {
	r := make([]bigint.Int, p.n * 2)
	coeffs := make([]bigint.Int, p.n)
	coeffs1 := p1.GetCoefficients()
	coeffs2 := p2.GetCoefficients()
	tmp := new(bigint.Int)
	for i := uint32(0); i < p1.n; i++ {
		for j := uint32(0); j < p2.n; j++ {
			tmp.Mul(&coeffs1[i], &coeffs2[j])
			r[i+j].Add(&r[i+j], tmp)
			r[i+j].Mod(&r[i+j], &p.q)
		}
	}
	for i := p.n; i < 2*p.n-1; i++ {
		r[i-p.n].Sub(&r[i-p.n], &r[i])
		r[i-p.n].Mod(&r[i-p.n], &p.q)
	}
	for i := uint32(0); i < p.n; i++ {
		coeffs[i].SetBigInt(&r[i])
	}
	p.SetCoefficients(coeffs)
	return p, nil
}

// Div divides each coefficient of p1 by scalar, and sets p to the floor division results
func (p *Poly) Div(p1 *Poly, scalar bigint.Int) (*Poly, error) {
	if p.n != p1.n || !p.q.EqualTo(&p1.q) {
		return nil, errors.New("unmatched degree or module")
	}
	if scalar.EqualTo(bigint.NewInt(int64(0))) {
		return nil, errors.New("divisor cannot be zero")
	}
	for i := range p.coeffs {
		p.coeffs[i].Div(&p1.coeffs[i], &scalar)
	}
	return p, nil
}

// DivRound divides each coefficient of p1 by scalar, and sets p to the round division results
func (p *Poly) DivRound(p1 *Poly, scalar bigint.Int) (*Poly, error) {
	if p.n != p1.n || !p.q.EqualTo(&p1.q) {
		return nil, errors.New("unmatched degree or module")
	}
	if scalar.EqualTo(bigint.NewInt(int64(0))) {
		return nil, errors.New("divisor cannot be zero")
	}
	for i := range p.coeffs {
		p.coeffs[i].DivRound(&p1.coeffs[i], &scalar)
	}
	return p, nil
}

// Mod sets p to p1 mod m
func (p *Poly) Mod(p1 *Poly, m bigint.Int) (*Poly, error) {
	if p.n != p1.n || !p.q.EqualTo(&p1.q) {
		return nil, errors.New("unmatched degree or module")
	}
	for i := range p.coeffs {
		p.coeffs[i].Mod(&p1.coeffs[i], &m)
	}
	return p, nil
}

// And sets p to p1&m
func (p *Poly) And(p1 *Poly, m bigint.Int) (*Poly, error) {
	if p.n != p1.n || !p.q.EqualTo(&p1.q) {
		return nil, errors.New("unmatched degree or module")
	}
	for i := range p.coeffs {
		p.coeffs[i].And(&p1.coeffs[i], &m)
	}
	return p, nil
}

// Lsh sets p to p1 << m
func (p *Poly) Lsh(p1 *Poly, m uint32) (*Poly, error) {
	if p.n != p1.n || !p.q.EqualTo(&p1.q) {
		return nil, errors.New("unmatched degree or module")
	}
	for i := range p.coeffs {
		p.coeffs[i].Lsh(&p1.coeffs[i], m)
	}
	return p, nil
}

// Rsh sets p to p1 >> m
func (p *Poly) Rsh(p1 *Poly, m uint32) (*Poly, error) {
	if p.n != p1.n || !p.q.EqualTo(&p1.q) {
		return nil, errors.New("unmatched degree or module")
	}
	for i := range p.coeffs {
		p.coeffs[i].Rsh(&p1.coeffs[i], m)
	}
	return p, nil
}
