package models

type Preferences struct {
	PostingDefaultVisibility string  `json:"posting:default:visibility"`
	PostingDefaultSensitive  bool    `json:"posting:default:sensitive"`
	PostingDefaultLanguage   *string `json:"posting:default:language"`
	ReadingExpandMedia       string  `json:"reading:expand:media"`
	ReadingExpandSpoilers    bool    `json:"reading:expand:spoilers"`
}
