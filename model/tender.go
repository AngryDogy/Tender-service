package model

import (
	"time"
)

type Tender struct {
	Id              string      `json:"id"`
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	ServiceType     ServiceType `json:"serviceType"`
	Status          Status      `json:"status"`
	OrganizationId  string      `json:"organizationID"`
	CreatorUsername string      `json:"creatorUsername"`
	Version         int32       `json:"version"`
	CreatedAt       time.Time   `json:"createdAt"`
}

type CreateTenderDto struct {
	Name            string      `validate:"max=100" json:"name"`
	Description     string      `validate:"max=500" json:"description"`
	ServiceType     ServiceType `json:"serviceType"`
	OrganizationId  string      `validate:"max=100" json:"organizationId"`
	CreatorUsername string      `json:"creatorUsername"`
}

type PatchTenderDto struct {
	Name        string      `validate:"max=100" json:"name"`
	Description string      `validate:"max=500" json:"description"`
	ServiceType ServiceType `json:"serviceType"`
}
