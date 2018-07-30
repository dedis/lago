package crypto

import (
	"github.com/dedis/student_18_lattices/ring"
	"github.com/dedis/student_18_lattices/bigint"
)

type Ciphertext struct {
	c0 *ring.Ring
	c1 *ring.Ring
}

func (c *Ciphertext) GetCiphertext() ([]bigint.Int, []bigint.Int){
	return c.c0.GetCoefficients(), c.c1.GetCoefficients()
}
