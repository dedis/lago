package crypto

import (
	"testing"
	"github.com/dedis/lago/bigint"
	"io/ioutil"
	"fmt"
	"strings"
	"strconv"
)

func TestFVContext(t *testing.T) {
	for i := 0; i <=1; i++ {
		testfile, err := ioutil.ReadFile(fmt.Sprintf("test_data/testvector_fv_%d", i))
		if err != nil {
			t.Errorf("Failed to open file: %s", err.Error())
		}
		filecontent := strings.TrimSpace(string(testfile))
		vs := strings.Split(filecontent, "\n")
		if len(vs) != 9 {
			t.Errorf("Error in data read from test_data: len(vs) = %d", len(vs))
		}

		// load BigQ
		BigQ := vs[0]

		// load Q
		Q := vs[1]

		// load T
		T, err := strconv.Atoi(vs[2])
		if err != nil {
			t.Errorf("Invalid integer: %v", vs[2])
		}

		// load N
		N, err := strconv.Atoi(vs[3])
		if err != nil {
			t.Errorf("Invalid integer: %v", vs[3])
		}

		// create FV context
		fv := NewFVContext(uint32(N), *bigint.NewInt(int64(T)), *bigint.NewIntFromString(Q), *bigint.NewIntFromString(BigQ))
		// generate new keys
		key := GenerateKey(fv)

		// load first plaintext
		plaintext1String := strings.Split(strings.TrimSpace(vs[4]), ", ")
		plaintext1Coeffs := make([]bigint.Int, N)
		for i := range plaintext1Coeffs {
			tmp, err := strconv.Atoi(plaintext1String[i])
			if err != nil {
				t.Errorf("Invalid integer of p1 coeffs: %v", plaintext1String[i])
			}
			plaintext1Coeffs[i].SetInt(int64(tmp))
		}
		plaintext1 := NewPlaintext(fv.N, fv.Q, fv.NttParams)
		plaintext1.Value.Poly.SetCoefficients(plaintext1Coeffs)

		// load second plaintext
		plaintext2String := strings.Split(strings.TrimSpace(vs[5]), ", ")
		plaintext2Coeffs := make([]bigint.Int, N)
		for i := range plaintext2Coeffs {
			tmp, err := strconv.Atoi(plaintext2String[i])
			if err != nil {
				t.Errorf("Invalid integer of p2 coeffs: %v", plaintext2String[i])
			}
			plaintext2Coeffs[i].SetInt(int64(tmp))
		}
		plaintext2 := NewPlaintext(fv.N, fv.Q, fv.NttParams)
		plaintext2.Value.Poly.SetCoefficients(plaintext2Coeffs)

		// load add plaintext
		plaintextAddString := strings.Split(strings.TrimSpace(vs[6]), ", ")
		plaintextAddCoeffs := make([]bigint.Int, N)
		for i := range plaintextAddCoeffs {
			tmp, err := strconv.Atoi(plaintextAddString[i])
			if err != nil {
				t.Errorf("Invalid integer of add coeffs: %v", plaintextAddString[i])
			}
			plaintextAddCoeffs[i].SetInt(int64(tmp))
		}

		// load sub plaintext
		plaintextSubString := strings.Split(strings.TrimSpace(vs[7]), ", ")
		plaintextSubCoeffs := make([]bigint.Int, N)
		for i := range plaintextSubCoeffs {
			tmp, err := strconv.Atoi(plaintextSubString[i])
			if err != nil {
				t.Errorf("Invalid integer of sub coeffs: %v", plaintextSubString[i])
			}
			plaintextSubCoeffs[i].SetInt(int64(tmp))
		}

		// load mult plaintext
		plaintextMultString := strings.Split(strings.TrimSpace(vs[8]), ", ")
		plaintextMultCoeffs := make([]bigint.Int, N)
		for i := range plaintextMultCoeffs {
			tmp, err := strconv.Atoi(plaintextMultString[i])
			if err != nil {
				t.Errorf("Invalid integer of mult coeffs: %v", plaintextMultString[i])
			}
			plaintextMultCoeffs[i].SetInt(int64(tmp))
		}

		// encrypt
		encryptor := NewEncryptor(fv, &key.PubKey)
		ciphertext1 := encryptor.Encrypt(plaintext1)
		ciphertext2 := encryptor.Encrypt(plaintext2)
		// decrypt
		decryptor := NewDecryptor(fv, &key.SecKey)
		new_plaintext1 := decryptor.Decrypt(ciphertext1)

		// test encrypt and decrypt
		new_msg1 := new_plaintext1.Value.GetCoefficients()
		for i := uint32(0); i < fv.N; i++ {
			if ! new_msg1[i].EqualTo(&plaintext1Coeffs[i]) {
				t.Errorf("Error in enc/dec, expected %v, got %v", new_msg1[i].Int64(), plaintext1Coeffs[i].Int64())
			}
		}

		// generate new evaluator for addition, subtraction, multiplication
		evaluator := NewEvaluator(fv, &key.EvaKey, key.EvaSize)

		// add
		add_cipher := evaluator.Add(ciphertext1, ciphertext2)
		add_plain := decryptor.Decrypt(add_cipher)
		// test add
		add_msg := add_plain.Value.GetCoefficients()
		for i := uint32(0); i < fv.N; i++ {
			if ! add_msg[i].EqualTo(&plaintextAddCoeffs[i]) {
				t.Errorf("Error in add, expected %v, got %v", plaintextAddCoeffs[i].Int64(), add_msg[i].Int64())
			}
		}

		// sub
		sub_cipher := evaluator.Sub(ciphertext1, ciphertext2)
		sub_plain := decryptor.Decrypt(sub_cipher)
		// test sub
		sub_msg := sub_plain.Value.GetCoefficients()
		for i := uint32(0); i < fv.N; i++ {
			if ! sub_msg[i].EqualTo(&plaintextSubCoeffs[i]) {
				t.Errorf("Error in sub, expected %v, got %v", plaintextSubCoeffs[i].Int64(), sub_msg[i].Int64())
			}
		}

		// multiply
		multiply_cipher := evaluator.Multiply(ciphertext1, ciphertext2)
		multiply_plain := decryptor.Decrypt(multiply_cipher)
		// test multiply
		multiply_msg := multiply_plain.Value.GetCoefficients()
		for i := uint32(0); i < fv.N; i++ {
			if ! multiply_msg[i].EqualTo(&plaintextMultCoeffs[i]) {
				t.Errorf("Error in multiply, expected %v, got %v", plaintextMultCoeffs[i].Int64(), multiply_msg[i].Int64())
			}
		}
	}
}

