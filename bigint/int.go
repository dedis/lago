package bigint

import "github.com/dedis/student_18_lattices/big"

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
	zero := NewInt(0)
	if a.EqualTo(zero) {
		return zero
	}
	_a := NewInt(1)
	_a.SetBigInt(a)
	_b := NewInt(1)
	_b.SetBigInt(b)
	r := NewInt(1)
	i.Value.Quo(&_a.Value, &_b.Value)
	r.Value.Rem(&_a.Value, &_b.Value)
	midValue := NewInt(1)
	midValue.Value.Quo(&_b.Value, &NewInt(2).Value)
	if !(NewInt(1).Value.Abs(&r.Value).Cmp(NewInt(1).Value.Abs(&midValue.Value)) == -1.0) {
		if i.Value.Cmp(&zero.Value) == -1.0 {
			i.Sub(i, NewInt(1))
		} else {
			i.Add(i, NewInt(1))
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
