package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// build as ricks.colors.and.commas to be used in .bash_profile

/*
This program is meant to accept standard output via a pipe from the gls command.
And, it is very picky re the exact number of fields etc. per line of input. i.e. :
gls -oFGtpsrh --color=auto --time-style=full-iso --block-size=1 --group-directories-first | thisProgram
*/

// Constants for colors
const colorReset = "\033[0m"
const colorRed = "\033[31m"
const colorGreen = "\033[32m"
const colorCyan = "\033[36m"
const colorPurple = "\033[35m"
const colorYellow = "\033[33m"
const cyanBack = "\033[46m"

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

func printFileInfo(path string) {
	fi, err := os.Stat(path) // fi being a bunch of file information codes
	// fmt.Printf("fi is %s\n", fi)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// used to formulate creationTime -> birth
	stat := fi.Sys().(*syscall.Stat_t) // stat being a bunch [slice/map or what have you] of numbers
	// fmt.Printf("stat is %d\n", stat)

	// birthdate and time :: [from:  fi.Sys().(*syscall.Stat_t)  ]
	creationTime := time.Unix(stat.Birthtimespec.Sec, stat.Birthtimespec.Nsec)
	// creationTime being info such as :: 2024-06-18 09:59:58.624201859 -0700 PDT // note the fancy inclusion of PDT
	// fmt.Printf("creationTime is %s\n", creationTime)

	// permissions changed
	lastStatusChangeTime := time.Unix(stat.Ctimespec.Sec, stat.Ctimespec.Nsec)

	// fmt.Println("File:", path, "  Size:", fi.Size(), "bytes")
	if creationTime.Format(time.RFC1123) == fi.ModTime().Format(time.RFC1123) {
		// do not print redundant creation and modification info if the gls has already provided same
	} else {
		fmt.Println("                                   originally created:", creationTime.Format(time.RFC1123))
		fmt.Printf("                                 touched, or %smodified: %s%s\n", colorCyan, fi.ModTime().Format(time.RFC1123), colorReset)
		if fi.ModTime().Format(time.RFC1123) == lastStatusChangeTime.Format(time.RFC1123) {
			// if redundant, do not print lastStatusChangeTime
		} else {
			fmt.Printf("                           %spermissions|status changed: %s%s\n", colorRed, lastStatusChangeTime.Format(time.RFC1123), colorReset)
		}
	}
	fmt.Println()
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
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line) // trim spaces

		lineCount++

		if strings.HasPrefix(line, "total") {
			// fmt.Println(line) // print "total" line as-is
			fmt.Printf("%s  x4096  attributes     links    owner          size        date         time        fileName  %s\n", cyanBack, colorReset)
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
		dateMonth = translateMonth(dateMonth) // change numeric month to Jan-Dec
		dateDay := matches[8]
		dateDay = replaceLeading0withSpace(dateDay) // strip leading 0s from single digit days
		time := matches[9]
		// nanoseconds := matches[10] // no one really wants to see nanoseconds
		timeZone := matches[11] // actually an offset, but will use this to effect a DST/STD time attribution
		fileName := matches[12]

		// Determine color based on file name suffix
		color := getColor(fileName)

		// Print formatted line with colorized file name and aligned columns
		if timeZone == "-0700" {
			fmt.Printf("%7s %11s %4s %slinks%s %8s %12s %sbytes%s  %s-%s-%s %s%s %s%s%s\n",
				blocks, permissions, links, colorCyan, colorReset, owner, formattedSize, colorCyan, colorReset, dateDay, dateMonth, dateYear, time, "-dst",
				color, fileName, "\033[0m", // reset color
			)
			printFileInfo(fileName)
		} else {
			fmt.Printf("%7s %11s %4s %slinks%s %8s %12s %sbytes%s  %s-%s-%s %s%s %s%s%s\n",
				blocks, permissions, links, colorCyan, colorReset, owner, formattedSize, colorCyan, colorReset, dateDay, dateMonth, dateYear, time, "-std",
				color, fileName, "\033[0m", // reset color
			)
			printFileInfo(fileName)
		}

		if lineCount > 36 {
			lineCount = 0
			fmt.Printf("%s  x4096  attributes     links    owner          size        date         time        fileName  %s\n", cyanBack, colorReset)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	fmt.Printf("%s  x4096  attributes     links    owner          size        date         time        fileName  %s\n", cyanBack, colorReset)
}

func replaceLeading0withSpace(dateDay string) string {
	if dateDay == "01" {
		dateDay = " 1"
	}
	if dateDay == "02" {
		dateDay = " 2"
	}
	if dateDay == "03" {
		dateDay = " 3"
	}
	if dateDay == "04" {
		dateDay = " 4"
	}
	if dateDay == "05" {
		dateDay = " 5"
	}
	if dateDay == "06" {
		dateDay = " 6"
	}
	if dateDay == "07" {
		dateDay = " 7"
	}
	if dateDay == "08" {
		dateDay = " 8"
	}
	if dateDay == "09" {
		dateDay = " 9"
	} else {
		// return dateDay unchanged
	}

	return dateDay
}

func translateMonth(dateM string) string {
	if dateM == "01" {
		dateM = "Jan"
	}
	if dateM == "02" {
		dateM = "Feb"
	}
	if dateM == "03" {
		dateM = "Mar"
	}
	if dateM == "04" {
		dateM = "Apr"
	}
	if dateM == "05" {
		dateM = "May"
	}
	if dateM == "06" {
		dateM = "Jun"
	}
	if dateM == "07" {
		dateM = "Jul"
	}
	if dateM == "08" {
		dateM = "Aug"
	}
	if dateM == "09" {
		dateM = "Sep"
	}
	if dateM == "10" {
		dateM = "Oct"
	}
	if dateM == "11" {
		dateM = "Nov"
	}
	if dateM == "12" {
		dateM = "Dec"
	}
	return dateM
}
