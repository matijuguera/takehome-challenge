package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"takehome-challenge/downloader"
	"takehome-challenge/houseresponse"
	"time"
)

var url = "https://app-homevision-staging.herokuapp.com/api_project/houses"
var totalpages = 10 // Requests the first 10 pages of results as requested
var maxRetries = 5  // Max retries for backoff

// execute the program
func execute() {
	//Create the waitgroup
	var wg sync.WaitGroup

	// Iterates the pages
	for i := 1; i <= totalpages; i++ {
		// Add one to the queue
		wg.Add(1)

		// We pass the "i" value so it knows which page it is
		go func(i int) {

			// Initializes an empty result
			var result houseresponse.HouseResponse

			// Checks the page we attempt to fetch is reachable
			fetch(i, &result)

			// Downloads all the houses from the page
			downloadImages(&result, &wg)

			// Notifies we completed downloading the page
			fmt.Println("Finished downloading page:", i)
			// Substracts one from the queue
			wg.Done()

		}(i)
	}
	// Wait until queue is completed before continuing
	wg.Wait()

	fmt.Printf("Finished downloading all photos for the first %d pages\n", totalpages)
}

// Downloads the images
func downloadImages(result *houseresponse.HouseResponse, wg *sync.WaitGroup) {

	// Iterates each house from the page
	for _, house := range result.Houses {

		// Add one to the queue
		wg.Add(1)

		// We pass the house
		go func(house houseresponse.Houses) {

			// Prepares the strings required for download the file
			fileExtension := "." + strings.Split(house.PhotoURL, ".")[len(strings.Split(house.PhotoURL, "."))-1]
			fileName := "images/" + strconv.Itoa(house.Id) + "-" + house.Address + fileExtension
			imageUrl := house.PhotoURL

			// Download the file
			downloader.DownloadFile(imageUrl, fileName)

			// Substracts one from the queue
			wg.Done()

		}(house)
	}
}

// Fetches and checks if the page works
func fetch(i int, result *houseresponse.HouseResponse) {

	// Initialize variables
	correctStatusCode := false                  // If the response code is 200 it modifies to true
	backoff := 0.5                              // Time it will wait after a failure
	retries := 0                                // current retries
	pageUrl := url + "?page=" + strconv.Itoa(i) // Specific page URL to fetch

	// Tries multiple times (up to the maxRetries) until it gets a 200 response code
	for !correctStatusCode && retries < maxRetries {

		// Get url
		resp, err := http.Get(pageUrl)

		// add error if failed
		if err != nil {
			errors.New("No response from request")
		}

		// Make sure the response will be closed
		defer resp.Body.Close()

		// Checks the responseCode is 200
		if resp.StatusCode != 200 {

			// Notifies the issue
			errors.New("Received non 200 response code")

			// Waits so it doesn't overload the page
			time.Sleep(time.Duration(backoff) * time.Second)

			// updates the backoff and retries
			backoff = backoff * 2
			retries = retries + 1

		} else {
			// if the response is 200, then continue with the program
			correctStatusCode = true
		}

		// Read the response
		body, err := ioutil.ReadAll(resp.Body)

		// Parses the JSON
		if err := json.Unmarshal(body, &result); err != nil {
			log.Fatal("Can not unmarshal JSON")
		}
	}

	// If it tried too many times, continue with the program and ignore this page (we assume the site is not working)
	if retries >= 5 {
		fmt.Printf("Too many retries (%d) when trying to fetch %s, ignoring the URL", retries, pageUrl)
	}
}

func main() {

	start := time.Now()
	execute()
	elapsed := time.Since(start)
	fmt.Printf("Total execution time: %v", elapsed)

}
