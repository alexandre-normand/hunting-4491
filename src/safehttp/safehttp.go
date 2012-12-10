/**
  Wraps the http package to execute Get and Post requests with retries
*/
package safehttp

import (
	"io"
	"log"
	"net/http"
	"time"
)

func Get(url string, maxAttempts int, retryDelayInSeconds int) (resp *http.Response, err error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return doRequestWithRetry(request, maxAttempts, retryDelayInSeconds)
}

func Post(url string, bodyType string, body io.Reader, maxAttempts int, retryDelayInSeconds int) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	if http.DefaultClient.Jar != nil {
		for _, cookie := range http.DefaultClient.Jar.Cookies(req.URL) {
			req.AddCookie(cookie)
		}
	}

	if err != nil {
		return nil, err
	}

	return doRequestWithRetry(req, maxAttempts, retryDelayInSeconds)
}

func Put(url string, bodyType string, body io.Reader, maxAttempts int, retryDelayInSeconds int) (resp *http.Response, err error) {
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	if http.DefaultClient.Jar != nil {
		for _, cookie := range http.DefaultClient.Jar.Cookies(req.URL) {
			req.AddCookie(cookie)
		}
	}

	if err != nil {
		return nil, err
	}

	return doRequestWithRetry(req, maxAttempts, retryDelayInSeconds)
}

func doRequestWithRetry(request *http.Request, maxAttempts int, retryDelayInSeconds int) (resp *http.Response, err error) {
	resp, err = http.DefaultClient.Do(request)

	for attemptCount := 1; resp == nil && attemptCount < maxAttempts; {
		if err != nil {
			log.Printf("Attempt #%d of %d of doing %s(%s) failed with: %s.\nRetrying in %d seconds", request.Method,
				request.URL, attemptCount, maxAttempts, err, retryDelayInSeconds)
		}

		time.Sleep(time.Duration(retryDelayInSeconds) * time.Second)

		resp, err = http.DefaultClient.Do(request)
		attemptCount++
	}

	return resp, err
}
