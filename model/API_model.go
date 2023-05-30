package model

import "time"

type API_Request struct {
	National_id    string `json:"national_id"`
	Country        string `json:"country"`
	Entity_type    string `json:"type"`
	UserAuthorized bool   `json:"userAuthorized"`
}

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
