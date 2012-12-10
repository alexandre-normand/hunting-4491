/*

*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"safehttp"
	"time"
)

type EchoResponse struct {
	Method         ContentAttribute  `json:"method,omitempty"`
	ApiKey         ContentAttribute  `json:"apiKey,omitempty"`
	NoJsonCallback ContentAttribute  `json:"nojsoncallback,omitempty"`
	Format         ContentAttribute  `json:"format,omitempty"`
	Stat           string            `json:"stat,omitempty"`
}

type ContentAttribute struct {
	Content      string              `json:"_content,omitempty"`
}

const (
	BASE_URL               = "http://api.flickr.com/services/rest/?method=flickr.test.echo&api_key=%s&format=json&nojsoncallback=1"
	API_KEY                = "REPLACE_WITH_YOURS"
)

/**
*/
func main() {
	channelOne := make(chan *EchoResponse)
	go ExecutePollingJob(channelOne, 60, 30)

	channelTwo := make(chan *EchoResponse)
	go ExecutePollingJob(channelTwo, 60, 20)

	channelThree := make(chan *EchoResponse)
    go ExecutePollingJob(channelThree, 60, 10)

	jobOne, jobTwo, jobThree := <- channelOne, <- channelTwo, <- channelThree

	HandleJobsTermination("Finished phase 1 with statuses [%s, %s, %s]", jobOne, jobTwo, jobThree)

	channelFour := make(chan *EchoResponse)
	go ExecutePollingJob(channelFour, 65, 500)

	channelFive := make(chan *EchoResponse)
	go ExecutePollingJob(channelFive, 60, 550)

	jobFour, jobFive := <- channelFour, <- channelFive
	HandleJobsTermination("Finished phase 2 with statuses: [%s, %s]\n",
		jobFour, jobFive)
}

/**
Prints the final state of the jobs and exits if any one of them didn't complete successfully
*/
func HandleJobsTermination(logformat string, responses ...*EchoResponse) {
	logValues := make([]interface{}, len(responses))
	success := true
	for i := range responses {
		if responses[i] == nil {
			logValues[i] = "nil"
			success = false
		} else {
			logValues[i] = responses[i].Stat
			if responses[i].Stat != "ok" {
				success = false
			}
		}
	}

	log.Printf(logformat, logValues...)

	if !success {
		log.Fatal("Aborting test since not all jobs completed successfully")
	}
}

/**
Makes a call to the flickr api at the given interval
*/
func ExecutePollingJob(destinationChannel chan *EchoResponse, intervalInSeconds int, numberOfPolls int) {
	response := ExecutePollingTask(intervalInSeconds, numberOfPolls)

	destinationChannel <- response
}

/**
Polls for a job to Complete or fail. When done, it sends the last jobExecution to the returnChannel
*/
func ExecutePollingTask(intervalInSeconds int, numberOfPolls int) (lastResponse *EchoResponse) {
	echoResponse, err := ExecuteSinglePoll()
	log.Printf("Poll with API_KEY [%s] with stat [%s]\n", API_KEY, echoResponse.Stat)

	for count := 0; count < numberOfPolls && err == nil; count++ {
		time.Sleep(time.Duration(intervalInSeconds) * time.Second)
		echoResponse, err = ExecuteSinglePoll()
		if echoResponse != nil {
			log.Printf("Polled with API_KEY[%s] and got stat [%s]\n", API_KEY, echoResponse.Stat)
		} else {
			log.Printf("Response is nil, err is %s", err)
		}
	}

	return echoResponse
}

/**
Gets the latest state of a job given its batch id
*/
func ExecuteSinglePoll() (response *EchoResponse, err error) {
	url := fmt.Sprintf(BASE_URL, API_KEY)
	resp, err := safehttp.Get(url, 10, 10)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var echoResponse EchoResponse
	err = json.Unmarshal(content, &echoResponse)
	if err != nil {
		return nil, err
	}

	return &echoResponse, nil
}
