package yesno_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tamnd/yesno-cli/yesno"
)

const fakeAnswerJSON = `{"answer":"yes","forced":false,"image":"https://yesno.wtf/assets/yes/1.gif"}`

func newTestClient(ts *httptest.Server) *yesno.Client {
	cfg := yesno.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.Rate = 0
	return yesno.NewClient(cfg)
}

func TestRandomAnswerParsesFields(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, fakeAnswerJSON)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	a, err := c.RandomAnswer(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if a.Answer != "yes" {
		t.Errorf("Answer = %q, want yes", a.Answer)
	}
	if a.Forced {
		t.Error("Forced = true, want false")
	}
	if a.Image != "https://yesno.wtf/assets/yes/1.gif" {
		t.Errorf("Image = %q, unexpected", a.Image)
	}
}

func TestRandomAnswerSendsUserAgent(t *testing.T) {
	var gotUA string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		_, _ = fmt.Fprint(w, fakeAnswerJSON)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.RandomAnswer(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if gotUA == "" {
		t.Error("request carried no User-Agent")
	}
}

func TestRandomAnswerRetriesOn503(t *testing.T) {
	var hits int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		if hits < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		_, _ = fmt.Fprint(w, fakeAnswerJSON)
	}))
	defer ts.Close()

	cfg := yesno.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.Rate = 0
	cfg.Retries = 3
	c := yesno.NewClient(cfg)

	_, err := c.RandomAnswer(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if hits != 3 {
		t.Errorf("server saw %d hits, want 3", hits)
	}
}

func TestRandomAnswerNon200Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.RandomAnswer(context.Background())
	if err == nil {
		t.Error("expected error for 404, got nil")
	}
}

func TestRandomAnswerHitsAPIPath(t *testing.T) {
	var gotPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = fmt.Fprint(w, fakeAnswerJSON)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	_, err := c.RandomAnswer(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/api" {
		t.Errorf("path = %q, want /api", gotPath)
	}
}
