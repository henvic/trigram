package server

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p := Params{
		Address: "127.0.0.1:9375",
	}

	var w sync.WaitGroup
	w.Add(1)

	go runServer(ctx, t, p, w.Done)

	testGenerateBeforeAvailability(t, p)
	testGenerateInvalidMethod(t, p)
	testLearnInvalidMethod(t, p)
	testLearnContentTypeNotAccepted(t, p)
	testLearn(t, p)
	testGenerate(t, p)

	w.Wait()
}

func testGenerateBeforeAvailability(t *testing.T, p Params) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/generate", p.Address), nil)

	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status to be %v, got %v instead", http.StatusNotFound, resp.StatusCode)
	}
}

func testGenerateInvalidMethod(t *testing.T, p Params) {
	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("http://%s/generate", p.Address), nil)

	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status to be %v, got %v instead", http.StatusMethodNotAllowed, resp.StatusCode)
	}
}

func testLearnInvalidMethod(t *testing.T, p Params) {
	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("http://%s/learn", p.Address), nil)

	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status to be %v, got %v instead", http.StatusMethodNotAllowed, resp.StatusCode)
	}
}

func testLearnContentTypeNotAccepted(t *testing.T, p Params) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/learn", p.Address), nil)

	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	if resp.StatusCode != http.StatusNotAcceptable {
		t.Errorf("Expected status to be %v, got %v instead", http.StatusNotAcceptable, resp.StatusCode)
	}
}

func testLearn(t *testing.T, p Params) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/learn", p.Address), nil)

	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	req.Header.Set("Content-Type", "text/plain; charset=UTF-8")

	req.Body = ioutil.NopCloser(bytes.NewBufferString("To be, or not to be, that is the question."))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status to be %v, got %v instead", http.StatusOK, resp.StatusCode)
	}
}

func testGenerate(t *testing.T, p Params) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/generate", p.Address), nil)

	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status to be %v, got %v instead", http.StatusOK, resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("Expected no error, got %v instead", err)
	}

	if len(b) == 0 {
		t.Error("No text received")
	}
}

func runServer(ctx context.Context, t *testing.T, p Params, done func()) {
	defer done()

	if err := Run(ctx, p); err != http.ErrServerClosed {
		t.Errorf("Expected error %v, got %v instead", http.ErrServerClosed, err)
	}
}
