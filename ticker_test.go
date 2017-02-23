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
			input: "-span=1d X",
			valid: true,
		},
		{
			input: "-span=5d X",
			valid: true,
		},
		{
			input: "-span=1m X",
			valid: true,
		},
		{
			input: "-span=3m X",
			valid: true,
		},
		{
			input: "-span=6m X",
			valid: true,
		},
		{
			input: "-span=1y X",
			valid: true,
		},
		{
			input: "-span=2y X",
			valid: true,
		},
		{
			input: "-span=5y X",
			valid: true,
		},
		{
			input: "-span=my X",
			valid: true,
		},
		{
			input: "-span=X X",
			valid: false,
		},
		{
			input: "-type=l X",
			valid: true,
		},
		{
			input: "-type=b X",
			valid: true,
		},
		{
			input: "-type=c X",
			valid: true,
		},
		{
			input: "-type=X X",
			valid: false,
		},
		{
			input: "-log=true X",
			valid: true,
		},
		{
			input: "-log=false X",
			valid: true,
		},
		{
			input: "-log=X X",
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
