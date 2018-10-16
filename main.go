package main

import (
	"github.com/dedis/lago/bigint"
	"github.com/dedis/lago/crypto"
	"github.com/dedis/lago/encoding"
	"fmt"
)

func main() {
	msg1 := bigint.NewInt(10)
	msg2 := bigint.NewInt(8)

	N := uint32(32)  // polynomial degree
	T := bigint.NewInt(10)  // plaintext moduli
	Q := bigint.NewInt(8380417)  // ciphertext moduli
	BigQ := bigint.NewIntFromString("4611686018326724609")  // big ciphertext moduli, used in homomorphic multiplication and should be greater than q^2

	// create FV context and generate keys
	fv := crypto.NewFVContext(N, *T, *Q, *BigQ)
	key := crypto.GenerateKey(fv)

	// encode messages
	encoder := encoding.NewEncoder(fv)
	plaintext1 := crypto.NewPlaintext(N, *Q, fv.NttParams)
	plaintext2 := crypto.NewPlaintext(N, *Q, fv.NttParams)
	encoder.Encode(msg1, plaintext1)
	encoder.Encode(msg2, plaintext2)

	// encrypt plainetexts
	encryptor := crypto.NewEncryptor(fv, &key.PubKey)
	ciphertext1 := encryptor.Encrypt(plaintext1)
	ciphertext2 := encryptor.Encrypt(plaintext2)

	// evaluate ciphertexts
	evaluator := crypto.NewEvaluator(fv, &key.EvaKey, key.EvaSize)
	add_cipher := evaluator.Add(ciphertext1, ciphertext2)
	mul_cipher := evaluator.Multiply(add_cipher, ciphertext2)

	// decrypt ciphertexts
	decryptor := crypto.NewDecryptor(fv, &key.SecKey)
	new_plaintext1 := decryptor.Decrypt(ciphertext1)
	new_plaintext2 := decryptor.Decrypt(ciphertext2)
	add_plaintext := decryptor.Decrypt(add_cipher)
	mul_plaintext := decryptor.Decrypt(mul_cipher)

	// decode messages
	new_msg1 := new(bigint.Int)
	new_msg2 := new(bigint.Int)
	add_msg := new(bigint.Int)
	mul_msg := new(bigint.Int)
	encoder.Decode(new_msg1, new_plaintext1)
	encoder.Decode(new_msg2, new_plaintext2)
	encoder.Decode(add_msg, add_plaintext)
	encoder.Decode(mul_msg, mul_plaintext)

	fmt.Printf("%v + %v = %v\n", msg1.Int64(), msg2.Int64(), add_msg.Int64())
	fmt.Printf("(%v + %v) * %v = %v\n", msg1.Int64(), msg2.Int64(), msg2.Int64(), mul_msg.Int64())
}
