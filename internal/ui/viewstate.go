package ui

// ViewState はUIの表示状態を表す型です
type ViewState int

const (
	// BucketsView はバケット一覧表示状態
	BucketsView ViewState = iota
	// ObjectsView はオブジェクト一覧表示状態
	ObjectsView
)

// String はViewStateを文字列で返します
func (v ViewState) String() string {
	switch v {
	case BucketsView:
		return "buckets"
	case ObjectsView:
		return "objects"
	default:
		return "unknown"
	}
}
