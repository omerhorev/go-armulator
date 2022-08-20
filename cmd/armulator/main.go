package main

import (
	"fmt"
	"log"
	"os"

	"github.com/omerhorev/goarmulator"
)

func main() {
	if len(os.Args) < 2 {
		MsgError("Usage: armulator <command> [args...]\n")
		return
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		MsgError("ERROR: %v\n", err)
		return
	}
	defer f.Close()

	a, err := goarmulator.NewArmulator(f)
	if err != nil {
		MsgError("ERROR: %v\n", err)
		return
	}
	defer a.Close()

	a.Log = log.Default()

	if err := a.Run(); err != nil {
		MsgError("ERROR: %v\n", err)
		return
	}
}

func MsgError(format string, args ...any) (int, error) {
	return fmt.Fprintf(os.Stderr, format, args...)
}
