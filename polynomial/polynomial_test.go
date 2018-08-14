package polynomial

import (
	"testing"
	"github.com/dedis/student_18_lattices/bigint"
)

func TestPoly_MulPoly(t *testing.T) {
	q := int64(7681)
	n := uint32(256)
	p, _ := NewPolynomial(n, *bigint.NewInt(q))
	p1, _ := NewPolynomial(n, *bigint.NewInt(q))
	p2, _ := NewPolynomial(n, *bigint.NewInt(q))
	coeffs := make([]bigint.Int, n)
	for i := range coeffs {
		coeffs[i].SetInt(0)
	}
	coeffs[0].SetInt(1)
	p1.SetCoefficients(coeffs)
	for i := range coeffs {
		coeffs[i].SetInt(0)
	}
	coeffs[0].SetInt(1)
	p2.SetCoefficients(coeffs)
	p.MulPoly(p1, p2)
}
