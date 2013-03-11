/*

*/
package main

import (
	"fmt"
	"io/ioutil"
	"safehttp"
	"regexp"
	"sync"
	"os"
	"log"
	"strconv"
)

const (
	// generally 7 digits
	URL_FORMAT = "http://oxforddictionaries.com/view/entry/m_en_us%d"
)

var REGEXP = regexp.MustCompile("Definition of (.+) in Oxford")

type Entry struct {
	id int
}

/**
Gets the word for that entry id
*/
func FetchWord(entry Entry) (word *string) {
	url := fmt.Sprintf(URL_FORMAT, entry.id)
	resp, err := safehttp.Get(url, 10, 10)
	error := "error"
	if err != nil {
		return &error
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &error
	}

	//log.Printf(string(content))

	matches := REGEXP.FindSubmatch(content)
	noMatch := "no match"
	if len(matches) != 2 {
		return &noMatch
		return
	}

	wordString := string(matches[1])

	word = &wordString
	return word
}

func consumeEntryId(queue chan Entry) {
	defer consumer_wg.Done()

	for item := range queue {
		fmt.Printf("m_en_us%d\t%s\n", item.id, *FetchWord(item))
	}
}

var resultingChannel = make(chan Entry)
var consumer_wg sync.WaitGroup

func main() {
	args := os.Args
	if (len(args) < 3) {
		log.Printf("Invalid syntax, supply start and end entry ids");
		os.Exit(1)
	}

	for c := 0; c < 20; c++ {
		consumer_wg.Add(1)
		go consumeEntryId(resultingChannel)
	}

	start, _ := strconv.ParseInt(args[1], 10, 0)
	end, _ := strconv.ParseInt(args[2], 10, 0)
	for i := int(start); i < int(end); i++ {
		resultingChannel <- Entry{id: i}
	}

	close(resultingChannel)

	consumer_wg.Wait()
}
