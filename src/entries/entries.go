package entries

import (
	"fmt"
	"log"
	"os"
	"reflect"
)

// func GetFiles() ([]string, error) {
func GetFiles() {
	// first two paths are system-wide, last one is user-specific
	var paths = [3]string{
		"/usr/share/applications",
		"/usr/local/share/applications",
		"~/.local/share/applications",
	}

	var entries []os.DirEntry

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
		entries = append(entries, dir...)
	}

	fmt.Println(entries)

}
