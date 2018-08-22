package bigint

import (
	"github.com/dedis/student_18_lattices/big"
)

// Int is a generic implementation of natural arithmetic on integers,
// built using Go's built-in "math/big.Int"
type Int struct {
	Value big.Int // Integer value, theoretically ranging from -infinite to +infinite
}

// NewInt creates a new Int with a given int64 value.
func NewInt(v int64) *Int {
	i := new(Int)
	i.Value.SetInt64(v)
	return i
}

// NewIntFromString creates a new Int from a string.
// A prefix of ``0x'' or ``0X'' selects base 16;
// the ``0'' prefix selects base 8, and
// a ``0b'' or ``0B'' prefix selects base 2.
// Otherwise the selected base is 10.
func NewIntFromString(s string) *Int {
	i := new(Int)
	i.Value.SetString(s, 0)
	return i
}

// SetInt sets Int i with value v
func (i *Int) SetInt(v int64) {
	i.Value.SetInt64(v)
}

// SetBigInt sets Int i with bigint.Int
func (i *Int) SetBigInt(v *Int) {
	i.Value.Set(&v.Value)
}

// SetString sets the value of i from a string
// A prefix of ``0x'' or ``0X'' selects base 16;
// the ``0'' prefix selects base 8, and
// a ``0b'' or ``0B'' prefix selects base 2.
// Otherwise the selected base is 10.
func (i *Int) SetString(s string) {
	i.Value.SetString(s, 0)
}

// Add sets the target i to a + b.
func (i *Int) Add(a, b *Int) *Int {
	i.Value.Add(&a.Value, &b.Value)
	return i
}

// Sub sets the target i to a - b.
func (i *Int) Sub(a, b *Int) *Int {
	i.Value.Sub(&a.Value, &b.Value)
	return i
}

// Mul sets the target i to a * b.
func (i *Int) Mul(a, b *Int) *Int {
	i.Value.Mul(&a.Value, &b.Value)
	return i
}

// Div sets the target i to ceil(a / b), which is the closest integer to zero for a/b
func (i *Int) Div(a, b *Int) *Int {
	i.Value.Quo(&a.Value, &b.Value)
	return i
}

// DivRound sets the target i to the integer closest to a / b .
func (i *Int) DivRound(a, b *Int) *Int {
	_a := NewInt(1)
	_a.SetBigInt(a)
	i.Value.Quo(&_a.Value, &b.Value)
	r := NewInt(1)
	r.Value.Rem(&_a.Value, &b.Value)
	r2 := NewInt(1).Mul(r, NewInt(2))
	if r2.Value.CmpAbs(&b.Value) != -1.0 {
		if _a.Value.Sign() == b.Value.Sign() {
			i.Add(i, NewInt(1))
		} else {
			i.Sub(i, NewInt(1))
		}
	}
	return i
}

// Exp sets the target i to a^b mod m
func (i *Int) Exp(a , b, m *Int) *Int {
	i.Value.Exp(&a.Value, &b.Value, &m.Value)
	return i
}

// Mod sets the target i to a mod m.
func (i *Int) Mod(a, m *Int) *Int {
	i.Value.Mod(&a.Value, &m.Value)
	return i
}

// Inv sets the target i to a^-1 mod m.
func (i *Int) Inv(a, m *Int) *Int {
	i.Value.ModInverse(&a.Value, &m.Value)
	return i
}

// Neg sets the target i to -a mod m.
func (i *Int) Neg(a, m *Int) *Int {
	i.Value.Neg(&a.Value)
	i.Mod(i, m)
	return i
}

// Lsh sets the target i to a << m.
func (i *Int) Lsh(a *Int, m uint32) *Int {
	i.Value.Lsh(&a.Value, uint(m))
	return i
}

// Rsh sets the target i to a >> m.
func (i *Int) Rsh(a *Int, m uint32) *Int {
	i.Value.Rsh(&a.Value, uint(m))
	return i
}

// And sets the target i to a & b.
func (i *Int) And(a, b *Int) *Int {
	i.Value.And(&a.Value, &b.Value)
	return i
}

// EqualTo judges if i and i2 have the same value.
func (i *Int) EqualTo(i2 *Int) bool {
	r := i.Value.Cmp(&i2.Value)
	if r == 0 {
		return true
	} else {
		return false
	}
}

// Cmp compares i and i2 and returns:
//
//   -1 if i <  i2
//    0 if i == i2
//   +1 if i >  i2
//
func (i *Int) Compare(i2 *Int) int{
	return i.Value.Cmp(&i2.Value)
}

// Bits returns the bit stream and bit length of i's absolute value.
// For example, 6=110, this function will return ([0, 1, 1], 3)
func (i *Int) Bits() ([]uint, uint) {
	var z Int
	z.Value.Abs(&i.Value)
	n := z.Value.BitLen()
	bits := make([]uint, n)
	for j := 0; j < n; j++ {
		bits[j] = z.Value.Bit(j)
	}
	return bits, uint(n)
}

// Uint32 returns the low 32 bits of i as uint32
func (i *Int) Uint32() uint32 {
	return uint32(i.Value.Int64())
}

// Int64 returns the low 64 bits of i as int64
func (i *Int) Int64() int64 {
	return i.Value.Int64()
}


func MontgomeryReduce(x, q, qInv, montgomeryMod *Int, bitLen uint32) *Int{
	u := new(Int).Mul(x, qInv)
	u.And(u, montgomeryMod)
	u.Mul(u, q)
	x.Add(x, u)
	x.Rsh(x, bitLen)
	if x.Compare(q) != -1.0 {
		return x.Sub(x, q)
	}
	return x
}

func montgomeryReduce(a int64) int64 {
	u := a * 7679
	u &= (1 << 18) - 1
	u *= 7681
	a += u
	a = int64(a >> 18)
	if a >= 7681 {
		return a - 7681
	}
	return a
}

func BarrettReduce(x, q *Int) *Int {
	u := new(Int).Rsh(x, uint32(13)) // ((uint32_t) a * sinv) >> 16
	u.Mul(u, q)
	x.Sub(x, u)
	return x
}

func barrettReduce(a int64) int64 {
	u := int64(a >> 13) // ((uint32_t) a * sinv) >> 16
	u *= 7681
	a -= int64(u)
	return a
}