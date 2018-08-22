package polynomial

import (
	"testing"
	"github.com/dedis/student_18_lattices/bigint"
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
)

type argFactors struct {
	q bigint.Int
	factors []bigint.Int
}
var factorVec = []argFactors {
	{*bigint.NewInt(7680), []bigint.Int{*bigint.NewInt(2), *bigint.NewInt(3), *bigint.NewInt(5)}},
	{*bigint.NewInt(1152921504382476288), []bigint.Int{*bigint.NewInt(2), *bigint.NewInt(3), *bigint.NewInt(131), *bigint.NewInt(358110657323)}},
}

func TestGetFactors(t *testing.T) {
	for i, testPair := range factorVec {
		factors := getFactors(&testPair.q)
		for j := range factors {
			if ! factors[j].EqualTo(&testPair.factors[j]) {
				t.Errorf("factor not match in test pair %v", i)
			}
		}
	}
}

type argRoots struct {
	q, root bigint.Int
}
var rootsVec = []argRoots {
	{*bigint.NewInt(7681), *bigint.NewInt(17)},
	{*bigint.NewInt(1152921504382476289), *bigint.NewInt(11)},
}
func TestPrimitiveRoot(t *testing.T) {
	for i, testPair := range rootsVec {
		root := primitiveRoot(&testPair.q)
		if !root.EqualTo(&testPair.root) {
			t.Errorf("primitive root not match %v", i)
		}
	}
}

func TestGenerateNTTParameters(t *testing.T) {
	for i := 0; i <=0; i++ {
		testfile, err := ioutil.ReadFile(fmt.Sprintf("test_data/testvector_params_%d", i))
		if err != nil {
			t.Errorf("Failed to open file: %s", err.Error())
		}
		filecontent := strings.TrimSpace(string(testfile))
		vs := strings.Split(filecontent, "\n")
		if len(vs) != 4 {
			t.Errorf("Error in data read from test_data: len(vs) = %d", len(vs))
		}
		q, err := strconv.Atoi(vs[0])
		if err != nil {
			t.Errorf("Invalid integer: %v", vs[0])
		}
		n, err := strconv.Atoi(vs[1])
		if err != nil {
			t.Errorf("Invalid integer: %v", vs[1])
		}

		params, _ := GenerateNTTParameters(uint32(n), *bigint.NewInt(int64(q)))

		psiReverseString := strings.Split(strings.TrimSpace(vs[2]), ", ")
		psiReverse := make([]bigint.Int, n)
		for i := range psiReverse {
			tmp, err := strconv.Atoi(psiReverseString[i])
			if err != nil {
				t.Errorf("Invalid integer: %v", psiReverseString[i])
			}
			psiReverse[i].SetInt(int64(tmp))
			if !psiReverse[i].EqualTo(&params.PsiReverse[i]) {
				t.Errorf("psi unmatch : %v", i)
			}
		}
		psiInvReverseString := strings.Split(strings.TrimSpace(vs[3]), ", ")
		psiInvReverse := make([]bigint.Int, n)
		for i := range psiInvReverse {
			tmp, err := strconv.Atoi(psiInvReverseString[i])
			if err != nil {
				t.Errorf("Invalid integer: %v", psiInvReverseString[i])
			}
			psiInvReverse[i].SetInt(int64(tmp))
			if !psiInvReverse[i].EqualTo(&params.PsiInvReverse[i]) {
				t.Errorf("psiInv unmatch : %v", i)
			}
		}
	}
}
