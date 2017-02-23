// Copyright 2015, 2016, 2017 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the COPYING file.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestStatusError(t *testing.T) {
	se := StatusError{200, errors.New("OK")}
	if se.Status() != 200 {
		t.Error("expected 200, got", se.Status())
	}
	if se.Error() != "OK" {
		t.Error("expected OK, got", se.Error())
	}
}

func GetRIDTestHandler() http.HandlerFunc {
	fn := func(rw http.ResponseWriter, req *http.Request) {
		if rid := RequestID(req.Context()); rid != currentID {
			http.Error(rw, fmt.Sprintf(
				"expected request ID %d, got %d", currentID,
				rid), 500)
		}
	}
	return http.HandlerFunc(fn)
}

func TestRequestIDMiddleware(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	ctx := NewContext(req.Context(), req)
	ridm := RequestIDMiddleware(GetRIDTestHandler())
	ridm.ServeHTTP(w, req.WithContext(ctx))
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler status code: got %v, want %v: %s",
			status, http.StatusOK, w.Body.String())
	}
}

func TestErrorHandlerLogging(t *testing.T) {
	var str bytes.Buffer
	log.SetOutput(&str)
	defer log.SetOutput(os.Stderr)
	log.SetFlags(0)
	defer log.SetFlags(log.Ldate | log.Ltime)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	ctx := NewContext(req.Context(), req)
	ridm := ErrorHandler(func(rw http.ResponseWriter, req *http.Request) error {
		return nil
	})
	ridm.ServeHTTP(w, req.WithContext(ctx))
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler status code: got %v, want %v: %s",
			status, http.StatusOK, w.Body.String())
	}
	expected := fmt.Sprintf("[%d] %s - %s %s %s\n", currentID, req.RemoteAddr,
		req.Method, req.URL.RequestURI(), req.Proto)
	if str.String() != expected {
		t.Errorf("expected '%s', got '%s'", expected, str.String())
	}
}

func TestErrorHandlerLoggingXRealIP(t *testing.T) {
	var str bytes.Buffer
	log.SetOutput(&str)
	defer log.SetOutput(os.Stderr)
	log.SetFlags(0)
	defer log.SetFlags(log.Ldate | log.Ltime)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Real-IP", "1.2.3.4")
	ctx := NewContext(req.Context(), req)
	ridm := ErrorHandler(func(rw http.ResponseWriter, req *http.Request) error {
		return nil
	})
	ridm.ServeHTTP(w, req.WithContext(ctx))
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler status code: got %v, want %v: %s",
			status, http.StatusOK, w.Body.String())
	}
	expected := fmt.Sprintf("[%d] 1.2.3.4 (%s) - %s %s %s\n", currentID,
		req.RemoteAddr, req.Method, req.URL.RequestURI(), req.Proto)
	if str.String() != expected {
		t.Errorf("expected '%s', got '%s'", expected, str.String())
	}
}

func TestErrorHandlerLoggingXForwardedFor(t *testing.T) {
	var str bytes.Buffer
	log.SetOutput(&str)
	defer log.SetOutput(os.Stderr)
	log.SetFlags(0)
	defer log.SetFlags(log.Ldate | log.Ltime)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-For", "1.2.3.4")
	ctx := NewContext(req.Context(), req)
	ridm := ErrorHandler(func(rw http.ResponseWriter, req *http.Request) error {
		return nil
	})
	ridm.ServeHTTP(w, req.WithContext(ctx))
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler status code: got %v, want %v: %s",
			status, http.StatusOK, w.Body.String())
	}
	expected := fmt.Sprintf("[%d] 1.2.3.4 (%s) - %s %s %s\n", currentID,
		req.RemoteAddr, req.Method, req.URL.RequestURI(), req.Proto)
	if str.String() != expected {
		t.Errorf("expected '%s', got '%s'", expected, str.String())
	}
}

func TestErrorHandlerLoggingErrorFailure(t *testing.T) {
	var str bytes.Buffer
	log.SetOutput(&str)
	defer log.SetOutput(os.Stderr)
	log.SetFlags(0)
	defer log.SetFlags(log.Ldate | log.Ltime)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	ctx := NewContext(req.Context(), req)
	ridm := ErrorHandler(func(rw http.ResponseWriter, req *http.Request) error {
		return StatusError{420, errors.New("error")}
	})
	ridm.ServeHTTP(w, req.WithContext(ctx))
	if status := w.Code; status != 420 {
		t.Errorf("handler status code: got %v, want %v: %s",
			status, 420, w.Body.String())
	}
	expected := fmt.Sprintf("[%d] %s - %s %s %s\n[%d] %s\n", currentID,
		req.RemoteAddr, req.Method, req.URL.RequestURI(), req.Proto,
		currentID, "error")
	if str.String() != expected {
		t.Errorf("expected '%s', got '%s'", expected, str.String())
	}
	expected = "error\n"
	if w.Body.String() != expected {
		t.Errorf("expected '%s', got '%s'", expected, w.Body.String())
	}
}

func TestErrorHandlerLoggingGeneralFailure(t *testing.T) {
	var str bytes.Buffer
	log.SetOutput(&str)
	defer log.SetOutput(os.Stderr)
	log.SetFlags(0)
	defer log.SetFlags(log.Ldate | log.Ltime)

	msg := "error"

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	ctx := NewContext(req.Context(), req)
	ridm := ErrorHandler(func(rw http.ResponseWriter, req *http.Request) error {
		return errors.New(msg)
	})
	ridm.ServeHTTP(w, req.WithContext(ctx))
	if status := w.Code; status != http.StatusInternalServerError {
		t.Errorf("handler status code: got %v, want %v: %s",
			status, http.StatusInternalServerError,
			w.Body.String())
	}
	http500 := http.StatusText(http.StatusInternalServerError)
	expected := fmt.Sprintf("[%d] %s - %s %s %s\n[%d] %s - %s\n", currentID,
		req.RemoteAddr, req.Method, req.URL.RequestURI(), req.Proto,
		currentID, http500, msg)
	if str.String() != expected {
		t.Errorf("expected '%s', got '%s'", expected, str.String())
	}
	expected = fmt.Sprintf("%s - %s\n", http500, msg)
	if w.Body.String() != expected {
		t.Errorf("expected '%s', got '%s'", expected, w.Body.String())
	}
}
