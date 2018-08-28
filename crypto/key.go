package crypto

import (
	"github.com/dedis/student_18_lattices/ring"
	"github.com/dedis/student_18_lattices/bigint"
)

type Key struct {
	PubKey PublicKey
	SecKey SecretKey
	EvaKey EvaluationKey
	EvaSize uint32
}

type PublicKey = [2]*ring.Ring

type SecretKey = *ring.Ring

type EvaluationKey = [][2]*ring.Ring

// KeyGenerator generates the public key and secret key of given FV context
func GenerateKey(fv *FVContext) *Key {
	key := new(Key)
	err := *new(error)
	// generate secret key
	key.SecKey, err = ring.NewUniformPoly(fv.N, fv.Q, fv.NttParams, *bigint.NewInt(int64(2)))
	if err != nil {
		panic(err)
	}
	key.SecKey.Poly.NTT()  // store secret key in NTT form

	// generate public key: PubKey[0] = e - a * sk, PubKey[1] = a
	key.PubKey[1], err = ring.NewGaussPoly(fv.N, fv.Q, fv.NttParams, fv.Sigma)
	if err != nil {
		panic(err)
	}
	key.PubKey[1].Poly.NTT()

	a_sk, err := ring.NewRing(fv.N, fv.Q, fv.NttParams)
	if err != nil {
		panic(err)
	}
	_, err = a_sk.MulCoeffs(key.PubKey[1], key.SecKey)
	if err != nil {
		panic(err)
	}

	key.PubKey[0], err = ring.NewGaussPoly(fv.N, fv.Q, fv.NttParams, fv.Sigma)
	if err != nil {
		panic(err)
	}
	key.PubKey[0].Poly.NTT()

	_, err = key.PubKey[0].Sub(key.PubKey[0], a_sk)
	if err != nil {
		panic(err)
	}

	// generate evaluation key
	l := fv.Q.Value.BitLen()
	key.EvaKey = make([][2]*ring.Ring, l)
	key.EvaSize = 16

	w := bigint.NewInt(1)  // decomposition base, corresponding to T^i in the paper, here we choose T=2
	for i := 0; i < l ; i++ {
		// evaluationKey[i][1] = a_i, where a_i sampled from R_q
		key.EvaKey[i][1], err = ring.NewUniformPoly(fv.N, fv.Q, fv.NttParams, fv.Q)
		if err != nil {
			panic(err)
		}
		key.EvaKey[i][1].Poly.NTT()

		// evaluationKey[i][0] = -(a_i * s + e_i) + T^i * s * s mod q
		key.EvaKey[i][0], err = ring.NewGaussPoly(fv.N, fv.Q, fv.NttParams, fv.Sigma)
		if err != nil {
			panic(err)
		}
		key.EvaKey[i][0].Poly.NTT()

		tmp1, err := ring.NewRing(fv.N, fv.Q, fv.NttParams)
		if err != nil {
			panic(err)
		}
		tmp1.MulCoeffs(key.EvaKey[i][1], key.SecKey)

		tmp2, err := ring.NewRing(fv.N, fv.Q, fv.NttParams)
		if err != nil {
			panic(err)
		}
		tmp2.MulCoeffs(key.SecKey, key.SecKey)
		tmp2.MulScalar(tmp2, *w)

		key.EvaKey[i][0].Sub(key.EvaKey[i][0], tmp1)
		key.EvaKey[i][0].Add(key.EvaKey[i][0], tmp2)

		w.Lsh(w, key.EvaSize)
	}

	return key
}