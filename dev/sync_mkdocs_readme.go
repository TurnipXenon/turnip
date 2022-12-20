// sync_mkdocs_readme.go copies the mkdocs index.md file and overwrites over the README.md. This script
// handles transforming the relative references into the correct references.

package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func main() {
	// reference: nishanthshanmugham @ https://stackoverflow.com/a/74107689/17836168
	sourceFile, err := os.Open("./docs/docs/index.md")
	if err != nil {
		log.Fatalln(err)
	}
	defer func(sourceFile *os.File) {
		_ = sourceFile.Close()
	}(sourceFile) // ok to ignore error: file was opened read-only.

	destinationFile, err := os.Create("./README.md")
	if err != nil {
		log.Fatalln(err)
	}

	dataWriter := bufio.NewWriter(destinationFile)
	defer func() {
		err := dataWriter.Flush()
		if err != nil {
			log.Printf("Error flushing writer: %s\n", err)
		}

		c := destinationFile.Close()
		// Report the error from Close, if any.
		// But do so only if there isn't already
		// an outgoing error.
		if c != nil && err == nil {
			err = c
		}
	}()

	// reference: Roy Lee @ https://stackoverflow.com/a/36111861/17836168
	// todo: replace image references properly
	dataReader := bufio.NewScanner(sourceFile)
	for dataReader.Scan() {
		line := dataReader.Text()

		// transform references
		if strings.Contains(line, "](./") {
			line = strings.Replace(line, "](./", "](./docs/docs/", 1)
		}

		_, _ = dataWriter.WriteString(line + "\n")
	}
}
