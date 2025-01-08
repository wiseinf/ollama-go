package main

import (
	"context"
	"fmt"
	"log"
	"time"

	ollama "github.com/wiseinf/ollama-go"
)

func main() {
	ctx := context.Background()
	// Using a different base URL
	// client := ollama.NewClient(ollama.WithBaseURL("http://127.0.0.1:11435"))
	client := ollama.NewClient()

	// List local models
	models, err := client.ListModels(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, model := range models {
		fmt.Printf("Model: %s, Size: %d\n", model.Name, model.Size)
	}

	// Pull model
	pullChan, err := client.PullModel(ctx, &ollama.PullModelRequest{
		Name: "llama3.2:1b",
	})
	if err != nil {
		log.Fatal(err)
	}
	for status := range pullChan {
		fmt.Println("Pull status:", status.Status)
	}

	// Generate embeddings
	embeddings, err := client.Embeddings(ctx, &ollama.EmbeddingRequest{
		Model:  "llama3.2:1b",
		Prompt: "Hello world",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Embeddings: %v\n", embeddings.Embedding[:5]) // 只打印前5个值

	// Copy model
	err = client.CopyModel(ctx, &ollama.CopyModelRequest{
		Source:      "llama3.2:1b",
		Destination: "llama3.2:1b-copy",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Relist local models
	models, err = client.ListModels(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, model := range models {
		fmt.Printf("Model: %s, Size: %d\n", model.Name, model.Size)
	}

	// Show model
	modelInfo, err := client.ShowModel(ctx, "llama3.2:1b", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Model info: %+v\n", modelInfo)

	// List running model
	runningModels, err := client.ListRunningModels(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, model := range runningModels {
		fmt.Printf("Running model: %s\n", model.Name)
	}

	resp, err := client.Chat(ctx, &ollama.ChatRequest{
		Model: "llama3.2:1b",
		Messages: []ollama.ChatMessage{
			{
				Role:    "user",
				Content: "why is the sky blue?",
			},
		},
		KeepAlive: ollama.Duration(5 * time.Minute),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Chat response: %+v\n", resp)

	generateRes, err := client.Generate(ctx, &ollama.GenerateRequest{
		Model:     "llama3.2:1b",
		Prompt:    "Hello world",
		KeepAlive: ollama.Duration(5 * time.Minute),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generate response: %+v\n", generateRes)

	if err := client.DeleteModel(ctx, "llama3.2:1b-copy"); err != nil {
		log.Fatal(err)
	}

	// Show model should report error
	if _, err := client.ShowModel(ctx, "llama3.2:1b-copy", nil); err == nil {
		log.Fatal("model still exists after deletion")
	}

	// Stream chat
	req := &ollama.ChatRequest{
		Model: "llama3.2:1b",
		Messages: []ollama.ChatMessage{
			{
				Role:    "user",
				Content: "why is the sky blue?",
			},
		},
	}

	stream, err := client.ChatStream(ctx, req)

	if err != nil {
		log.Fatalf("chat stream error: %v", err)
	}
	for item := range stream {
		if item.Error != nil {
			fmt.Printf("Chat stream response: %+v\n", item.Error)
		} else {
			fmt.Printf("Chat stream response: %+v\n", item.ChatResponse)
		}
	}

	generateStream, err := client.GenerateStream(ctx, &ollama.GenerateRequest{
		Model:     "llama3.2:1b",
		Prompt:    "Hello world",
		KeepAlive: ollama.Duration(5 * time.Minute),
	})
	if err != nil {
		log.Fatalf("generate stream error: %v", err)
	}
	for item := range generateStream {
		if item.Error != nil {
			fmt.Printf("Generate stream response: %+v\n", item.Error)
		} else {
			fmt.Printf("Generate stream response: %+v\n", item.GenerateResponse)
		}
	}

	fmt.Printf("Finished.")
}
