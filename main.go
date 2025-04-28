package main

import (
	"github.com/charmbracelet/log"
)

func main() {
	log.SetReportCaller(true)
	log.Info("Hello, world!")
}
