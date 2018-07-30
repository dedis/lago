package crypto

import (
	"github.com/dedis/student_18_lattices/ring"
	"github.com/dedis/student_18_lattices/bigint"
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
	fv.PrivateKey.s = ring.NewUniformPoly(fv.N, *bigint.NewInt(int64(2)))
	fv.PrivateKey.s.Q = fv.Q
	fv.PublicKey.pub1 = ring.NewUniformPoly(fv.N, fv.Q)
	e := ring.NewGaussPolyFromBLISS(fv.N, fv.Q)
	fv.PublicKey.pub0 = ring.NewUniformPoly(fv.N, fv.Q)
	fv.PublicKey.pub0.MulPoly(fv.PublicKey.pub1, fv.PrivateKey.s)
	fv.PublicKey.pub0.Add(fv.PublicKey.pub0, e)
	fv.PublicKey.pub0.Neg(fv.PublicKey.pub0)
}

func (fv *FV) Encrypt(m *Plaintext) *Ciphertext {
	var delta bigint.Int
	newM := ring.NewRing(fv.N, fv.Q)
	delta.Div(&fv.Q, &fv.T)
	newM.MulScalar(m.Msg, delta)
	u := ring.NewUniformPoly(fv.N, *bigint.NewInt(int64(2)))
	u.Q = fv.Q
	e1 := ring.NewGaussPolyFromBLISS(fv.N, fv.Q)
	e2 := ring.NewGaussPolyFromBLISS(fv.N, fv.Q)
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
	m.Msg.MulPoly(c.c1, fv.PrivateKey.s)
	m.Msg.Add(m.Msg, c.c0)
	m.Msg.MulScalar(m.Msg, fv.T)
	m.Msg.DivRound(m.Msg, fv.Q)
	m.Msg.Mod(m.Msg, fv.T)
	return m
}
