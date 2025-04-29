package bgw

// import (
// 	ffarith "github.com/fprasx/secrets-and-spies/ff"
// 	"github.com/fprasx/secrets-and-spies/utils"
// )

// // returns a matrix where each row represents shares of the ith element of location vector
// func ShareLocation(location int, noCities int, t int, n int, p int) ([][][2]ffarith.FFNum, error) {
// 	shares := make([][][2]ffarith.FFNum, noCities)
// 	for i := 0; i < noCities; i++ {
// 		if i != location {
// 			shares[i], _ = ShareSecret(ffarith.NewFFNum(p, 0), t, n)

// 		} else {
// 			shares[i], _ = ShareSecret(ffarith.NewFFNum(p, 1), t, n)
// 		}
// 	}
// 	return shares, nil
// }

// // outputs a share of the validated move computation
// func ValidateMoveShares(graph [][]ffarith.FFNum, noCities int, prevLoc [][2]ffarith.FFNum, newLoc [][2]ffarith.FFNum, party int) ([2]ffarith.FFNum, [][2]ffarith.FFNum) {
// 	newShare := make([][2]ffarith.FFNum, noCities)
// 	for i := 0; i < noCities; i++ {
// 		newShare[i] = DotProductConstant(graph[i], prevLoc, party)
// 	}
// 	return DotProductShares(newShare, newLoc, party), newShare

// }

// // returns shares of a 2t degree polynomial
// func DotProductShares(a [][2]ffarith.FFNum, b [][2]ffarith.FFNum, party int) [2]ffarith.FFNum {
// 	p := a[0][0].P()
// 	utils.Assert(a[0][0].P() == b[0][0].P(), "mismatched prime")
// 	sum := ffarith.NewFFNum(p, 0)
// 	for i := 0; i < len(a); i++ {
// 		sum = sum.Plus(a[i][1].Times(b[i][1]))
// 	}
// 	return [2]ffarith.FFNum{ffarith.NewFFNum(p, party+1), sum}
// }

// func DotProductConstant(a []ffarith.FFNum, b [][2]ffarith.FFNum, party int) [2]ffarith.FFNum {
// 	p := a[0].P()
// 	utils.Assert(a[0].P() == b[0][0].P(), "mismatched prime")
// 	sum := ffarith.NewFFNum(p, 0)
// 	for i := 0; i < len(a); i++ {
// 		//	fmt.Printf("%d:%d ", a[i], b[i][1])
// 		//	fmt.Printf("\n")
// 		sum = sum.Plus(a[i].Times(b[i][1]))
// 	}
// 	return [2]ffarith.FFNum{ffarith.NewFFNum(p, party+1), sum}
// }
