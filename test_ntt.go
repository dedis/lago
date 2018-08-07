package main

import (
	"github.com/dedis/student_18_lattices/bigint"
	"fmt"
	"github.com/dedis/student_18_lattices/kyber"
	"github.com/LoCCS/bliss/poly"
	"github.com/dedis/student_18_lattices/ring"
)

func main() {
	N := uint32(256)
	Q := bigint.NewInt(7681)
	// Two set of coefficients (polynomials)
	p1 := ring.NewGaussPoly(N, *Q, *Q)
	coeffs1 := p1.GetCoefficientsInt64()
	p2 := ring.NewUniformPoly(N, *Q, *bigint.NewInt(2))
	coeffs2 := p2.GetCoefficientsInt64()


	// kyber.NTT
	coeffs1Kyber := [256]uint16{}
	coeffs2Kyber:= [256]uint16{}
	for i := range coeffs1Kyber {
		coeffs1Kyber[i] = uint16(coeffs1[i])
		coeffs2Kyber[i] = uint16(coeffs2[i])
	}

	kyber.NttRef(&coeffs1Kyber)
	kyber.NttRef(&coeffs2Kyber)
	for i := range coeffs1 {
		coeffs1Kyber[i] = coeffs1Kyber[i] * coeffs2Kyber[i]
		coeffs1Kyber[i] = coeffs1Kyber[i] % uint16(Q.Int64())
	}
	kyber.InvnttRef(&coeffs1Kyber)

	nttResultKyber := coeffs1Kyber
	fmt.Printf("Kyber : %v\n", nttResultKyber)


	//bliss.NTT
	coeffs1Bliss := make([]int32, N)
	coeffs2Bliss := make([]int32, N)
	for i := range coeffs1Bliss {
		coeffs1Bliss[i] = int32(coeffs1[i])
		coeffs2Bliss[i] = int32(coeffs2[i])
	}
	p1Bliss, _ := poly.New(0)
	p2Bliss, _ := poly.New(0)
	p1Bliss.SetData(coeffs1Bliss)
	p2Bliss.SetData(coeffs2Bliss)

	p1BlissNTT, _ := p1Bliss.NTT()
	nttResultBliss, _ := p2Bliss.MultiplyNTT(p1BlissNTT)
	fmt.Printf("Bliss : %v\n", nttResultBliss)


	//our.NTT
	p := ring.NewRing(N, *Q)
	p.MulPoly(p1, p2)
	nttResultOur := p.GetCoefficientsInt64()
	fmt.Printf("Our   : %v\n", nttResultOur)

}
