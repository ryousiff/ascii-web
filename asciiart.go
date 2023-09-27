package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// WhatToPrint takes an index (i), a word, a filename, and a result string.
// It retrieves the ASCII representation of the word from the specified filename,
// based on the index, and appends it to the result string.
func WhatToPrint(i int, word, filename, res string) (string, error) {
	for _, letter := range word {
		line, err := GetLine(1+(int(letter)-32)*9+i, filename)
		res += line
		if err != nil {
			return "", err
		}
	}
	return res, nil
}

// GetLine retrieves a specific line (num) from a file (filename).
// It returns the line as a string and any error encountered during file operations.
func GetLine(num int, filename string) (string, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil { // Check if there was an error opening the file
		return "", err
	}
	defer file.Close()

	str := ""
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	// Scan through the file until the desired line number is reached
	for scanner.Scan() {
		if lineNumber == num {
			str = scanner.Text()
		}
		lineNumber++
	}

	return str, nil
}

// handleASCII handles the HTTP request for ASCII art generation.
// It reads input from the request, selects the banner option, and converts the input to ASCII art.
// The generated ASCII art is then returned as a string.
func handleASCII(w http.ResponseWriter, r *http.Request) string {

	input := r.FormValue("fname")
	// Retrieve the selected banner option from the request
	bannerOption  := r.FormValue("banner")
	if bannerOption == ""{
		// Return an error message if the banner file cannot be read
		http.ServeFile(w, r, "template/error3.html")
		return ""	
	}

	banner := ""
	switch bannerOption {
	case "shadow":
		banner = "shadow.txt"
	case "thinkertoy":
		banner = "thinkertoy.txt"
	case "standard":
		banner = "standard.txt"
	}

	if banner == ""{
		http.Error(w,"asdf",http.StatusInternalServerError)
	}
	input = strings.ReplaceAll(input, "\r\n", "\\n")
	lines := strings.Split(input, "\\n")

	// Create a strings.Builder to store the output
	var output string

	output, err := ConvertToASCIIArt(lines, banner)
	if err != nil {
		fmt.Println("error")
		os.Exit(0)
	}
	return output
}

// ConvertToASCIIArt converts a slice of input lines into ASCII art using a specified banner.
// It returns the generated ASCII art as a string and any error encountered during conversion.
func ConvertToASCIIArt(lines []string, banner string) (string, error) {
	output := ""
	res := ""
	for _, word := range lines {
		if word == "" {
			output = output + "\n"
			continue
		}

		for i := 0; i < 8; i++ {
			res, err := WhatToPrint(i, word, banner, res)
			if err != nil {
				os.Exit(0)
			}
			output = output + "\n" + res
			res = ""
		}
	}

	return output, nil
}
