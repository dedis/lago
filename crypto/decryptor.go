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
	plaintext.Value.MulCoeffs(ciphertext.value[1], *decryptor.secretkey)
	plaintext.Value.Add(plaintext.Value, ciphertext.value[0])
	plaintext.Value.Poly.InverseNTT()
	center(plaintext.Value)
	plaintext.Value.MulScalar(plaintext.Value, decryptor.ctx.T)
	plaintext.Value.DivRound(plaintext.Value, decryptor.ctx.Q)
	plaintext.Value.Mod(plaintext.Value, decryptor.ctx.T)
	return plaintext
}
