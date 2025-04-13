/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package utils

import (
	"errors"
	"strings"
	"unicode"

	"github.com/microcosm-cc/bluemonday"
)

// TextSanitizer provides methods to sanitize text input
type TextSanitizer struct {
	policy *bluemonday.Policy
}

// NewTextSanitizer creates a new text sanitizer
func NewTextSanitizer() *TextSanitizer {
	// Create a policy that allows some basic formatting but removes potentially dangerous content
	policy := bluemonday.UGCPolicy()

	// Allow some additional tags that might be useful for chat messages
	policy.AllowStandardAttributes()
	policy.AllowStandardURLs()

	return &TextSanitizer{
		policy: policy,
	}
}

// SanitizeText cleans user input text to prevent XSS and other injection attacks
func (s *TextSanitizer) SanitizeText(input string) string {
	// Apply HTML sanitization
	sanitized := s.policy.Sanitize(input)

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	return sanitized
}

// ValidateMessageContent checks if a message content is valid
func ValidateMessageContent(content string) (bool, string) {
	// Check if content is empty after trimming
	if strings.TrimSpace(content) == "" {
		return false, "Message content cannot be empty"
	}

	// Check message length (e.g., max 10000 characters)
	if len(content) > 10000 {
		return false, "Message content is too long (maximum 10000 characters)"
	}

	// Check for other invalid patterns (example: control characters)
	for _, r := range content {
		if unicode.IsControl(r) && !unicode.IsSpace(r) {
			return false, "Message contains invalid control characters"
		}
	}

	return true, ""
}

// ValidateAndSanitizeMessage validates and sanitizes a message
func ValidateAndSanitizeMessage(content string) (string, error) {
	// Validate the message content
	valid, errMsg := ValidateMessageContent(content)
	if !valid {
		return "", errors.New(errMsg)
	}

	// Sanitize the content
	sanitizer := NewTextSanitizer()
	sanitized := sanitizer.SanitizeText(content)

	return sanitized, nil
}
