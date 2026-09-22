package typesafe

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestClient(server *httptest.Server, logger *slog.Logger) *Client {
	return NewClientWithOptions("test-key", logger, WithBaseURL(server.URL))
}

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestCall_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("missing or invalid authorization header")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"model": "jev-1.13.0",
			"answers": {"test": {"type": "choice", "choice": "success"}},
			"usage": {"input_tokens": 100, "output_tokens": 50},
			"request_id": "test-123",
			"evaluation_time_ms": 45.5
		}`))
	}))
	defer server.Close()

	client := newTestClient(server, newTestLogger())
	result, err := client.Call(context.Background(), &Request{
		State: map[string]any{"test": "state"},
		Model: "jev-latest",
		Questions: map[string]any{"test": "question"},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if result.Model != "jev-1.13.0" {
		t.Errorf("expected model='jev-1.13.0', got %v", result.Model)
	}

	if len(result.Answers) == 0 {
		t.Error("expected non-empty answers")
	}
}

func TestCall_RetryOn500(t *testing.T) {
	attempt := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt++
		if attempt == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"server error"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"model": "jev-1.13.0",
			"answers": {"test": {"type": "choice", "choice": "success"}},
			"usage": {"input_tokens": 100, "output_tokens": 50},
			"request_id": "test-456",
			"evaluation_time_ms": 50.0
		}`))
	}))
	defer server.Close()

	client := newTestClient(server, newTestLogger())
	result, err := client.Call(context.Background(), &Request{
		State: map[string]any{"test": "state"},
		Model: "jev-latest",
		Questions: map[string]any{"test": "question"},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if attempt != 2 {
		t.Errorf("expected 2 attempts, got %d", attempt)
	}

	if result == nil || result.Model != "jev-1.13.0" {
		t.Error("expected successful result after retry")
	}
}

func TestCall_NoRetryOn400(t *testing.T) {
	attempt := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt++
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"bad request"}`))
	}))
	defer server.Close()

	client := newTestClient(server, newTestLogger())
	_, err := client.Call(context.Background(), &Request{
		State: map[string]any{"test": "state"},
		Model: "jev-latest",
		Questions: map[string]any{"test": "question"},
	})

	if err == nil {
		t.Fatal("expected error for 400 response")
	}

	if attempt != 1 {
		t.Errorf("expected 1 attempt (no retries), got %d", attempt)
	}
}

func TestCall_RetryOn429(t *testing.T) {
	attempt := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt++
		if attempt == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":"rate limited"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"model": "jev-1.13.0",
			"answers": {"test": {"type": "choice", "choice": "success"}},
			"usage": {"input_tokens": 100, "output_tokens": 50},
			"request_id": "test-789",
			"evaluation_time_ms": 55.0
		}`))
	}))
	defer server.Close()

	client := newTestClient(server, newTestLogger())
	result, err := client.Call(context.Background(), &Request{
		State: map[string]any{"test": "state"},
		Model: "jev-latest",
		Questions: map[string]any{"test": "question"},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if attempt != 2 {
		t.Errorf("expected 2 attempts, got %d", attempt)
	}

	if result == nil || result.Model != "jev-1.13.0" {
		t.Error("expected successful result after 429 retry")
	}
}

func TestCall_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"model": "jev-1.13.0",
			"answers": {},
			"usage": {"input_tokens": 0, "output_tokens": 0},
			"request_id": "test",
			"evaluation_time_ms": 0
		}`))
	}))
	defer server.Close()

	client := newTestClient(server, newTestLogger())

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := client.Call(ctx, &Request{
		State: map[string]any{"test": "state"},
		Model: "jev-latest",
		Questions: map[string]any{"test": "question"},
	})

	if err == nil {
		t.Fatal("expected error due to context cancellation")
	}
}

func TestCall_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{invalid json}`))
	}))
	defer server.Close()

	client := newTestClient(server, newTestLogger())
	_, err := client.Call(context.Background(), &Request{
		State: map[string]any{"test": "state"},
		Model: "jev-latest",
		Questions: map[string]any{"test": "question"},
	})

	if err == nil {
		t.Fatal("expected error for invalid JSON response")
	}
}
