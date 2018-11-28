package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Defined as environment variables so that private info isn't stored on a public repo
var (
	userid    string = os.Getenv("SPEECHMATICSUSERID")
	authtoken string = os.Getenv("SPEECHMATICSAUTHTOKEN")
)

// Helper function to download files to disk
func downloadFromUrl(url string) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	fmt.Println("Downloading " + url + "...")
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	fmt.Println("downloading complete")
	fmt.Println()

}

func main() {

	// Make sure userid and authtoken exist
	if userid == "" || authtoken == "" {
		// panic prints out the output and quits the program with status 1
		panic("Please make sure you defined the variables SPEECHMATICSUSERID and SPEECHMATICSAUTHTOKEN in your environment variables")
	}

	// Download the example audio file
	downloadFromUrl("https://s3.amazonaws.com/voiceit-api2-testing-files/test-data/enrollmentStephen1.wav")

	// open file and assign to a file variable of type File
	file, err := os.Open("./enrollmentStephen1.wav")
	if err != nil {
		os.Remove("./enrollmentStephen1.wav")
		panic(err)
	}

	// read file into io.Reader to  make it compatible with the writer interface (for multipart-form data writing)
	fileContents, _ := ioutil.ReadAll(file)
	file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file as data_file multipart field
	part, _ := writer.CreateFormFile("data_file", filepath.Base(file.Name()))
	// write the contents of the file into the data_file field (-F "data_file=@/filepath/enrollmentStephen1.wav" in curl)
	part.Write(fileContents)
	// also add model=en-US (-F "model=en-US" in curl)
	writer.WriteField("model", "en-US")
	writer.Close()

	// define URL
	url := "https://api.speechmatics.com/v1.0/user/" + userid + "/jobs/?auth_token=" + authtoken
	fmt.Println("URL to hit:", url)
	fmt.Println()

	// Define post request to the jobs request endpoint with the multipart body attached
	req, _ := http.NewRequest("POST", url, body)

	// Create a http client with a timeout of 5 minutes
	client := &http.Client{Timeout: 5 * time.Minute}
	// Execute the POST request and save the result to resp variable
	resp, err := client.Do(req)
	if err != nil {
		os.Remove("./enrollmentStephen1.wav")
		fmt.Println("resp, err := client.Do(req)")
		panic(err)
	}

	// Print the status
	fmt.Println("Status:", resp.Status)
	defer resp.Body.Close()

	// Read the result of the resp body into a byte array
	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		os.Remove("./enrollmentStephen1.wav")
		fmt.Println("reply, err := ioutil.ReadAll(resp.Body)")
		panic(err)
	}

	// Print the string casted reply
	fmt.Println(string(reply))

	// Remove file at the end of the program
}
