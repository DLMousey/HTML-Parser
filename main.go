package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Tag struct {
	Name string
	LineNumber int
	TagStartColumn int
	TagEndColumn int
	IsOpening bool
	IsClosing bool
}

func main() {
	start := time.Now()
	html, err := os.Open("index.html")

	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(html)
	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	var tags []Tag
	for lineNumber, line := range lines {
		lineSlice := strings.Split(line, "")

		// Iterate over characters in line
		for column, char := range lineSlice {

			// If we find an opening tag marker
			if char == "<" {

				// Set up a few tracking variables for the loop
				locatedTag := ""
				locatedEnd := false
				lastIndex := column

				// Do a separate iteration moving up the indexes (without moving lineNumber)
				// to find the closing brace. @TODO - move onto next line if closing brace is not on the same line
				for locatedEnd == false {
					indexSafe := (lastIndex + 1) < len(lineSlice)

					if !indexSafe {
						panic("Lookahead would go past array boundaries on line " + strconv.Itoa(lineNumber) + " while reading line; " + line)
					}

					isClosingBrace := lineSlice[lastIndex + 1] == ">"
					isAttributeSpace := lineSlice[lastIndex + 1] == " "

					if !isClosingBrace && !isAttributeSpace {
						locatedTag += lineSlice[lastIndex + 1]
					} else {
						locatedEnd = true
					}

					lastIndex++
				}

				// We now know what tag we've found and what it's line number, start column and end column are.
				// @todo find a better way of determining if it's a closing tag, this is gross
				isOpening := !strings.Contains(locatedTag, "/")
				tags = append(tags, Tag {
					Name: locatedTag,
					LineNumber: lineNumber,
					TagStartColumn: column + 1, // +1 due to columns not being 0 indexed
					TagEndColumn: lastIndex + 2, // +2 due to columns not being 0 indexed + lookahead in loop
					IsOpening: isOpening,
					IsClosing: !isOpening,
				})
			}
		}
	}

	// Marshall the tags slice into json, use MarshalIndent so it can be pretty printed with tabs
	out, err := json.MarshalIndent(tags, "", "\t")
	fmt.Println(string(out))

	// Print profiling information since profiling in goland doesn't seem to be a thing on linux
	executionTime := time.Since(start)
	fmt.Println("Took " + executionTime.String() + " to parse " + strconv.Itoa(len(tags)) + " tags")
}
