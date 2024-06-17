package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
This variant prints Dir/Folder entries only
//
This program is meant to accept standard output via a pipe from the gls command.
And, it is very picky re the exact number of fields etc. per line of input. i.e. :
gls -oFGtpsrh --color=auto --time-style=full-iso --block-size=1 --group-directories-first | thisProgram
*/

// Constants for colors
const colorReset = "\033[0m"
const colorCyan = "\033[36m"
const cyanBack = "\033[46m"

// Colors for file types
var (
	colors = map[string]string{
		"/":       "\033[91m", // Slightly-brighter Red for directories
		"default": "\033[0m",  // reset
	}
)

func dirOnly(fileName string) bool {
	if strings.HasSuffix(fileName, "/") {
		// print it
		return true
	} else {
		// do not print it
		return false
	}
}

// Function to format file sizes with commas
func formatSize(size int64) string {
	str := strconv.FormatInt(size, 10)
	n := len(str)
	if n <= 3 {
		return str
	}
	start := n % 3
	if start == 0 {
		start = 3
	}
	parts := []string{str[:start]}
	for i := start; i < n; i += 3 {
		parts = append(parts, str[i:i+3])
	}
	return strings.Join(parts, ",")
}

// Function to get color (a specified color) for each type of file
func getColor(fileName string) string {
	for suffix, color := range colors {
		if strings.HasSuffix(fileName, suffix) {
			return color
		}
	}
	return colors["default"]
}

func main() {
	// Regular expression to parse the input line
	lineRegex := regexp.MustCompile(`^(\d+)\s+(\S+)\s+(\d+)\s+(\S+)\s+(\d+)\s+(\d{4})-(\d{2})-(\d{2})\s+(\d{2}:\d{2}:\d{2})+(\.\d{9})\s+([+-]\d{4})\s+(.+)$`)
	//                                      1        2      3       4       5         6 date 7     8       9 time h:m:s        10 nano s    11 offset      12 fileName
	//
	// Read from standard input (assumes specific format of input from `gls`, see prior comment).
	scanner := bufio.NewScanner(os.Stdin) // (os.Stdin) reads from standard input
	lineCount := 0
	dirLinesCounter := 0
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line) // trim spaces

		lineCount++

		if strings.HasPrefix(line, "total") {
			continue // get next line
		}

		// Parse the line using regular expression
		matches := lineRegex.FindStringSubmatch(line)
		if len(matches) != 13 {
			continue // skip lines that don't match expected format, i.e., lines with more or less than 11 fields
		}

		// Extract parts from the parsed line
		blocks := matches[0] // matches[0] is apparently some sort of line/carriage control char, discard it elegantly on next line
		blocks = matches[1]  // see above comment
		// Convert len as blocks to number of blocks (one block being 4096 bytes)
		if blocks != "0" {
			blocksAsNumber, err1 := strconv.Atoi(blocks)
			blocks = strconv.Itoa(blocksAsNumber / 4096)
			if err1 != nil {
				fmt.Println("there was an error while doing blocks conversion")
			}
		}
		permissions := matches[2]
		links := matches[3]
		owner := matches[4]
		size, err := strconv.ParseInt(matches[5], 10, 64)
		if err != nil {
			continue // skip lines where size cannot be parsed
		}
		formattedSize := formatSize(size) // create the comma-formatted string to print
		dateYear := matches[6]
		dateMonth := matches[7]
		dateDay := matches[8]
		time := matches[9]
		// nanoseconds := matches[10] // no one really wants to see nanoseconds
		timeZone := matches[11] // actually an offset, but will use this to effect a DST/STD time attribution
		fileName := matches[12]

		// Determine color based on file name suffix
		color := getColor(fileName)

		// Print formatted line with colorized file name and aligned columns
		if dirOnly(fileName) {
			dirLinesCounter++
			if timeZone == "-0700" {
				fmt.Printf("%7s %11s %4s %slinks%s %8s %12s %sbytes%s  %s-%s-%s %s%s %s%s%s\n",
					blocks, permissions, links, colorCyan, colorReset, owner, formattedSize, colorCyan, colorReset, dateDay, dateMonth, dateYear, time, "-dst",
					color, fileName, "\033[0m", // reset color
				)
			} else {
				fmt.Printf("%7s %11s %4s %slinks%s %8s %12s %sbytes%s  %s-%s-%s %s%s %s%s%s\n",
					blocks, permissions, links, colorCyan, colorReset, owner, formattedSize, colorCyan, colorReset, dateDay, dateMonth, dateYear, time, "-std",
					color, fileName, "\033[0m", // reset color
				)
			}

			if lineCount > 36 {
				lineCount = 0
				fmt.Printf("%s  x4096  attributes     links    owner          size        date         time        fileName%s\n", cyanBack, colorReset)
			}
		}
	}
	if dirLinesCounter < 1 {
		fmt.Println("The current dir has no sub directories")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	// fmt.Printf("%s  x4096  attributes     links    owner          size        date         time        fileName%s\n", cyanBack, colorReset)
}
