package osfiles

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

var dir string

func GetFiles(args []string, showDotFiles bool) ([]os.DirEntry, string) {
	if len(args) > 0 {
		dir = args[0] // Get the first non-flag argument passed to the program
	} else {
		dir = "." // Default to current directory if no non-flag argument is provided
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading directory: %v\n", err)
		os.Exit(1)
	}
	if !showDotFiles {
		var filteredFiles []os.DirEntry
		for _, file := range files {
			if !strings.HasPrefix(file.Name(), ".") {
				filteredFiles = append(filteredFiles, file)
			}
		}
		files = filteredFiles
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })

	return files, dir
}
