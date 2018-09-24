package encoding

import (
	"github.com/dedis/student_18_lattices/bigint"
	"github.com/dedis/student_18_lattices/crypto"
)

type Encoder struct {
	N uint32
	Base uint32
}

// NewEncoder creates an Encoder.
func NewEncoder(ctx *crypto.FVContext) *Encoder {
	encoder := new(Encoder)
	encoder.N = ctx.N
	encoder.Base = 2  // base is set to 2.
	return encoder
}

// Encode encodes an integer to a plaintext (polynomial ring).
func (encoder *Encoder) Encode(msg *bigint.Int, plaintext *crypto.Plaintext) {
	msgBits, bitLen := msg.Bits()
	coeffs := make([]bigint.Int, encoder.N)
	for i := uint(0); i < bitLen; i++ {
		coeffs[i].SetInt(int64(msgBits[i]))
	}
	plaintext.Value.Poly.SetCoefficients(coeffs)
}

// Decode decodes an integer from a plaintext (polynomial ring).
func (encoder *Encoder) Decode(msg *bigint.Int, plaintext *crypto.Plaintext) {
	msg.SetInt(0)
	coeffs := plaintext.Value.GetCoefficients()
	tmp := new(bigint.Int)
	for i := uint32(0); i < encoder.N; i++ {
		tmp.Mul(&coeffs[i], bigint.NewInt(int64(1 << i)))
		msg.Add(msg, tmp)
	}
}
