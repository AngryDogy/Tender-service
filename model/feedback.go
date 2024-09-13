package model

import "time"

type Feedback struct {
	Id          string    `json:"id"`
	BidId       string    `json:"bid_id" `
	Description string    `json:"description"`
	Username    string    `json:"username"`
	CreatedAt   time.Time `json:"created_at"`
}
