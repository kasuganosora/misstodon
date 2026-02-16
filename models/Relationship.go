package models

type Relationship struct {
	ID                  string   `json:"id"`
	Following           bool     `json:"following"`
	ShowingReblogs      bool     `json:"showing_reblogs"`
	Notifying           bool     `json:"notifying"`
	FollowedBy          bool     `json:"followed_by"`
	Blocking            bool     `json:"blocking"`
	BlockedBy           bool     `json:"blocked_by"`
	Muting              bool     `json:"muting"`
	MutingNotifications bool     `json:"muting_notifications"`
	Requested           bool     `json:"requested"`
	RequestedBy         bool     `json:"requested_by"`
	DomainBlocking      bool     `json:"domain_blocking"`
	Endorsed            bool     `json:"endorsed"`
	Languages           []string `json:"languages"`
	Note                string   `json:"note"`
}
