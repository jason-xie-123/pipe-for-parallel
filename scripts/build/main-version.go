package main

import (
	"fmt"
	packageVersion "windows-safe-pipe/version"
)

func main() {
	fmt.Printf("%v\n", packageVersion.Version)
}
