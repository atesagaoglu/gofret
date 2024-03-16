package desktopentry

import (
	"log"
	"os"
	"os/exec"
	"strings"
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

func (de DesktopEntry) Run() {

	if de.Terminal {

		envTerm := os.Getenv("TERM")
		if envTerm == "" {
			// if a terminal is not specified, use xterm
			envTerm = "xterm"
		}

		execField := de.Exec
		execField = envTerm + " -e " + execField
		parts := strings.Fields(execField)
		var args []string

		// remove freedesktop reletad arguments
		for _, part := range parts {
			if strings.HasPrefix(part, "%") {
				continue
			}
			args = append(args, part)
		}

		exec.Command(args[0], args[1:]...).Run()

	} else {
		parts := strings.Fields(de.Exec)
		var args []string

		// remove freedesktop reletad arguments
		for _, part := range parts {
			if strings.HasPrefix(part, "%") {
				continue
			}
			args = append(args, part)
		}

		exec.Command(args[0], args[1:]...).Run()
	}
}