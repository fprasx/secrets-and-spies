package ff

import (
	"math/big"

	"github.com/fprasx/secrets-and-spies/utils"
)

type Num struct {
	Inner *big.Int
}

var p *big.Int = nil

func FieldPrime() *big.Int {
	if p == nil {
		// https://datatracker.ietf.org/doc/html/rfc3526#section-2
		s := "FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD1" +
			"29024E088A67CC74020BBEA63B139B22514A08798E3404DD" +
			"EF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245" +
			"E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7ED" +
			"EE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3D" +
			"C2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F" +
			"83655D23DCA3AD961C62F356208552BB9ED529077096966D" +
			"670C354E4ABC9804F1746C08CA237327FFFFFFFFFFFFFFFF"
		p = new(big.Int)
		p.SetString(s, 16)
	}
	return p
}

// modify returns num mod p. Does not actually modify num.
func modify(num *big.Int) *big.Int {
	mod := new(big.Int)
	quot := new(big.Int)
	quot.DivMod(num, FieldPrime(), mod)
	return mod
}

func New(val int64) Num {
	return Num{Inner: big.NewInt(int64(val))}
}

func (a Num) Plus(b Num) Num {
	aa := a.Inner
	bb := b.Inner
	sum := new(big.Int)
	sum.Add(aa, bb)
	return Num{Inner: modify(sum)}
}

func (a Num) Minus(b Num) Num {
	aa := a.Inner
	bb := b.Inner
	diff := new(big.Int)
	diff.Sub(aa, bb)
	return Num{Inner: modify(diff)}
}

func (a Num) Times(b Num) Num {
	aa := a.Inner
	bb := b.Inner
	prod := new(big.Int)
	prod.Mul(aa, bb)
	return Num{Inner: modify(prod)}
}

func (a Num) Div(b Num) Num {
	return a.Times(b.Inv())
}

func (a Num) Inv() Num {
	aa := a.Inner
	inv := new(big.Int)
	inv.ModInverse(aa, FieldPrime())
	return Num{Inner: inv}
}

func (a Num) Pow(exponent *big.Int) Num {
	aa := a.Inner
	res := new(big.Int)
	res.Exp(aa, exponent, FieldPrime())
	return Num{Inner: res}
}

func (a Num) IsZero() bool {
	return a.Inner.Sign() == 0
}

func (a Num) IsNonZero() bool {
	utils.Assert(a.Inner.Sign() != -1, "field element shouldn't be negative")
	return a.Inner.Sign() == 1
}

func (a Num) Eq(b Num) bool {
	return a.Inner.Cmp(b.Inner) == 0
}

func (a Num) Neq(b Num) bool {
	return a.Inner.Cmp(b.Inner) != 0
}

func (a Num) BigInt() *big.Int {
	return a.Inner
}
