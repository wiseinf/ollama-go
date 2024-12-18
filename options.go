package ollama

import (
	"net/http"
	"time"
)

// ClientOptions includes options for client configuration
type ClientOptions struct {
	BaseURL          string
	HTTPClient       *http.Client
	MaxRetries       int
	RetryWaitTime    time.Duration
	RetryMaxWaitTime time.Duration
	RateLimit        int // 每秒请求数
	Timeout          time.Duration
	Debug            bool
	Logger           Logger
}

// default options
func defaultOptions() *ClientOptions {
	return &ClientOptions{
		BaseURL:          "http://localhost:11434",
		MaxRetries:       3,
		RetryWaitTime:    time.Second,
		RetryMaxWaitTime: time.Second * 30,
		RateLimit:        10,
		Timeout:          time.Minute * 5,
		Debug:            false,
		Logger:           newDefaultLogger(),
	}
}

func WithBaseURL(url string) func(*ClientOptions) {
	return func(o *ClientOptions) {
		o.BaseURL = url
	}
}

func WithHTTPClient(client *http.Client) func(*ClientOptions) {
	return func(o *ClientOptions) {
		o.HTTPClient = client
	}
}

func WithMaxRetries(retries int) func(*ClientOptions) {
	return func(o *ClientOptions) {
		o.MaxRetries = retries
	}
}

func WithRetryWaitTime(duration time.Duration) func(*ClientOptions) {
	return func(o *ClientOptions) {
		o.RetryWaitTime = duration
	}
}

func WithRateLimit(rps int) func(*ClientOptions) {
	return func(o *ClientOptions) {
		o.RateLimit = rps
	}
}

func WithDebug(debug bool) func(*ClientOptions) {
	return func(o *ClientOptions) {
		o.Debug = debug
	}
}

func WithLogger(logger Logger) func(*ClientOptions) {
	return func(o *ClientOptions) {
		o.Logger = logger
	}
}
