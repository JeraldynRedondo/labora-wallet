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

	"github.com/joho/godotenv"
)

type TruoraPostResponse struct {
	Check struct {
		CheckID string `json:"check_id"`
	} `json:"check"`
}

type TruoraGetResponse struct {
	Check struct {
		CheckID        string `json:"check_id"`
		CompanySummary struct {
			CompanyStatus string `json:"company_status"`
			Result        string `json:"result"`
		} `json:"company_summary"`
		Country      string    `json:"country"`
		CreationDate time.Time `json:"creation_date"`
		NameScore    int       `json:"name_score"`
		IDScore      int       `json:"id_score"`
		Score        int       `json:"score"`
	}
}

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
	urlFull, data := getPostUrl(nationalID, country, entity_type, userAuthorized)
	//Create the request input variables
	urlStr := urlFull.String() // "https://api.checks.truora.com/v1/checks"
	payload := strings.NewReader(data.Encode())
	method := "POST"

	body, err := request(method, urlStr, payload)
	if err != nil {

		return "", fmt.Errorf("Error, failed to make POST request to API: %w", err)
	}

	var Response TruoraPostResponse
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
	url := "https://api.checks.truora.com/v1/checks/" + checkID
	method := "GET"
	payload := strings.NewReader("")

	body, err := request(method, url, payload)
	if err != nil {

		return -1, fmt.Errorf("Error, failed to make GET request to API: %w", err)
	}

	var Response TruoraGetResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {

		return -1, fmt.Errorf("Error decoding the GET response API: %w", err)
	}

	// We get the score of the object
	score := Response.Check.Score

	return score, nil
}

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

// getPostUrl is a function that returns the url of the post service with the related body data.
func getPostUrl(nationalID, country, entity_type string, userAuthorized bool) (*url.URL, url.Values) {
	apiUrl := "https://api.checks.truora.com"
	resource := "/v1/checks"
	data := url.Values{}
	data.Set("national_id", nationalID)
	data.Set("country", country)
	data.Set("type", entity_type)
	data.Set("user_authorized", strconv.FormatBool(userAuthorized))

	// Create the url type variable
	urlFull, _ := url.ParseRequestURI(apiUrl)
	urlFull.Path = resource

	return urlFull, data
}
