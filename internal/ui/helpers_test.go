package ui

import (
	"reflect"
	"testing"
)

func TestFilterItems(t *testing.T) {
	// テストケースの定義
	testCases := []struct {
		name     string
		items    []string
		filter   string
		expected []string
	}{
		{
			name:     "空のフィルター",
			items:    []string{"bucket1", "bucket2", "bucket3"},
			filter:   "",
			expected: []string{"bucket1", "bucket2", "bucket3"},
		},
		{
			name:     "一致する要素がある場合",
			items:    []string{"bucket1", "bucket2", "my-bucket", "bucket3"},
			filter:   "my",
			expected: []string{"my-bucket"},
		},
		{
			name:     "大文字小文字の区別なし",
			items:    []string{"Bucket1", "bucket2", "MY-bucket", "bucket3"},
			filter:   "my",
			expected: []string{"MY-bucket"},
		},
		{
			name:     "複数の要素が一致する場合",
			items:    []string{"test1", "test2", "test3", "other"},
			filter:   "test",
			expected: []string{"test1", "test2", "test3"},
		},
		{
			name:     "一致する要素がない場合",
			items:    []string{"bucket1", "bucket2", "bucket3"},
			filter:   "object",
			expected: []string{},
		},
		{
			name:     "空の入力スライス",
			items:    []string{},
			filter:   "test",
			expected: []string{},
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filterItems(tc.items, tc.filter)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("期待結果 %v, 実際の結果 %v", tc.expected, result)
			}
		})
	}
}

func TestMin(t *testing.T) {
	// テストケースの定義
	testCases := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{
			name:     "aがbより小さい",
			a:        5,
			b:        10,
			expected: 5,
		},
		{
			name:     "bがaより小さい",
			a:        10,
			b:        5,
			expected: 5,
		},
		{
			name:     "aとbが等しい",
			a:        7,
			b:        7,
			expected: 7,
		},
		{
			name:     "負の数を含む場合",
			a:        -3,
			b:        2,
			expected: -3,
		},
		{
			name:     "両方が負の数の場合",
			a:        -3,
			b:        -5,
			expected: -5,
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := min(tc.a, tc.b)
			if result != tc.expected {
				t.Errorf("期待結果 %d, 実際の結果 %d", tc.expected, result)
			}
		})
	}
}
