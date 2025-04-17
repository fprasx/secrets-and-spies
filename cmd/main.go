package main

import (
	"fmt"

	"github.com/fprasx/secrets-and-spies/internal/ff_arith"
)

func main() {
	a := ffarith.NewFFNum(71, 3)

	fmt.Printf("%v\n", a.ToThe(4))
}
