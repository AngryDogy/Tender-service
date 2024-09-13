package model

import "time"

type Bid struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	TenderId    string     `json:"tender_id"`
	AuthorType  AuthorType `json:"author_type"`
	AuthorId    string     `json:"author_id"`
	Version     int        `json:"version"`
	CreatedAt   time.Time  `json:"created_at"`
}

type CreateBidDto struct {
	Name        string     `json:"name" validate:"max=100"`
	Description string     `json:"description" validate:"max=500"`
	TenderId    string     `json:"tenderId"`
	AuthorType  AuthorType `json:"authorType"`
	AuthorId    string     `json:"authorId"`
}

type PatchBidDto struct {
	Name        string `json:"name" validate:"max=100"`
	Description string `json:"description" validate:"max=500"`
}
