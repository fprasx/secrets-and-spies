package ff_test

import (
	"testing"

	"github.com/fprasx/secrets-and-spies/ff"
	"github.com/fprasx/secrets-and-spies/utils"
)

const p = 1000000007

func TestPow(t *testing.T) {
	num := ff.New(p, 0xc0ffee)
	utils.Assert(num.Pow(0) == ff.New(p, 1), "exponentiation by 0 incorrect")
}
