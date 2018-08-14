package bigint

import "testing"


// test vectors for function Add
type argAdd struct {
	x, y, want *Int
}
var addVec = []argAdd{
	{NewInt(0), NewInt(0), NewInt(0)},
	{NewInt(0), NewInt(1), NewInt(1)},
	{NewInt(0), NewInt(-1), NewInt(-1)},
	{NewInt(1), NewInt(-1), NewInt(0)},
	{NewInt(123456789), NewInt(987654321), NewInt(1111111110)},
	{NewInt(-123456789), NewInt(987654321), NewInt(864197532)},
	{NewInt(-123456789), NewInt(-987654321), NewInt(-1111111110)},
	{NewInt(123456789), NewIntFromString("1234567899999999999999"), NewIntFromString("1234567900000123456788")},
	{NewInt(-123456789), NewIntFromString("1234567899999999999999"), NewIntFromString("1234567899999876543210")},
	{NewIntFromString("123456789123456789123456789123456789"), NewIntFromString("-123456789123456789123456789123456789"), NewInt(0)},
	{NewIntFromString("123456789123456789123456789123456789"), NewIntFromString("987654321987654321987654321987654321"), NewIntFromString("1111111111111111111111111111111111110")},
}

func TestAdd(t *testing.T) {
	var z Int
	for i, testPair := range addVec {
		if !z.Add(testPair.x, testPair.y).EqualTo(testPair.want) {
			t.Errorf("Error Add test pair %v", i)
		}
	}
}

func BenchmarkAdd(b *testing.B) {
	var z Int
	x := NewIntFromString("123456789123456789123456789123456789")
	y := NewIntFromString("-987654321987654321987654321987654321")
	for i := 0; i < b.N; i++ {
		z.Add(x, y)
	}
}


// test vectors for function Sub
type argSub struct {
	x, y, want *Int
}
var subVec = []argSub{
	{NewInt(0), NewInt(0), NewInt(0)},
	{NewInt(0), NewInt(1), NewInt(-1)},
	{NewInt(0), NewInt(-1), NewInt(1)},
	{NewInt(-1), NewInt(-1), NewInt(0)},
	{NewInt(123456789), NewInt(987654321), NewInt(-864197532)},
	{NewInt(-123456789), NewInt(987654321), NewInt(-1111111110)},
	{NewInt(-123456789), NewInt(-987654321), NewInt(864197532)},
	{NewInt(123456789), NewIntFromString("1234567899999999999999"), NewIntFromString("-1234567899999876543210")},
	{NewInt(-123456789), NewIntFromString("1234567899999999999999"), NewIntFromString("-1234567900000123456788")},
	{NewIntFromString("123456789123456789123456789123456789"), NewIntFromString("987654321987654321987654321987654321"), NewIntFromString("-864197532864197532864197532864197532")},
	{NewIntFromString("-123456789123456789123456789123456789"), NewIntFromString("-987654321987654321987654321987654321"), NewIntFromString("864197532864197532864197532864197532")},
}

func TestSub(t *testing.T) {
	var z Int
	for i, testPair := range subVec {
		if !z.Sub(testPair.x, testPair.y).EqualTo(testPair.want) {
			t.Errorf("Error Sub test pair %v", i)
		}
	}
}

func BenchmarkSub(b *testing.B) {
	var z Int
	x := NewIntFromString("123456789123456789123456789123456789")
	y := NewIntFromString("-987654321987654321987654321987654321")
	for i := 0; i < b.N; i++ {
		z.Sub(x, y)
	}
}


