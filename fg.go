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
var maxLines = flag.Int("lines", 5, "max lines of the fortune")
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

	log.Println("Found", len(files), "files:", files)

	// create output directory
	err = os.MkdirAll(*outputDirectory, 0755)
	if err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}
	for _, f := range files {
		path := filepath.Join(*sourceDirectory, f)
		// read contents
		contents, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		newFile, err := filterFile(f, string(contents))
		if err != nil {
			fmt.Println("Error filtering file:", err)
			return
		}

		// write file
		err = os.WriteFile(filepath.Join(*outputDirectory, f), []byte(newFile), 0644)
		if err != nil {
			fmt.Println("Error writing file:", err)
			return
		}
	}

}

func isExecutable(fileInfo os.DirEntry) bool {
	return (fileInfo.Type().Perm() & 0111) != 0
}

const fortuneSeparator = "\n%\n"

func filterFile(fileName string, originalContent string) (string, error) {
	newFile := ""
	fortunes := strings.Split(originalContent, fortuneSeparator)
	log.Println("["+fileName+"]", "Found", len(fortunes), "fortunes in file")

	newCount := 0
	for _, line := range fortunes {
		if line == "" {
			continue
		}
		if strings.Count(line, "\n") > *maxLines {
			continue
		}
		if len(line) > *maxLength {
			continue
		}
		newFile += line + fortuneSeparator
		newCount++
	}
	newFile = strings.TrimSuffix(newFile, fortuneSeparator)
	log.Println("["+fileName+"]", "Filtered to", newCount, "fortunes")

	return newFile, nil
}
