package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Constants for colors
const colorReset = "\033[0m"
const colorRed = "\033[31m"
const colorGreen = "\033[32m"
const colorCyan = "\033[36m"
const colorPurple = "\033[35m"
const colorYellow = "\033[33m"

// Colors for file types
var (
	colors = map[string]string{
		".tiff":   "\033[92m", // Slightly-brighter Green
		"/":       "\033[91m", // Slightly-brighter Red for directories
		".jpeg":   "\033[32m", // Green
		".jpg":    "\033[35m", // Magenta
		".txt":    "\033[36m", // Cyan
		".go":     "\033[32m", // Green
		".py":     "\033[34m", // Blue
		".sh":     "\033[33m", // Gold
		".pages":  "\033[96m", // Light Cyan
		"default": "\033[0m",  // reset
	}
)

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

// Function to get color for the file name
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
	lineRegex := regexp.MustCompile(`^(\d+)\s+(\S+)\s+(\d+)\s+(\S+)\s+(\d+)\s+(\d{4}-\d{2}-\d{2})\s+(\d{2}:\d{2}:\d{2})+(\.\d{9})\s+([+-]\d{4})\s+(.+)$`)

	// Read lines from standard input (assumes input from `ls` or `gls`)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line) // trim spaces

		if strings.HasPrefix(line, "total") {
			fmt.Println(line) // print "total" line as-is
			continue
		}

		// Parse the line using regular expression
		matches := lineRegex.FindStringSubmatch(line)
		if len(matches) != 11 {
			continue // skip lines that don't match expected format
		}

		// Extract parts from the parsed line
		blocks := matches[1]
		permissions := matches[2]
		links := matches[3]
		owner := matches[4]
		size, err := strconv.ParseInt(matches[5], 10, 64)
		if err != nil {
			continue // skip lines where size cannot be parsed
		}
		formattedSize := formatSize(size) // this is the string to print
		date := matches[6]
		time := matches[7]
		nanoseconds := matches[8]
		timeZone := matches[9]
		fileName := matches[10]

		// Determine color based on file name suffix
		color := getColor(fileName)

		if nanoseconds == "had to use it" {
		}

		// Print formatted line with colorized file name and aligned columns
		if timeZone == "-0700" {
			fmt.Printf("%10s %11s %4s %slinks%s %8s %12s %sbytes%s  %s %s %s %s%s%s\n",
				blocks, permissions, links, colorCyan, colorReset, owner, formattedSize, colorCyan, colorReset, date, time, "dst",
				color, fileName, "\033[0m", // reset color
			)
		} else {
			fmt.Printf("%10s %11s %4s %slinks%s %8s %12s %sbytes%s  %s %s %s %s%s%s\n",
				blocks, permissions, links, colorCyan, colorReset, owner, formattedSize, colorCyan, colorReset, date, time, "std",
				color, fileName, "\033[0m", // reset color
			)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
