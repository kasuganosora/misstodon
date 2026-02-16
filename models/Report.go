package models

type Report struct {
	ID            string  `json:"id"`
	ActionTaken   bool    `json:"action_taken"`
	ActionTakenAt *string `json:"action_taken_at"`
	Category      string  `json:"category"`
	Comment       string  `json:"comment"`
	Forwarded     bool    `json:"forwarded"`
	CreatedAt     string  `json:"created_at"`
	StatusIDs     []string `json:"status_ids"`
	RuleIDs       []string `json:"rule_ids"`
	TargetAccount Account  `json:"target_account"`
}
