package ui

import (
	"fmt"
	"strings"
)

// View はUIの現在の状態を表示します
func (m UIModel) View() string {
	if m.msg != "" {
		return fmt.Sprintf("%s\n\nCtrl+Cで終了してください。", m.msg)
	}

	if m.state == BucketsView {
		return m.renderBucketView()
	} else {
		return m.renderObjectView()
	}
}

// renderBucketView はバケット一覧ビューを描画します
func (m UIModel) renderBucketView() string {
	var profile, endpoint string
	if m.s3Client != nil {
		profile = m.s3Client.GetProfile()
		if m.s3Client.GetEndpointURL() != "" {
			endpoint = m.s3Client.GetEndpointURL()
		}
	}

	// ヘッダー部分（常に表示）
	header := fmt.Sprintf("Profile: %s\nEndpoint url: %s\n\n", profile, endpoint)
	header += m.filterInput.View() + "\n\n"

	// リスト部分（共通関数を使用）
	listView := m.renderList(
		m.bucketModel.FilteredBuckets,
		m.bucketModel.Cursor,
		"条件に一致するバケットが見つかりません",
	)

	// フッター部分（常に表示）
	footer := "\n(↑/↓: 移動, Enter: 選択, Ctrl+C: 終了)"

	return header + listView + footer
}

// renderObjectView はオブジェクト一覧ビューを描画します
func (m UIModel) renderObjectView() string {
	var profile, endpoint string
	if m.s3Client != nil {
		profile = m.s3Client.GetProfile()
		if m.s3Client.GetEndpointURL() != "" {
			endpoint = m.s3Client.GetEndpointURL()
		}
	}

	// ヘッダー部分（常に表示）
	header := fmt.Sprintf("Profile: %s\nEndpoint url: %s\nBucket: %s\n\n", profile, endpoint, m.objectModel.BucketName)
	header += m.filterInput.View() + "\n\n"

	// リスト部分（共通関数を使用）
	listView := m.renderList(
		m.objectModel.FilteredObjects,
		m.objectModel.Cursor,
		"条件に一致するオブジェクトが見つかりません",
	)

	// フッター部分（常に表示）
	footer := "\n(↑/↓: 移動, Enter: ダウンロード, Esc: バケット一覧に戻る, Ctrl+C: 終了)"

	return header + listView + footer
}

// renderList はリスト部分を描画する共通関数です
func (m UIModel) renderList(items []string, cursor int, emptyMessage string) string {
	if len(items) == 0 {
		return emptyMessage
	}

	// 表示可能な最大行数を計算
	maxVisibleItems := m.height - 10 // ヘッダーとフッターのスペースを考慮
	if maxVisibleItems < 1 {
		maxVisibleItems = 1 // 最低でも1行は表示
	}

	// 表示範囲を計算
	startIdx, endIdx := m.calculateVisibleRange(items, cursor, maxVisibleItems)

	// 表示する範囲のアイテムを描画
	resultItems := make([]string, 0, endIdx-startIdx)
	for i := startIdx; i < endIdx; i++ {
		cursorMark := " "
		if i == cursor {
			cursorMark = ">"
		}
		resultItems = append(resultItems, fmt.Sprintf("%s %s", cursorMark, items[i]))
	}

	// スクロールインジケータを表示
	var result string
	if startIdx > 0 {
		result += "↑ (more)\n"
	}

	result += strings.Join(resultItems, "\n")

	if endIdx < len(items) {
		result += "\n↓ (more)"
	}

	return result
}

// calculateVisibleRange は表示する項目の範囲を計算します
func (m UIModel) calculateVisibleRange(items []string, cursor int, maxVisibleItems int) (int, int) {
	startIdx := 0
	if len(items) > maxVisibleItems {
		// カーソルが画面外に出ないように調整
		if cursor >= maxVisibleItems {
			startIdx = cursor - maxVisibleItems + 1
			if startIdx+maxVisibleItems > len(items) {
				startIdx = len(items) - maxVisibleItems
			}
		}
	}

	endIdx := startIdx + maxVisibleItems
	if endIdx > len(items) {
		endIdx = len(items)
	}

	return startIdx, endIdx
}
