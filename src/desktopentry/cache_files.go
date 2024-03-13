package desktopentry

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

// Just to see the entries better
func PrintEntries(entries []DesktopEntry) {
	for _, entry := range entries {
		fmt.Println(entry.Path)
	}
}

func CacheEntries() {

	entries := ParseEntries()

	// PrintEntries(entries)
	var buff bytes.Buffer

	encoder := gob.NewEncoder(&buff)

	err := encoder.Encode(entries)
	if err != nil {
		log.Fatal(err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(homeDir)

	cacheFile, err := os.Create(homeDir + "/.cache/gofret")
	if err != nil {
		log.Fatal(err)
	}

	cacheFile.Write(buff.Bytes())
}
