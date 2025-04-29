package commitmment_test

import (
	"testing"

	"github.com/fprasx/secrets-and-spies/commitmment"
	"github.com/fprasx/secrets-and-spies/utils"
)

func Test1(t *testing.T) {
	message := uint(0xdecaf)
	commitment, nonce := commitmment.Commit(message)
	utils.Assert(
		commitmment.Verify(commitment, message, nonce),
		"failed to verify commitmment",
	)
}
