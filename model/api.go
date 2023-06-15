package model

import "time"

//API_Request is a structure that represents the body of the API request.
type API_Request struct {
	NationalId     string `json:"national_id"`
	Country        string `json:"country"`
	EntityType     string `json:"type"`
	UserAuthorized bool   `json:"user_authorized"`
}

// CheckPost is a structure that represents a part of the response of the POST request to the Background Checks PI.
type CheckPost struct {
	CheckID string `json:"check_id"`
}

// TruoraPostResponse is a structure that represents the response of the POST request to the Background Checks API.
type TruoraPostResponse struct {
	Check CheckPost `json:"check"`
}

// CompanySummary is a function that represents a part of the structure Check of the GET request response to the Background Checks API.
type CompanySummary struct {
	CompanyStatus string `json:"company_status"`
	Result        string `json:"result"`
}

// CheckGet is a functon that represents a part of the structure Check of the GET request response to the Background Checks API.
type CheckGet struct {
	CheckID        string         `json:"check_id"`
	CompanySummary CompanySummary `json:"company_smmary"`
	Country        string         `json:"country"`
	CreationDate   time.Time      `json:"creation_dat"`
	NameScore      int            `json:"name_score`
	IDScore        int            `json:"id_scor"`
	Score          int            `json:"score"`
}

// TruoraGetResponse is a strucure that represents the response of the GET request to the Background Checks API.
type TruoraGetResponse struct {
	Check CheckGet `json:"check"`
}
