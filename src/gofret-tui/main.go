package main

import (
	"fmt"
	"log"

	"github.com/atesagaoglu/gofret/src/desktopentry"
)

func main() {
	fmt.Println("Hello, World!")
	desktopentry.CacheEntries()
	entries, err := desktopentry.ReadCache()
	if err != nil {
		log.Fatal(err)
	} else {
		desktopentry.PrintEntries(entries)
	}

}
