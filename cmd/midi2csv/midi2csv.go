package main

import (
	"fmt"
	"os"

	"github.com/poolpOrg/go-midicsv/encoding"
)

func main() {
	fp, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	records, err := encoding.NewEncoder(fp).Encode()
	if err != nil {
		panic(err)
	}
	for _, record := range records {
		fmt.Println(record)
	}
}
