# Ollama Go Library

The Ollama Go library provides the easiest way to integrate your JavaScript project with [Ollama](https://github.com/ollama/ollama).

## Usage

```go

import (
    ollama "github.com/wiseinf/ollama-go"
)

...

client := ollama.NewClient(
    ollama.WithBaseURL("http://localhost:11434"),
    ollama.WithMaxRetries(3),
    ollama.WithRateLimit(10),
    ollama.WithDebug(true),
)
client.Generate(context.Background(), &GenerateRequest{
    Model:  "llama2",
    Prompt: "Hello",
})
```

## API

The Ollama Go library's API is designed around the [Ollama REST API](https://github.com/ollama/ollama/blob/main/docs/api.md).