// test vectors for function Mul
type argMul struct {
	x, y, want *Int
}
var mulVec = []argMul{
	{NewInt(0), NewInt(0), NewInt(0)},
	{NewInt(1), NewInt(0), NewInt(0)},
	{NewInt(1), NewInt(-1), NewInt(-1)},
	{NewInt(-1), NewInt(-1), NewInt(1)},
	{NewInt(123456789), NewInt(987654321), NewInt(121932631112635269)},
	{NewInt(-123456789), NewInt(987654321), NewInt(-121932631112635269)},
	{NewInt(-123456789), NewInt(-987654321), NewInt(121932631112635269)},
	{NewInt(123456789), NewIntFromString("1234567899999999999999"), NewIntFromString("152415788736473099999876543211")},
	{NewInt(-123456789), NewIntFromString("1234567899999999999999"), NewIntFromString("-152415788736473099999876543211")},
	{NewIntFromString("123456789123456789123456789123456789"), NewIntFromString("-987654321987654321987654321987654321"), NewIntFromString("-121932631356500531591068431825636331816338969581771069347203169112635269")},
	{NewIntFromString("-123456789123456789123456789123456789"), NewIntFromString("-987654321987654321987654321987654321"), NewIntFromString("121932631356500531591068431825636331816338969581771069347203169112635269")},
}

func TestMul(t *testing.T) {
	var z Int
	for i, testPair := range mulVec {
		if !z.Mul(testPair.x, testPair.y).EqualTo(testPair.want) {
			t.Errorf("Error Mul test pair %v", i)
		}
	}
}

func BenchmarkMul(b *testing.B) {
	var z Int
	x := NewIntFromString("123456789123456789123456789123456789")
	y := NewIntFromString("-987654321987654321987654321987654321")
	for i := 0; i < b.N; i++ {
		z.Mul(x, y)
	}
}


// test vectors for function Div
type argDiv struct {
	x, y, want *Int
}
var divVec = []argDiv{
	{NewInt(0), NewInt(1), NewInt(0)},
	{NewInt(1), NewInt(2), NewInt(0)},
	{NewInt(5), NewInt(2), NewInt(2)},
	{NewInt(17), NewInt(-2), NewInt(-8)},
	{NewInt(987654321), NewInt(123456789), NewInt(8)},
	{NewInt(-987654320), NewInt(123456789), NewInt(-8)},
	{NewInt(-121932631112635269), NewInt(-987654321), NewInt(123456789)},
	{NewIntFromString("123456789123456789123456789123456789"), NewInt(123456789), NewIntFromString("1000000001000000001000000001")},
	{NewIntFromString("123456789123456789123456789123456789"), NewInt(-123456789), NewIntFromString("-1000000001000000001000000001")},
	{NewIntFromString("987654321987654321987654321987654321"), NewIntFromString("123456789123456789123456789123456789"), NewInt(8)},
	{NewIntFromString("-987654321987654321987654321987654321"), NewIntFromString("-123456789123456789123456789123456789"), NewInt(8)},
}

func TestDiv(t *testing.T) {
	var z Int
	for i, testPair := range divVec {
		if !z.Div(testPair.x, testPair.y).EqualTo(testPair.want) {
			t.Errorf("Error Div test pair %v", i)
		}
	}
}

func BenchmarkDiv(b *testing.B) {
	var z Int
	x := NewIntFromString("987654321987654321987654321987654321")
	y := NewIntFromString("-123456789123456789123456789123456789")
	for i := 0; i < b.N; i++ {
		z.Div(x, y)
	}
}


// test vectors for function DivRound
type argDivRound struct {
	x, y, want *Int
}
var divRoundVec = []argDivRound{
	{NewInt(0), NewInt(1), NewInt(0)},
	{NewInt(1), NewInt(2), NewInt(1)},
	{NewInt(5), NewInt(2), NewInt(3)},
	{NewInt(5), NewInt(3), NewInt(2)},
	{NewInt(5), NewInt(-2), NewInt(-3)},
	{NewInt(-5), NewInt(2), NewInt(-3)},
	{NewInt(-5), NewInt(-2), NewInt(3)},
	{NewInt(987654321), NewInt(123456789), NewInt(8)},
	{NewInt(-987654320), NewInt(123456789), NewInt(-8)},
	{NewInt(-121932631112635269), NewInt(-987654321), NewInt(123456789)},
	{NewIntFromString("123456789123456789123456789123456789"), NewInt(123456789), NewIntFromString("1000000001000000001000000001")},
	{NewIntFromString("987654321987654321987654321987654321"), NewIntFromString("123456789123456789123456789123456789"), NewInt(8)},
	{NewIntFromString("-987654321987654321987654321987654321"), NewIntFromString("-123456789123456789123456789123456789"), NewInt(8)},
}

