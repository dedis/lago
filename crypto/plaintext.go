package crypto

import (
	"github.com/dedis/student_18_lattices/ring"
	"github.com/dedis/student_18_lattices/bigint"
	"github.com/dedis/student_18_lattices/polynomial"
)

type Plaintext struct {
	Value *ring.Ring
}

// NewPlaintext creates a Plaintext with given parameters.
func NewPlaintext(n uint32, q bigint.Int, nttParams *polynomial.NttParams) *Plaintext {
	plaintext := new(Plaintext)
	err := *new(error)
	plaintext.Value, err = ring.NewRing(n, q, nttParams)
	if err != nil {
		panic(err)
	}
	return plaintext
}
