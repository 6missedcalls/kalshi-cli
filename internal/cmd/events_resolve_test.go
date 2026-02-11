package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/6missedcalls/kalshi-cli/internal/api"
	"github.com/6missedcalls/kalshi-cli/internal/config"
	"github.com/6missedcalls/kalshi-cli/pkg/models"
)

func newCmdTestClient(t *testing.T, serverURL string) *api.Client {
	t.Helper()

	cfg := &config.Config{
		API: config.APIConfig{
			Production: false,
			Timeout:    5 * time.Second,
		},
	}

	client := api.NewClient(cfg, nil)
	client.SetBaseURL(serverURL)

	return client
}

func TestResolveSeriesTicker_ExplicitSeriesUsedDirectly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("API should not be called when --series is provided explicitly")
	}))
	defer server.Close()

	client := newCmdTestClient(t, server.URL)
	ctx := context.Background()

	series, err := resolveSeriesTicker(ctx, client, "KXELONMARS-99", "KXELONMARS")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if series != "KXELONMARS" {
		t.Errorf("expected series ticker %q, got %q", "KXELONMARS", series)
	}
}

func TestResolveSeriesTicker_AutoResolvesFromEvent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := api.TradeAPIPrefix + "/events/KXELONMARS-99"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.EventResponse{
			Event: models.Event{
				EventTicker:  "KXELONMARS-99",
				SeriesTicker: "KXELONMARS",
				Title:        "Will Elon reach Mars by 2099?",
			},
		})
	}))
	defer server.Close()

	client := newCmdTestClient(t, server.URL)
	ctx := context.Background()

	series, err := resolveSeriesTicker(ctx, client, "KXELONMARS-99", "")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if series != "KXELONMARS" {
		t.Errorf("expected auto-resolved series ticker %q, got %q", "KXELONMARS", series)
	}
}

func TestResolveSeriesTicker_EventNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"code":    "not_found",
			"message": "event not found",
		})
	}))
	defer server.Close()

	client := newCmdTestClient(t, server.URL)
	ctx := context.Background()

	_, err := resolveSeriesTicker(ctx, client, "NONEXISTENT-EVENT", "")

	if err == nil {
		t.Fatal("expected error when event is not found, got nil")
	}
}

func TestResolveSeriesTicker_EventHasEmptySeriesTicker(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.EventResponse{
			Event: models.Event{
				EventTicker:  "BROKEN-EVENT",
				SeriesTicker: "",
				Title:        "Event with no series",
			},
		})
	}))
	defer server.Close()

	client := newCmdTestClient(t, server.URL)
	ctx := context.Background()

	_, err := resolveSeriesTicker(ctx, client, "BROKEN-EVENT", "")

	if err == nil {
		t.Fatal("expected error when event has empty series ticker, got nil")
	}
}

func TestCandlesticksCmd_SeriesFlagIsOptional(t *testing.T) {
	flag := eventsCandlesticksCmd.Flags().Lookup("series")
	if flag == nil {
		t.Fatal("expected --series flag to be registered")
	}

	// Check that the flag is not required by looking for Cobra's required annotation
	annotations := flag.Annotations
	if annotations != nil {
		if _, ok := annotations["cobra_annotation_bash_completion_one_required_flag"]; ok {
			t.Error("--series flag should NOT be marked as required")
		}
	}
}

func TestCandlesticksCmd_SeriesFlagDescription(t *testing.T) {
	flag := eventsCandlesticksCmd.Flags().Lookup("series")
	if flag == nil {
		t.Fatal("expected --series flag to be registered")
	}

	if flag.Usage == "series ticker (required for candlesticks)" {
		t.Error("flag description still says 'required'; should indicate auto-resolution")
	}
}
