package main

import (
	"os"

	"github.com/poolpOrg/go-midicsv/encoding"
)

func main() {
	fp, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	b, err := encoding.NewDecoder(fp).Decode()
	if err != nil {
		panic(err)
	}

	fp2, err := os.Create("/tmp/foo")
	if err != nil {
		panic(err)
	}
	defer fp2.Close()
	fp2.Write(b)

}
