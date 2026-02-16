package models

type MkRelation struct {
	ID                             string `json:"id"`
	IsFollowing                    bool   `json:"isFollowing"`
	IsFollowed                     bool   `json:"isFollowed"`
	HasPendingFollowRequestFromYou bool   `json:"hasPendingFollowRequestFromYou"`
	HasPendingFollowRequestToYou   bool   `json:"hasPendingFollowRequestToYou"`
	IsBlocking                     bool   `json:"isBlocking"`
	IsBlocked                      bool   `json:"isBlocked"`
	IsMuted                        bool   `json:"isMuted"`
	IsRenoteMuted                  bool   `json:"isRenoteMuted"`
}

func (r MkRelation) ToRelationship() Relationship {
	return Relationship{
		ID:             r.ID,
		Following:      r.IsFollowing,
		ShowingReblogs: !r.IsRenoteMuted,
		FollowedBy:     r.IsFollowed,
		Requested:      r.HasPendingFollowRequestFromYou,
		RequestedBy:    r.HasPendingFollowRequestToYou,
		Languages:      []string{},
		Blocking:       r.IsBlocking,
		BlockedBy:      r.IsBlocked,
		Muting:         r.IsMuted,
	}
}
