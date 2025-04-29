package bgw

import (
	"fmt"
	"testing"

	"github.com/fprasx/secrets-and-spies/ff"
)

func TestMoveCheck(t *testing.T) {
	fmt.Println("Running test MoveCheck")
	tVal := 2
	n := 5
	noCities := 4
	tempshares := make([][][2]ff.Num, noCities)
	for i := 0; i < noCities; i++ {
		tempshares[i] = make([][2]ff.Num, n)
	}
	oldloc, err := ShareLocation(3, noCities, tVal, n)
	if err != nil {
		t.Fatalf("ShareLocation failed: %v", err)
	}
	newloc, err := ShareLocation(3, noCities, tVal, n)
	if err != nil {
		t.Fatalf("ShareLocation failed: %v", err)
	}
	endshares := make([][2]ff.Num, n)

	for i := 0; i < n; i++ {
		columnold := make([][2]ff.Num, noCities)
		for j := 0; j < noCities; j++ {
			columnold[j] = oldloc[j][i]
		}
		columnnew := make([][2]ff.Num, noCities)
		for j := 0; j < noCities; j++ {
			columnnew[j] = newloc[j][i]
		}
		endshares[i] = DotProductShares(columnold, columnnew, i)
	}

	res, err := ReconstructSecret(endshares)
	if err != nil {
		t.Fatalf("ReconstructSecret failed: %v", err)
	}

	if res.Neq(ff.New(1)) {
		t.Fatalf("expected result 1, got %d", res.BigInt())
	}
}

func TestMoveValidation(t *testing.T) {
	fmt.Println("Running test MoveValidation")
	tVal := 2
	n := 5
	noCities := 4
	g := [][]int{
		{1, 0, 0, 1},
		{0, 1, 1, 0},
		{0, 1, 1, 0},
		{1, 0, 0, 1},
	}
	graph := make([][]ff.Num, noCities)
	tempshares := make([][][2]ff.Num, noCities)
	for i := 0; i < noCities; i++ {
		tempshares[i] = make([][2]ff.Num, n)
		graph[i] = make([]ff.Num, noCities)
		for j := 0; j < noCities; j++ {
			graph[i][j] = ff.New(int64(g[i][j]))
		}
	}
	oldloc, err := ShareLocation(1, noCities, tVal, n)
	if err != nil {
		t.Fatalf("ShareLocation failed: %v", err)
	}
	newloc, err := ShareLocation(2, noCities, tVal, n)
	if err != nil {
		t.Fatalf("ShareLocation failed: %v", err)
	}
	endshares := make([][2]ff.Num, n)

	for i := 0; i < n; i++ {
		columnold := make([][2]ff.Num, noCities)
		for j := 0; j < noCities; j++ {
			columnold[j] = oldloc[j][i]
		}
		columnnew := make([][2]ff.Num, noCities)
		for j := 0; j < noCities; j++ {
			columnnew[j] = newloc[j][i]
		}
		tempcolumn := make([][2]ff.Num, n)
		endshares[i], tempcolumn = ValidateMoveShares(graph, noCities, columnold, columnnew, i)
		for j := 0; j < noCities; j++ {
			tempshares[j][i] = tempcolumn[j]
		}
	}

	res, err := ReconstructSecret(endshares)
	if err != nil {
		t.Fatalf("ReconstructSecret failed: %v", err)
	}

	if res.Neq(ff.New(1)) {
		t.Fatalf("expected result 1, got %d", res.BigInt())
	}
}
func TestReconstructSecret(t *testing.T) {
	points := []Share{
		{ff.New(1), ff.New(10)},
		{ff.New(2), ff.New(21)},
		{ff.New(3), ff.New(9)},
	}

	secret, err := ReconstructSecret(points)
	if err != nil {
		t.Fatalf("ReconstructSecret failed: %v", err)
	}

	if secret.Neq(ff.New(5)) {
		t.Fatalf("expected secret 5, got %d", secret.BigInt())
	}
}

func TestShareSecret(t *testing.T) {
	secret := ff.New(28)

	shares, err := ShareSecret(secret, 10, 20)
	if err != nil {
		t.Fatalf("ShareSecret failed: %v", err)
	}

	recon, err := ReconstructSecret(shares)
	if err != nil {
		t.Fatalf("ReconstructSecret failed: %v", err)
	}

	if recon.Neq(secret) {
		t.Fatalf("expected secret %d, got %d", secret.BigInt(), recon.BigInt())
	}
}

func TestDotProduct(t *testing.T) {
	tVal := 2
	n := 5
	noCities := 3
	a := []ff.Num{
		ff.New(0),
		ff.New(1),
		ff.New(0),
	}

	shares, err := ShareLocation(2, noCities, tVal, n)
	if err != nil {
		t.Fatalf("ShareLocation failed: %v", err)
	}

	newShares := make([]Share, n)
	for i := 0; i < n; i++ {
		column := make([]Share, noCities)
		for j := 0; j < noCities; j++ {
			column[j] = shares[j][i]
		}
		newShares[i] = DotProductConstant(a, column, i)
	}

	secret, err := ReconstructSecret(newShares)
	if err != nil {
		t.Fatalf("ReconstructSecret failed: %v", err)
	}

	if secret.IsNonZero() {
		t.Fatalf("expected dot product 0, got %d", secret.BigInt())
	}
}

func TestDegreeReduce(t *testing.T) {
	poly := []ff.Num{
		ff.New(1),
		ff.New(3),
		ff.New(2),
		ff.New(1),
		ff.New(1),
	}
	reducedPoly := []ff.Num{
		ff.New(1),
		ff.New(3),
		ff.New(2),
	}

	// Create original shares
	g := []ff.Num{
		evaluatePolynomial(poly, ff.New(1)),
		evaluatePolynomial(poly, ff.New(2)),
		evaluatePolynomial(poly, ff.New(3)),
		evaluatePolynomial(poly, ff.New(4)),
		evaluatePolynomial(poly, ff.New(5)),
	}

	newg, err := DegreeReduce(g, 2) // Reduce to degree 2
	if err != nil {
		t.Fatalf("DegreeReduce failed: %v", err)
	}

	// Verify new shares match evaluation of reduced polynomial
	for i, val := range newg {
		expected := evaluatePolynomial(reducedPoly, ff.New(int64(i+1)))
		if val.Neq(expected) {
			t.Fatalf("party %d: expected %d, got %d", i+1, expected.BigInt(), val.BigInt())
		}
	}
}
