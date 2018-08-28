package crypto

import (
	"testing"
	"github.com/dedis/student_18_lattices/bigint"
	"github.com/dedis/student_18_lattices/ring"
)

func TestFVContext(t *testing.T) {
	N := uint32(256)
	Q := bigint.NewInt(int64(8380417))
	T := bigint.NewInt(int64(30))
	fv := NewFVContext(N, *Q, *T)
	key := GenerateKey(fv)
	plaintext1 := new(Plaintext)
	plaintext1.value, _ = ring.NewUniformPoly(N, *Q, fv.NttParams, *T)
	msg1 := plaintext1.value.GetCoefficientsInt64()
	plaintext1.value.Poly.NTT()

	encryptor := NewEncryptor(fv, &key.PubKey)
	ciphertext := encryptor.Encrypt(plaintext1)

	decryptor := NewDecryptor(fv, &key.SecKey)
	plaintext2 := decryptor.Decrypt(ciphertext)

	msg2 := plaintext2.value.GetCoefficientsInt64()
	for i := range msg1 {
		if msg1[i] != msg2[i] {
			t.Errorf("Error in FV: expected %v, got %v", msg1[i], msg2[i])
		}
	}
}