func BenchmarkFVContext(b *testing.B) {
	for i := 0; i <=0; i++ {
		testfile, _ := ioutil.ReadFile(fmt.Sprintf("test_data/testvector_fv_%d", i))
		filecontent := strings.TrimSpace(string(testfile))
		vs := strings.Split(filecontent, "\n")

		// load BigQ
		BigQ := vs[0]

		// load Q
		Q := vs[1]

		// load T
		T, _ := strconv.Atoi(vs[2])

		// load N
		N, _ := strconv.Atoi(vs[3])

		// create FV context
		fv := NewFVContext(uint32(N), *bigint.NewInt(int64(T)), *bigint.NewIntFromString(Q), *bigint.NewIntFromString(BigQ))
		// generate new keys
		key := GenerateKey(fv)

		// load first plaintext
		plaintext1String := strings.Split(strings.TrimSpace(vs[4]), ", ")
		plaintext1Coeffs := make([]bigint.Int, N)
		for i := range plaintext1Coeffs {
			tmp, _ := strconv.Atoi(plaintext1String[i])
			plaintext1Coeffs[i].SetInt(int64(tmp))
		}
		plaintext1 := NewPlaintext(fv.N, fv.Q, fv.NttParams)
		plaintext1.Value.Poly.SetCoefficients(plaintext1Coeffs)
		plaintext1.Value.Poly.NTT()

		// load second plaintext
		plaintext2String := strings.Split(strings.TrimSpace(vs[5]), ", ")
		plaintext2Coeffs := make([]bigint.Int, N)
		for i := range plaintext2Coeffs {
			tmp, _ := strconv.Atoi(plaintext2String[i])
			plaintext2Coeffs[i].SetInt(int64(tmp))
		}
		plaintext2 := NewPlaintext(fv.N, fv.Q, fv.NttParams)
		plaintext2.Value.Poly.SetCoefficients(plaintext2Coeffs)
		plaintext2.Value.Poly.NTT()

		// encrypt
		encryptor := NewEncryptor(fv, &key.PubKey)
		ciphertext1 := encryptor.Encrypt(plaintext1)
		ciphertext2 := encryptor.Encrypt(plaintext2)

		evaluator := NewEvaluator(fv, &key.EvaKey, key.EvaSize)
		b.ResetTimer()
		for j := 0; j < b.N ;j++ {
			evaluator.Multiply(ciphertext1, ciphertext2)
		}
	}
}
