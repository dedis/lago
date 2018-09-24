package crypto

import (
	"github.com/dedis/student_18_lattices/bigint"
	"github.com/dedis/student_18_lattices/polynomial"
	"github.com/dedis/student_18_lattices/ring"
)


// This code implements the text-FV scheme in paper https://eprint.iacr.org/2012/144.pdf,
// including key generation, encryption, and decryption.

type FVContext struct {
	N uint32  // polynomial degree
	T bigint.Int  // plaintext modulus
	Q bigint.Int  // ciphertext modulus
	BigQ bigint.Int  // big ciphertext modulus, used for relinearisation
	Delta bigint.Int  // floor(ciphertext modulus / plaintext modulus)
	InvDelta bigint.Int
	Sigma float64
	NttParams *polynomial.NttParams
	BigNttParams *polynomial.NttParams
}

// NewFVContext creates a new FV context containing all required parameters.
func NewFVContext(N uint32, T, Q, BigQ bigint.Int) *FVContext {
	fv := new(FVContext)
	fv.N = N
	fv.T = T
	fv.Q = Q
	fv.BigQ = BigQ
	fv.Delta.Div(&Q, &T)
	fv.InvDelta.Inv(&fv.Delta, &Q)
	fv.Sigma = 3.19  // distributed gaussian noise parameter, suggested by SEAL library.
	fv.NttParams = polynomial.GenerateNTTParams(N, Q)
	fv.BigNttParams = polynomial.GenerateNTTParams(N, BigQ)
	return fv
}

// center shifts r from [0, q) to (-q/2, q/2]
func center(r *ring.Ring) {
	coeffs := r.GetCoefficients()
	qDiv2 := bigint.NewInt(1)
	qDiv2.Div(&r.Q, bigint.NewInt(2))
	for i := range coeffs {
		if coeffs[i].Compare(qDiv2) == 1.0 {
			coeffs[i].Sub(&coeffs[i], &r.Q)
		}
	}
	r.Poly.SetCoefficients(coeffs)
}
