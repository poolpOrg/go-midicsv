package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/poolpOrg/go-midicsv/encoding"
)

func main() {
	var opt_file string

	flag.StringVar(&opt_file, "out", "", "output file (.csv)")
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("need a source file to process")
	}

	ifp, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer ifp.Close()

	records, err := encoding.NewEncoder(ifp).Encode()
	if err != nil {
		log.Fatal(err)
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
	for _, record := range records {
		fmt.Fprintf(ofp, "%s\n", record)
	}
	os.Exit(0)
}
