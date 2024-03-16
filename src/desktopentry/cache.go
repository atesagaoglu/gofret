package desktopentry

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"sync"
	"time"

	"gopkg.in/ini.v1"
)

type Cache struct {
	Entries     []DesktopEntry
	LatestCache time.Time
}

func (c Cache) write() {
	log.Println("writing")
	var buff bytes.Buffer

	encoder := gob.NewEncoder(&buff)

	err := encoder.Encode(c)
	if err != nil {
		log.Fatal(err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	cacheFile, err := os.Create(homeDir + "/.cache/gofret")
	if err != nil {
		log.Fatal(err)
	}

	cacheFile.Write(buff.Bytes())
}

func (c Cache) read() (Cache, error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Cache{}, err
	}

	var cachePath = homeDir + "/.cache/gofret"

	data, err := os.ReadFile(cachePath)
	if err != nil {
		return Cache{}, err
	}

	var buff = bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buff)

	err = decoder.Decode(&c)
	if err != nil {
		return Cache{}, err
	}

	return c, nil
}

// gets last modification time of the directories
// it compares dates from all paths and returns the most recent one
func lastMod() time.Time {
	var paths = [3]string{
		"/usr/share/applications",
		"/usr/local/share/applications",
		"~/.local/share/applications",
	}

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	var times []time.Time

	for _, path := range paths {
		wg.Add(1)

		go func() {
			defer wg.Done()

			info, err := os.Lstat(path)
			if err != nil {
				log.Println(err)
			} else {
				mu.Lock()
				times = append(times, info.ModTime())
				mu.Unlock()
			}

		}()
	}

	wg.Wait()
	log.Println("latest mod time: ", maxTime(times))
	return maxTime(times)
}

// helper function to get the most recent date in a slice
func maxTime(t []time.Time) time.Time {
	if len(t) == 0 {
		return time.Time{}
	}
	latest := t[0]
	for _, t := range t[1:] {
		if t.After(latest) {
			latest = t
		}
	}
	return latest
}

// checks cache writing/reading sequence
func CacheEntries() ([]DesktopEntry, error) {
	var cache Cache
	// initial cache read
	cache, err := cache.read()
	if err != nil {
		log.Println(err)
	}

	// if any of the 3 paths has been modified cache them again
	if cache.LatestCache.Before(lastMod()) {
		cache, err = cache.createCache()
	} else {
		cache, err = cache.read()
	}

	return cache.Entries, err
}

// reads every .desktop file concurrently
func (c Cache) createCache() (Cache, error) {
	var paths = [3]string{
		"/usr/share/applications",
		"/usr/local/share/applications",
		"~/.local/share/applications",
	}
	
	// for outer goroutine
	var entries []DesktopEntry
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	// check each directory concurrently
	for _, path := range paths {
		wg.Add(1)

		go func() {
			defer wg.Done()

			dir, err := os.ReadDir(path)
			if err != nil {
				if reflect.TypeOf(err) == reflect.TypeOf(&os.PathError{}) {
					// we can ignore if the directory does not exist
					log.Println(path, " does not exist")
				} else {
					log.Fatal(err)
				}
			}

			// these are for the inner goroutine
			var entriesin []DesktopEntry
			wgin := sync.WaitGroup{}
			muin := sync.Mutex{}

			// read each .desktop file inside a directory concurrently
			for _, dirEntry := range dir {
				wgin.Add(1)

				go func() {
					defer wgin.Done()

					var entry DesktopEntry
					entry.Path = path + "/" + dirEntry.Name()

					file, err := ini.Load(entry.Path)
					if err != nil {
						log.Println(err, err.Error())
						// return
					}

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

					entry.Print()
					muin.Lock()
					entriesin = append(entriesin, entry)
					muin.Unlock()
				}()

			}
			wgin.Wait()

			mu.Lock()
			entries = append(entries, entriesin...)
			mu.Unlock()

		}()

	}
	wg.Wait()

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	var newcache = Cache{
		Entries:     entries,
		LatestCache: lastMod(),
	}

	newcache.write()
	return newcache, nil
}

// Just to see the entries better
func PrintEntries(entries []DesktopEntry) {
	for _, entry := range entries {
		fmt.Println(entry.Path)
	}
}