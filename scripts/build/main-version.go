package main

import (
	"fmt"
	packageVersion "pipe-for-parallel/version"
)

func main() {
	fmt.Printf("%v\n", packageVersion.Version)
}
