package bgw

import ffarith "github.com/fprasx/secrets-and-spies/internal/ff_arith"

func multiplyMatrices(A, B [][]ffarith.FFNum) [][]ffarith.FFNum {
	n := len(A)
	m := len(B[0])
	p := A[0][0].P()

	result := make([][]ffarith.FFNum, n)
	for i := 0; i < n; i++ {
		result[i] = make([]ffarith.FFNum, m)
		for j := 0; j < m; j++ {
			sum := ffarith.NewFFNum(p, 0)
			for k := 0; k < len(B); k++ {
				sum = sum.Plus(A[i][k].Times(B[k][j]))
			}
			result[i][j] = sum
		}
	}
	return result
}
func invertMatrix(A [][]ffarith.FFNum) ([][]ffarith.FFNum, error) {
	n := len(A)
	p := A[0][0].P()

	// Create augmented matrix (A | I)
	M := make([][]ffarith.FFNum, n)
	for i := 0; i < n; i++ {
		M[i] = make([]ffarith.FFNum, 2*n)
		for j := 0; j < n; j++ {
			M[i][j] = A[i][j]
		}
		for j := n; j < 2*n; j++ {
			if j-n == i {
				M[i][j] = ffarith.NewFFNum(p, 1)
			} else {
				M[i][j] = ffarith.NewFFNum(p, 0)
			}
		}
	}

	// Forward elimination
	for i := 0; i < n; i++ {
		// Find pivot
		if M[i][i].Int() == 0 {
			for k := i + 1; k < n; k++ {
				if M[k][i].Int() != 0 {
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
	invA := make([][]ffarith.FFNum, n)
	for i := 0; i < n; i++ {
		invA[i] = make([]ffarith.FFNum, n)
		for j := 0; j < n; j++ {
			invA[i][j] = M[i][j+n]
		}
	}

	return invA, nil
}
