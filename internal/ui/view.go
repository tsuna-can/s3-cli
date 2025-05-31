package ui

import (
	"fmt"
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

	view := fmt.Sprintf("Profile: %s, Endpoint url: %s\n\n",
		profile, endpoint)
	view += m.filterInput.View() + "\n\n"

	if len(m.bucketModel.FilteredBuckets) == 0 {
		view += "条件に一致するバケットが見つかりません"
	} else {
		for i, bucket := range m.bucketModel.FilteredBuckets {
			cursor := " "
			if i == m.bucketModel.Cursor {
				cursor = ">"
			}
			view += fmt.Sprintf("%s %s\n", cursor, bucket)
		}
	}
	view += "\n(↑/↓: 移動, Enter: 選択, Ctrl+C: 終了)"

	return view
}

// renderObjectView はオブジェクト一覧ビューを描画します
func (m UIModel) renderObjectView() string {
	view := fmt.Sprintf("%s内のオブジェクト\n\n", m.objectModel.BucketName)
	view += m.filterInput.View() + "\n\n"

	if len(m.objectModel.FilteredObjects) == 0 {
		view += "条件に一致するオブジェクトが見つかりません"
	} else {
		for i, object := range m.objectModel.FilteredObjects {
			cursor := " "
			if i == m.objectModel.Cursor {
				cursor = ">"
			}
			view += fmt.Sprintf("%s %s\n", cursor, object)
		}
	}
	view += "\n(↑/↓: 移動, Enter: ダウンロード, Esc: バケット一覧に戻る, Ctrl+C: 終了)"

	return view
}
