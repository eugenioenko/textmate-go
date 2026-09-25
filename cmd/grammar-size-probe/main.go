// Command grammar-size-probe is the minimal executable used for reproducible
// embedded-grammar binary-size measurements.
package main

import (
	"fmt"
	"os"

	"github.com/eugenioenko/textmate-go/grammars"
)

func main() {
	grammar, err := grammars.Load("source.go")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if grammar == nil {
		fmt.Fprintln(os.Stderr, "source.go is not embedded")
		os.Exit(1)
	}
	fmt.Printf("%d embedded scopes\n", len(grammars.Scopes()))
}
