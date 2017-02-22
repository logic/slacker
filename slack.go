// Copyright 2015 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the COPYING file.

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

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
	if err := h(w, req); err != nil {
		addr := req.Header.Get("X-Real-IP")
		if addr == "" {
			addr = req.Header.Get("X-Forwarded-For")
			if addr == "" {
				addr = req.RemoteAddr
			}
		}
		switch e := err.(type) {
		case Error:
			log.Printf("%s - HTTP %d - %s", addr, e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			msg := fmt.Sprintf("%s - %s",
				http.StatusText(http.StatusInternalServerError),
				err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
		}
	}
}

// SlackDispatcher does basic validation, then routes Slack slash commands
// to an appropriate handler.
func SlackDispatcher(w http.ResponseWriter, req *http.Request) error {
	if req.Method != "POST" {
		return StatusError{http.StatusBadRequest,
			errors.New("Only POST is supported")}
	}
	if len(Config.Tokens) > 0 {
		token := req.FormValue("token")
		found := false
		for _, t := range Config.Tokens {
			if token == t {
				found = true
				break
			}
		}
		if !found {
			return StatusError{http.StatusBadRequest,
				fmt.Errorf("Token '%s' is invalid", token)}
		}
	}
	command := strings.ToLower(req.FormValue("command"))
	switch command {
	case "/ticker":
		return Ticker(w, req)
	default:
		return StatusError{http.StatusBadRequest,
			fmt.Errorf("Command '%s' is invalid", command)}
	}
}
