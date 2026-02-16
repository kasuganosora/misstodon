package models

type Translation struct {
	Content                string `json:"content"`
	DetectedSourceLanguage string `json:"detected_source_language"`
	Provider               string `json:"provider"`
}
