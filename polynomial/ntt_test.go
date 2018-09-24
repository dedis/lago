package polynomial

import (
	"testing"
	"github.com/dedis/student_18_lattices/bigint"
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
	"math/rand"
	"github.com/dedis/student_18_lattices/kyber"
	"github.com/LoCCS/bliss/poly"
)

// Test the correctness of NTT and InverseNTT functions with different params from test_data/testvector_ntt_i
func TestNTT(t *testing.T) {
	for i := 0; i <=2; i++ {
		testfile, err := ioutil.ReadFile(fmt.Sprintf("test_data/testvector_ntt_%d", i))
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
		coeffsString := strings.Split(strings.TrimSpace(vs[2]), " ")
		coeffs := make([]bigint.Int, n)
		for i := range coeffs {
			tmp, err := strconv.Atoi(coeffsString[i])
			if err != nil {
				t.Errorf("Invalid integer: %v", coeffsString[i])
			}
			coeffs[i].SetInt(int64(tmp))
		}
		nttCoeffsString := strings.Split(strings.TrimSpace(vs[3]), " ")
		nttCoeffs := make([]bigint.Int, n)
		for i := range nttCoeffs {
			tmp, err := strconv.Atoi(nttCoeffsString[i])
			if err != nil {
				t.Errorf("Invalid integer: %v", nttCoeffsString[i])
			}
			nttCoeffs[i].SetInt(int64(tmp))
		}
		nttParams := GenerateNTTParams(uint32(n), *bigint.NewInt(int64(q)))
		p, err := NewPolynomial(uint32(n), *bigint.NewInt(int64(q)), nttParams)
		if err != nil {
			t.Error("Error in creating new polynomial")
		}
		p.SetCoefficients(coeffs)

		// Test the correctness of NTT
		p.NTT()
		pNttCoeffs := p.GetCoefficients()
		for i := range pNttCoeffs {
			if !pNttCoeffs[i].EqualTo(&nttCoeffs[i]) {
				continue
				t.Errorf("NTT Error in ntt coeffs: %v, %v", nttCoeffs[i].Int64(), pNttCoeffs[i])
			}
		}

		//Test the correctness of inverse NTT
		p.InverseNTT()
		pInverseNttCoeffs := p.GetCoefficients()
		for i := range pInverseNttCoeffs {
			if !pInverseNttCoeffs[i].EqualTo(coeffs[i].Mod(&coeffs[i], &p.q)) {
				t.Errorf("Inverse NTT Error in ntt coeffs: %v", coeffs[i].Int64())
			}
		}
	}
}

// Cross-verify the correctness of NTT with kyber.NTT and loccs.NTT
func TestNTTCross(t *testing.T) {
	q := bigint.NewInt(7681)
	n := uint32(256)
	// Two set of coefficients (polynomials)
	coeffs1 := make([]bigint.Int, n)
	for i := range coeffs1 {
		coeffs1[i].SetInt(rand.Int63n(10000))
		coeffs1[i].Mod(&coeffs1[i], q)
	}
	coeffs2 := make([]bigint.Int, n)
	for i := range coeffs2 {
		coeffs2[i].SetInt(rand.Int63n(20000))
		coeffs2[i].Mod(&coeffs2[i], q)
	}

	// kyber.NTT
	coeffs1Kyber := [256]uint16{}
	coeffs2Kyber:= [256]uint16{}
	for i := range coeffs1Kyber {
		coeffs1Kyber[i] = uint16(coeffs1[i].Int64())
		coeffs2Kyber[i] = uint16(coeffs2[i].Int64())
	}

	kyber.NttRef(&coeffs1Kyber)
	kyber.NttRef(&coeffs2Kyber)
	for i := range coeffs1 {
		coeffs1Kyber[i] = uint16(uint32(coeffs1Kyber[i]) * uint32(coeffs2Kyber[i]) % uint32(q.Int64()))
	}
	kyber.InvnttRef(&coeffs1Kyber)
	for i := range coeffs1 {
		coeffs1Kyber[i] = uint16(uint32(coeffs1Kyber[i]) % uint32(q.Int64()))
	}
	nttResultKyber := coeffs1Kyber

	//bliss.NTT
	coeffs1Bliss := make([]int32, n)
	coeffs2Bliss := make([]int32, n)
	for i := range coeffs1Bliss {
		coeffs1Bliss[i] = int32(coeffs1[i].Int64())
		coeffs2Bliss[i] = int32(coeffs2[i].Int64())
	}
	p1Bliss, _ := poly.New(0)
	p2Bliss, _ := poly.New(0)
	p1Bliss.SetData(coeffs1Bliss)
	p2Bliss.SetData(coeffs2Bliss)

	p1BlissNTT, _ := p1Bliss.NTT()
	nttResultBlissPoly, _ := p2Bliss.MultiplyNTT(p1BlissNTT)
	nttResultBliss := nttResultBlissPoly.GetData()

	// this.NTT
	nttParams := GenerateNTTParams(uint32(n), *q)
	p1, _ := NewPolynomial(n, *q, nttParams)
	p1.SetCoefficients(coeffs1)
	p2, _ := NewPolynomial(n, *q, nttParams)
	p2.SetCoefficients(coeffs2)
	p1.MulPoly(p1, p2)
	nttResultThis := p1.GetCoefficientsInt64()

	// verify if the outputs of three methods are the same
	for i := range nttResultThis {
		if nttResultKyber[i] != uint16(nttResultBliss[i]) {
			t.Errorf("Unmatch between Kyber and Bliss in coeffs: %v (Kyber), %v (Bliss)", nttResultKyber[i], nttResultBliss[i])
		}
		if nttResultKyber[i] != uint16(nttResultThis[i]) {
			t.Errorf("Unmatch between Kyber and This in coeffs: %v (Kyber), %v (This)", nttResultKyber[i], nttResultThis[i])
		}
		if nttResultBliss[i] != int32(nttResultThis[i]) {
			t.Errorf("Unmatch between Bliss and This in coeffs: %v (Bliss), %v (This)", nttResultBliss[i], nttResultThis[i])
		}
	}
}

