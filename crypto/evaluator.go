package crypto

import (
	"github.com/dedis/lago/ring"
	"github.com/dedis/lago/bigint"
	"math"
)

type Evaluator struct {
	ctx *FVContext	  // FV context
	evalkey *EvaluationKey
	evalsize uint32
}

// NewEvaluator creates a new evaluator for varies evaluation, e.g. add, sub, mul.
func NewEvaluator(ctx *FVContext, evalkey *EvaluationKey, evalsize uint32) *Evaluator {
	evaluator := new(Evaluator)
	evaluator.ctx = ctx
	evaluator.evalkey = evalkey
	evaluator.evalsize = evalsize
	return evaluator
}

func (evaluator *Evaluator)change2BigNttParams(r *ring.Ring) {
	r.Q = evaluator.ctx.BigQ
	r.Poly.SetNTTParams(evaluator.ctx.BigNttParams)
}

func (evaluator *Evaluator)change2NormalNttParams(r *ring.Ring) {
	r.Q = evaluator.ctx.Q
	r.Poly.SetNTTParams(evaluator.ctx.NttParams)
}

// Add conducts the homomorphic addition between ciphertexts c1 and c2
func (evaluator *Evaluator) Add(c1, c2 *Ciphertext) *Ciphertext {
	c := NewCiphertext(evaluator.ctx.N, evaluator.ctx.Q, evaluator.ctx.NttParams)
	c.value[0].Add(c1.value[0], c2.value[0])
	c.value[1].Add(c1.value[1], c2.value[1])
	return c
}

// Sub conducts the homomorphic subtraction between ciphertexts c1 and c2
func (evaluator *Evaluator) Sub(c1, c2 *Ciphertext) *Ciphertext {
	c := NewCiphertext(evaluator.ctx.N, evaluator.ctx.Q, evaluator.ctx.NttParams)
	c.value[0].Sub(c1.value[0], c2.value[0])
	c.value[1].Sub(c1.value[1], c2.value[1])
	return c
}

// Multiply conducts the homomorphic multiplication between ciphertexts c1 and c2
func (evaluator *Evaluator) Multiply(ct1, ct2 *Ciphertext) *Ciphertext {
	c0, err := ring.NewRing(evaluator.ctx.N, evaluator.ctx.BigQ, evaluator.ctx.BigNttParams)
	if err != nil {
		panic(err)
	}
	c1, err := ring.NewRing(evaluator.ctx.N, evaluator.ctx.BigQ, evaluator.ctx.BigNttParams)
	if err != nil {
		panic(err)
	}
	c2, err := ring.NewRing(evaluator.ctx.N, evaluator.ctx.BigQ, evaluator.ctx.BigNttParams)
	if err != nil {
		panic(err)
	}
	tmp, err := ring.NewRing(evaluator.ctx.N, evaluator.ctx.BigQ, evaluator.ctx.BigNttParams)
	if err != nil {
		panic(err)
	}

	// set c0, c1, c2
	ct1.value[0].Poly.InverseNTT()
	ct1.value[1].Poly.InverseNTT()
	ct2.value[0].Poly.InverseNTT()
	ct2.value[1].Poly.InverseNTT()

	evaluator.change2BigNttParams(ct1.value[0])
	evaluator.change2BigNttParams(ct1.value[1])
	evaluator.change2BigNttParams(ct2.value[0])
	evaluator.change2BigNttParams(ct2.value[1])

	c0.MulPoly(ct1.value[0], ct2.value[0])

	center(c0)
	c0.MulScalar(c0, evaluator.ctx.T)
	c0.DivRound(c0, evaluator.ctx.Q)
	c0.Mod(c0, evaluator.ctx.Q)
	evaluator.change2NormalNttParams(c0)
	c0.Poly.NTT()

	c1.MulPoly(ct1.value[0], ct2.value[1])
	tmp.MulPoly(ct1.value[1], ct2.value[0])
	c1.Add(c1, tmp)
	center(c1)
	c1.MulScalar(c1, evaluator.ctx.T)
	c1.DivRound(c1, evaluator.ctx.Q)
	c1.Mod(c1, evaluator.ctx.Q)
	evaluator.change2NormalNttParams(c1)
	c1.Poly.NTT()

	c2.MulPoly(ct1.value[1], ct2.value[1])
	center(c2)
	c2.MulScalar(c2, evaluator.ctx.T)
	c2.DivRound(c2, evaluator.ctx.Q)
	c2.Mod(c2, evaluator.ctx.Q)
	evaluator.change2NormalNttParams(c2)

	evaluator.change2NormalNttParams(ct1.value[0])
	evaluator.change2NormalNttParams(ct1.value[1])
	evaluator.change2NormalNttParams(ct2.value[0])
	evaluator.change2NormalNttParams(ct2.value[1])
	ct1.value[0].Poly.NTT()
	ct1.value[1].Poly.NTT()
	ct2.value[0].Poly.NTT()
	ct2.value[1].Poly.NTT()

	// relinearisation
	c2_i, err := ring.NewRing(evaluator.ctx.N, evaluator.ctx.Q, evaluator.ctx.NttParams)
	if err != nil {
		panic(err)
	}
	l := int(math.Floor(float64(evaluator.ctx.Q.Value.BitLen() - 1) / float64(evaluator.evalsize))) + 1
	mask := bigint.NewInt(1)
	mask.Lsh(mask, evaluator.evalsize)
	mask.Sub(mask, bigint.NewInt(1))
	for i := 0; i < l; i++ {
		c2_i.And(c2, *mask)
		c2_i.Poly.NTT()
		c2.Rsh(c2, evaluator.evalsize)

		tmp.MulCoeffs(c2_i, (*evaluator.evalkey)[i][0])
		c0.Add(c0, tmp)

		tmp.MulCoeffs(c2_i, (*evaluator.evalkey)[i][1])
		c1.Add(c1, tmp)
	}
	c0.Mod(c0, evaluator.ctx.Q)
	c1.Mod(c1, evaluator.ctx.Q)
	// construct result ciphertext
	newCiphertext := new(Ciphertext)
	newCiphertext.value[0] = c0
	newCiphertext.value[1] = c1
	return newCiphertext
}
