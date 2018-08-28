package crypto

type Decryptor struct {
	ctx *FVContext	  // FV context
	secretkey *SecretKey   // secret key
}

// NewDecryptor creates a new Decryptor for decryption
func NewDecryptor(ctx *FVContext, secretkey *SecretKey) *Decryptor {
	decryptor := new(Decryptor)
	decryptor.ctx = ctx
	decryptor.secretkey = secretkey
	return decryptor
}

// Decrypt decrypts ciphertext to plaintext with decryptor parameters,
// both ciphertext and plaintext are in NTT form.
func (decryptor *Decryptor) Decrypt(ciphertext *Ciphertext) *Plaintext {
	plaintext := NewPlaintext(decryptor.ctx.N, decryptor.ctx.Q, decryptor.ctx.NttParams)
	plaintext.value.MulCoeffs(ciphertext.value[1], *decryptor.secretkey)
	plaintext.value.Add(plaintext.value, ciphertext.value[0])
	plaintext.value.Poly.InverseNTT()
	plaintext.value.MulScalar(plaintext.value, decryptor.ctx.T)
	plaintext.value.DivRound(plaintext.value, decryptor.ctx.Q)
	plaintext.value.Mod(plaintext.value, decryptor.ctx.T)
	return plaintext
}
