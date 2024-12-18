package ollama

import "fmt"

// APIError represents an error returned by the Ollama API
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("ollama api error: %s (status code: %d)", e.Message, e.StatusCode)
}
