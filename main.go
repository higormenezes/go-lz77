package main

import (
	"github.com/higormenezes/lz77/lz77"
)

func main() {
	err := lz77.Compress("./lorem.txt", "./lorem-c.txt")
	if err != nil {
		panic(err)
	}

	err = lz77.Decompress("./lorem-c.txt", "./lorem-d.txt")
	if err != nil {
		panic(err)
	}
}
