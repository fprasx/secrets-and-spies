package bgw

import (
	"github.com/fprasx/secrets-and-spies/ff"
)

func multiplyMatrices(A, B [][]ff.Num) [][]ff.Num {
	n := len(A)
	m := len(B[0])

	result := make([][]ff.Num, n)
	for i := 0; i < n; i++ {
		result[i] = make([]ff.Num, m)
		for j := 0; j < m; j++ {
			sum := ff.New(0)
			for k := 0; k < len(B); k++ {
				sum = sum.Plus(A[i][k].Times(B[k][j]))
			}
			result[i][j] = sum
		}
	}
	return result
}
func invertMatrix(A [][]ff.Num) ([][]ff.Num, error) {
	n := len(A)

	// Create augmented matrix (A | I)
	M := make([][]ff.Num, n)
	for i := 0; i < n; i++ {
		M[i] = make([]ff.Num, 2*n)
		for j := 0; j < n; j++ {
			M[i][j] = A[i][j]
		}
		for j := n; j < 2*n; j++ {
			if j-n == i {
				M[i][j] = ff.New(1)
			} else {
				M[i][j] = ff.New(0)
			}
		}
	}

	// Forward elimination
	for i := 0; i < n; i++ {
		// Find pivot
		if M[i][i].IsZero() {
			for k := i + 1; k < n; k++ {
				if M[k][i].IsNonZero() {
					M[i], M[k] = M[k], M[i]
					break
				}
			}
		}

		// Normalize pivot row
		inv := M[i][i].Inv()
		for j := 0; j < 2*n; j++ {
			M[i][j] = M[i][j].Times(inv)
		}

		// Eliminate below and above
		for k := 0; k < n; k++ {
			if k == i {
				continue
			}
			factor := M[k][i]
			for j := 0; j < 2*n; j++ {
				M[k][j] = M[k][j].Minus(factor.Times(M[i][j]))
			}
		}
	}

	// Extract inverse matrix
	invA := make([][]ff.Num, n)
	for i := 0; i < n; i++ {
		invA[i] = make([]ff.Num, n)
		for j := 0; j < n; j++ {
			invA[i][j] = M[i][j+n]
		}
	}

	return invA, nil
}
