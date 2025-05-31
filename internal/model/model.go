package model

// BucketListModel represents the model for the bucket list view
type BucketListModel struct {
	Buckets         []string
	FilteredBuckets []string
	Cursor          int
	Filter          string
}

// ObjectListModel represents the model for the object list view
type ObjectListModel struct {
	BucketName      string
	Objects         []string
	FilteredObjects []string
	Cursor          int
	Filter          string
}
