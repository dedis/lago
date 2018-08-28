package main

import (
	"github.com/dedis/student_18_lattices/bigint"
	"math/rand"
	"fmt"
	"github.com/dedis/student_18_lattices/polynomial"
)

func main() {
	q := bigint.NewInt(7681)
	n := uint32(256)
	nttParams := polynomial.GenerateNTTParams(n, *q)
	p, _ := polynomial.NewPolynomial(n, *q, nttParams)
	coeffs := make([]bigint.Int, n)
	halfQ := bigint.NewInt(2).Div(q, bigint.NewInt(2))
	for i := range coeffs {
		coeffs[i].SetInt(int64(rand.Int63n(7681)))
		coeffs[i].Mod(&coeffs[i], q)
		coeffs[i].Sub(&coeffs[i], halfQ)

	}
	p.SetCoefficients(coeffs)
	fmt.Println(p.GetCoefficientsInt64())
	p.NTT()
	fmt.Println(p.GetCoefficientsInt64())
	p.InverseNTT()
	fmt.Println(p.GetCoefficientsInt64())

}
