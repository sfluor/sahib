package main

import (
	"fmt"
	"os"
	"sahib/clients"
)

func failIf(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
    client, err := clients.NewHansWehrClient(os.Args[1])
    failIf(err)

    entries, err := client.Query(os.Args[2])
    failIf(err)

    for _, e := range entries.Entries {
        fmt.Printf("entry: %+v\n", e)
    }

}
