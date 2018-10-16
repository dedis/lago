package crypto

import (
	"github.com/dedis/lago/ring"
	"github.com/dedis/lago/bigint"
	"github.com/dedis/lago/polynomial"
)

type Ciphertext struct {
	value [2]*ring.Ring
}

// NewCiphertext creates a new ciphertext
func NewCiphertext(n uint32, q bigint.Int, nttParams *polynomial.NttParams) *Ciphertext {
	ciphertext := new(Ciphertext)
	err := *new(error)
	ciphertext.value[0], err = ring.NewRing(n, q, nttParams)
	if err != nil {
		panic(err)
	}
	ciphertext.value[1], err = ring.NewRing(n, q, nttParams)
	if err != nil {
		panic(err)
	}
	return ciphertext
}
