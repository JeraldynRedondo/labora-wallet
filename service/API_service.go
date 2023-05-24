package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
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

// request is a function that makes an http request to a Truora API and returns the body of the response.
func request(method, url string, payload *strings.Reader) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, fmt.Errorf("Error creating request: %w", err)
	}
	req.Header.Add("Truora-API-Key", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50X2lkIjoiIiwiYWRkaXRpb25hbF9kYXRhIjoie30iLCJjbGllbnRfaWQiOiJUQ0lmMTk4Y2Y4NjYzNjk2ZGM1YWQ4MGY4N2U5NjQzODBmNiIsImV4cCI6MzI2MTYzMDE4OCwiZ3JhbnQiOiIiLCJpYXQiOjE2ODQ4MzAxODgsImlzcyI6Imh0dHBzOi8vY29nbml0by1pZHAudXMtZWFzdC0xLmFtYXpvbmF3cy5jb20vdXMtZWFzdC0xXzJSZ2ZJSmxQeCIsImp0aSI6IjdlZjIwZjJjLWUwOGUtNDI1Mi05OTc0LWIzMjUzZmQ1NmM5NCIsImtleV9uYW1lIjoiYXBpX3dhbGxldCIsImtleV90eXBlIjoiYmFja2VuZCIsInVzZXJuYW1lIjoiZ21haWxyZWRvbmRvaW5nLWFwaV93YWxsZXQifQ.NiFZeqqIP8fv4XgWVAXa0xr7Sx59feiJi1bg-3sRPB0")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %w", err)
	}
	defer resp.Body.Close()

	// Read API response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response: %w", err)
	}

	return body, nil
}

// postTruoraAPIRequest is a function that makes a POST request to Truora "Background Check" API and returns the person's checkID.
func postTruoraAPIRequest() (string, error) {
	url := "https://api.checks.truora.com/v1/checks"
	method := "POST"
	payload := strings.NewReader("national_id=74909799&country=PE&type=person&user_authorized=true")

	body, err := request(method, url, payload)
	if err != nil {
		return "", fmt.Errorf("Error making request: %w", err)
	}

	var Response TruoraPostResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {
		return "", fmt.Errorf("Error decoding response: %w", err)
	}

	// We get the score of the object
	checkID := Response.Check.CheckID

	// show api response
	//fmt.Println("Valor de check_id:", checkID)
	return checkID, nil
}

// getTruoraAPIRequest is a function that makes a GET request to Truora "Background Check" API and returns the person's score.
func getTruoraAPIRequest() (int, error) {
	checkID, err := postTruoraAPIRequest()
	if err != nil {
		return -1, fmt.Errorf("Post request failed: %w", err)
	}

	url := "https://api.checks.truora.com/v1/checks/" + checkID
	method := "GET"
	payload := strings.NewReader("")

	body, err := request(method, url, payload)
	if err != nil {
		return -1, fmt.Errorf("Error making request: %w", err)
	}

	var Response TruoraGetResponse
	err = json.Unmarshal(body, &Response)
	if err != nil {
		return -1, fmt.Errorf("Error decoding response: %w", err)
	}

	// We get the score of the object
	score := Response.Check.Score

	// show api response
	//fmt.Println("Score:", score)

	return score, nil
}

func GetApproval() (bool, error) {
	score, err := getTruoraAPIRequest()

	if score == 1 {
		if err != nil {
			return false, fmt.Errorf("Error,score request failed: %w", err)
		}
	}

	return true, nil
}
