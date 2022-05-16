package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var maxRetries = 5

func DownloadFile(URL, fileName string) {
	//Get the response bytes from the url
	connected := false
	retries := 0
	backoff := 0.5

	var response *http.Response

	// Retries until we connect  to the page url.
	for !connected && retries < maxRetries {

		var err error

		response, err = http.Get(URL)
		if err != nil {
			fmt.Printf("Couldn't get %s, retrying...\n", URL)
		}

		defer response.Body.Close()

		if response.StatusCode != 200 {
			fmt.Printf("Response code from %s is not 200, retrying...\n", URL)
			// Waits so it doesn't overload the page
			time.Sleep(time.Duration(backoff) * time.Second)
			// updates the backoff and retries
			backoff = backoff * 2
			retries = retries + 1
		} else {
			connected = true
		}

	}
	if retries >= 5 {
		fmt.Printf("Too many retries (%d) when trying to fetch %s (photoUrl), ignoring the URL", retries, URL)
	} else {
		//Create a empty file
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatal("Error when creating the file")
		}
		defer file.Close()

		//Write the bytes to the file
		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Fatal("Error when copying the file")
		}
	}
}
