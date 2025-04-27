package bgw

import (
	"crypto/rand"
	"fmt"
	"math/big"

	ffarith "github.com/fprasx/secrets-and-spies/internal/ff_arith"
)

// ShareSecret constructs shares according to Shamir's algorithm
func ShareSecret(secret ffarith.FFNum, t int, n int) ([][2]ffarith.FFNum, error) {
	if t > n {
		return nil, fmt.Errorf("threshold cannot be greater than number of parties")
	}

	p := secret.P()

	// Random coefficients: a_1 to a_{t-1}
	coeffs := make([]ffarith.FFNum, t)
	coeffs[0] = secret // constant term is the secret

	for i := 1; i < t; i++ {
		randVal, err := rand.Int(rand.Reader, big.NewInt(int64(p)))
		if err != nil {
			return nil, err
		}
		coeffs[i] = ffarith.NewFFNum(p, int(randVal.Int64()))
	}

	// Generate shares (i, f(i))
	shares := make([][2]ffarith.FFNum, n)
	for i := 1; i <= n; i++ {
		x := ffarith.NewFFNum(p, i)
		y := evaluatePolynomial(coeffs, x)
		shares[i-1] = [2]ffarith.FFNum{x, y}
	}

	return shares, nil
}

// DegreeReduce reduces the degree of shares using Vandermonde matrices
func DegreeReduce(g []ffarith.FFNum, t int) ([]ffarith.FFNum, error) {
	n := len(g)
	p := g[0].P()

	// Build V
	V := make([][]ffarith.FFNum, n)
	for i := 0; i < n; i++ {
		V[i] = make([]ffarith.FFNum, n)
		x := ffarith.NewFFNum(p, i+1) // Party indices are 1-based
		power := ffarith.NewFFNum(p, 1)
		for j := 0; j < n; j++ {
			V[i][j] = power
			power = power.Times(x)
		}
	}

	// Build P
	P := make([][]ffarith.FFNum, n)
	for i := 0; i < n; i++ {
		P[i] = make([]ffarith.FFNum, n)
		for j := 0; j < n; j++ {
			if i == j && i <= t {
				P[i][j] = ffarith.NewFFNum(p, 1)
			} else {
				P[i][j] = ffarith.NewFFNum(p, 0)
			}
		}
	}

	// Compute V inverse
	Vinv, err := invertMatrix(V)
	if err != nil {
		return nil, fmt.Errorf("failed to invert V: %v", err)
	}

	// Compute A = V * P * Vinv
	VP := multiplyMatrices(V, P)
	A := multiplyMatrices(VP, Vinv)

	// Represent g as a column vector
	gvec := make([][]ffarith.FFNum, n)
	for i := 0; i < n; i++ {
		gvec[i] = []ffarith.FFNum{g[i]}
	}

	// Compute new_g = A * g
	newgMat := multiplyMatrices(A, gvec)

	// Extract result
	newg := make([]ffarith.FFNum, n)
	for i := 0; i < n; i++ {
		newg[i] = newgMat[i][0]
	}

	return newg, nil
}
