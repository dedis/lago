package crypto

import (
	"github.com/dedis/student_18_lattices/ring"
	"github.com/dedis/student_18_lattices/bigint"
	"fmt"
)

type FV struct {
	N uint32
	Q bigint.Int
	T bigint.Int
	PublicKey *publicKey
	PrivateKey *privateKey
}

type publicKey struct {
	pub0 *ring.Ring
	pub1 *ring.Ring
}

type privateKey struct {
	s *ring.Ring
}

func NewFVContext(N uint32, Q, T bigint.Int) *FV {
	fv := new(FV)
	fv.N = N
	fv.Q = Q
	fv.T = T
	fv.PublicKey = new(publicKey)
	fv.PrivateKey = new(privateKey)
	return fv
}

func (fv *FV) KeyGenerate() {
	if fv.N == 0 || fv.Q.EqualTo(bigint.NewInt(0)) {
		panic("Invalid FV context")
	}
	// generate private key
	fv.PrivateKey.s = ring.NewUniformPoly(fv.N, fv.Q, *bigint.NewInt(int64(2)))

	//generate public key: pub0 = e - a * sk, pub1 = a
	a := ring.NewGaussPoly(fv.N, fv.Q, fv.Q)
	fv.PublicKey.pub1 = ring.NewRing(fv.N, fv.Q)
	fv.PublicKey.pub1.Poly.SetCoefficients(a.Poly.GetCoefficients())
	b := ring.NewGaussPoly(fv.N, fv.Q, *bigint.NewInt(int64(2)))

	a.MulPoly(a, fv.PrivateKey.s)

	fv.PublicKey.pub0, _ = b.Sub(b, a)
}

func (fv *FV) Encrypt(m *Plaintext) *Ciphertext {
	var delta bigint.Int
	newM := ring.NewRing(fv.N, fv.Q)
	delta.Div(&fv.Q, &fv.T)
	newM.MulScalar(m.Msg, delta)
	newM.Mod(newM, fv.Q)

	u := ring.NewUniformPoly(fv.N, fv.Q, *bigint.NewInt(int64(2)))
	e1 := ring.NewGaussPoly(fv.N, fv.Q, *bigint.NewInt(int64(2)))
	fmt.Println(e1.GetCoefficientsInt64())
	e2 := ring.NewGaussPoly(fv.N, fv.Q, *bigint.NewInt(int64(2)))

	c := new(Ciphertext)
	c.c0 = ring.NewRing(fv.N, fv.Q)
	c.c1 = ring.NewRing(fv.N, fv.Q)

	c.c0.MulPoly(fv.PublicKey.pub0, u)
	c.c0.Add(c.c0, e1)
	c.c0.Add(c.c0, newM)

	c.c1.MulPoly(fv.PublicKey.pub1, u)
	c.c1.Add(c.c1, e2)

	return c
}

func (fv *FV) Decrypt(c *Ciphertext) *Plaintext {
	m := new(Plaintext)
	m.Msg = ring.NewRing(fv.N, fv.Q)
	c.c1.MulPoly(c.c1, fv.PrivateKey.s)
	c.c0.Add(c.c0, c.c1)
	m.Msg.MulScalar(c.c0, fv.T)
	m.Msg.DivRound(m.Msg, fv.Q)
	m.Msg.Mod(m.Msg, fv.T)
	return m
}
