package ui

import (
	"strings"
)

// filterItems は文字列のスライスをフィルタリングします
func filterItems(items []string, filter string) []string {
	if filter == "" {
		return items
	}

	filtered := make([]string, 0)
	for _, item := range items {
		if strings.Contains(strings.ToLower(item), strings.ToLower(filter)) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// min は2つの整数の小さい方を返します
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
