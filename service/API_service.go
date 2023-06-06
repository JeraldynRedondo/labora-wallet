package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"my-labora-wallet-project/model"

	"github.com/joho/godotenv"
)

// Singers that relate to Truora API requests
const (
	BaseUrl     = "https://api.checks.truora.com/v1/checks"
	ContentType = "application/x-www-form-urlencoded"
	APITimeout  = 5 * time.Second
)

// request is a function that makes an http request to a Truora API and returns the body of the response.
func request(method, url string, payload *strings.Reader) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {

		return nil, fmt.Errorf("Error creating request to API: %w", err)
	}

	req.Header.Add("Truora-API-Key", getAPI_KEY())
	req.Header.Add("Content-Type", ContentType)

	resp, err := client.Do(req)
	if err != nil {

		return nil, fmt.Errorf("Error making request to API: %w", err)
	}
	defer resp.Body.Close()

	// Read API response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		return nil, fmt.Errorf("Error reading the response API: %w", err)
	}

	return body, nil
}

// postTruoraAPIRequest is a function that makes a POST request to Truora "Background Check" API and returns the person's checkID.
func postTruoraAPIRequest(nationalID, country, entity_type string, userAuthorized bool) (string, error) {

	//Getting the url with the request body
	data := makePostBody(nationalID, country, entity_type, userAuthorized)
	//Create the request input variables
	payload := strings.NewReader(data.Encode())
	method := "POST"

	body, err := request(method, BaseUrl, payload)
	if err != nil {

		return "", fmt.Errorf("Error, failed to make POST request to API: %w", err)
	}

	var Response model.TruoraPostResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {

		return "", fmt.Errorf("Error decoding the POST response API: %w", err)
	}

	// Get the ID of the object
	checkID := Response.Check.CheckID

	return checkID, nil
}

// getTruoraAPIRequest is a function that makes a GET request to Truora "Background Check" API and returns the person's score.
func getTruoraAPIRequest(checkID string) (int, error) {
	url := BaseUrl + "/" + checkID
	method := "GET"
	payload := strings.NewReader("")

	body, err := request(method, url, payload)
	if err != nil {

		return -1, fmt.Errorf("Error, failed to make GET request to API: %w", err)
	}

	var Response model.TruoraGetResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {

		return -1, fmt.Errorf("Error decoding the GET response API: %w", err)
	}

	// We get the score of the object
	score := Response.Check.Score

	return score, nil
}

// truoraAPIRequest is a function that executes the two requests to the truora API to get the score.
func truoraAPIRequest(national_id, country, entity_type string, userAuthorized bool) (int, error) {
	checkID, err := postTruoraAPIRequest(national_id, country, entity_type, userAuthorized)
	if err != nil {

		return -1, fmt.Errorf("Post request failed: %w", err)
	}

	time.Sleep(APITimeout) // Espera 2 segundos antes de llamar a getTruoraAPIRequest

	score, err := getTruoraAPIRequest(checkID)
	if err != nil {

		return -1, fmt.Errorf("Get request failed: %w", err)
	}

	return score, nil
}

// GetApproval is a function that decides if the creation of the wallet is approved according to the score received.
func GetApproval(national_id, country, entity_type string, userAuthorized bool) (bool, error) {
	score, err := truoraAPIRequest(national_id, country, entity_type, userAuthorized)

	if err != nil {

		return false, fmt.Errorf("Error,score request failed: %w", err)
	}

	return score == 1, nil
}

// getAPI_KEY is a function that returns the API key that it extracts from the .env file.
func getAPI_KEY() string {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error loading .env file")
	}
	API_KEY := string(os.Getenv("TRUORA_API_KEY"))

	return API_KEY
}

// makePostBody is a function that returns the url of the post service with the related body data.
func makePostBody(nationalID, country, entity_type string, userAuthorized bool) url.Values {
	data := url.Values{}
	data.Set("national_id", nationalID)
	data.Set("country", country)
	data.Set("type", entity_type)
	data.Set("user_authorized", strconv.FormatBool(userAuthorized))

	return data
}
