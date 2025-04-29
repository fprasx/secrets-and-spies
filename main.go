package main

import (
	"github.com/fprasx/secrets-and-spies/ui/lobby"
	"github.com/fprasx/secrets-and-spies/ui/menu"
	"github.com/fprasx/secrets-and-spies/utils"
)

func init() {
	utils.RegisterRpcTypes()
	utils.RegisterLogger()
}

func main() {
	menu.Show()
	lobby.Show()
}
