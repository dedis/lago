package polynomial

import (
	"testing"
	"github.com/dedis/student_18_lattices/bigint"
)

func TestPoly_NTT(t *testing.T) {
	q := int64(Q)
	p, _ := NewPolynomial(N, *bigint.NewInt(q))
	coeffs := make([]bigint.Int, N)
	for i := range coeffs {
		coeffs[i].SetString("10", 0)
	}
	p.SetCoefficients(coeffs)
	p.NTT()
	p.InverseNTT()
	for i := range coeffs {
		if ! coeffs[i].Mod(&coeffs[i], &p.q).EqualTo(&p.coeffs[i]) {
			t.Error("Coefficients not euqal!")
		}
	}

}

func BenchmarkPoly_NTT(b *testing.B) {
	q := int64(Q)
	p, _ := NewPolynomial(N, *bigint.NewInt(q))
	coeffs := make([]bigint.Int, N)
	for i := range coeffs {
		coeffs[i].SetString("123456789000", 0)
	}
	p.SetCoefficients(coeffs)
	for i := 0; i < b.N; i++ {
		p.NTT()
	}
}

func BenchmarkPoly_InverseNTT(b *testing.B) {
	q := int64(Q)
	p, _ := NewPolynomial(N, *bigint.NewInt(q))
	coeffs := make([]bigint.Int, N)
	for i := range coeffs {
		coeffs[i].SetString("123456789000", 0)
	}
	p.SetCoefficients(coeffs)
	for i := 0; i < b.N; i++ {
		p.NTT()
	}
}