// Benchmark this NTT
func BenchmarkNTT(b *testing.B) {
	testfile, err := ioutil.ReadFile(fmt.Sprint("test_data/testvector_ntt_0"))
	if err != nil {
		b.Errorf("Failed to open file: %s", err.Error())
	}
	filecontent := strings.TrimSpace(string(testfile))
	vs := strings.Split(filecontent, "\n")
	if len(vs) != 4 {
		b.Errorf("Error in data read from test_data: len(vs) = %d", len(vs))
	}
	q, err := strconv.Atoi(vs[0])
	if err != nil {
		b.Errorf("Invalid integer: %v", vs[0])
	}
	n, err := strconv.Atoi(vs[1])
	if err != nil {
		b.Errorf("Invalid integer: %v", vs[1])
	}
	coeffsString := strings.Split(strings.TrimSpace(vs[2]), " ")
	coeffs := make([]bigint.Int, n)
	for i := range coeffs {
		tmp, err := strconv.Atoi(coeffsString[i])
		if err != nil {
			b.Errorf("Invalid integer: %v", coeffsString[i])
		}
		coeffs[i].SetInt(int64(tmp))
	}
	nttParams := GenerateNTTParams(uint32(n), *bigint.NewInt(int64(q)))
	p, err := NewPolynomial(uint32(n), *bigint.NewInt(int64(q)), nttParams)
	if err != nil {
		b.Error("Error in creating new polynomial")
	}
	p.SetCoefficients(coeffs)
	b.ResetTimer()
	for i :=0; i < b.N; i++ {
		p.NTT()
	}
}

// Benchmark this NTT
func BenchmarkFastNTT(b *testing.B) {
	testfile, err := ioutil.ReadFile(fmt.Sprint("test_data/testvector_ntt_0"))
	if err != nil {
		b.Errorf("Failed to open file: %s", err.Error())
	}
	filecontent := strings.TrimSpace(string(testfile))
	vs := strings.Split(filecontent, "\n")
	if len(vs) != 4 {
		b.Errorf("Error in data read from test_data: len(vs) = %d", len(vs))
	}
	q, err := strconv.Atoi(vs[0])
	if err != nil {
		b.Errorf("Invalid integer: %v", vs[0])
	}
	n, err := strconv.Atoi(vs[1])
	if err != nil {
		b.Errorf("Invalid integer: %v", vs[1])
	}
	coeffsString := strings.Split(strings.TrimSpace(vs[2]), " ")
	coeffs := make([]bigint.Int, n)
	for i := range coeffs {
		tmp, err := strconv.Atoi(coeffsString[i])
		if err != nil {
			b.Errorf("Invalid integer: %v", coeffsString[i])
		}
		coeffs[i].SetInt(int64(tmp))
	}
	nttParams := GenerateNTTParams(uint32(n), *bigint.NewInt(int64(q)))
	p, err := NewPolynomial(uint32(n), *bigint.NewInt(int64(q)), nttParams)
	if err != nil {
		b.Error("Error in creating new polynomial")
	}
	p.SetCoefficients(coeffs)
	b.ResetTimer()
	for i :=0; i < b.N; i++ {
		p.NTTFast()
	}
}

