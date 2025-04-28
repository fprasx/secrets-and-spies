package main

import (
	"fmt"

	"github.com/fprasx/secrets-and-spies/bgw"
	ffarith "github.com/fprasx/secrets-and-spies/ff"
)

// example for checking if moved onto city
func testMoveCheck() int {
	p := 29
	t := 2
	n := 5
	noCities := 4
	tempshares := make([][][2]ffarith.Num, noCities)
	for i := 0; i < noCities; i++ {
		tempshares[i] = make([][2]ffarith.Num, n)
	}
	oldloc, err := bgw.ShareLocation(3, noCities, t, n, p)
	if err != nil {
		panic(err)
	}
	newloc, err := bgw.ShareLocation(3, noCities, t, n, p)
	if err != nil {
		panic(err)
	}
	endshares := make([][2]ffarith.Num, n)

	for i := 0; i < n; i++ {
		columnold := make([][2]ffarith.Num, noCities)
		for j := 0; j < noCities; j++ {
			columnold[j] = oldloc[j][i]
		}
		columnnew := make([][2]ffarith.Num, noCities)
		for j := 0; j < noCities; j++ {
			columnnew[j] = newloc[j][i]
		}
		fmt.Println()
		endshares[i] = bgw.DotProductShares(columnold, columnnew, i)

	}

	res, err := bgw.ReconstructSecret(endshares)

	if err != nil {
		panic(err)
	}
	fmt.Println("Verdict is:", res.Int())
	if res.Int() != 1 {
		return -1
	}
	return 0
}

// example of move validation
func testMoveValidation() int {
	p := 29
	t := 2
	n := 5
	noCities := 4
	//0 - 3
	//1 - 2
	g := [][]int{
		{1, 0, 0, 1},
		{0, 1, 1, 0},
		{0, 1, 1, 0},
		{1, 0, 0, 1},
	}
	graph := make([][]ffarith.Num, noCities)
	tempshares := make([][][2]ffarith.Num, noCities)
	for i := 0; i < noCities; i++ {
		tempshares[i] = make([][2]ffarith.Num, n)
		graph[i] = make([]ffarith.Num, noCities)
		for j := 0; j < noCities; j++ {
			graph[i][j] = ffarith.New(p, g[i][j])
		}
	}
	oldloc, err := bgw.ShareLocation(1, noCities, t, n, p)
	if err != nil {
		panic(err)
	}
	newloc, err := bgw.ShareLocation(2, noCities, t, n, p)
	if err != nil {
		panic(err)
	}
	endshares := make([][2]ffarith.Num, n)

	for i := 0; i < n; i++ {
		columnold := make([][2]ffarith.Num, noCities)
		for j := 0; j < noCities; j++ {
			columnold[j] = oldloc[j][i]
		}
		columnnew := make([][2]ffarith.Num, noCities)
		for j := 0; j < noCities; j++ {
			columnnew[j] = newloc[j][i]
		}
		fmt.Println()
		tempcolumn := make([][2]ffarith.Num, n)
		endshares[i], tempcolumn = bgw.ValidateMoveShares(graph, noCities, columnold, columnnew, i)
		for j := 0; j < noCities; j++ {
			tempshares[j][i] = tempcolumn[j]
		}
	}

	for i := 0; i < noCities; i++ {
		res, _ := bgw.ReconstructSecret(tempshares[i])
		fmt.Printf("%d ", res.Int())
	}
	fmt.Println()
	res, err := bgw.ReconstructSecret(endshares)
	if err != nil {
		panic(err)
	}
	fmt.Println("Verdict is:", res.Int())
	if res.Int() != 1 {
		return -1
	}
	return 0
}

// test for reconstruction
func test1() int {
	p := 29
	points := [][2]ffarith.Num{
		{ffarith.New(p, 1), ffarith.New(p, 10)},
		{ffarith.New(p, 2), ffarith.New(p, 21)},
		{ffarith.New(p, 3), ffarith.New(p, 9)},
	}

	secret, err := bgw.ReconstructSecret(points)
	if err != nil {
		panic(err)
	}

	fmt.Println("Secret is:", secret.Int())
	if secret.Int() != 5 {
		return secret.Int()
	}
	return 0
}

