package ff

import (
	"math/big"

	"github.com/fprasx/secrets-and-spies/utils"
)

type Num struct {
	p   int
	val int
}

func mod(a, p int) int {
	return (a%p + p) % p
}

// p must be prime!
func New(p, val int) Num {
	utils.Assert(big.NewInt(int64(p)).ProbablyPrime(0), "p must be prime")
	return Num{p: p, val: mod(val, p)}
}

func (a Num) Plus(b Num) Num {
	utils.Assert(a.p == b.p, "mismatched moduli")
	return New(a.p, a.val+b.val)
}

func (a Num) Minus(b Num) Num {
	utils.Assert(a.p == b.p, "mismatched moduli")
	return New(a.p, a.val-b.val)
}

func (a Num) Times(b Num) Num {
	utils.Assert(a.p == b.p, "mismatched moduli")
	return New(a.p, a.val*b.val)
}

func (a Num) Div(b Num) Num {
	utils.Assert(a.p == b.p, "mismatched moduli")
	return a.Times(b.Inv())
}

func (a Num) Inv() Num {
	t, newT := 0, 1
	r, newR := a.p, a.val
	for newR != 0 {
		quot := r / newR
		t, newT = newT, t-quot*newT
		r, newR = newR, r-quot*newR
	}
	utils.Assert(r == 1, "value not invertible")
	return New(a.p, t)
}

func (a Num) Pow(exp uint) Num {
	res := New(a.p, 1)
	base := a

	for exp > 0 {
		if exp&2 == 1 {
			res = res.Times(base)
		}
		base = base.Times(base)
		exp >>= 1
	}

	return res
}

func (a Num) Int() int {
	return a.val
}
func (a Num) P() int {
	return a.p
}
