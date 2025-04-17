package ffarith

import "github.com/fprasx/secrets-and-spies/internal/utils"

type FFNum struct {
	p   int
	val int
}

func mod(a, p int) int {
	a = a % p
	if a < 0 {
		a += p
	}
	return a
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i < n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// p must be prime!
func NewFFNum(p, val int) FFNum {
	utils.Assert(isPrime(p), "p must be prime")
	return FFNum{p: p, val: mod(val, p)}
}

func (a FFNum) Plus(b FFNum) FFNum {
	utils.Assert(a.p == b.p, "mismatched moduli")
	return NewFFNum(a.p, a.val+b.val)
}

func (a FFNum) Minus(b FFNum) FFNum {
	utils.Assert(a.p == b.p, "mismatched moduli")
	return NewFFNum(a.p, a.val-b.val)
}

func (a FFNum) Times(b FFNum) FFNum {
	utils.Assert(a.p == b.p, "mismatched moduli")
	return NewFFNum(a.p, a.val*b.val)
}

func (a FFNum) Over(b FFNum) FFNum {
	utils.Assert(a.p == b.p, "mismatched moduli")
	return a.Times(b.Inv())
}

func (a FFNum) Inv() FFNum {
	t, newT := 0, 1
	r, newR := a.p, a.val
	for newR != 0 {
		quot := r / newR
		t, newT = newT, t-quot*newT
		r, newR = newR, r-quot*newR
	}
	utils.Assert(r == 1, "value not invertible")
	return NewFFNum(a.p, t)
}

func (a FFNum) ToThe(power uint) FFNum {
	res := NewFFNum(a.p, 1)
	for i := uint(0); i < power; i++ {
		res = res.Times(a)
	}
	return res
}

func (a FFNum) Int() int {
	return a.val
}