func TestDivRound(t *testing.T) {
	var z Int
	for i, testPair := range divRoundVec {
		if !z.DivRound(testPair.x, testPair.y).EqualTo(testPair.want) {
			t.Errorf("Error DivRound test pair %v", i)
		}
	}
}

func BenchmarkDivRound(b *testing.B) {
	var z Int
	x := NewIntFromString("987654321987654321987654321987654321")
	y := NewIntFromString("-123456789123456789123456789123456789")
	for i := 0; i < b.N; i++ {
		z.DivRound(x, y)
	}
}


// test vectors for function Exp
type argExp struct {
	x, y, m, want *Int
}
var expVec = []argExp{
	{NewInt(0), NewInt(1), NewInt(2), NewInt(0)},
	{NewInt(1), NewInt(0), NewInt(2), NewInt(1)},
	{NewInt(-1), NewInt(0), NewInt(2), NewInt(1)},
	{NewInt(123456789), NewInt(12345), NewInt(987654321), NewInt(658957095)},
	{NewInt(-123456789), NewInt(12345), NewInt(987654321), NewInt(328697226)},
	{NewInt(123456789), NewInt(12345), NewIntFromString("123456789123456789"), NewIntFromString("87718977473362236")},
	{NewInt(-123456789), NewInt(12345), NewIntFromString("123456789123456789"), NewIntFromString("35737811650094553")},
	{NewIntFromString("123456789123456789"), NewInt(12345), NewIntFromString("987654321987654321"), NewIntFromString("313081623313081623")},
	{NewIntFromString("-123456789123456789"), NewInt(12345), NewIntFromString("987654321987654321"), NewIntFromString("674572698674572698")},
}

func TestExp(t *testing.T) {
	var z Int
	for i, testPair := range expVec {
		if !z.Exp(testPair.x, testPair.y, testPair.m).EqualTo(testPair.want) {
			t.Errorf("Error Exp test pair %v", i)
		}
	}
}

func BenchmarkExp(b *testing.B) {
	var z Int
	x := NewIntFromString("987654321987654321987654321987654321")
	y := NewIntFromString("12345")
	m := NewIntFromString("6780883635459973527839456474")
	for i := 0; i < b.N; i++ {
		z.Exp(x, y, m)
	}
}


// test vectors for function Mod
type argMod struct {
	x, m, want *Int
}
var modVec = []argMod{
	{NewInt(0), NewInt(1), NewInt(0)},
	{NewInt(1), NewInt(1), NewInt(0)},
	{NewInt(1), NewInt(2), NewInt(1)},
	{NewInt(5), NewInt(2), NewInt(1)},
	{NewInt(5), NewInt(3), NewInt(2)},
	{NewInt(-5), NewInt(2), NewInt(1)},
	{NewInt(-5), NewInt(4), NewInt(3)},
	{NewInt(987654321), NewInt(123456789), NewInt(9)},
	{NewInt(-987654321), NewInt(123456789), NewInt(123456780)},
	{NewIntFromString("123456789123456789123456789123456789"), NewInt(123456789), NewIntFromString("0")},
	{NewInt(123456789), NewIntFromString("123456789123456789123456789123456789"), NewIntFromString("123456789")},
	{NewIntFromString("123456789123456789123456789123456789"), NewInt(987654321), NewInt(246911409)},
	{NewIntFromString("98765432198765432198734567654567876789654321987654321"), NewIntFromString("123456789123456789123456789123456789"), NewIntFromString("41821061497654370731506248029506247")},
	{NewIntFromString("-98765432198765432198734567654567876789654321987654321"), NewIntFromString("123456789123456789123456789123456789"), NewIntFromString("81635727625802418391950541093950542")},
}

func TestMod(t *testing.T) {
	var z Int
	for i, testPair := range modVec {
		if !z.Mod(testPair.x, testPair.m).EqualTo(testPair.want) {
			t.Errorf("Error Mod test pair %v", i)
		}
	}
}

