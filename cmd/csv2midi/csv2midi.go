package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/poolpOrg/go-midicsv/encoding"
)

func main() {
	var opt_file string

	flag.StringVar(&opt_file, "out", "", "output file (.mid)")
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("need a source file to process")
	}

	ifp, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer ifp.Close()

	b, err := encoding.NewDecoder(ifp).Decode()
	if err != nil {
		panic(err)
	}

	var ofp io.Writer
	if opt_file == "" {
		ofp = os.Stdout
	} else {
		ofp, err = os.Create(opt_file)
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err = ofp.Write(b)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
