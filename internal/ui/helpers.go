package ui

import (
	"strings"
)

// filterItems は文字列のスライスをフィルタリングします（部分一致・大文字小文字無視）
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