func BenchmarkMod(b *testing.B) {
	var z Int
	x := NewIntFromString("-98765432198765432198734567654567876789654321987654321")
	y := NewIntFromString("123456789123456789123456789123456789")
	for i := 0; i < b.N; i++ {
		z.Mod(x, y)
	}
}


// test vectors for function Inv
type argInv struct {
	x, m *Int
}
var invVec = []argInv{
	{NewInt(1), NewInt(2)},
	{NewInt(-1), NewInt(2)},
	{NewInt(12345), NewInt(10001473)},
	{NewInt(123456789123456789), NewIntFromString("1152921504382476289")},
	{NewIntFromString("123456789123456789"), NewIntFromString("1152921504382476289")},
	{NewIntFromString("-123456789123456789"), NewIntFromString("1152921504382476289")},
}

func TestInv(t *testing.T) {
	var z Int
	one := NewInt(1)
	for i, testPair := range invVec {
		z.Inv(testPair.x, testPair.m)
		z.Mul(&z, testPair.x)
		z.Mod(&z, testPair.m)
		if !z.EqualTo(one) {
			t.Errorf("Error Inv test pair %v", i)
		}
	}
}

func BenchmarkInv(b *testing.B) {
	var z Int
	x := NewIntFromString("123456789123456789")
	y := NewIntFromString("1152921504382476289")
	for i := 0; i < b.N; i++ {
		z.Inv(x, y)
	}
}


// test vectors for function Neg
type argNeg struct {
	x, m *Int
}
var negVec = []argNeg{
	{NewInt(1), NewInt(2)},
	{NewInt(-1), NewInt(2)},
	{NewInt(12345), NewInt(10001473)},
	{NewInt(123456789123456789), NewIntFromString("1152921504382476289")},
	{NewIntFromString("123456789123456789"), NewIntFromString("1152921504382476289")},
	{NewIntFromString("-123456789123456789"), NewIntFromString("1152921504382476289")},
}

func TestNeg(t *testing.T) {
	var z Int
	zero := NewInt(0)
	for i, testPair := range negVec {
		z.Neg(testPair.x, testPair.m)
		z.Add(&z, testPair.x)
		z.Mod(&z, testPair.m)
		if !z.EqualTo(zero) {
			t.Errorf("Error Neg test pair %v", i)
		}
	}
}

func BenchmarkNeg(b *testing.B) {
	var z Int
	x := NewIntFromString("123456789123456789")
	y := NewIntFromString("1152921504382476289")
	for i := 0; i < b.N; i++ {
		z.Neg(x, y)
	}
}


// test vectors for function Bits
type argBits struct {
	x *Int
	xBits []uint
	xBitLen uint
}
var bitsVec = []argBits{
	{NewInt(0), []uint{0}, uint(0)},
	{NewInt(1), []uint{1}, uint(1)},
	{NewInt(-1), []uint{1}, uint(1)},
	{NewInt(2), []uint{0, 1}, uint(2)},
	{NewInt(-2), []uint{0, 1}, uint(2)},
	{NewInt(123456789), []uint{1, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 1, 1, 1, 1, 0, 1, 1, 0, 1, 0, 1, 1, 1}, uint(27)},
	{NewInt(-123456789), []uint{1, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 1, 1, 1, 1, 0, 1, 1, 0, 1, 0, 1, 1, 1}, uint(27)},
	{NewIntFromString("1152921504382476289"), []uint{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, uint(60)},
	{NewIntFromString("-1152921504382476289"), []uint{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, uint(60)},
}

func TestBits(t *testing.T) {
	var zBits []uint
	var zBitLen uint
	for i, testPair := range bitsVec {
		zBits, zBitLen = testPair.x.Bits()
		if zBitLen != testPair.xBitLen {
			t.Errorf("Error Bits test pair %v", i)
		}
		for j := range zBits {
			if zBits[j] != testPair.xBits[j] {
				t.Errorf("Error Bits test pair %v", i)
			}
		}
	}
}

func BenchmarkBits(b *testing.B) {
	x := NewIntFromString("123456789123456789123456789123456789")
	for i := 0; i < b.N; i++ {
		x.Bits()
	}
}
