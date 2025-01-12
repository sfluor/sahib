package main

import (
	"fmt"
	"sahib/clients"
)

func noErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	resp, err := clients.QueryMaany("%D8%AA%D9%81%D8%A7%D9%88%D8%B6")
	noErr(err)

	for _, r := range resp {
		fmt.Printf("row : %v\n", r)
	}
}
