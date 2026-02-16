package models

type Marker struct {
	LastReadID string `json:"last_read_id"`
	Version    int    `json:"version"`
	UpdatedAt  string `json:"updated_at"`
}
