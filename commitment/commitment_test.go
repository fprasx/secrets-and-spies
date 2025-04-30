package commitment_test

import (
	"testing"

	"github.com/fprasx/secrets-and-spies/commitment"
	"github.com/fprasx/secrets-and-spies/utils"
)

func TestSuccess(t *testing.T) {
	message := uint(0xdecaf)
	commitment, nonce := commitment.Commit(message)
	utils.Assert(
		commitment.Verify(commitment, message, nonce),
		"failed to verify commitmment",
	)
}

func TestFail(t *testing.T) {
	message := uint(0xdecaf)
	commitment, nonce := commitment.Commit(message)
	utils.Assert(
		!commitment.Verify(commitment, 0x4311, nonce),
		"commitmment should not have verified",
	)
}
