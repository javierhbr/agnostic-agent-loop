package token

import "strings"

// CountTokens provides a rough estimation of token count.
// For MVP, we assume ~4 characters per token average for code/english.
// This avoids heavy dependencies for now.
func CountTokens(text string) int {
	if text == "" {
		return 0
	}
	// A slightly more robust heuristic?
	// 1 token ~= 4 chars in English
	return len(text) / 4
}

// CountTokensRough counts tokens very roughly based on words.
// Just another heuristic option.
func CountTokensWords(text string) int {
	words := strings.Fields(text)
	// 0.75 words per token is a common approximation
	return int(float64(len(words)) * 1.3)
}
