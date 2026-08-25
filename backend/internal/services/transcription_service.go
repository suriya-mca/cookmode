package services

// TranscriptionService turns recipe videos into text with word-level
// timestamps so anchor points can be inferred for each cooking step.
type TranscriptionService struct {
	apiKey string
}

func NewTranscriptionService(apiKey string) *TranscriptionService {
	return &TranscriptionService{apiKey: apiKey}
}

type TranscriptSegment struct {
	StartSeconds float64 `json:"start_seconds"`
	EndSeconds   float64 `json:"end_seconds"`
	Text         string  `json:"text"`
}

// TODO: Transcribe(videoHLSOrFileURL string) ([]TranscriptSegment, error)
