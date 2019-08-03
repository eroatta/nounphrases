package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jdkato/prose.v2"
)

func TestFilterInsignificantSuffixes_ShouldReturnFilteredArray(t *testing.T) {
	tests := []struct {
		name   string
		tokens []prose.Token
		want   []prose.Token
	}{
		{"empty_array", []prose.Token{}, []prose.Token{}},
		{"no_filtered_array", []prose.Token{{Tag: "NNI", Text: "a", Label: "NNI"}},
			[]prose.Token{{Tag: "NNI", Text: "a", Label: "NNI"}}},
		{"filtered_array", []prose.Token{
			{Tag: "DT", Text: "a", Label: "DT"},
			{Tag: "CC", Text: "b", Label: "CC"},
			{Tag: "NNI", Text: "a", Label: "NNI"},
			{Tag: "PRP$", Text: "c", Label: "PRP$"},
			{Tag: "PRP", Text: "d", Label: "PRP"}},
			[]prose.Token{{Tag: "NNI", Text: "a", Label: "NNI"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterInsignificantSuffixes(tt.tokens)

			assert.Equal(t, tt.want, got, "elements should match in number and order")
		})
	}
}

func TestNormalizeTokens_ShouldReturnNormalizedTokens(t *testing.T) {
	tests := []struct {
		name  string
		chunk []prose.Token
		want  []prose.Token
	}{
		{"empty_array", []prose.Token{}, []prose.Token{}},
		{"unnormalized_array", []prose.Token{{Tag: "NN", Text: "a", Label: "NN"}},
			[]prose.Token{{Tag: "NN", Text: "a", Label: "NN"}}},
		{"normalized_array", []prose.Token{
			{Tag: "NP-TL", Text: "a", Label: ""},
			{Tag: "NP", Text: "b", Label: ""},
			{Tag: "NNP-TL", Text: "c", Label: ""},
			{Tag: "NNS", Text: "d", Label: ""},
			{Tag: "NN", Text: "e", Label: ""}},
			[]prose.Token{
				{Tag: "NNP", Text: "a", Label: ""},
				{Tag: "NNP", Text: "b", Label: ""},
				{Tag: "NNP", Text: "c", Label: ""},
				{Tag: "NN", Text: "d", Label: ""},
				{Tag: "NN", Text: "e", Label: ""}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeTokens(tt.chunk)

			assert.Equal(t, tt.want, got, "elements should match in number and order")
		})
	}
}
