package ui

import (
	"context"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Update はUIイベントを処理し、モデルを更新します
func (m UIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

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
		m.err = msg.err
	}

	var cmd tea.Cmd
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
		if m.state == "objects" {
			m.state = "buckets"
			m.filterInput.Reset()
			m.filterInput.Placeholder = "Filter buckets..."
			return m, nil
		}

	case tea.KeyEnter:
		if m.state == "buckets" && len(m.bucketModel.FilteredBuckets) > 0 {
			selectedBucket := m.bucketModel.FilteredBuckets[m.bucketModel.Cursor]
			m.state = "objects"
			m.objectModel.BucketName = selectedBucket
			m.filterInput.Reset()
			m.filterInput.Placeholder = "Filter objects..."
			return m, m.fetchObjects(selectedBucket)
		}

	case tea.KeyUp:
		if m.state == "buckets" {
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

	case tea.KeyDown:
		if m.state == "buckets" {
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
	}

	return m, nil
}

// applyFilter はフィルターを適用します
func (m *UIModel) applyFilter() {
	if m.state == "buckets" {
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
