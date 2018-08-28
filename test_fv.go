package main

import (
	"github.com/dedis/student_18_lattices/crypto"
	"github.com/dedis/student_18_lattices/bigint"
	"github.com/dedis/student_18_lattices/ring"
	"fmt"
	"os"
)

func main() {
	N := uint32(256)
	Q := bigint.NewInt(int64(8380417))
	T := bigint.NewInt(int64(30))
	fv := crypto.NewFVContext(N, *Q, *T)
	fv.KeyGenerate()
	plain := new(crypto.Plaintext)
	plain.Msg = ring.NewRing(N, *Q)
	coeffs := make([]bigint.Int, N)
	for i := range coeffs {
		coeffs[i].SetInt(int64(2))
	}
	plain.Msg.Poly.SetCoefficients(coeffs)
	cipher := fv.Encrypt(plain)
	newPlain := fv.Decrypt(cipher)

	msg1 := plain.Msg.GetCoefficientsInt64()
	msg2 := newPlain.Msg.GetCoefficientsInt64()
	fmt.Printf("msg1: %v\n", msg1)
	fmt.Printf("msg2: %v\n", msg2)
	for i := range msg1 {
		if msg1[i] != msg2[i] {
			os.Exit(0)
		}
	}
	fmt.Println("PASS")

}
