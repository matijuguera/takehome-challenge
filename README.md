# HomeVision Backend Take Home
## _House image downloader_

[![Build Status](https://travis-ci ?branch=master)](https://github.com/urizz4n/takehome-challenge)

This is a small project in which given a specific API, I wrote a script in Go that does the following:

- Requests the first 10 pages of results from the API
- Parses the JSON returned by the API
- Downloads the photo for each house and saves it in a file with the name formatted as: `[id]-[address].[ext]`
- Uses concurrency to download the photos


## Additional features

- Backoff: With configurable max amount of retries (default: 5)
- Configurable total amount of pages to fetch
- Logs: Notifies when a page was completely downloaded and total execution time.

## Instructions

To run this program, we first need to have Go installed. If you do not have Go installed in your computer, you can go the link below:
https://go.dev/doc/install

Then, follow these steps:

- Open cmd : (WIN + R -> type cmd and open)
- Enter "cd <path to the project folder (takehome-challenge)>" For example: "cd C:\Users\<User>\Desktop\Go project\takehome-challenge"
- Enter "go run main.go"

The current status of the download will be shown in the terminal and the images will be downloaded to the "images" folder

Please let me know if you have any questions!

----
Thanks,
Uriel Sánchez
uriel.sanchez@alu.bue.edu.ar
