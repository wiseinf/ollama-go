package ollama

import "time"

// GenerateRequest represents a request to generate a completion
type GenerateRequest struct {
	Model  string   `json:"model"`
	Prompt string   `json:"prompt"`
	Suffix string   `json:"suffix,omitempty"`
	Images []string `json:"images,omitempty"`

	//Advanced parameters (optional)
	// json or json schema
	Format   interface{}            `json:"format,omitempty"`
	Options  map[string]interface{} `json:"options,omitempty"`
	System   string                 `json:"system,omitempty"`
	Template string                 `json:"template,omitempty"`

	Stream    bool          `json:"stream,omitempty"`
	Raw       bool          `json:"raw,omitempty"`
	KeepAlive time.Duration `json:"keep_alive,omitempty"`

	// Deprecated
	Context []int `json:"context,omitempty"`
}

// GenerateResponse represents a response from the generate endpoint
type GenerateResponse struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Response           string    `json:"response"`
	Done               bool      `json:"done"`
	Context            []int     `json:"context,omitempty"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int64     `json:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count"`
	PromptEvalDuration int64     `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"`
	EvalDuration       int64     `json:"eval_duration"`
}

// ChatMessage represents a chat message
type ChatMessage struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	Images    []string   `json:"images,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	Function struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	} `json:"function"`
}

// ChatRequest represents a request to the chat endpoint
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Tools    []Tool        `json:"tools,omitempty"`
	// Advanced parameters (optional)
	Format    interface{}            `json:"format,omitempty"`
	Stream    bool                   `json:"stream,omitempty"`
	Options   map[string]interface{} `json:"options,omitempty"`
	KeepAlive time.Duration          `json:"keep_alive,omitempty"`
}

type Tool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Parameters  struct {
			Type       string                   `json:"type"`
			Required   []string                 `json:"required"`
			Properties map[string]PropertyField `json:"properties"`
		} `json:"parameters"`
	} `json:"function"`
}

type PropertyField struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"` // Optional field
}

// ChatResponse represents a response from the chat endpoint
type ChatResponse struct {
	Model              string      `json:"model"`
	CreatedAt          time.Time   `json:"created_at"`
	Message            ChatMessage `json:"message"`
	Done               bool        `json:"done"`
	TotalDuration      int64       `json:"total_duration"`
	LoadDuration       int64       `json:"load_duration"`
	PromptEvalCount    int         `json:"prompt_eval_count"`
	PromptEvalDuration int64       `json:"prompt_eval_duration"`
	EvalCount          int         `json:"eval_count"`
	EvalDuration       int64       `json:"eval_duration"`
}

// ModelInfo represents information about a model
type ModelInfo struct {
	Name       string                 `json:"name"`
	Modified   time.Time              `json:"modified"`
	Size       int64                  `json:"size"`
	Digest     string                 `json:"digest"`
	Details    map[string]interface{} `json:"details,omitempty"`
	License    string                 `json:"license,omitempty"`
	Modelfile  string                 `json:"modelfile,omitempty"`
	Parameters string                 `json:"parameters,omitempty"`
	Template   string                 `json:"template,omitempty"`
}

// CreateModelRequest represents a request to create a model
type CreateModelRequest struct {
	Name      string `json:"name"`
	Path      string `json:"path,omitempty"`
	Modelfile string `json:"modelfile"`
}

// CopyModelRequest represents a request to copy a model
type CopyModelRequest struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

// PullModelRequest represents a request to pull a model
type PullModelRequest struct {
	Name     string `json:"name"`
	Insecure bool   `json:"insecure,omitempty"`
	Stream   bool   `json:"stream,omitempty"`
}

// PushModelRequest represents a request to push a model
type PushModelRequest struct {
	Name     string `json:"name"`
	Insecure bool   `json:"insecure,omitempty"`
	Stream   bool   `json:"stream,omitempty"`
}

// EmbeddingRequest represents a request to generate embeddings
type EmbeddingRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// EmbeddingResponse represents a response from the embedding endpoint
type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

// ModelResponse represents a response containing model status
type ModelResponse struct {
	Status string `json:"status"`
}
