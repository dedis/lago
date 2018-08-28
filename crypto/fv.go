package crypto

import (
	"github.com/dedis/student_18_lattices/bigint"
	"github.com/dedis/student_18_lattices/polynomial"
)


// This code implements the text-FV scheme in paper https://eprint.iacr.org/2012/144.pdf,
// including key generation, encryption, and decryption.

type FVContext struct {
	N uint32
	Q bigint.Int
	T bigint.Int
	Delta bigint.Int  // floor(ciphertext modulus / plaintext modulus)
	Sigma float64
	NttParams *polynomial.NttParams
}

// NewFVContext creates a new FV context containing all required parameters.
func NewFVContext(N uint32, Q, T bigint.Int) *FVContext {
	fv := new(FVContext)
	fv.N = N  // polynomial degree
	fv.Q = Q  // ciphertext modulus
	fv.T = T  // plaintext modulus
	fv.Delta.Div(&Q, &T)
	fv.Sigma = 3.19  // distributed gaussian noise parameter
	fv.NttParams = polynomial.GenerateNTTParams(N, Q)
	return fv
}
