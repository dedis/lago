package crypto

import (
	"github.com/dedis/student_18_lattices/ring"
	"github.com/dedis/student_18_lattices/bigint"
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
	c0, err := ring.NewRing(evaluator.ctx.N, evaluator.ctx.Q, evaluator.ctx.NttParams)
	if err != nil {
		panic(err)
	}
	c1, err := ring.NewRing(evaluator.ctx.N, evaluator.ctx.Q, evaluator.ctx.NttParams)
	if err != nil {
		panic(err)
	}
	c2, err := ring.NewRing(evaluator.ctx.N, evaluator.ctx.Q, evaluator.ctx.NttParams)
	if err != nil {
		panic(err)
	}
	tmp, err := ring.NewRing(evaluator.ctx.N, evaluator.ctx.Q, evaluator.ctx.NttParams)
	if err != nil {
		panic(err)
	}
	// set c0, c1, c2
	c0.MulCoeffs(ct1.value[0], ct2.value[0])
	c0.Poly.InverseNTT()
	c0.MulScalar(c0, evaluator.ctx.T)
	c0.DivRound(c0, evaluator.ctx.Q)
	c0.Mod(c0, evaluator.ctx.Q)

	c1.MulCoeffs(ct1.value[0], ct2.value[1])
	tmp.MulCoeffs(ct1.value[1], ct2.value[0])
	c1.Add(c1, tmp)
	c1.Poly.InverseNTT()
	c1.MulScalar(c1, evaluator.ctx.T)
	c1.DivRound(c1, evaluator.ctx.Q)
	c1.Mod(c1, evaluator.ctx.Q)

	c2.MulCoeffs(ct1.value[1], ct2.value[1])
	c2.Poly.InverseNTT()
	c2.MulScalar(c2, evaluator.ctx.T)
	c2.DivRound(c2, evaluator.ctx.Q)
	c2.Mod(c2, evaluator.ctx.Q)

	// relinearisation
	c2_i, err := ring.NewRing(evaluator.ctx.N, evaluator.ctx.Q, evaluator.ctx.NttParams)
	if err != nil {
		panic(err)
	}
	l := evaluator.ctx.Q.Value.BitLen()
	mask := bigint.NewInt(1)
	mask.Lsh(mask, evaluator.evalsize)
	for i := 0; i < l; i++ {
		c2_i.And(c2, *mask)
		c2.Rsh(c2, evaluator.evalsize)
		c2_i.MulCoeffs(c2_i, evaluator.evalkey[i][0])
		c0.Add(c0, evaluator.evalkey[i][0])
	}
}
