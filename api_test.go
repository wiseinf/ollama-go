package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

// Setup test server and client
func setupTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	server := httptest.NewServer(handler)
	client := NewClient(WithBaseURL(server.URL))
	return server, client
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		name     string
		request  *GenerateRequest
		response *GenerateResponse
		wantErr  bool
	}{
		{
			name: "successful generation",
			request: &GenerateRequest{
				Model:  "llama3.2:1b",
				Prompt: "Hello",
			},
			response: &GenerateResponse{
				Model:    "llama3.2:1b",
				Response: "Hi there!",
				Done:     true,
			},
			wantErr: false,
		},
		{
			name: "error response",
			request: &GenerateRequest{
				Model: "invalid-model",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(map[string]string{"error": "invalid model"})
					return
				}
				json.NewEncoder(w).Encode(tt.response)
			})
			defer server.Close()

			resp, err := client.Generate(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(resp, tt.response) {
				t.Errorf("Generate() got = %v, want %v", resp, tt.response)
			}
		})
	}
}

func TestGenerateStream(t *testing.T) {
	responses := []GenerateResponse{
		{Response: "Hello", Done: false},
		{Response: " World", Done: false},
		{Response: "!", Done: true},
	}

	server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("Expected http.Flusher")
		}

		for _, resp := range responses {
			json.NewEncoder(w).Encode(resp)
			flusher.Flush()
		}
	})
	defer server.Close()

	stream, err := client.GenerateStream(context.Background(), &GenerateRequest{
		Model:  "llama3.2:1b",
		Prompt: "Hello",
	})
	if err != nil {
		t.Fatalf("GenerateStream() error = %v", err)
	}

	var received []GenerateResponse
	for response := range stream {
		if response.Error != nil {
			received = append(received, *response.GenerateResponse)
		}
	}

	if !reflect.DeepEqual(received, responses) {
		t.Errorf("GenerateStream() got = %v, want %v", received, responses)
	}
}

func TestChat(t *testing.T) {
	tests := []struct {
		name     string
		request  *ChatRequest
		response *ChatResponse
		wantErr  bool
	}{
		{
			name: "successful chat",
			request: &ChatRequest{
				Model: "llama3.2:1b",
				Messages: []ChatMessage{
					{Role: "user", Content: "Hi"},
				},
			},
			response: &ChatResponse{
				Model: "llama3.2:1b",
				Message: ChatMessage{
					Role:    "assistant",
					Content: "Hello! How can I help you?",
				},
				Done: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				json.NewEncoder(w).Encode(tt.response)
			})
			defer server.Close()

			resp, err := client.Chat(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(resp, tt.response) {
				t.Errorf("Chat() got = %v, want %v", resp, tt.response)
			}
		})
	}
}

func TestListModels(t *testing.T) {
	expectedModels := []ModelInfo{
		{Name: "llama3.2:1b", Size: 1000},
		{Name: "mistral", Size: 2000},
	}

	server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string][]ModelInfo{"models": expectedModels})
	})
	defer server.Close()

	models, err := client.ListModels(context.Background())
	if err != nil {
		t.Fatalf("ListModels() error = %v", err)
	}

	if !reflect.DeepEqual(models, expectedModels) {
		t.Errorf("ListModels() got = %v, want %v", models, expectedModels)
	}
}

func TestShowModel(t *testing.T) {
	tests := []struct {
		name         string
		modelName    string
		opts         *ShowModelOptions
		expectedBody map[string]interface{}
		response     *ModelInfo
		wantErr      bool
	}{
		{
			name:      "basic show model",
			modelName: "llama2",
			opts:      nil,
			expectedBody: map[string]interface{}{
				"model": "llama2",
			},
			response: &ModelInfo{
				Name: "llama2",
				Size: 1000,
			},
			wantErr: false,
		},
		{
			name:      "show model with verbose",
			modelName: "llama2",
			opts: &ShowModelOptions{
				Verbose: true,
			},
			expectedBody: map[string]interface{}{
				"model":   "llama2",
				"verbose": true,
			},
			response: &ModelInfo{
				Name:      "llama2",
				Size:      1000,
				Modelfile: "FROM llama2\nPARAMETER temperature 0.7",
				Template:  "<prompt>",
				License:   "MIT",
				Details: map[string]interface{}{
					"format":       "gguf",
					"quantization": "q4_0",
				},
			},
			wantErr: false,
		},
		{
			name:      "model not found",
			modelName: "nonexistent",
			opts:      nil,
			expectedBody: map[string]interface{}{
				"model": "nonexistent",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST request, got %s", r.Method)
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				var reqBody map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("Failed to decode request body: %v", err)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if !reflect.DeepEqual(reqBody, tt.expectedBody) {
					t.Errorf("Request body mismatch\ngot:  %v\nwant: %v", reqBody, tt.expectedBody)
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				if tt.wantErr {
					w.WriteHeader(http.StatusNotFound)
					json.NewEncoder(w).Encode(map[string]string{
						"error": "model not found",
					})
					return
				}

				json.NewEncoder(w).Encode(tt.response)
			})
			defer server.Close()

			model, err := client.ShowModel(context.Background(), tt.modelName, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("ShowModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(model, tt.response) {
				t.Errorf("ShowModel() got = %v, want %v", model, tt.response)
			}
		})
	}
}

