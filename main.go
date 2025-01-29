/**
dd is a markdown viewer created specifically for my personal use.
Every setting is adjusted to fit my workflow.
**/

package main

import (
	"os"

	"github.com/cjaewon/dd/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
