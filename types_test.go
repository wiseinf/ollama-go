package ollama

import (
	"encoding/json"
	"testing"
	"time"
)

func TestChatRequestMarshalling(t *testing.T) {
	req := ChatRequest{
		Model: "llama2",
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: "Hello",
			},
		},
		KeepAlive: Duration(5 * time.Minute),
	}

	got, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal ChatRequest: %v", err)
	}

	want := `{"model":"llama2","messages":[{"role":"user","content":"Hello"}],"keep_alive":"5m"}`

	if string(got) != want {
		t.Errorf("json.Marshal(req) = %v, want %v", string(got), want)
	}
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Duration
		wantErr bool
	}{
		{
			name:    "empty string",
			input:   `""`,
			want:    Duration(0),
			wantErr: false,
		},
		{
			name:    "simple minutes",
			input:   `"5m"`,
			want:    Duration(5 * time.Minute),
			wantErr: false,
		},
		{
			name:    "hours and minutes",
			input:   `"1h30m"`,
			want:    Duration(90 * time.Minute),
			wantErr: false,
		},
		{
			name:    "complex duration",
			input:   `"2h45m30s"`,
			want:    Duration(2*time.Hour + 45*time.Minute + 30*time.Second),
			wantErr: false,
		},
		{
			name:    "days",
			input:   `"2d"`,
			want:    Duration(48 * time.Hour),
			wantErr: false,
		},
		{
			name:    "days and hours",
			input:   `"2d12h"`,
			want:    Duration(60 * time.Hour),
			wantErr: false,
		},
		{
			name:    "invalid unit",
			input:   `"5x"`,
			wantErr: true,
		},
		{
			name:    "invalid format",
			input:   `"5m4"`,
			wantErr: true,
		},
		{
			name:    "invalid number",
			input:   `"am"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Duration
			err := json.Unmarshal([]byte(tt.input), &d)

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && d != tt.want {
				t.Errorf("UnmarshalJSON() = %v, want %v", d, tt.want)
			}
		})
	}
}

func TestDuration_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		d       Duration
		want    string
		wantErr bool
	}{
		{
			name: "zero duration",
			d:    Duration(0),
			want: `""`,
		},
		{
			name: "simple seconds",
			d:    Duration(45 * time.Second),
			want: `"45s"`,
		},
		{
			name: "simple minutes",
			d:    Duration(5 * time.Minute),
			want: `"5m"`,
		},
		{
			name: "simple hours",
			d:    Duration(2 * time.Hour),
			want: `"2h"`,
		},
		{
			name: "simple days",
			d:    Duration(3 * 24 * time.Hour),
			want: `"3d"`,
		},
		{
			name: "minutes and seconds",
			d:    Duration(5*time.Minute + 30*time.Second),
			want: `"5m30s"`,
		},
		{
			name: "hours and minutes",
			d:    Duration(2*time.Hour + 30*time.Minute),
			want: `"2h30m"`,
		},
		{
			name: "hours, minutes and seconds",
			d:    Duration(2*time.Hour + 30*time.Minute + 45*time.Second),
			want: `"2h30m45s"`,
		},
		{
			name: "hours, and seconds",
			d:    Duration(2*time.Hour + 45*time.Second),
			want: `"2h45s"`,
		},
		{
			name: "days and hours",
			d:    Duration(2*24*time.Hour + 5*time.Hour),
			want: `"2d5h"`,
		},
		{
			name: "days, hours and minutes",
			d:    Duration(2*24*time.Hour + 5*time.Hour + 30*time.Minute),
			want: `"2d5h30m"`,
		},
		{
			name: "complex full duration",
			d:    Duration(2*24*time.Hour + 5*time.Hour + 30*time.Minute + 15*time.Second),
			want: `"2d5h30m15s"`,
		},
		{
			name: "59 seconds",
			d:    Duration(59 * time.Second),
			want: `"59s"`,
		},
		{
			name: "59 minutes",
			d:    Duration(59 * time.Minute),
			want: `"59m"`,
		},
		{
			name: "23 hours",
			d:    Duration(23 * time.Hour),
			want: `"23h"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("MarshalJSON() = %v, want %v", string(got), tt.want)
			}

			if !tt.wantErr {
				var decoded Duration
				err = json.Unmarshal(got, &decoded)
				if err != nil {
					t.Errorf("Failed to unmarshal marshaled value: %v", err)
				}
				if decoded != tt.d {
					t.Errorf("Marshal/Unmarshal roundtrip failed: got %v, want %v", decoded, tt.d)
				}
			}
		})
	}
}
