package observatory

// Recommendation is the explainable advice object (ADR-0006).
type Recommendation struct {
	ID               string   `json:"id"`
	Product          string   `json:"product"`
	Title            string   `json:"title"`
	Action           string   `json:"action"`
	Evidence         []string `json:"evidence"`
	Confidence       float64  `json:"confidence"` // 0–1
	ExpectedImpact   string   `json:"expected_impact"`
	EstimatedEffort  string   `json:"estimated_effort"`
	SupportingScores []string `json:"supporting_scores,omitempty"`
	FiveQuestions    []string `json:"five_questions_addressed,omitempty"`
}
