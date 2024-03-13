package main

import (
	// "log"

	"github.com/atesagaoglu/gofret/src/desktopentry"
)

func main() {
	desktopentry.CacheEntries()
	
	/* For now, don't read from the cache

	entries, err := desktopentry.ReadCache()
	if err != nil {
		log.Fatal(err)
	} else {
		desktopentry.PrintEntries(entries)
	}
	*/

}
