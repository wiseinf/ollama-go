package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// Generate sends a generation request to the Ollama API
func (c *Client) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	req.Stream = false
	resp, err := c.sendRequest(ctx, "POST", "/api/generate", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GenerateStream sends a streaming generation request to the Ollama API
func (c *Client) GenerateStream(ctx context.Context, req *GenerateRequest) (<-chan GenerateStreamResponse, error) {
	req.Stream = true

	resp, err := c.sendRequest(ctx, "POST", "/api/generate", req)
	if err != nil {
		return nil, err
	}

	ch := make(chan GenerateStreamResponse)
	go func() {
		defer close(ch)
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		for {
			var response GenerateResponse
			if err := decoder.Decode(&response); err != nil {
				if err != io.EOF {
					ch <- GenerateStreamResponse{
						Error: err,
					}
				}
				return
			}
			ch <- GenerateStreamResponse{
				GenerateResponse: &response,
			}
		}
	}()

	return ch, nil
}

// Chat sends a chat request to the Ollama API
func (c *Client) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	resp, err := c.sendRequest(ctx, "POST", "/api/chat", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListModels returns a list of local models
func (c *Client) ListModels(ctx context.Context) ([]ModelInfo, error) {
	resp, err := c.sendRequest(ctx, "GET", "/api/tags", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Models []ModelInfo `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Models, nil
}

// ShowModelOptions contains optional parameters for ShowModel
type ShowModelOptions struct {
	Verbose bool `json:"verbose,omitempty"`
}

// ShowModel returns information about a specific model
func (c *Client) ShowModel(ctx context.Context, name string, opts *ShowModelOptions) (*ModelInfo, error) {
	// Create request
	req := struct {
		Model   string `json:"model"`
		Verbose bool   `json:"verbose,omitempty"`
	}{
		Model: name,
	}

	if opts != nil {
		req.Verbose = opts.Verbose
	}

	// Send POST request
	resp, err := c.sendRequest(ctx, "POST", "/api/show", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ModelInfo
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateModel creates a new model
func (c *Client) CreateModel(ctx context.Context, req *CreateModelRequest) error {
	resp, err := c.sendRequest(ctx, "POST", "/api/create", req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// CopyModel copies a model
func (c *Client) CopyModel(ctx context.Context, req *CopyModelRequest) error {
	resp, err := c.sendRequest(ctx, "POST", "/api/copy", req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// DeleteModel deletes a model from the Ollama server
func (c *Client) DeleteModel(ctx context.Context, name string) error {
	reqBody := struct {
		Model string `json:"model"`
	}{
		Model: name,
	}

	_, err := c.sendRequest(ctx, "DELETE", "/api/delete", reqBody)
	if err != nil {
		return fmt.Errorf("failed to delete model: %w", err)
	}

	return nil
}

// PullModel pulls a model from a registry
func (c *Client) PullModel(ctx context.Context, req *PullModelRequest) (<-chan ModelResponse, error) {
	req.Stream = true
	resp, err := c.sendRequest(ctx, "POST", "/api/pull", req)
	if err != nil {
		return nil, err
	}

	ch := make(chan ModelResponse)
	go func() {
		defer resp.Body.Close()
		defer close(ch)

		decoder := json.NewDecoder(resp.Body)
		for {
			var response ModelResponse
			if err := decoder.Decode(&response); err != nil {
				if err != io.EOF {
					// Handle error
				}
				return
			}
			ch <- response
		}
	}()

	return ch, nil
}

// PushModel pushes a model to a registry
func (c *Client) PushModel(ctx context.Context, req *PushModelRequest) (<-chan ModelResponse, error) {
	req.Stream = true
	resp, err := c.sendRequest(ctx, "POST", "/api/push", req)
	if err != nil {
		return nil, err
	}

	ch := make(chan ModelResponse)
	go func() {
		defer resp.Body.Close()
		defer close(ch)

		decoder := json.NewDecoder(resp.Body)
		for {
			var response ModelResponse
			if err := decoder.Decode(&response); err != nil {
				if err != io.EOF {
					// Handle error
				}
				return
			}
			ch <- response
		}
	}()

	return ch, nil
}

// Embeddings generates embeddings for the given input
func (c *Client) Embeddings(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	resp, err := c.sendRequest(ctx, "POST", "/api/embeddings", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ListRunningModels returns a list of currently running models
func (c *Client) ListRunningModels(ctx context.Context) ([]ModelInfo, error) {
	resp, err := c.sendRequest(ctx, "GET", "/api/ps", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Models []ModelInfo `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Models, nil
}

// ChatStream sends a streaming chat request to the Ollama API
func (c *Client) ChatStream(ctx context.Context, req *ChatRequest) (<-chan ChatStreamResponse, error) {
	req.Stream = true
	resp, err := c.sendRequest(ctx, "POST", "/api/chat", req)
	if err != nil {
		return nil, err
	}

	ch := make(chan ChatStreamResponse)
	go func() {
		defer resp.Body.Close()
		defer close(ch)
		decoder := json.NewDecoder(resp.Body)
		for {
			var response ChatResponse
			if err := decoder.Decode(&response); err != nil {
				if err != io.EOF {
					ch <- ChatStreamResponse{
						Error: err,
					}
				}
				return
			}
			ch <- ChatStreamResponse{
				ChatResponse: &response,
			}
		}
	}()

	return ch, nil
}
