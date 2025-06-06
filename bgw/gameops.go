package bgw

import (
	"github.com/fprasx/secrets-and-spies/ff"
)

// returns a matrix where each row represents shares of the ith element of location vector
func ShareLocation(location int, noCities int, t int, n int) ([][][2]ff.Num, error) {
	shares := make([][][2]ff.Num, noCities)
	for i := 0; i < noCities; i++ {
		if i != location {
			shares[i], _ = ShareSecret(ff.New(0), t, n)

		} else {
			shares[i], _ = ShareSecret(ff.New(1), t, n)
		}
	}
	return shares, nil
}

// outputs a share of the validated move computation
func ValidateMoveShares(graph [][]ff.Num, noCities int, prevLoc [][2]ff.Num, newLoc [][2]ff.Num, party int) ([2]ff.Num, [][2]ff.Num) {
	newShare := make([][2]ff.Num, noCities)
	for i := 0; i < noCities; i++ {
		newShare[i] = DotProductConstant(graph[i], prevLoc, party)
	}
	return DotProductShares(newShare, newLoc, party), newShare

}

// returns shares of a 2t degree polynomial
func DotProductShares(a [][2]ff.Num, b [][2]ff.Num, party int) [2]ff.Num {
	sum := ff.New(0)
	for i := 0; i < len(a); i++ {
		sum = sum.Plus(a[i][1].Times(b[i][1]))
	}
	return [2]ff.Num{ff.New(int64(party + 1)), sum}
}

func DotProductConstant(a []ff.Num, b [][2]ff.Num, party int) [2]ff.Num {
	sum := ff.New(0)
	for i := 0; i < len(a); i++ {
		//	fmt.Printf("%d:%d ", a[i], b[i][1])
		//	fmt.Printf("\n")
		sum = sum.Plus(a[i].Times(b[i][1]))
	}
	return [2]ff.Num{ff.New(int64(party + 1)), sum}
}
