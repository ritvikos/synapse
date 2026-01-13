// Copyright 2025-2026 Ritvik Gupta
// SPDX-License-Identifier: Apache-2.0

package transform

import "strings"

type Transformer func(string) string

var (
	TrimSpace Transformer = strings.TrimSpace
	ToLower   Transformer = strings.ToLower
	ToUpper   Transformer = strings.ToUpper
)

func ApplyTransformations(text string, transformers ...Transformer) string {
	for _, transformer := range transformers {
		text = transformer(text)
	}
	return text
}

func OnlyDigits(text string) string {
	return strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' || r == '.' {
			return r
		}
		return -1
	}, text)
}

func NormalizeWhitespace(text string) string {
	return strings.Join(strings.Fields(text), " ")
}

// func CleanHTML(html string) string {}
