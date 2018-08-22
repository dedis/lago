package crypto

import (
	"testing"
	"github.com/dedis/student_18_lattices/bigint"
	"github.com/dedis/student_18_lattices/ring"
)

func TestFV(t *testing.T) {
	N := uint32(256)
	Q := bigint.NewInt(int64(8380417))
	T := bigint.NewInt(int64(30))
	fv := NewFVContext(N, *Q, *T)
	fv.KeyGenerate()
	plain := new(Plaintext)
	plain.Msg = ring.NewUniformPoly(N, *Q, *T)
	cipher := fv.Encrypt(plain)
	newPlain := fv.Decrypt(cipher)

	msg1 := plain.Msg.GetCoefficientsInt64()
	msg2 := newPlain.Msg.GetCoefficientsInt64()
	for i := range msg1 {
		if msg1[i] != msg2[i] {
			t.Errorf("Error in FV: expected %v, got %v", msg1[i], msg2[i])
		}
	}
}
