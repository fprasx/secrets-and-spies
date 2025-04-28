package bgw

import (
	"testing"

	"github.com/fprasx/secrets-and-spies/ff"
)

func TestReconstructSecret(t *testing.T) {
	p := 29
	points := []Share{
		{ff.New(p, 1), ff.New(p, 10)},
		{ff.New(p, 2), ff.New(p, 21)},
		{ff.New(p, 3), ff.New(p, 9)},
	}

	secret, err := ReconstructSecret(points)
	if err != nil {
		t.Fatalf("ReconstructSecret failed: %v", err)
	}

	if secret.Int() != 5 {
		t.Fatalf("expected secret 5, got %d", secret.Int())
	}
}

func TestShareSecret(t *testing.T) {
	p := 29
	secret := ff.New(p, 28)

	shares, err := ShareSecret(secret, 10, 20)
	if err != nil {
		t.Fatalf("ShareSecret failed: %v", err)
	}

	recon, err := ReconstructSecret(shares)
	if err != nil {
		t.Fatalf("ReconstructSecret failed: %v", err)
	}

	if recon.Int() != secret.Int() {
		t.Fatalf("expected secret %d, got %d", secret.Int(), recon.Int())
	}
}

func TestDotProduct(t *testing.T) {
	p := 29
	tVal := 2
	n := 5
	noCities := 3
	a := []ff.Num{
		ff.New(p, 0),
		ff.New(p, 1),
		ff.New(p, 0),
	}

	shares, err := ShareLocation(2, noCities, tVal, n, p)
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

	if secret.Int() != 0 {
		t.Fatalf("expected dot product 0, got %d", secret.Int())
	}
}

func TestDegreeReduce(t *testing.T) {
	p := 29
	poly := []ff.Num{
		ff.New(p, 1),
		ff.New(p, 3),
		ff.New(p, 2),
		ff.New(p, 1),
		ff.New(p, 1),
	}
	reducedPoly := []ff.Num{
		ff.New(p, 1),
		ff.New(p, 3),
		ff.New(p, 2),
	}

	// Create original shares
	g := []ff.Num{
		ff.New(p, evaluatePolynomial(poly, ff.New(p, 1)).Int()),
		ff.New(p, evaluatePolynomial(poly, ff.New(p, 2)).Int()),
		ff.New(p, evaluatePolynomial(poly, ff.New(p, 3)).Int()),
		ff.New(p, evaluatePolynomial(poly, ff.New(p, 4)).Int()),
		ff.New(p, evaluatePolynomial(poly, ff.New(p, 5)).Int()),
	}

	newg, err := DegreeReduce(g, 2) // Reduce to degree 2
	if err != nil {
		t.Fatalf("DegreeReduce failed: %v", err)
	}

	// Verify new shares match evaluation of reduced polynomial
	for i, val := range newg {
		expected := evaluatePolynomial(reducedPoly, ff.New(p, i+1)).Int()
		if val.Int() != expected {
			t.Fatalf("party %d: expected %d, got %d", i+1, expected, val.Int())
		}
	}
}
