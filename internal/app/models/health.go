package models

type HealthOutput struct {
	Body struct {
		Ok       bool              `json:"ok" example:"true" doc:"Is app ok"`
		Database map[string]string `json:"database" example:"{\"status\": \"up\"}" doc:"database status"`
	}
}

type HealthInput struct{}
