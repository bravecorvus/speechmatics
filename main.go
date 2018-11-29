package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gilgameshskytrooper/speechmatics/structs"
)

// Defined as environment variables so that private info isn't stored on a public repo
var (
	userid    string = os.Getenv("SPEECHMATICSUSERID")
	authtoken string = os.Getenv("SPEECHMATICSAUTHTOKEN")
)

// Helper function to download files to disk
func downloadFromUrl(url string) (string, error) {
	tokens := strings.Split(url, "/")
	filename := tokens[len(tokens)-1]

	fmt.Println("Downloading " + url + "...")
	output, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error while creating", filename, "-", err)
		return "", err
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return "", err
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return "", err
	}
	fmt.Println("downloading complete")
	fmt.Println()

	return filename, nil
}

func main() {

	// Make sure userid and authtoken exist
	if userid == "" || authtoken == "" {
		// log.Fatal prints out the output and quits the program with status 1
		log.Fatal("Please make sure you defined the variables \nSPEECHMATICSUSERID\n\nand\n\nSPEECHMATICSAUTHTOKEN in your environment variables")
	}

	// Make sure user passes link to download
	if len(os.Args) == 1 {
		log.Fatal("Please pass the link to download as the argument")
	}

	// Download the example audio file using the link provided in
	filename, err := downloadFromUrl(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// open file and assign to a file variable of type File
	file, err := os.Open(filename)
	if err != nil {
		os.Remove(filename)
		log.Fatal(err)
	}

	// read file into io.Reader to  make it compatible with the writer interface (for multipart-form data writing)
	fileContents, _ := ioutil.ReadAll(file)
	file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file as data_file multipart field
	part, _ := writer.CreateFormFile("data_file", filepath.Base(file.Name()))
	// write the contents of the file into the data_file field (-F "data_file=@/filepath/filename" in curl)
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
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}

	start := time.Now()

	// Execute the POST request and save the result to resp variable
	resp, err := client.Do(req)
	if err != nil {
		os.Remove(filename)
		fmt.Println("resp, err := client.Do(req)")
		log.Fatal(err)
	}

	// Print the status
	fmt.Println("Status:", resp.Status)
	defer resp.Body.Close()

	// Read the result of the resp body into a byte array
	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		os.Remove(filename)
		fmt.Println("reply, err := ioutil.ReadAll(resp.Body)")
		log.Fatal(err)
	}

	var initialresponsereply structs.Initial
	json.Unmarshal(reply, &initialresponsereply)

	fmt.Println(initialresponsereply)

	result, end := job(initialresponsereply.ID, filename)

	sttresult := make([]string, 0)
	for _, word := range result.Words {
		sttresult = append(sttresult, word.Name)
	}

	fmt.Println("STT Result: \"" + strings.Join(sttresult, " ") + "\"")
	fmt.Println("Call took " + strconv.Itoa(int(end.Sub(start).Seconds())) + " seconds.")

	os.Remove(filename)
}

func job(id int, filename string) (*structs.Transcript, time.Time) {
	url := "https://api.speechmatics.com/v1.0/user/" + userid + "/jobs/" + strconv.Itoa(id) + "/transcript?format=json&auth_token=" + authtoken
	resp, err := http.Get(url)
	if err != nil {
		os.Remove(filename)
		fmt.Println("resp, err := http.Get(url)")
		log.Fatal(err)
	}

	defer resp.Body.Close()

	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		os.Remove(filename)
		fmt.Println("reply, err := ioutil.ReadAll(resp.Body)")
		log.Fatal(err)
	}

	var checkstatus structs.CheckStatus
	json.Unmarshal(reply, &checkstatus)

	if checkstatus.Code == 404 {
		time.Sleep(1 * time.Millisecond)
		return job(id, filename)
	} else {
		var transcript structs.Transcript
		json.Unmarshal(reply, &transcript)
		return &transcript, time.Now()
	}
}
