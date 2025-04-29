package main

import (
	"encoding/gob"
	"log"
	"net"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/fprasx/secrets-and-spies/service"
	"github.com/fprasx/secrets-and-spies/ui/menu"
)

func init() {
	gob.Register(&net.UnixAddr{})
	gob.Register(&net.TCPAddr{})
	gob.Register(&net.UDPAddr{})
	gob.Register(&net.IPAddr{})

	f, err := tea.LogToFile("debug.log", "debug")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
}

func main() {
	menu.Show()

	if !menu.Host {
		menu.ShowJoin()
	}

	srv := service.
		New(menu.Name, menu.Address).
		WithHost(menu.Host).
		Join(menu.HostAddress)

	for {
		_ = srv
	}
}
