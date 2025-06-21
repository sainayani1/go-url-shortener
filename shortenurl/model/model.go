package model

import "time"

type ShortURL struct {
	ID          int       `json:"id"`
	URL         string    `json:"url"`
	ShortCode   string    `json:"shortCode"`
	AccessCount int       `json:"accessCount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ShortURLCreateReq struct {
	URL string `json:"url"`
}
