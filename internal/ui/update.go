package ui

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Update はUIイベントを処理し、モデルを更新します
func (m UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		model, cmd := m.handleKeyMsg(msg)
		if model != nil {
			return model, cmd
		}
		// handleKeyMsgが処理しなかった場合（nilを返した場合）は、以下の処理に進む

	case s3ClientInitMsg:
		m.s3Client = msg.client
		return m, m.fetchBuckets

	case bucketsMsg:
		m.bucketModel.Buckets = msg.buckets
		m.bucketModel.FilteredBuckets = msg.buckets
		m.bucketModel.Cursor = 0

	case objectsMsg:
		m.objectModel.Objects = msg.objects
		m.objectModel.FilteredObjects = msg.objects
		m.objectModel.Cursor = 0

	case errorMsg:
		m.err = nil
		m.msg = fmt.Sprintf("エラー: %v", msg.err)

	case downloadedMsg:
		m.err = nil
		m.msg = fmt.Sprintf("ダウンロード完了: %s/%s → %s", msg.bucket, msg.key, msg.outputDir)
		return m, tea.Quit
	}

	var cmd tea.Cmd
	// filterInputの値が変わるたびにapplyFilterが呼ばれ、部分一致で絞り込みされる
	m.filterInput, cmd = m.filterInput.Update(msg)

	// フィルター変更時に適用
	m.applyFilter()

	return m, cmd
}

// handleKeyMsg はキーボード入力を処理します
func (m UIModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit

	case tea.KeyEsc:
		if m.state == ObjectsView { // 文字列比較から定数比較に変更
			m.state = BucketsView // 文字列からViewState型に変更
			m.filterInput.Reset()
			m.filterInput.Placeholder = "Filter buckets..."
			return m, nil
		}

	case tea.KeyEnter:
		if m.state == BucketsView && len(m.bucketModel.FilteredBuckets) > 0 { // 文字列比較から定数比較に変更
			selectedBucket := m.bucketModel.FilteredBuckets[m.bucketModel.Cursor]
			m.state = ObjectsView // 文字列からViewState型に変更
			m.objectModel.BucketName = selectedBucket
			m.filterInput.Reset()
			m.filterInput.Placeholder = "Filter objects..."
			return m, m.fetchObjects(selectedBucket)
		}
		if m.state == ObjectsView && len(m.objectModel.FilteredObjects) > 0 {
			selectedObject := m.objectModel.FilteredObjects[m.objectModel.Cursor]
			bucket := m.objectModel.BucketName
			outputDir := m.outputDir
			return m, m.downloadObject(bucket, selectedObject, outputDir)
		}

	case tea.KeyUp:
		if m.state == BucketsView {
			m.bucketModel.Cursor--
			if m.bucketModel.Cursor < 0 {
				m.bucketModel.Cursor = 0
			}
		} else {
			m.objectModel.Cursor--
			if m.objectModel.Cursor < 0 {
				m.objectModel.Cursor = 0
			}
		}
		return m, nil

	case tea.KeyDown:
		if m.state == BucketsView {
			m.bucketModel.Cursor++
			if m.bucketModel.Cursor >= len(m.bucketModel.FilteredBuckets) {
				m.bucketModel.Cursor = len(m.bucketModel.FilteredBuckets) - 1
			}
		} else {
			m.objectModel.Cursor++
			if m.objectModel.Cursor >= len(m.objectModel.FilteredObjects) {
				m.objectModel.Cursor = len(m.objectModel.FilteredObjects) - 1
			}
		}
		return m, nil
	}
	// ここでnil,nilを返すことで、通常の文字入力はfilterInputに渡される
	return nil, nil
}

// applyFilter はフィルターを適用します
func (m *UIModel) applyFilter() {
	if m.state == BucketsView {
		filter := strings.ToLower(m.filterInput.Value())
		m.bucketModel.FilteredBuckets = filterItems(m.bucketModel.Buckets, filter)
		if len(m.bucketModel.FilteredBuckets) > 0 {
			m.bucketModel.Cursor = min(m.bucketModel.Cursor, len(m.bucketModel.FilteredBuckets)-1)
		} else {
			m.bucketModel.Cursor = 0
		}
	} else {
		filter := strings.ToLower(m.filterInput.Value())
		m.objectModel.FilteredObjects = filterItems(m.objectModel.Objects, filter)
		if len(m.objectModel.FilteredObjects) > 0 {
			m.objectModel.Cursor = min(m.objectModel.Cursor, len(m.objectModel.FilteredObjects)-1)
		} else {
			m.objectModel.Cursor = 0
		}
	}
}

// fetchBuckets はS3バケット一覧を取得します
func (m UIModel) fetchBuckets() tea.Msg {
	buckets, err := m.s3Client.ListBuckets(context.Background())
	if err != nil {
		return errorMsg{err}
	}
	return bucketsMsg{buckets}
}

// fetchObjects はバケット内のオブジェクト一覧を取得します
func (m UIModel) fetchObjects(bucketName string) tea.Cmd {
	return func() tea.Msg {
		objects, err := m.s3Client.ListObjects(context.Background(), bucketName)
		if err != nil {
			return errorMsg{err}
		}
		return objectsMsg{objects}
	}
}

// downloadObject はオブジェクトをダウンロードするCmdを返します
func (m UIModel) downloadObject(bucket, key, outputDir string) tea.Cmd {
	return func() tea.Msg {
		err := m.s3Client.DownloadObject(context.Background(), bucket, key, outputDir)
		if err != nil {
			return errorMsg{err}
		}
		return downloadedMsg{bucket: bucket, key: key, outputDir: outputDir}
	}
}
