package specview

// TokenUsage represents token consumption for a single AI call.
type TokenUsage struct {
	CandidatesTokens int32
	Model            string
	PromptTokens     int32
	TotalTokens      int32
}

// Add combines two TokenUsage values.
func (t TokenUsage) Add(other TokenUsage) TokenUsage {
	return TokenUsage{
		CandidatesTokens: t.CandidatesTokens + other.CandidatesTokens,
		Model:            t.Model,
		PromptTokens:     t.PromptTokens + other.PromptTokens,
		TotalTokens:      t.TotalTokens + other.TotalTokens,
	}
}
