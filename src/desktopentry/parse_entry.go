package desktopentry

import (
	"log"
	"os"
	"reflect"

	"gopkg.in/ini.v1"
)

func ParseEntries() []DesktopEntry {
	// var execs []string
	var entries []DesktopEntry
	addPaths(&entries)

	for i := range entries {
		parseEntry(&entries[i])
		entries[i].Print()
	}

	return entries
}

func parseEntry(entry *DesktopEntry) {

	file, err := ini.Load(entry.Path)
	if err != nil {
		log.Println(err, err.Error())
		// return
	}

	// log.Println(file.Section("Des	ktop Entry").Key("Exec").String())
	var isHidden bool
	var isNoDisplay bool
	hidden := file.Section("Desktop Entry").Key("Hidden").String()
	if hidden == "" {
		isHidden = false
	} else {
		isHidden = true
	}

	noDisplay := file.Section("Desktop Entry").Key("NoDisplay").String()
	if noDisplay == "" || noDisplay == "false" {
		isNoDisplay = false
	} else {
		isNoDisplay = true
	}

	isTerminal := file.Section("Desktop Entry").Key("Terminal").String()
	if isTerminal == "" || isTerminal == "false" {
		entry.Terminal = false
	} else {
		entry.Terminal = true
	}

	entry.Hide = (isHidden || isNoDisplay)
	entry.Exec = file.Section("Desktop Entry").Key("Exec").String()
	entry.Name = file.Section("Desktop Entry").Key("Name").String()
}

func addPaths(entries *[]DesktopEntry) {
	// first two paths are system-wide, last one is user-specific
	var paths = [3]string{
		"/usr/share/applications",
		"/usr/local/share/applications",
		"~/.local/share/applications",
	}

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
			*entries = append(*entries, DesktopEntry{
				Path: path + "/" + entry.Name(),
			})
		}
	}
}
