package service

import (
	"testing"
)

func TestDotProductShares(t *testing.T) {
	// Create two Spies objects
	s1 := New("unix:/tmp/spies/a").WithName("s1").WithHost(true)
	s2 := New("unix:/tmp/spies/b").WithName("s2").WithHost(false)
	s2.Join("unix:/tmp/spies/a")
	s1.HostStart()

}
