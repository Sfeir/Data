package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndpoint(t *testing.T) {
	assert.Equal(t, base+sample, endpoint().URL.String())

	Language = []string{"en"}
	assert.Equal(t, base+sample+"?language=en", endpoint().URL.String())
	Language = []string{}

	Track = []string{"twitter", "golang"}
	assert.Equal(t, base+filter+"?track=twitter%2Cgolang", endpoint().URL.String())
	Track = []string{}

	Follow = []string{"12345"}
	assert.Equal(t, base+filter+"?follow=12345", endpoint().URL.String())
	Follow = []string{}

	Locations = []string{"-122.75", "36.8", "-121.75", "37.8"}
	assert.Equal(t, base+filter+"?locations=-122.75%2C36.8%2C-121.75%2C37.8", endpoint().URL.String())
	Locations = []string{}
}

var rootstests = []struct {
	in  []byte
	out map[string]bool
}{
	{[]byte(`{ "hello": "goodbye" }`),
		map[string]bool{"hello": true}},
	{[]byte(`{ "You": ["say", "hello"] }`),
		map[string]bool{"You": true}},
	{[]byte(`{ "I": "say", "Goodbye": "..." }`),
		map[string]bool{"I": true, "Goodbye": true}},
}

func TestRoots(t *testing.T) {
	for _, tt := range rootstests {
		r := roots(tt.in)
		assert.Equal(t, tt.out, r)
	}
}
