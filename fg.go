package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var maxLength = flag.Int("max", 160, "max length of the fortune")
var maxLines = flag.Int("lines", 15, "max lines of the fortune")
var sourceDirectory = flag.String("dir", "/usr/share/games/fortunes", "source directory of the fortunes")
var outputDirectory = flag.String("out", "", "output directory of the fortunes (defaults for dir/dat)")

func main() {
	flag.Parse()

	*sourceDirectory = "./test-fortunes"

	*sourceDirectory = strings.TrimSuffix(*sourceDirectory, "/")
	if *outputDirectory == "" {
		*outputDirectory = *sourceDirectory + "/dat"
	}

	var files []string

	// list all files in source dir
	allFiles, err := os.ReadDir(*sourceDirectory)
	if err != nil {
		fmt.Println("Error opening directory:", err)
		return
	}

	for _, f := range allFiles {
		if !f.IsDir() {
			if !isExecutable(f) && filepath.Ext(f.Name()) == "" {
				files = append(files, f.Name())
			}
		}
	}

	log.Println("Found", len(files), "files")
	log.Println(files)
}

func isExecutable(fileInfo os.DirEntry) bool {
	return (fileInfo.Type().Perm() & 0111) != 0
}
