package polynomial

import (
	"testing"
	"github.com/dedis/student_18_lattices/bigint"
	"fmt"
)

func TestPoly_MulPoly(t *testing.T) {
	q := int64(Q)
	p, _ := NewPolynomial(N, *bigint.NewInt(q))
	p1, _ := NewPolynomial(N, *bigint.NewInt(q))
	p2, _ := NewPolynomial(N, *bigint.NewInt(q))
	coeffs := make([]bigint.Int, N)
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
	fmt.Println(p.GetCoefficients())
}
