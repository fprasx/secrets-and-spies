package ff_test

import (
	"math/big"
	"testing"

	"github.com/fprasx/secrets-and-spies/ff"
)

func TestPow0(t *testing.T) {
	num := ff.New(0xc0ffee)
	res := num.Pow(big.NewInt(0))
	if !res.Eq(ff.New(1)) {
		t.Errorf("got %v, which is not 1", res.BigInt())
	}
}

func TestSub(t *testing.T) {
	x := ff.New(100)
	y := ff.New(200)
	xmy := x.Minus(y)
	shouldbex := xmy.Plus(y)
	if !shouldbex.Eq(x) {
		t.Errorf("expected %v, got %v", x.BigInt(), shouldbex.BigInt())
	}
}

func TestDivInvRoundtrip(t *testing.T) {
	y := ff.New(200)
    roundtrip1 := y.Inv().Times(y)
    if !roundtrip1.Eq(ff.New(1)) {
        t.Errorf("expected 1, got %v", roundtrip1.BigInt())
    }

	x := ff.New(100)
    roundtrip2 := x.Div(y).Times(y)
    if !roundtrip2.Eq(x) {
        t.Errorf("expected 100, got %v", roundtrip2.BigInt())
    }

}
