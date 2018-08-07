package main

import (
	"github.com/LoCCS/bliss/sampler"
	"github.com/LoCCS/bliss/poly"
	"fmt"
	"github.com/dedis/student_18_lattices/ring"
	"github.com/dedis/student_18_lattices/bigint"
)

func main() {
	seed := make([]uint8, sampler.SHA_512_DIGEST_LENGTH)
	for i := 0; i < len(seed); i++ {
		seed[i] = uint8(i % 8)
	}
	entropy, _ := sampler.NewEntropy(seed)
	mySampler, _ := sampler.New(0, entropy)
	gaussPoly := poly.GaussPoly(0, mySampler)
	fmt.Println(gaussPoly)

	myGaussPoly := ring.NewGaussPolyFromBLISS(256, *bigint.NewInt(7681))
	fmt.Println(myGaussPoly.GetCoefficientsInt64())

	myUniformPoly := ring.NewUniformPoly(256, *bigint.NewInt(7681))
	fmt.Println(myUniformPoly.GetCoefficientsInt64())
}
