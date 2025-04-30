package commitment

import (
	crand "crypto/rand"
	"math/big"

	"github.com/fprasx/secrets-and-spies/ff"
)

func PedersonParams() (g ff.Num, h ff.Num) {
	g = ff.New(2)
	h = ff.New(0xb1eed).Pow(ff.New(0xb100d).BigInt())
	return
}

func Commit(value uint) (commitment *big.Int, nonce *big.Int) {
	g, h := PedersonParams()
	nonce, err := crand.Int(crand.Reader, ff.FieldPrime())
	if err != nil {
		panic("failed to generate nonce")
	}
	commitment = g.Pow(big.NewInt(int64(value))).Times(h.Pow(nonce)).BigInt()
	return
}

func Verify(commitment *big.Int, value uint, nonce *big.Int) bool {
	g, h := PedersonParams()
	recalculated := g.Pow(big.NewInt(int64(value))).Times(h.Pow(nonce))
	return recalculated.BigInt().Cmp(commitment) == 0
}
