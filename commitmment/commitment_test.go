package commitmment_test

import (
	"testing"

	"github.com/fprasx/secrets-and-spies/commitmment"
	"github.com/fprasx/secrets-and-spies/utils"
)

func TestSuccess(t *testing.T) {
	message := uint(0xdecaf)
	commitment, nonce := commitmment.Commit(message)
	utils.Assert(
		commitmment.Verify(commitment, message, nonce),
		"failed to verify commitmment",
	)
}

func TestFail(t *testing.T) {
	message := uint(0xdecaf)
	commitment, nonce := commitmment.Commit(message)
	utils.Assert(
		!commitmment.Verify(commitment, 0x4311, nonce),
		"commitmment should not have verified",
	)
}