// Benchmark this DebugNTT
func BenchmarkDebugNTT(b *testing.B) {
	testfile, err := ioutil.ReadFile(fmt.Sprint("test_data/testvector_ntt_0"))
	if err != nil {
		b.Errorf("Failed to open file: %s", err.Error())
	}
	filecontent := strings.TrimSpace(string(testfile))
	vs := strings.Split(filecontent, "\n")
	if len(vs) != 4 {
		b.Errorf("Error in data read from test_data: len(vs) = %d", len(vs))
	}
	q, err := strconv.Atoi(vs[0])
	if err != nil {
		b.Errorf("Invalid integer: %v", vs[0])
	}
	n, err := strconv.Atoi(vs[1])
	if err != nil {
		b.Errorf("Invalid integer: %v", vs[1])
	}
	coeffsString := strings.Split(strings.TrimSpace(vs[2]), " ")
	coeffs := make([]int64, 256)
	for i := range coeffs {
		tmp, err := strconv.Atoi(coeffsString[i])
		if err != nil {
			b.Errorf("Invalid integer: %v", coeffsString[i])
		}
		coeffs[i] = int64(tmp)
	}
	nttParams := GenerateNTTParams(uint32(n), *bigint.NewInt(int64(q)))
	p, err := NewPolynomial(uint32(n), *bigint.NewInt(int64(q)), nttParams)
	if err != nil {
		b.Error("Error in creating new polynomial")
	}
	psiReverse := make([]int64, n)
	for i := range psiReverse {
		psiReverse[i] = p.nttParams.PsiReverse[i].Int64()
	}
	b.ResetTimer()
	for i :=0; i < b.N; i++ {
		NTTFastInt64(coeffs, psiReverse, int64(q), int64(n))
	}
}

// Benchmark this DebugNTT2
func BenchmarkDebugNTT2(b *testing.B) {
	testfile, err := ioutil.ReadFile(fmt.Sprint("test_data/testvector_ntt_0"))
	if err != nil {
		b.Errorf("Failed to open file: %s", err.Error())
	}
	filecontent := strings.TrimSpace(string(testfile))
	vs := strings.Split(filecontent, "\n")
	if len(vs) != 4 {
		b.Errorf("Error in data read from test_data: len(vs) = %d", len(vs))
	}
	q, err := strconv.Atoi(vs[0])
	if err != nil {
		b.Errorf("Invalid integer: %v", vs[0])
	}
	n, err := strconv.Atoi(vs[1])
	if err != nil {
		b.Errorf("Invalid integer: %v", vs[1])
	}
	coeffsString := strings.Split(strings.TrimSpace(vs[2]), " ")
	coeffs := make([]int64, 256)
	for i := range coeffs {
		tmp, err := strconv.Atoi(coeffsString[i])
		if err != nil {
			b.Errorf("Invalid integer: %v", coeffsString[i])
		}
		coeffs[i] = int64(tmp)
	}
	nttParams := GenerateNTTParams(uint32(n), *bigint.NewInt(int64(q)))
	p, err := NewPolynomial(uint32(n), *bigint.NewInt(int64(q)), nttParams)
	if err != nil {
		b.Error("Error in creating new polynomial")
	}
	psiReverse := make([]int64, n)
	for i := range psiReverse {
		psiReverse[i] = p.nttParams.PsiReverse[i].Int64()
	}
	b.ResetTimer()
	for i :=0; i < b.N; i++ {
		NTTInt64(coeffs, psiReverse, int64(q), int64(n))
	}
}

// Benchmark the Kyber NTT
func BenchmarkKyberNTT(b *testing.B) {
	testfile, err := ioutil.ReadFile(fmt.Sprint("test_data/testvector_ntt_0"))
	if err != nil {
		b.Errorf("Failed to open file: %s", err.Error())
	}
	filecontent := strings.TrimSpace(string(testfile))
	vs := strings.Split(filecontent, "\n")
	if len(vs) != 4 {
		b.Errorf("Error in data read from test_data: len(vs) = %d", len(vs))
	}
	coeffsString := strings.Split(strings.TrimSpace(vs[2]), " ")
	coeffs := [256]uint16{}
	for i := range coeffs {
		tmp, err := strconv.Atoi(coeffsString[i])
		if err != nil {
			b.Errorf("Invalid integer: %v", coeffsString[i])
		}
		coeffs[i] = uint16(tmp)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kyber.NttRef(&coeffs)
	}
}

// Benchmark the Bliss NTT
func BenchmarkBlissNTT(b *testing.B) {
	testfile, err := ioutil.ReadFile(fmt.Sprint("test_data/testvector_ntt_0"))
	if err != nil {
		b.Errorf("Failed to open file: %s", err.Error())
	}
	filecontent := strings.TrimSpace(string(testfile))
	vs := strings.Split(filecontent, "\n")
	if len(vs) != 4 {
		b.Errorf("Error in data read from test_data: len(vs) = %d", len(vs))
	}
	coeffsString := strings.Split(strings.TrimSpace(vs[2]), " ")
	coeffs := make([]int32, 256)
	for i := range coeffs {
		tmp, err := strconv.Atoi(coeffsString[i])
		if err != nil {
			b.Errorf("Invalid integer: %v", coeffsString[i])
		}
		coeffs[i] = int32(tmp)
	}
	p, _ := poly.New(0)
	p.SetData(coeffs)

	for i :=0; i < b.N; i++ {
		p.NTT()
	}
}
