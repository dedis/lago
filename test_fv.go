package main

import (
	"github.com/dedis/student_18_lattices/crypto"
	"github.com/dedis/student_18_lattices/bigint"
	"github.com/dedis/student_18_lattices/ring"
	"fmt"
)

func main() {
	N := uint32(256)
	Q := bigint.NewInt(int64(7681))
	T := bigint.NewInt(int64(10))
	fv := crypto.NewFVContext(N, *Q, *T)
	fv.KeyGenerate()
	plain := new(crypto.Plaintext)
	plain.Msg = ring.NewRing(N, *Q)
	coeffs := make([]bigint.Int, N)
	for i := range coeffs {
		coeffs[i].SetInt(int64(9))
	}
	plain.Msg.Poly.SetCoefficients(coeffs)
	cipher := fv.Encrypt(plain)
	newPlain := fv.Decrypt(cipher)

	fmt.Println(plain.Msg.GetCoefficientsInt64())
	fmt.Println(newPlain.Msg.GetCoefficientsInt64())

}
