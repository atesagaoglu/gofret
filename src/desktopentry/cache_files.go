package desktopentry

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
	"reflect"
)

// This only have one field for now, but it will be extended
type DesktopEntry struct {
	Path string
}

// Just to see the entries better
func PrintEntries(entries []DesktopEntry) {
	for _, entry := range entries {
		log.Println(entry.Path)
	}
}

func CacheEntries() {
	// first two paths are system-wide, last one is user-specific
	var paths = [3]string{
		"/usr/share/applications",
		"/usr/local/share/applications",
		"~/.local/share/applications",
	}

	var entries []DesktopEntry

	// add .desktop files from each path to entries
	for _, path := range paths {
		dir, err := os.ReadDir(path)
		if err != nil {
			if reflect.TypeOf(err) == reflect.TypeOf(&os.PathError{}) {
				// we can ignore if the directory does not exist
				log.Println(path, " does not exist")
			} else {
				log.Fatal(err)
			}

		}

		// concatinate with the path to get absolute path
		for _, entry := range dir {
			entries = append(entries, DesktopEntry{
				Path: path + "/" + entry.Name(),
			})
		}

	}

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
