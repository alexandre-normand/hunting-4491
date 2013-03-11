/**
 * Created with IntelliJ IDEA.
 * User: skippyjon
 * Date: 2013-03-10
 * Time: 4:22 PM
 * To change this template use File | Settings | File Templates.
 */
package main

import (
	"fmt"
	"os"
	"log"
	"bufio"
	"strings"
)

const (
	DICTIONARY = "./words.txt"
)

func main() {
	args := os.Args
	if (len(args) < 2) {
		log.Printf("Invalid syntax, supply list of words file path");
		os.Exit(1)
	}

	path := args[1]
	if inputFile, err := os.Open(path); err != nil {
		log.Fatal(err)
	} else {
		defer inputFile.Close()

		reader := bufio.NewReader(inputFile)
		dictionary := loadDictionary(DICTIONARY)
		fmt.Fprintf(os.Stderr, "Entries in dictionary: %d\n", len(dictionary))

		for word, err := reader.ReadString('\n'); err == nil; word, err = reader.ReadString('\n') {
			word = strings.TrimSpace(word)
			if entryId, exists := dictionary[word]; exists {
				fmt.Printf("%s\t%s\n", entryId, word)
			} else {
				fmt.Fprintf(os.Stderr, "Can't find %s in dictionary\n", word)
			}
		}
	}

}

func loadDictionary(path string) (map[string] string) {
	dictionary := make(map[string] string)

	if file, err := os.Open(path); err != nil {
		log.Fatal(err)
	} else {
		defer file.Close()
		reader := bufio.NewReader(file)

		for line, err := reader.ReadString('\n'); err == nil; line, err = reader.ReadString('\n') {
			parts := strings.Split(line, "\t")
			if (len(parts) == 2) {
				entryId := strings.TrimSpace(parts[0])
				word := strings.TrimSpace(parts[1])
				dictionary[word] = entryId
			}
		}
	}

	return dictionary
}

