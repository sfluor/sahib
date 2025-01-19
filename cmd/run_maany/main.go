package main

import (
	"fmt"
	"os"
	"sahib/clients"
)

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	resp, err := clients.QueryMaany(os.Args[1])
	noErr(err)

	for _, r := range resp.List {
		fmt.Printf("row : %+v\n", r)
	}
}
