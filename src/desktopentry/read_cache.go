package desktopentry

import(
	"os"
	"encoding/gob"
	"bytes"
)

func ReadCache() ([]DesktopEntry, error) {

	var entries []DesktopEntry

	homeDir, err := os.UserHomeDir()
	if err != nil {
		// log.Fatal(err)
		return nil,err
	}

	var cachePath = homeDir + "/.cache/gofret"

	data, err := os.ReadFile(cachePath)
	if err != nil {
		// log.Fatal(err)
		return nil,err
	}

	var buff = bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buff)

	err = decoder.Decode(&entries)
	if err != nil {
		// log.Fatal(err)
		return nil,err
	}

	return entries,nil
}