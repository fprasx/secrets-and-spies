package bgw

import (
	"errors"
	"fmt"

	"github.com/fprasx/secrets-and-spies/ff"
)

type Share = [2]ff.Num

// evaluatePolynomial evaluates the polynomial at point x
func evaluatePolynomial(coeffs []ff.Num, x ff.Num) ff.Num {
	result := ff.New(0)
	power := ff.New(1) // x^0

	for _, coeff := range coeffs {
		result = result.Plus(coeff.Times(power))
		power = power.Times(x) // next power of x
	}

	return result
}

// assumes t-1 degree poly
func ReconstructSecret(points []Share) (ff.Num, error) {
	t := len(points)
	if t == 0 {
		return ff.Num{}, fmt.Errorf("no points provided")
	}

	secret := ff.New(0)

	for k := 0; k < t; k++ {
		lambda := ff.New(1)
		xk := points[k][0]

		for j := 0; j < t; j++ {
			if j == k {
				continue
			}
			xj := points[j][0]
			num := xj
			den := xj.Minus(xk)
			lambda = lambda.Times(num.Times(den.Inv()))
            // fmt.Printf("den: %v\n", den.BigInt())
		}

		secret = secret.Plus(lambda.Times(points[k][1]))
        // fmt.Printf("secret: %v\n", secret.BigInt())
        // fmt.Printf("close: %v\n", ff.New(5).Minus(secret).BigInt())
	}

	return secret, nil
}

// SolveLinearSystemFF solves A * x = b over a finite field
func SolveLinearSystemFF(A [][]ff.Num, b []ff.Num) ([]ff.Num, error) {
	n := len(A)

	// Augment A with b
	for i := 0; i < n; i++ {
		A[i] = append(A[i], b[i])
	}

	// Forward elimination
	for i := 0; i < n; i++ {
		fmt.Printf("\nMatrix after elimination step %d:\n", i)
		for _, row := range A {
			for _, elem := range row {
				fmt.Printf("%d ", elem.BigInt())
			}
			fmt.Println()
		}
		// Find pivot row
		pivotRow := i
		for k := i + 1; k < n; k++ {
			if A[k][i].IsNonZero() {
				pivotRow = k
				break
			}
		}
		// Swap rows if needed
		A[i], A[pivotRow] = A[pivotRow], A[i]

		// Check for singularity
		if A[i][i].IsZero() {
			return nil, errors.New("singular matrix: no solution")
		}

		// Normalize pivot row
		invPivot := A[i][i].Inv()
		for j := i; j <= n; j++ {
			A[i][j] = A[i][j].Times(invPivot)
		}

		// Eliminate below
		for k := i + 1; k < n; k++ {
			factor := A[k][i]
			for j := i; j <= n; j++ {
				A[k][j] = A[k][j].Minus(factor.Times(A[i][j]))
			}
		}

	}

	// Back substitution
	x := make([]ff.Num, n)
	for i := n - 1; i >= 0; i-- {
		x[i] = A[i][n]
		for j := i + 1; j < n; j++ {
			x[i] = x[i].Minus(A[i][j].Times(x[j]))
		}
	}

	return x, nil
}
