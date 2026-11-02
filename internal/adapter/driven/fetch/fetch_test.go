// Package fetch: unit tests of the driven GET adapter against a local
// recording server (net/http/httptest runs on 127.0.0.1, no external
// network). The adapter is the single net/http site of the project
// (application-no-net); these tests pin its error classes: a 2xx body, a
// non-2xx status and a transport failure are all errors the caller
// of the use case classifies (SPEC-054).
package fetch

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGet_OKReturnsBody(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("BODY"))
	}))
	defer srv.Close()
	c := New(5 * time.Second)
	body, err := c.Get(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if string(body) != "BODY" {
		t.Errorf("Get: body = %q, want BODY", body)
	}
}

func TestGet_Non2xxIsAnError(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not the file", http.StatusNotFound)
	}))
	defer srv.Close()
	c := New(5 * time.Second)
	_, err := c.Get(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("Get of a 404 source returned no error, want the unreadable class")
	}
	if !strings.Contains(err.Error(), "unexpected status 404") {
		t.Errorf("err = %v, want the status in the message", err)
	}
}

func TestGet_UnreachableIsAnError(t *testing.T) {
	t.Parallel()
	c := New(1 * time.Second)
	// A port of the loopback range nothing should listen on for this test;
	// the refused connection is the transport failure the read carries.
	_, err := c.Get(context.Background(), "http://127.0.0.1:1/x")
	if err == nil {
		t.Fatal("Get of a dead address returned no error")
	}
}

func TestGet_ContextCancellation(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-time.After(5 * time.Second):
		}
	}))
	defer srv.Close()
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-time.After(100 * time.Millisecond)
		cancel()
	}()
	_, err := New(5*time.Second).Get(ctx, srv.URL)
	if err == nil {
		t.Fatal("Get on a cancelled context returned no error")
	}
	if !errors.Is(err, context.Canceled) && !strings.Contains(err.Error(), "canceled") {
		t.Errorf("err = %v, want the cancellation reason", err)
	}
}
