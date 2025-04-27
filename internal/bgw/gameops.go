package bgw

import (
	"fmt"

	ffarith "github.com/fprasx/secrets-and-spies/internal/ff_arith"
)

// returns a matrix where each row represents shares of the ith element of location vector
func ShareLocation(location int, noCities int, t int, n int, p int) ([][][2]ffarith.FFNum, error) {
	shares := make([][][2]ffarith.FFNum, noCities)
	for i := 0; i < noCities; i++ {
		if i != location {
			shares[i], _ = ShareSecret(ffarith.NewFFNum(p, 0), t, n)

		} else {
			shares[i], _ = ShareSecret(ffarith.NewFFNum(p, 1), t, n)
		}
	}
	return shares, nil
}
func DotProductConstant(a []ffarith.FFNum, b [][2]ffarith.FFNum, party int) [2]ffarith.FFNum {
	p := a[0].P()
	sum := ffarith.NewFFNum(p, 0)
	for i := 0; i < len(a); i++ {
		fmt.Printf("%d:%d ", a[i], b[i][1])
		fmt.Printf("\n")
		sum = sum.Plus(a[i].Times(b[i][1]))
	}
	return [2]ffarith.FFNum{ffarith.NewFFNum(p, party+1), sum}
}
