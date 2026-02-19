package models

type (
	Instance struct {
		Uri              string       `json:"uri"`
		Title            string       `json:"title"`
		ShortDescription string       `json:"short_description"`
		Description      string       `json:"description"`
		Email            string       `json:"email"`
		Version          string       `json:"version"`
		Urls             InstanceUrls `json:"urls"`
		Stats            struct {
			UserCount   int `json:"user_count"`
			StatusCount int `json:"status_count"`
			DomainCount int `json:"domain_count"`
		} `json:"stats"`
		Thumbnail        string `json:"thumbnail"`
		Languages        any    `json:"languages"`
		Registrations    bool   `json:"registrations"`
		ApprovalRequired bool   `json:"approval_required"`
		InvitesEnabled   bool   `json:"invites_enabled"`
		Configuration    struct {
			Accounts struct {
				MaxFeaturedTags int `json:"max_featured_tags"`
			} `json:"accounts"`
			Statuses struct {
				MaxCharacters            int `json:"max_characters"`
				MaxMediaAttachments      int `json:"max_media_attachments"`
				CharactersReservedPerUrl int `json:"characters_reserved_per_url"`
			} `json:"statuses"`
			MediaAttachments struct {
				SupportedMimeTypes  []string `json:"supported_mime_types"`
				ImageSizeLimit      int      `json:"image_size_limit"`
				ImageMatrixLimit    int      `json:"image_matrix_limit"`
				VideoSizeLimit      int      `json:"video_size_limit"`
				VideoFrameRateLimit int      `json:"video_frame_rate_limit"`
				VideoMatrixLimit    int      `json:"video_matrix_limit"`
			} `json:"media_attachments"`
			Polls struct {
				MaxOptions             int `json:"max_options"`
				MaxCharactersPerOption int `json:"max_characters_per_option"`
				MinExpiration          int `json:"min_expiration"`
				MaxExpiration          int `json:"max_expiration"`
			} `json:"polls"`
		} `json:"configuration"`
		ContactAccount *Account       `json:"contact_account"`
		Rules          []InstanceRule `json:"rules"`
	}
	InstanceUrls struct {
		StreamingApi string `json:"streaming_api"`
	}
	InstanceRule struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	}
	InstanceIcon struct {
		Src  string `json:"src"`
		Size string `json:"size"`
	}
	VapidConfig struct {
		PublicKey string `json:"public_key,omitempty"`
	}
	InstanceV2 struct {
		Domain        string   `json:"domain"`
		Title         string   `json:"title"`
		Version       string   `json:"version"`
		SourceURL     string   `json:"source_url"`
		Description   string   `json:"description"`
		Usage         struct {
			Users struct {
				ActiveMonth int `json:"active_month"`
			} `json:"users"`
		} `json:"usage"`
		Thumbnail struct {
			URL string `json:"url"`
		} `json:"thumbnail"`
		Icon        []InstanceIcon `json:"icon"`
		Languages   []string       `json:"languages"`
		Configuration struct {
			Urls struct {
				Streaming string `json:"streaming"`
			} `json:"urls"`
			Vapid             *VapidConfig `json:"vapid,omitempty"`
			Accounts struct {
				MaxFeaturedTags    int `json:"max_featured_tags"`
				MaxPinnedStatuses int `json:"max_pinned_statuses"`
			} `json:"accounts"`
			Statuses struct {
				MaxCharacters            int `json:"max_characters"`
				MaxMediaAttachments      int `json:"max_media_attachments"`
				CharactersReservedPerUrl int `json:"characters_reserved_per_url"`
			} `json:"statuses"`
			MediaAttachments struct {
				SupportedMimeTypes  []string `json:"supported_mime_types"`
				ImageSizeLimit      int      `json:"image_size_limit"`
				ImageMatrixLimit    int      `json:"image_matrix_limit"`
				VideoSizeLimit      int      `json:"video_size_limit"`
				VideoFrameRateLimit int      `json:"video_frame_rate_limit"`
				VideoMatrixLimit    int      `json:"video_matrix_limit"`
			} `json:"media_attachments"`
			Polls struct {
				MaxOptions             int `json:"max_options"`
				MaxCharactersPerOption int `json:"max_characters_per_option"`
				MinExpiration          int `json:"min_expiration"`
				MaxExpiration          int `json:"max_expiration"`
			} `json:"polls"`
			Translation struct {
				Enabled bool `json:"enabled"`
			} `json:"translation"`
		} `json:"configuration"`
		Registrations struct {
			Enabled          bool    `json:"enabled"`
			ApprovalRequired bool    `json:"approval_required"`
			Message          *string `json:"message"`
			Url              *string `json:"url"`
		} `json:"registrations"`
		ApiVersions struct {
			Mastodon int `json:"mastodon"`
		} `json:"api_versions"`
		Contact struct {
			Email   string   `json:"email"`
			Account *Account `json:"account"`
		} `json:"contact"`
		Rules []InstanceRule `json:"rules"`
	}
)
