package bgw

import (
	"fmt"

	"github.com/fprasx/secrets-and-spies/ff"
)

// returns a matrix where each row represents shares of the ith element of location vector
func ShareLocation(location int, noCities int, t int, n int, p int) ([][][2]ff.Num, error) {
	shares := make([][][2]ff.Num, noCities)
	for i := 0; i < noCities; i++ {
		if i != location {
			shares[i], _ = ShareSecret(ff.New(p, 0), t, n)

		} else {
			shares[i], _ = ShareSecret(ff.New(p, 1), t, n)
		}
	}
	return shares, nil
}
func DotProductConstant(a []ff.Num, b [][2]ff.Num, party int) [2]ff.Num {
	p := a[0].P()
	sum := ff.New(p, 0)
	for i := 0; i < len(a); i++ {
		fmt.Printf("%d:%d ", a[i], b[i][1])
		fmt.Printf("\n")
		sum = sum.Plus(a[i].Times(b[i][1]))
	}
	return [2]ff.Num{ff.New(p, party+1), sum}
}
