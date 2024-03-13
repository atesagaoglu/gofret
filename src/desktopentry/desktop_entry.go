package desktopentry

import (
	"log"
)

// might need to populate this struct later
type DesktopEntry struct {
	Path     string
	Exec     string
	Name     string
	Hide     bool
	Terminal bool
}

func (de DesktopEntry) Print() {
	log.Printf(
		"\nPath: %s\nExec: %s\nHide: %t\nName: %s\nTerminal: %t\n\n",
		de.Path, de.Exec, de.Hide, de.Name, de.Terminal,
	)
}
