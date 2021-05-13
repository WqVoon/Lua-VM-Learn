package main

import (
	"fmt"
	"io/ioutil"
	"lua-vm/binchunk"
	"os"
)

func main() {
	var filename string
	if len(os.Args) != 1 {
		filename = os.Args[1]
	} else {
		filename = "a.out"
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	bc := binchunk.Undump(data)
	fmt.Printf("%+v", bc)
}
