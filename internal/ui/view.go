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
	profile := "default"
	endpoint := "AWS本番環境"
	if m.s3Client != nil {
		profile = m.s3Client.GetProfile()
		if m.s3Client.GetEndpointURL() != "" {
			endpoint = m.s3Client.GetEndpointURL()
		}
	}

	// ヘッダー部分（常に表示）
	header := fmt.Sprintf("Profile: %s\nEndpoint url: %s\n\n", profile, endpoint)
	header += m.filterInput.View() + "\n\n"

	// リスト部分（スクロール可能）
	var listView string
	if len(m.bucketModel.FilteredBuckets) == 0 {
		listView = "条件に一致するバケットが見つかりません"
	} else {
		// 表示可能な最大行数を計算
		maxVisibleItems := m.height - 10 // ヘッダーとフッターのスペースを考慮
		if maxVisibleItems < 1 {
			maxVisibleItems = 1 // 最低でも1行は表示
		}

		// カーソル位置に基づいて表示範囲を計算
		startIdx := 0
		if len(m.bucketModel.FilteredBuckets) > maxVisibleItems {
			// カーソルが画面外に出ないように調整
			if m.bucketModel.Cursor >= maxVisibleItems {
				startIdx = m.bucketModel.Cursor - maxVisibleItems + 1
				if startIdx+maxVisibleItems > len(m.bucketModel.FilteredBuckets) {
					startIdx = len(m.bucketModel.FilteredBuckets) - maxVisibleItems
				}
			}
		}

		endIdx := startIdx + maxVisibleItems
		if endIdx > len(m.bucketModel.FilteredBuckets) {
			endIdx = len(m.bucketModel.FilteredBuckets)
		}

		// 表示する範囲のバケットを描画
		items := make([]string, 0, endIdx-startIdx)
		for i := startIdx; i < endIdx; i++ {
			cursor := " "
			if i == m.bucketModel.Cursor {
				cursor = ">"
			}
			items = append(items, fmt.Sprintf("%s %s", cursor, m.bucketModel.FilteredBuckets[i]))
		}

		// スクロールインジケータを表示
		if startIdx > 0 {
			listView += "↑ (more)\n"
		}

		listView += strings.Join(items, "\n")

		if endIdx < len(m.bucketModel.FilteredBuckets) {
			listView += "\n↓ (more)"
		}
	}

	// フッター部分（常に表示）
	footer := "\n(↑/↓: 移動, Enter: 選択, Ctrl+C: 終了)"

	return header + listView + footer
}

// renderObjectView はオブジェクト一覧ビューを描画します
func (m UIModel) renderObjectView() string {
	// ヘッダー部分（常に表示）
	header := fmt.Sprintf("%s内のオブジェクト\n\n", m.objectModel.BucketName)
	header += m.filterInput.View() + "\n\n"

	// リスト部分（スクロール可能）
	var listView string
	if len(m.objectModel.FilteredObjects) == 0 {
		listView = "条件に一致するオブジェクトが見つかりません"
	} else {
		// 表示可能な最大行数を計算
		maxVisibleItems := m.height - 10 // ヘッダーとフッターのスペースを考慮
		if maxVisibleItems < 1 {
			maxVisibleItems = 1 // 最低でも1行は表示
		}

		// カーソル位置に基づいて表示範囲を計算
		startIdx := 0
		if len(m.objectModel.FilteredObjects) > maxVisibleItems {
			// カーソルが画面外に出ないように調整
			if m.objectModel.Cursor >= maxVisibleItems {
				startIdx = m.objectModel.Cursor - maxVisibleItems + 1
				if startIdx+maxVisibleItems > len(m.objectModel.FilteredObjects) {
					startIdx = len(m.objectModel.FilteredObjects) - maxVisibleItems
				}
			}
		}

		endIdx := startIdx + maxVisibleItems
		if endIdx > len(m.objectModel.FilteredObjects) {
			endIdx = len(m.objectModel.FilteredObjects)
		}

		// 表示する範囲のオブジェクトを描画
		items := make([]string, 0, endIdx-startIdx)
		for i := startIdx; i < endIdx; i++ {
			cursor := " "
			if i == m.objectModel.Cursor {
				cursor = ">"
			}
			items = append(items, fmt.Sprintf("%s %s", cursor, m.objectModel.FilteredObjects[i]))
		}

		// スクロールインジケータを表示
		if startIdx > 0 {
			listView += "↑ (more)\n"
		}

		listView += strings.Join(items, "\n")

		if endIdx < len(m.objectModel.FilteredObjects) {
			listView += "\n↓ (more)"
		}
	}

	// フッター部分（常に表示）
	footer := "\n(↑/↓: 移動, Enter: ダウンロード, Esc: バケット一覧に戻る, Ctrl+C: 終了)"

	return header + listView + footer
}
