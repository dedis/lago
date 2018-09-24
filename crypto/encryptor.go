package crypto

import (
	"github.com/dedis/student_18_lattices/bigint"
	"github.com/dedis/student_18_lattices/ring"
)

type Encryptor struct {
	ctx *FVContext	  // FV context
	publickey *PublicKey  // public key
}

// NewEncryptor creates a new Encryptor for encryption
func NewEncryptor(ctx *FVContext, publickey *PublicKey) *Encryptor {
	encryptor := new(Encryptor)
	encryptor.ctx = ctx
	encryptor.publickey = publickey
	return encryptor
}

// Encrypt encrypts plaintext to ciphertext with encryptor parameters,
// both plaintext and ciphertext are in NTT form.
func (encryptor *Encryptor) Encrypt(plaintext *Plaintext) *Ciphertext {
	plaintext.Value.Poly.NTT()
	// deltaM = delta * m
	deltaM, err := ring.NewRing(encryptor.ctx.N, encryptor.ctx.Q, encryptor.ctx.NttParams)
	if err != nil {
		panic(err)
	}
	deltaM.MulScalar(plaintext.Value, encryptor.ctx.Delta)

	// u sampled from R_2, e1 and e2 sampled from gaussian
	u, err := ring.NewUniformPoly(encryptor.ctx.N, encryptor.ctx.Q, encryptor.ctx.NttParams, *bigint.NewInt(2))
	if err != nil {
		panic(err)
	}
	u.Poly.NTT() // turn u to NTT form for polynomial multiplication

	e1, err := ring.NewGaussPoly(encryptor.ctx.N, encryptor.ctx.Q, encryptor.ctx.NttParams, encryptor.ctx.Sigma)
	if err != nil {
		panic(err)
	}
	e1.Poly.NTT()
	e2, err := ring.NewGaussPoly(encryptor.ctx.N, encryptor.ctx.Q, encryptor.ctx.NttParams, encryptor.ctx.Sigma)
	if err != nil {
		panic(err)
	}
	e2.Poly.NTT()

	// Ciphertext = (c0, c1)
	// c0 = delta * m + publickey[0] * u + e1
	// c1 = publickey[1] * u + e2
	ciphertext := NewCiphertext(encryptor.ctx.N, encryptor.ctx.Q, encryptor.ctx.NttParams)
	ciphertext.value[0].MulCoeffs(encryptor.publickey[0], u)
	ciphertext.value[0].Add(ciphertext.value[0], deltaM)
	ciphertext.value[0].Add(ciphertext.value[0], e1)

	ciphertext.value[1].MulCoeffs(encryptor.publickey[1], u)
	ciphertext.value[1].Add(ciphertext.value[1], e2)

	plaintext.Value.Poly.InverseNTT()
	return ciphertext
}
