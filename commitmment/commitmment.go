package commitmment

import (
	// "math/rand"

	// "github.com/fprasx/secrets-and-spies/ff"
)

// Completely secure implementation of Pederson commitmments in the group Z_(10**9+7)

// https://datatracker.ietf.org/doc/html/rfc3526#page-3

const p = 1000000007

// func PedersonParams() (g ff.Num, h ff.Num) {
// 	g = ff.New(p, 0xdecaf).Pow(0xc0ffee)
// 	h = ff.New(p, 0xb1eed).Pow(0xb100d)
// 	return
// }
//
// func Commit(value uint) (commitment ff.Num, nonce uint) {
// 	g, h := PedersonParams()
// 	nonce = uint(rand.Intn(p))
// 	commitment = g.Pow(value).Times(h.Pow(nonce))
// 	return
// }
//
// func MustValidateCommitment(commitment ff.Num, value uint, nonce uint) {
// 	g, h := PedersonParams()
// 	if g.Pow(value).Times(h.Pow(nonce)) != commitment {
// 		panic("invalid commitment")
// 	}
// }
