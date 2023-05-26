package model

type API_Request struct {
	National_id    string `json:"national_id"`
	Country        string `json:"country"`
	Entity_type    string `json:"type"`
	UserAuthorized bool   `json:"userAuthorized"`
}
