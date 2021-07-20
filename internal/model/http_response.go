package model

import "time"

type HttpResponse struct {
	Success   bool `json:"success"`
	Timestamp time.Time
	Message   string `json:"message"`
	Data      string `json:"data"`
}
