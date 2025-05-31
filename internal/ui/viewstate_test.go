package ui

import "testing"

// TestViewStateValues はViewState型の定数値をテストします
func TestViewStateValues(t *testing.T) {
	// iotaを使用しているので、順序通りに値が設定されていることを確認
	if BucketsView != 0 {
		t.Errorf("BucketsViewの値が期待と異なります: 期待値=%d, 実際値=%d", 0, BucketsView)
	}

	if ObjectsView != 1 {
		t.Errorf("ObjectsViewの値が期待と異なります: 期待値=%d, 実際値=%d", 1, ObjectsView)
	}
}

// TestViewStateString はViewStateのString()メソッドをテストします
func TestViewStateString(t *testing.T) {
	testCases := []struct {
		name     string
		state    ViewState
		expected string
	}{
		{
			name:     "BucketsViewの文字列表現",
			state:    BucketsView,
			expected: "buckets",
		},
		{
			name:     "ObjectsViewの文字列表現",
			state:    ObjectsView,
			expected: "objects",
		},
		{
			name:     "未定義の状態の文字列表現",
			state:    ViewState(99), // 未定義の値
			expected: "unknown",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.state.String()
			if result != tc.expected {
				t.Errorf("期待結果 %q, 実際の結果 %q", tc.expected, result)
			}
		})
	}
}
