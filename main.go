package main

import (
	"io/ioutil"
	"lua-vm/binchunk"
	"os"
)

func main() {
	var filename string
	if len(os.Args) != 1 {
		filename = os.Args[1]
	} else {
		filename = "luac.out"
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	binchunk.List(binchunk.Undump(data))
}
