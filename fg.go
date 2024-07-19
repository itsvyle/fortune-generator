package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const schemaVersion = 1

var maxLength = flag.Int("max", 160, "max length of the fortune")
var maxLines = flag.Int("lines", 5, "max lines of the fortune")
var sourceDirectory = flag.String("dir", "/usr/share/games/fortunes", "source directory of the fortunes")
var outputFile = flag.String("out", "", "output directory of the fortunes (defaults to fortunes.vyle)")

func main() {
	flag.Parse()

	*sourceDirectory = "./test-fortunes"

	*sourceDirectory = strings.TrimSuffix(*sourceDirectory, "/")
	if *outputFile == "" {
		*outputFile = *sourceDirectory + "/fortunes.vyle"
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

	output := []byte{}
	output = append(output, byte(schemaVersion))
	output = append(output, byte(len(files)))

	for i, f := range files {
		name := []byte(f)
		file := []byte{
			byte(i),
			byte(1),
			byte(len(name)),
		}
		file = append(file, name...)
		output = append(output, file...)
	}
	const empty = 10
	output = append(output, make([]byte, empty)...)

	numberOfEntriesByteLocation := len(output)
	output = append(output, make([]byte, 4)...)

	totalFortunes := int32(0)

	for fileID, f := range files {
		path := filepath.Join(*sourceDirectory, f)
		// read contents
		contents, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}

		// filter file
		lines := strings.Split(string(contents), fortuneSeparator)
		var totalOffset int32 = 0
		var sepLength int32 = int32(len(fortuneSeparator))
		fortunesCount := 0
		for _, line := range lines {
			if line == "" {
				continue
			}
			if strings.Count(line, "\n") > *maxLines {
				continue
			}
			if len(line) > *maxLength {
				continue
			}
			l := int32(len(line))
			output = append(output, []byte{
				byte(fileID),
				byte(totalOffset >> 24),
				byte(totalOffset >> 16),
				byte(totalOffset >> 8),
				byte(totalOffset),
				byte(l >> 24),
				byte(l >> 16),
				byte(l >> 8),
				byte(l),
				byte(0),
			}...)
			totalOffset += l + sepLength
			fortunesCount++
		}
		totalFortunes += int32(fortunesCount)
	}
	fillInt32(&output, numberOfEntriesByteLocation, totalFortunes)

	log.Println("Total fortunes:", totalFortunes)

	// Save it to a file:
	err = os.WriteFile(*outputFile, output, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

}

func fillInt32(array *[]byte, index int, value int32) {
	(*array)[index] = byte(value >> 24)
	(*array)[index+1] = byte(value >> 16)
	(*array)[index+2] = byte(value >> 8)
	(*array)[index+3] = byte(value)
}

func isExecutable(fileInfo os.DirEntry) bool {
	return (fileInfo.Type().Perm() & 0111) != 0
}