func TestCreateModel(t *testing.T) {
	tests := []struct {
		name    string
		request *CreateModelRequest
		wantErr bool
	}{
		{
			name: "successful creation",
			request: &CreateModelRequest{
				Name:      "custom-model",
				Modelfile: "FROM llama3.2:1b",
			},
			wantErr: false,
		},
		{
			name: "invalid request",
			request: &CreateModelRequest{
				Name: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				if tt.wantErr {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				w.WriteHeader(http.StatusOK)
			})
			defer server.Close()

			err := client.CreateModel(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateModel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCopyModel(t *testing.T) {
	server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		var req CopyModelRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if req.Source == "" || req.Destination == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	err := client.CopyModel(context.Background(), &CopyModelRequest{
		Source:      "llama3.2:1b",
		Destination: "llama3.2:1b-copy",
	})
	if err != nil {
		t.Errorf("CopyModel() error = %v", err)
	}
}

func TestDeleteModel(t *testing.T) {
	server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Query().Get("name") == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	err := client.DeleteModel(context.Background(), "llama3.2:1b")
	if err != nil {
		t.Errorf("DeleteModel() error = %v", err)
	}
}

func TestPullModel(t *testing.T) {
	responses := []ModelResponse{
		{Status: "downloading manifest"},
		{Status: "downloading weights"},
		{Status: "verifying checksums"},
		{Status: "completed"},
	}

	server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("Expected http.Flusher")
		}

		for _, resp := range responses {
			json.NewEncoder(w).Encode(resp)
			flusher.Flush()
		}
	})
	defer server.Close()

	stream, err := client.PullModel(context.Background(), &PullModelRequest{
		Name: "llama3.2:1b",
	})
	if err != nil {
		t.Fatalf("PullModel() error = %v", err)
	}

	var received []ModelResponse
	for response := range stream {
		received = append(received, response)
	}

	if !reflect.DeepEqual(received, responses) {
		t.Errorf("PullModel() got = %v, want %v", received, responses)
	}
}

func TestEmbeddings(t *testing.T) {
	expectedEmbedding := &EmbeddingResponse{
		Embedding: []float32{0.1, 0.2, 0.3},
	}

	server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(expectedEmbedding)
	})
	defer server.Close()

	embedding, err := client.Embeddings(context.Background(), &EmbeddingRequest{
		Model:  "llama3.2:1b",
		Prompt: "Hello world",
	})
	if err != nil {
		t.Fatalf("Embeddings() error = %v", err)
	}

	if !reflect.DeepEqual(embedding, expectedEmbedding) {
		t.Errorf("Embeddings() got = %v, want %v", embedding, expectedEmbedding)
	}
}

func TestListRunningModels(t *testing.T) {
	expectedModels := []ModelInfo{
		{Name: "llama3.2:1b", Size: 1000},
	}

	server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string][]ModelInfo{"models": expectedModels})
	})
	defer server.Close()

	models, err := client.ListRunningModels(context.Background())
	if err != nil {
		t.Fatalf("ListRunningModels() error = %v", err)
	}

	if !reflect.DeepEqual(models, expectedModels) {
		t.Errorf("ListRunningModels() got = %v, want %v", models, expectedModels)
	}
}

func TestContextCancellation(t *testing.T) {
	server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := client.Generate(ctx, &GenerateRequest{
		Model:  "llama3.2:1b",
		Prompt: "test",
	})

	if err == nil {
		t.Error("Expected context deadline exceeded error, got nil")
	}
}

func TestConcurrentRequests(t *testing.T) {
	server, client := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	defer server.Close()

	concurrency := 10
	errors := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(i int) {
			_, err := client.Generate(context.Background(), &GenerateRequest{
				Model:  "llama3.2:1b",
				Prompt: fmt.Sprintf("test-%d", i),
			})
			errors <- err
		}(i)
	}

	for i := 0; i < concurrency; i++ {
		if err := <-errors; err != nil {
			t.Errorf("Concurrent request %d failed: %v", i, err)
		}
	}
}
