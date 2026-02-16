package models

type Announcement struct {
	ID          string     `json:"id"`
	Content     string     `json:"content"`
	StartsAt    *string    `json:"starts_at"`
	EndsAt      *string    `json:"ends_at"`
	Published   bool       `json:"published"`
	AllDay      bool       `json:"all_day"`
	PublishedAt string     `json:"published_at"`
	UpdatedAt   string     `json:"updated_at"`
	Read        bool       `json:"read,omitempty"`
	Mentions    []struct{} `json:"mentions"`
	Statuses    []struct{} `json:"statuses"`
	Tags        []Tag      `json:"tags"`
	Emojis      []struct{} `json:"emojis"`
	Reactions   []struct{} `json:"reactions"`
}
