package zoho_cliq

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
)

func fixXML(filePatch string) {

	// Open the input file
	inputFile, err := os.Open(filePatch)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer inputFile.Close()

	// Open the output file
	outputFile, err := os.Create(filePatch + ".tmp")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outputFile.Close()

	// Create a regular expression pattern to match numeric tags
	tagPattern := regexp.MustCompile(`<(\d+|/\d+)>`)

	// Read the input file line by line
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		line := scanner.Text()

		// Replace the numeric tags with modified keys
		modifiedLine := tagPattern.ReplaceAllStringFunc(line, func(match string) string {
			// Extract the key value from the tag
			keyValue := tagPattern.FindStringSubmatch(match)[1]
			var modifiedKey string
			if string(rune(keyValue[0])) == "/" {
				modifiedKey = "</newkey>"
			} else {
				modifiedKey = "<newkey>"
			}

			return modifiedKey
		})

		// Write the modified line to the output file
		_, err := outputFile.WriteString(modifiedLine + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	//fmt.Println("Keys replaced successfully.")
	cmd := exec.Command("mv", filePatch+".tmp", filePatch)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
