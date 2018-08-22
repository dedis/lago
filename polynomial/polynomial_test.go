package polynomial

import (
	"testing"
	"io/ioutil"
	"fmt"
	"strings"
	"strconv"
	"github.com/dedis/student_18_lattices/bigint"
)

func TestPolynomial(t *testing.T) {
	for i := 0; i <=2; i++ {
		testfile, err := ioutil.ReadFile(fmt.Sprintf("test_data/testvector_polynomial_%d", i))
		if err != nil {
			t.Errorf("Failed to open file: %s", err.Error())
		}
		filecontent := strings.TrimSpace(string(testfile))
		vs := strings.Split(filecontent, "\n")
		if len(vs) != 14 {
			t.Errorf("Error in data read from test_data: len(vs) = %d", len(vs))
		}

		// load q
		q, err := strconv.Atoi(vs[0])
		if err != nil {
			t.Errorf("Invalid integer: %v", vs[0])
		}

		// load n
		n, err := strconv.Atoi(vs[1])
		if err != nil {
			t.Errorf("Invalid integer: %v", vs[1])
		}

		// load first polynomial
		p1String := strings.Split(strings.TrimSpace(vs[2]), ", ")
		p1Coeffs := make([]bigint.Int, n)
		for i := range p1Coeffs {
			tmp, err := strconv.Atoi(p1String[i])
			if err != nil {
				t.Errorf("Invalid integer of p1 coeffs: %v", p1String[i])
			}
			p1Coeffs[i].SetInt(int64(tmp))
		}
		p1, _ := NewPolynomial(uint32(n), *bigint.NewInt(int64(q)))
		p1.SetCoefficients(p1Coeffs)

		// load second polynomial
		p2String := strings.Split(strings.TrimSpace(vs[3]), ", ")
		p2Coeffs := make([]bigint.Int, n)
		for i := range p2Coeffs {
			tmp, err := strconv.Atoi(p2String[i])
			if err != nil {
				t.Errorf("Invalid integer of p2 coeffs: %v", p2String[i])
			}
			p2Coeffs[i].SetInt(int64(tmp))
		}
		p2, _ := NewPolynomial(uint32(n), *bigint.NewInt(int64(q)))
		p2.SetCoefficients(p2Coeffs)

		// load add coefficients
		addString := strings.Split(strings.TrimSpace(vs[4]), ", ")
		addCoeffs := make([]bigint.Int, n)
		for i := range addCoeffs {
			tmp, err := strconv.Atoi(addString[i])
			if err != nil {
				t.Errorf("Invalid integer of add coeffs: %v", addString[i])
			}
			addCoeffs[i].SetInt(int64(tmp))
		}

		// load sub coefficients
		subString := strings.Split(strings.TrimSpace(vs[5]), ", ")
		subCoeffs := make([]bigint.Int, n)
		for i := range subCoeffs {
			tmp, err := strconv.Atoi(subString[i])
			if err != nil {
				t.Errorf("Invalid integer of sub coeffs: %v", subString[i])
			}
			subCoeffs[i].SetInt(int64(tmp))
		}

		// load neg coefficients
		negString := strings.Split(strings.TrimSpace(vs[6]), ", ")
		negCoeffs := make([]bigint.Int, n)
		for i := range negCoeffs {
			tmp, err := strconv.Atoi(negString[i])
			if err != nil {
				t.Errorf("Invalid integer of neg coeffs: %v", negString[i])
			}
			negCoeffs[i].SetInt(int64(tmp))
		}

		// load mulCoeffs coefficients
		mulCoeffsString := strings.Split(strings.TrimSpace(vs[7]), ", ")
		mulCoeffsCoeffs := make([]bigint.Int, n)
		for i := range mulCoeffsCoeffs {
			tmp, err := strconv.Atoi(mulCoeffsString[i])
			if err != nil {
				t.Errorf("Invalid integer of mulCoeffs coeffs: %v", mulCoeffsString[i])
			}
			mulCoeffsCoeffs[i].SetInt(int64(tmp))
		}

		// load mulScalar coefficients
		scalar, err := strconv.Atoi(vs[8])
		if err != nil {
			t.Errorf("Invalid integer of scalar: %v", vs[8])
		}
		mulScalarString := strings.Split(strings.TrimSpace(vs[9]), ", ")
		mulScalarCoeffs := make([]bigint.Int, n)
		for i := range mulScalarCoeffs {
			tmp, err := strconv.Atoi(mulScalarString[i])
			if err != nil {
				t.Errorf("Invalid integer of mulScalar coeffs: %v", mulScalarString[i])
			}
			mulScalarCoeffs[i].SetInt(int64(tmp))
		}

		// load mulPoly coefficients
		mulPolyString := strings.Split(strings.TrimSpace(vs[10]), ", ")
		mulPolyCoeffs := make([]bigint.Int, n)
		for i := range mulPolyCoeffs {
			tmp, err := strconv.Atoi(mulPolyString[i])
			if err != nil {
				t.Errorf("Invalid integer of mulPoly coeffs: %v", mulScalarString[i])
			}
			mulPolyCoeffs[i].SetInt(int64(tmp))
		}

		// load divisor for div and divRound
		divisor, err := strconv.Atoi(vs[11])
		if err != nil {
			t.Errorf("Invalid integer of divisor: %v", vs[11])
		}

		// load div coefficients
		divString := strings.Split(strings.TrimSpace(vs[12]), ", ")
		divCoeffs := make([]bigint.Int, n)
		for i := range divCoeffs {
			tmp, err := strconv.Atoi(divString[i])
			if err != nil {
				t.Errorf("Invalid integer of div coeffs: %v", divString[i])
			}
			divCoeffs[i].SetInt(int64(tmp))
		}

		// load divRound coefficients
		divRoundString := strings.Split(strings.TrimSpace(vs[13]), ", ")
		divRoundCoeffs := make([]bigint.Int, n)
		for i := range divRoundCoeffs {
			tmp, err := strconv.Atoi(divRoundString[i])
			if err != nil {
				t.Errorf("Invalid integer of divRound coeffs: %v", divRoundString[i])
			}
			divRoundCoeffs[i].SetInt(int64(tmp))
		}

		pTest, _ := NewPolynomial(uint32(n), *bigint.NewInt(int64(q)))
		// Test add
		pTest.AddMod(p1, p2)
		for i := range pTest.coeffs {
			if !pTest.coeffs[i].EqualTo(&addCoeffs[i]) {
				t.Errorf("Error in add coeffs: index %v, value %v", i, pTest.coeffs[i])
			}
		}
		// Test sub
		pTest.SubMod(p1, p2)
		for i := range pTest.coeffs {
			if !pTest.coeffs[i].EqualTo(&subCoeffs[i]) {
				t.Errorf("Error in sub coeffs: index %v, value %v", i, pTest.coeffs[i])
			}
		}
		// Test neg
		pTest.Neg(p1)
		for i := range pTest.coeffs {
			if !pTest.coeffs[i].EqualTo(&negCoeffs[i]) {
				t.Errorf("Error in neg coeffs: index %v, value %v", i, pTest.coeffs[i])
			}
		}
		// Test mulCoeffs
		pTest.MulCoeffs(p1, p2)
		for i := range pTest.coeffs {
			if !pTest.coeffs[i].EqualTo(&mulCoeffsCoeffs[i]) {
				t.Errorf("Error in mulCoeffs coeffs: index %v, value %v", i, pTest.coeffs[i])
			}
		}
		// Test mulScalar
		pTest.MulScalar(p1, *bigint.NewInt(int64(scalar)))
		for i := range pTest.coeffs {
			if !pTest.coeffs[i].EqualTo(&mulScalarCoeffs[i]) {
				t.Errorf("Error in mulScalar coeffs: index %v, value %v", i, pTest.coeffs[i])
			}
		}
		// Test mulPoly
		pTest.MulPoly(p1, p2)
		for i := range pTest.coeffs {
			if !pTest.coeffs[i].EqualTo(&mulPolyCoeffs[i]) {
				t.Errorf("Error in mulPoly coeffs: index %v, value %v", i, pTest.coeffs[i])
			}
		}
		// Test div
		pTest.Div(p1, *bigint.NewInt(int64(divisor)))
		for i := range pTest.coeffs {
			if !pTest.coeffs[i].EqualTo(&divCoeffs[i]) {
				t.Errorf("Error in div coeffs: index %v, value %v", i, pTest.coeffs[i])
			}
		}
		// Test divRound
		pTest.DivRound(p1, *bigint.NewInt(int64(divisor)))
		for i := range pTest.coeffs {
			if !pTest.coeffs[i].EqualTo(&divRoundCoeffs[i]) {
				t.Errorf("Error in divRound coeffs: DivRound(%v, %v), expected %v, got, %v", p1Coeffs[i].Int64(), divisor, divRoundCoeffs[i], pTest.coeffs[i])
			}
		}
	}
}