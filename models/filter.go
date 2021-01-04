package models

type Filter struct {
	ID       string `json:"_id"`
	DateFrom string `json:"dateFrom"`
	DateTo   string `json:"dateTo"`
	Data     string `json:"data"`
}
