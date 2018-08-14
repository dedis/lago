package polynomial

import (
	"github.com/dedis/student_18_lattices/bigint"
)

// NTT performs the number theoretic transform on polynomial p's coefficients
// The implementation is based on the description from https://eprint.iacr.org/2016/504.pdf,
// while the underlying algorithm originates from
// https://www.usenix.org/system/files/conference/usenixsecurity16/sec16_paper_alkim.pdf
func (p *Poly) NTT() (*Poly, error) {
	var j1, j2 uint32
	var U, V, T bigint.Int
	var S *bigint.Int
	t := p.n
	for m := uint32(1); m < p.n; m <<= 1 {
		t >>= 1
		for i := uint32(0); i < m; i++ {
			j1 = 2 * i * t
			j2 = j1 + t - 1
			S = &p.psiReverse[m+i]
			for j := j1; j <= j2; j++ {
				// TODO: implement fast reduction algorithms
				U.SetBigInt(&p.coeffs[j])
				V.Mul(&p.coeffs[j+t], S)
				V.Mod(&V, &p.q)
				T.Add(&U, &V)
				p.coeffs[j].Mod(&T, &p.q)
				T.Sub(&U, &V)
				p.coeffs[j+t].Mod(&T, &p.q)
			}
		}
	}
	return p, nil
}

// NTT performs the inverse number theoretic transform on polynomial p's coefficients
func (p *Poly) InverseNTT() (*Poly, error) {
	var j1, j2, h uint32
	var U, V, T bigint.Int
	var S *bigint.Int
	t := uint32(1)
	for m := p.n; m > 1; m >>= 1 {
		j1 = 0
		h = m >> 1
		for i := uint32(0); i < h; i++ {
			j2 = j1 + t - 1
			S = &p.psiInvReverse[h+i]
			for j := j1; j <= j2; j++ {
				U.SetBigInt(&p.coeffs[j])
				V.SetBigInt(&p.coeffs[j+t])
				T.Add(&U, &V)
				p.coeffs[j].Mod(&T, &p.q)
				T.Sub(&U, &V)
				T.Mul(&T, S)
				p.coeffs[j+t].Mod(&T, &p.q)
			}
			j1 = j1 + (t << 1)
		}
		t <<= 1
	}
	var n_reverse bigint.Int
	n_reverse.Inv(bigint.NewInt(int64(p.n)), &p.q)
	for j := uint32(0); j < p.n; j++ {
		p.coeffs[j].Mod(p.coeffs[j].Mul(&p.coeffs[j], &n_reverse), &p.q)
	}
	return p, nil
}
