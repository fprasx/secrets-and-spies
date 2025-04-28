package main

import (
	"github.com/charmbracelet/log"
)

func main() {
	log.SetReportCaller(true)
	log.SetReportTimestamp(false)

	log.Info("Hello, world!")
}
