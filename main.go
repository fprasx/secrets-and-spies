package main

import (
	"fmt"

	"github.com/fprasx/secrets-and-spies/bgw"
	"github.com/fprasx/secrets-and-spies/ff"
)

// test for reconstruction
func test1() int {
	p := 29
	points := [][2]ff.Num{
		{ff.New(p, 1), ff.New(p, 10)},
		{ff.New(p, 2), ff.New(p, 21)},
		{ff.New(p, 3), ff.New(p, 9)},
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
	secret := ff.New(p, 28)
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
func evaluatePolynomial(coeffs []ff.Num, x ff.Num) ff.Num {
	p := x.P()
	result := ff.New(p, 0)
	power := ff.New(p, 1) // x^0

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
	a := []ff.Num{
		ff.New(p, 0),
		ff.New(p, 1),
		ff.New(p, 0),
	}
	shares, err := bgw.ShareLocation(2, noCities, t, n, p)
	if err != nil {
		panic(err)
	}
	newShares := make([][2]ff.Num, n)
	for i := 0; i < n; i++ {
		column := make([][2]ff.Num, noCities)
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
	poly := []ff.Num{
		ff.New(p, 1),
		ff.New(p, 3),
		ff.New(p, 2),
		ff.New(p, 1),
		ff.New(p, 1),
	}
	reducedpoly := []ff.Num{
		ff.New(p, 1),
		ff.New(p, 3),
		ff.New(p, 2),
	}
	g := []ff.Num{
		ff.New(p, evaluatePolynomial(poly, ff.New(p, 1)).Int()),
		ff.New(p, evaluatePolynomial(poly, ff.New(p, 2)).Int()),
		ff.New(p, evaluatePolynomial(poly, ff.New(p, 3)).Int()),
		ff.New(p, evaluatePolynomial(poly, ff.New(p, 4)).Int()),
		ff.New(p, evaluatePolynomial(poly, ff.New(p, 5)).Int()),
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
		if val.Int() != evaluatePolynomial(reducedpoly, ff.New(p, i+1)).Int() {
			fmt.Printf("FAIL, %d expected %d, got %d\n", i+1, evaluatePolynomial(reducedpoly, ff.New(p, i+1)).Int(), val.Int())
			return val.Int()
		}
	}
	return 0
}
func main() {
	fmt.Println("Running test...")
	// p := 29
	// points := [][2]ffarith.FFNum{
	// 	{ffarith.NewFFNum(p, 1), ffarith.NewFFNum(p, 2)},
	// 	{ffarith.NewFFNum(p, 2), ffarith.NewFFNum(p, 0)},
	// 	{ffarith.NewFFNum(p, 3), ffarith.NewFFNum(p, 3)},
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
	val := testDotProduct()
	if val != 0 {
		fmt.Printf("FAIL bad expected %d \n", val)
		return
	}
	fmt.Println("PASS")
}
