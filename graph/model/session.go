package model

import "time"

// Session as used in API
type Session struct {
	Token          string    `json:"token"`
	Note           string    `json:"note"`
	UserID         uint      `json:"-"`
	CreatedAt      time.Time `json:"createdAt"`
	LastActivityAt time.Time `json:"lastActivityAt"`
}
