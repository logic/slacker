// Copyright 2015, 2016, 2017 Ed Marshall. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the COPYING file.

package main

import (
	"testing"
)

func TestParseTickerCommand(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{
			input: "",
			valid: false,
		},
		{
			input: "X",
			valid: true,
		},
		{
			input: "X Y",
			valid: false,
		},
		{
			input: "invalid-ticker",
			valid: false,
		},
		{
			input: "-period=1d X",
			valid: true,
		},
		{
			input: "-period=5d X",
			valid: true,
		},
		{
			input: "-period=1Y X",
			valid: true,
		},
		{
			input: "-period=5Y X",
			valid: true,
		},
		{
			input: "-period=0 X",
			valid: false,
		},
		{
			input: "-period=0d X",
			valid: false,
		},
		{
			input: "-period=0Y X",
			valid: false,
		},
		{
			input: "-period=X X",
			valid: false,
		},
		{
			input: "-interval 0 X",
			valid: true,
		},
		{
			input: "-interval X X",
			valid: false,
		},
	}

	for i, test := range tests {
		_, err := ParseTickerCommand(test.input)
		if test.valid && err != nil {
			t.Errorf("%d. expected input to be valid, got %s", i, err)
		} else if !test.valid && err == nil {
			t.Errorf("%d. expected input to be invalid", i)
		}
	}
}
