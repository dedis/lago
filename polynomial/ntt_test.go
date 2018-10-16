package polynomial

import (
	"testing"
	"github.com/dedis/lago/bigint"
	"fmt"
	"io/ioutil"
	"strings"
	"strconv"
	"math/rand"
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

	// kyber.NTT, from https://github.com/Yawning/kyber
	coeffs1Kyber := [256]uint16{}
	coeffs2Kyber:= [256]uint16{}
	for i := range coeffs1Kyber {
		coeffs1Kyber[i] = uint16(coeffs1[i].Int64())
		coeffs2Kyber[i] = uint16(coeffs2[i].Int64())
	}

	NttRef(&coeffs1Kyber)
	NttRef(&coeffs2Kyber)
	for i := range coeffs1 {
		coeffs1Kyber[i] = uint16(uint32(coeffs1Kyber[i]) * uint32(coeffs2Kyber[i]) % uint32(q.Int64()))
	}
	InvnttRef(&coeffs1Kyber)
	for i := range coeffs1 {
		coeffs1Kyber[i] = uint16(uint32(coeffs1Kyber[i]) % uint32(q.Int64()))
	}
	nttResultKyber := coeffs1Kyber

	//bliss.NTT, from https://github.com/LoCCS/bliss
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

// Benchmark NTTFat
func BenchmarkNTTFast(b *testing.B) {
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

// Benchmark NTTFastInt64
func BenchmarkNTTFastInt64(b *testing.B) {
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

// Benchmark NTTInt64
func BenchmarkNTTInt64(b *testing.B) {
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
		NttRef(&coeffs)
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

// The following codes implement the kyber.Ntt and kyber.Invntt.
// for more details, please visit https://github.com/Yawning/kyber

func NttRef(p *[256]uint16) {
	var j int
	k := 1
	for level := 7; level >= 0; level-- {
		distance := 1 << uint(level)
		for start := 0; start < 256; start = j + distance {
			zeta := zetas[k]
			k++
			for j = start; j < start+distance; j++ {
				t := uint16(montgomeryReduce(int64(uint32(zeta) * uint32(p[j+distance]))))
				p[j+distance] = uint16(barrettReduce(int64(p[j] + 4*7681 - t)))

				if level&1 == 1 { // odd level
					p[j] = p[j] + t // Omit reduction (be lazy)
				} else {
					p[j] = uint16(barrettReduce(int64(p[j] + t)))
				}
			}
		}
	}
}

func InvnttRef(a *[256]uint16) {
	for level := 0; level < 8; level++ {
		distance := 1 << uint(level)
		for start := 0; start < distance; start++ {
			var jTwiddle int
			for j := start; j < 256-1; j += 2 * distance {
				w := uint32(omegasInvBitrevMontgomery[jTwiddle])
				jTwiddle++

				temp := a[j]

				if level&1 == 1 { // odd level
					a[j] = uint16(barrettReduce(int64(temp + a[j+distance])))
				} else {
					a[j] = temp + a[j+distance] // Omit reduction (be lazy)
				}

				t := w * (uint32(temp) + 4*7681 - uint32(a[j+distance]))

				a[j+distance] = uint16(montgomeryReduce(int64(t)))
			}
		}
	}

	for i, v := range psisInvMontgomery {
		a[i] = uint16(montgomeryReduce(int64(uint32(a[i]) * uint32(v))))
	}
}

var zetas = [256]uint16{
	990, 7427, 2634, 6819, 578, 3281, 2143, 1095, 484, 6362, 3336, 5382, 6086, 3823, 877, 5656,
	3583, 7010, 6414, 263, 1285, 291, 7143, 7338, 1581, 5134, 5184, 5932, 4042, 5775, 2468, 3,
	606, 729, 5383, 962, 3240, 7548, 5129, 7653, 5929, 4965, 2461, 641, 1584, 2666, 1142, 157,
	7407, 5222, 5602, 5142, 6140, 5485, 4931, 1559, 2085, 5284, 2056, 3538, 7269, 3535, 7190, 1957,
	3465, 6792, 1538, 4664, 2023, 7643, 3660, 7673, 1694, 6905, 3995, 3475, 5939, 1859, 6910, 4434,
	1019, 1492, 7087, 4761, 657, 4859, 5798, 2640, 1693, 2607, 2782, 5400, 6466, 1010, 957, 3851,
	2121, 6392, 7319, 3367, 3659, 3375, 6430, 7583, 1549, 5856, 4773, 6084, 5544, 1650, 3997, 4390,
	6722, 2915, 4245, 2635, 6128, 7676, 5737, 1616, 3457, 3132, 7196, 4702, 6239, 851, 2122, 3009,
	7613, 7295, 2007, 323, 5112, 3716, 2289, 6442, 6965, 2713, 7126, 3401, 963, 6596, 607, 5027,
	7078, 4484, 5937, 944, 2860, 2680, 5049, 1777, 5850, 3387, 6487, 6777, 4812, 4724, 7077, 186,
	6848, 6793, 3463, 5877, 1174, 7116, 3077, 5945, 6591, 590, 6643, 1337, 6036, 3991, 1675, 2053,
	6055, 1162, 1679, 3883, 4311, 2106, 6163, 4486, 6374, 5006, 4576, 4288, 5180, 4102, 282, 6119,
	7443, 6330, 3184, 4971, 2530, 5325, 4171, 7185, 5175, 5655, 1898, 382, 7211, 43, 5965, 6073,
	1730, 332, 1577, 3304, 2329, 1699, 6150, 2379, 5113, 333, 3502, 4517, 1480, 1172, 5567, 651,
	925, 4573, 599, 1367, 4109, 1863, 6929, 1605, 3866, 2065, 4048, 839, 5764, 2447, 2022, 3345,
	1990, 4067, 2036, 2069, 3567, 7371, 2368, 339, 6947, 2159, 654, 7327, 2768, 6676, 987, 2214,
}

var omegasInvBitrevMontgomery = [256 / 2]uint16{
	990, 254, 862, 5047, 6586, 5538, 4400, 7103, 2025, 6804, 3858, 1595, 2299, 4345, 1319, 7197,
	7678, 5213, 1906, 3639, 1749, 2497, 2547, 6100, 343, 538, 7390, 6396, 7418, 1267, 671, 4098,
	5724, 491, 4146, 412, 4143, 5625, 2397, 5596, 6122, 2750, 2196, 1541, 2539, 2079, 2459, 274,
	7524, 6539, 5015, 6097, 7040, 5220, 2716, 1752, 28, 2552, 133, 4441, 6719, 2298, 6952, 7075,
	4672, 5559, 6830, 1442, 2979, 485, 4549, 4224, 6065, 1944, 5, 1553, 5046, 3436, 4766, 959,
	3291, 3684, 6031, 2137, 1597, 2908, 1825, 6132, 98, 1251, 4306, 4022, 4314, 362, 1289, 5560,
	3830, 6724, 6671, 1215, 2281, 4899, 5074, 5988, 5041, 1883, 2822, 7024, 2920, 594, 6189, 6662,
	3247, 771, 5822, 1742, 4206, 3686, 776, 5987, 8, 4021, 38, 5658, 3017, 6143, 889, 4216,
}

var psisInvMontgomery = [256]uint16{
	1024, 4972, 5779, 6907, 4943, 4168, 315, 5580, 90, 497, 1123, 142, 4710, 5527, 2443, 4871,
	698, 2489, 2394, 4003, 684, 2241, 2390, 7224, 5072, 2064, 4741, 1687, 6841, 482, 7441, 1235,
	2126, 4742, 2802, 5744, 6287, 4933, 699, 3604, 1297, 2127, 5857, 1705, 3868, 3779, 4397, 2177,
	159, 622, 2240, 1275, 640, 6948, 4572, 5277, 209, 2605, 1157, 7328, 5817, 3191, 1662, 2009,
	4864, 574, 2487, 164, 6197, 4436, 7257, 3462, 4268, 4281, 3414, 4515, 3170, 1290, 2003, 5855,
	7156, 6062, 7531, 1732, 3249, 4884, 7512, 3590, 1049, 2123, 1397, 6093, 3691, 6130, 6541, 3946,
	6258, 3322, 1788, 4241, 4900, 2309, 1400, 1757, 400, 502, 6698, 2338, 3011, 668, 7444, 4580,
	6516, 6795, 2959, 4136, 3040, 2279, 6355, 3943, 2913, 6613, 7416, 4084, 6508, 5556, 4054, 3782,
	61, 6567, 2212, 779, 632, 5709, 5667, 4923, 4911, 6893, 4695, 4164, 3536, 2287, 7594, 2848,
	3267, 1911, 3128, 546, 1991, 156, 4958, 5531, 6903, 483, 875, 138, 250, 2234, 2266, 7222,
	2842, 4258, 812, 6703, 232, 5207, 6650, 2585, 1900, 6225, 4932, 7265, 4701, 3173, 4635, 6393,
	227, 7313, 4454, 4284, 6759, 1224, 5223, 1447, 395, 2608, 4502, 4037, 189, 3348, 54, 6443,
	2210, 6230, 2826, 1780, 3002, 5995, 1955, 6102, 6045, 3938, 5019, 4417, 1434, 1262, 1507, 5847,
	5917, 7157, 7177, 6434, 7537, 741, 4348, 1309, 145, 374, 2236, 4496, 5028, 6771, 6923, 7421,
	1978, 1023, 3857, 6876, 1102, 7451, 4704, 6518, 1344, 765, 384, 5705, 1207, 1630, 4734, 1563,
	6839, 5933, 1954, 4987, 7142, 5814, 7527, 4953, 7637, 4707, 2182, 5734, 2818, 541, 4097, 5641,
}