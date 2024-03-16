package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/atesagaoglu/gofret/src/desktopentry"
)

func main() {

	start := time.Now()
	entries, _ := desktopentry.CacheEntries()
	fmt.Println("Caching took: ", time.Since(start))
	fmt.Println("Found ", len(entries), " entries across the system.")
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter app name: ")
	scanner.Scan()
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}
	name := scanner.Text()

	for _, entry := range entries {
		if entry.Name == name {
			entry.Run()
		}
	}

}