// test for sharing a secret
func testShare() int {
	p := 29
	secret := ffarith.New(p, 28)
	shares, err := bgw.ShareSecret(secret, 10, 20)
	if err != nil {
		panic(err)
	}
	recon, err := bgw.ReconstructSecret(shares)
	if err != nil {
		panic(err)
	}
	fmt.Println("Secret is:", recon.Int())
	if secret.Int() != secret.Int() {
		return secret.Int()
	}
	return 0
}
func evaluatePolynomial(coeffs []ffarith.Num, x ffarith.Num) ffarith.Num {
	p := x.P()
	result := ffarith.New(p, 0)
	power := ffarith.New(p, 1) // x^0

	for _, coeff := range coeffs {
		result = result.Plus(coeff.Times(power))
		power = power.Times(x) // next power of x
	}

	return result
}

func testDotProduct() int {
	p := 29
	t := 2
	n := 5
	noCities := 3
	a := []ffarith.Num{
		ffarith.New(p, 0),
		ffarith.New(p, 1),
		ffarith.New(p, 0),
	}
	shares, err := bgw.ShareLocation(2, noCities, t, n, p)
	if err != nil {
		panic(err)
	}
	newShares := make([][2]ffarith.Num, n)
	for i := 0; i < n; i++ {
		column := make([][2]ffarith.Num, noCities)
		for j := 0; j < noCities; j++ {
			column[j] = shares[j][i]
		}
		fmt.Println()
		newShares[i] = bgw.DotProductConstant(a, column, i)
	}
	secret, err := bgw.ReconstructSecret(newShares)
	if err != nil {
		panic(err)
	}

	fmt.Println("Secret is:", secret.Int())
	if secret.Int() != 0 {
		return secret.Int()
	}
	return 0
}

// Polynomial is x^4+x^3+2x^2+3x+1
// Reduced is 2x^2+3x+1
func testDegreeReduce() int {
	p := 29
	poly := []ffarith.Num{
		ffarith.New(p, 1),
		ffarith.New(p, 3),
		ffarith.New(p, 2),
		ffarith.New(p, 1),
		ffarith.New(p, 1),
	}
	reducedpoly := []ffarith.Num{
		ffarith.New(p, 1),
		ffarith.New(p, 3),
		ffarith.New(p, 2),
	}
	g := []ffarith.Num{
		ffarith.New(p, evaluatePolynomial(poly, ffarith.New(p, 1)).Int()),
		ffarith.New(p, evaluatePolynomial(poly, ffarith.New(p, 2)).Int()),
		ffarith.New(p, evaluatePolynomial(poly, ffarith.New(p, 3)).Int()),
		ffarith.New(p, evaluatePolynomial(poly, ffarith.New(p, 4)).Int()),
		ffarith.New(p, evaluatePolynomial(poly, ffarith.New(p, 5)).Int()),
	}
	fmt.Println("old shares:")
	for i, val := range g {
		fmt.Printf("Party %d: %d\n", i+1, val.Int())
	}
	newg, err := bgw.DegreeReduce(g, 2) // Reduce to degree 2
	if err != nil {
		panic(err)
	}
	fmt.Println("New reduced shares:")
	for i, val := range newg {
		fmt.Printf("Party %d: %d\n", i+1, val.Int())
		if val.Int() != evaluatePolynomial(reducedpoly, ffarith.New(p, i+1)).Int() {
			fmt.Printf("FAIL, %d expected %d, got %d\n", i+1, evaluatePolynomial(reducedpoly, ffarith.New(p, i+1)).Int(), val.Int())
			return val.Int()
		}
	}
	return 0
}
func main() {
	fmt.Println("Running test...")
	// p := 29
	// points := [][2]ffarith.Num{
	// 	{ffarith.New(p, 1), ffarith.New(p, 2)},
	// 	{ffarith.New(p, 2), ffarith.New(p, 0)},
	// 	{ffarith.New(p, 3), ffarith.New(p, 3)},
	// }

	// secret, err := bgw.ReconstructSecret(points)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Secret is:", secret.Int())
	// val := testShare()
	// if val != 0 {
	// 	fmt.Printf("FAIL bad secret %d \n", val)
	// 	return
	// }
	// val = testDegreeReduce()
	// if val != 0 {
	// 	fmt.Printf("FAIL bad expected %d \n", val)
	// 	return
	// }
	val := testMoveValidation()
	if val != 0 {
		fmt.Printf("FAIL bad expected %d \n", val)
		return
	}
	fmt.Println("PASS")
}
