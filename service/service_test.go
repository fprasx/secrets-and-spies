package service

import (
	"testing"
	"time"

	"github.com/fprasx/secrets-and-spies/utils"
)

func TestBoth(t *testing.T) {
	// Create two Spies objects
	utils.RegisterRpcTypes()
	s1 := New("unix:/tmp/spies/a").WithName("s1").WithHost(true)
	s2 := New("unix:/tmp/spies/b").WithName("s2").WithHost(false)
	s2.Join("unix:/tmp/spies/a")
	//	go s1.PlayGame()
	//go s2.PlayGame()
	for {
		s1.Lock()
		if len(s1.peers) == 2 {
			s1.Unlock()
			break

		}
		s1.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
	s1.HostStart()
	for {
		time.Sleep(10 * time.Millisecond)
	}
}
func Test1(t *testing.T) {
	// Create two Spies objects
	utils.RegisterRpcTypes()
	s1 := New("unix:/tmp/spies/a").WithName("s1").WithHost(true)

	//go s1.PlayGame()

	for {
		s1.Lock()
		if len(s1.peers) == 2 {
			s1.Unlock()
			break

		}
		s1.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
	s1.HostStart()
	s1.PlayGame()
	for {
		time.Sleep(10 * time.Millisecond)
	}
}

func Test2(t *testing.T) {
	utils.RegisterRpcTypes()
	s2 := New("unix:/tmp/spies/b").WithName("s2").WithHost(false)
	go s2.PlayGame()
	s2.Join("unix:/tmp/spies/a")

	for {
		time.Sleep(10 * time.Millisecond)
	}
}
