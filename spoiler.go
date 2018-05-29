// Copyright 2015, 2016, 2017 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the COPYING file.

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type SpoilerAttachment struct {
	Color string `json:"color"`
	Text  string `json:"text"`
}

type SpoilerResponse struct {
	ResponseType string              `json:"response_type"`
	Attachments  []SpoilerAttachment `json:"attachments"`
}

// Spoiler is the handler for the "/spoiler" Slack slash command.
func Spoiler(w http.ResponseWriter, req *http.Request) error {
	text := req.FormValue("text")
	if len(text) == 0 {
		return StatusError{http.StatusInternalServerError,
			errors.New("usage: /spoiler <text to hide>")}
	}

	payload := SpoilerResponse{
		ResponseType: "in_channel",
		Attachments: []SpoilerAttachment{{
			Color: "danger",
			Text: fmt.Sprintf("%s posted a spoiler...\n\n\n\n\n%s",
				req.FormValue("user_name"), req.FormValue("text")),
		}},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return StatusError{http.StatusInternalServerError,
			fmt.Errorf("Could not marshal response: %+v", payload)}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonPayload)
	return nil
}

func init() {
	Commands["/spoiler"] = Spoiler
}
