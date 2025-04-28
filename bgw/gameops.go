package bgw

import (
	"github.com/fprasx/secrets-and-spies/ff"
)

// returns a matrix where each row represents shares of the ith element of location vector
func ShareLocation(location int, noCities int, t int, n int, p int) ([][]Share, error) {
	shares := make([][]Share, noCities)
	for i := 0; i < noCities; i++ {
		if i != location {
			shares[i], _ = ShareSecret(ff.New(p, 0), t, n)

		} else {
			shares[i], _ = ShareSecret(ff.New(p, 1), t, n)
		}
	}
	return shares, nil
}

func DotProductConstant(a []ff.Num, b []Share, party int) Share {
	p := a[0].P()
	sum := ff.New(p, 0)
	for i := 0; i < len(a); i++ {
		sum = sum.Plus(a[i].Times(b[i][1]))
	}
	return Share{ff.New(p, party+1), sum}
}
