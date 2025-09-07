package model

import "time"

type Account struct {
	ID              string    `json:"id"`
	Active          bool      `json:"active"`
	Name            string    `json:"name"`
	Address         string    `json:"address"`
	Packages        []string  `json:"packages"`
	BookID          string    `json:"book_id"`
	StaffingId      string    `json:"staffing_id"`
	CommunicationId string    `json:"communication_id"`
	Created         time.Time `json:"since"`
	Updated         time.Time `json:"updated"`
}
