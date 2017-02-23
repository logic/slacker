// Copyright 2015 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the COPYING file.

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

const requestIDKey uint64 = 0

var currentID uint64

// NewContext creates a context with a request ID attached to it.
func NewContext(ctx context.Context, req *http.Request) context.Context {
	id := atomic.AddUint64(&currentID, 1)
	return context.WithValue(ctx, requestIDKey, id)
}

// RequestID retrieves a request ID from a context.
func RequestID(ctx context.Context) uint64 {
	return ctx.Value(requestIDKey).(uint64)
}

// RequestIDMiddleware adds a request id context to the request.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := NewContext(req.Context(), req)
		next.ServeHTTP(rw, req.WithContext(ctx))
	})
}

// Error extends the error interface to add HTTP status information.
type Error interface {
	error
	Status() int
}

// StatusError implements Error, storing both an error and an HTTP status code.
type StatusError struct {
	Code int
	Err  error
}

// Error returns the text of the returned error.
func (se StatusError) Error() string { return se.Err.Error() }

// Status returns the HTTP error code of our raised error.
func (se StatusError) Status() int { return se.Code }

// ErrorHandler extended the Handler interface to include an error return value.
type ErrorHandler func(http.ResponseWriter, *http.Request) error

// ErrorHandler.ServeHTTP adds logging and error handling to a standard Handler
func (h ErrorHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	addr := req.Header.Get("X-Real-IP")
	if addr == "" {
		addr = req.Header.Get("X-Forwarded-For")
	}
	if addr == "" {
		addr = req.RemoteAddr
	} else {
		addr = fmt.Sprintf("%s (%s)", addr, req.RemoteAddr)
	}
	rid := RequestID(req.Context())
	log.Printf("[%d] %s - %s %s %s", rid, addr, req.Method,
		req.URL.RequestURI(), req.Proto)

	if err := h(w, req); err != nil {
		switch e := err.(type) {
		case Error:
			log.Printf("[%d] %s", rid, e)
			http.Error(w, e.Error(), e.Status())
		default:
			msg := fmt.Sprintf("%s - %s",
				http.StatusText(http.StatusInternalServerError),
				err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
		}
	}
}
