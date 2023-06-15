package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"my-labora-wallet-project/model"

	"github.com/joho/godotenv"
)

var API_KEY string

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	API_KEY = string(os.Getenv("TRUORA_API_KEY"))
}

// Singers that relate to Truora API requests
const (
	BaseUrl        = "https://api.checks.truora.com/v1/checks"
	ContentType    = "application/x-www-form-urlencoded"
	APITimeout     = 5 * time.Second
	ScoreThreshold = 1
)

// request is a function that makes an http request to a Truora API and returns the body of the response.
func request(method, url string, payload *strings.Reader) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, fmt.Errorf("Error creating request to API: %w", err)
	}

	req.Header.Add("Truora-API-Key", API_KEY)
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

// postAPIRequest is a function that makes a POST request to Truora "Background Check" API and returns the person's checkID.
func postAPIRequest(nationalID, country, entity_type string, userAuthorized bool) (string, error) {
	data := buildBodyRequest(nationalID, country, entity_type, userAuthorized)
	payload := strings.NewReader(data.Encode())

	body, err := request(http.MethodPost, BaseUrl, payload)
	if err != nil {
		return "", fmt.Errorf("Error, failed to make POST request to API: %w", err)
	}

	var Response model.TruoraPostResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {
		return "", fmt.Errorf("Error decoding the POST response API: %w", err)
	}

	checkID := Response.Check.CheckID

	return checkID, nil
}

// getAPIRequest is a function that makes a GET request to Truora "Background Check" API and returns the person's score.
func getAPIRequest(checkID string) (int, error) {
	url := BaseUrl + "/" + checkID
	payload := strings.NewReader("")

	body, err := request(http.MethodGet, url, payload)
	if err != nil {
		return -1, fmt.Errorf("Error, failed to make GET request to API: %w", err)
	}

	var Response model.TruoraGetResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {
		return -1, fmt.Errorf("Error decoding the GET response API: %w", err)
	}

	score := Response.Check.Score

	return score, nil
}

// APIRequest is a function that executes the two requests to the truora API to get the score.
func APIRequest(national_id, country, entity_type string, userAuthorized bool) (int, error) {
	checkID, err := postAPIRequest(national_id, country, entity_type, userAuthorized)
	if err != nil {
		return -1, fmt.Errorf("Post request failed: %w", err)
	}

	time.Sleep(APITimeout)

	score, err := getAPIRequest(checkID)
	if err != nil {
		return -1, fmt.Errorf("Get request failed: %w", err)
	}

	return score, nil
}

// GetApproval is a function that decides if the creation of the wallet is approved according to the score received.
func GetApproval(national_id, country, entity_type string, userAuthorized bool) (bool, error) {
	score, err := APIRequest(national_id, country, entity_type, userAuthorized)

	if err != nil {
		return false, fmt.Errorf("Error,score request failed: %w", err)
	}

	return score >= ScoreThreshold, nil
}

// buildBodyRequest is a function that returns the url of the post service with the related body data.
func buildBodyRequest(nationalID, country, entityType string, userAuthorized bool) url.Values {
	data := url.Values{}
	data.Set("national_id", nationalID)
	data.Set("country", country)
	data.Set("type", entityType)
	data.Set("user_authorized", strconv.FormatBool(userAuthorized))

	return data
}
