package typesafe

// Request represents the TypeSafe SystemOne API request
type Request struct {
	Context   map[string]any `json:"context"`
	Questions map[string]any `json:"questions"`
}

// Response represents the TypeSafe SystemOne API response
type Response struct {
	Model          string         `json:"model"`
	Answers        map[string]Answer `json:"answers"`
	Usage          Usage          `json:"usage"`
	RequestID      string         `json:"request_id"`
	EvaluationTime float64        `json:"evaluation_time_ms"`
}

// Answer represents a single answer from TypeSafe
type Answer struct {
	Type          string             `json:"type"`
	Choice        string             `json:"choice,omitempty"`
	Score         *float64           `json:"score,omitempty"`
	Noul          *float64           `json:"noul,omitempty"`
	Legend        map[string]string  `json:"legend,omitempty"`
	Confidence    *float64           `json:"confidence,omitempty"`
	Probabilities map[string]float64 `json:"probabilities,omitempty"`
	Stats         map[string]any     `json:"stats,omitempty"`
}

// Usage contains token usage information
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}
