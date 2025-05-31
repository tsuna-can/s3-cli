package ui

import (
	"github.com/tsuna-can/s3-cli/internal/aws"
)

// s3ClientInitMsg はS3クライアントの初期化メッセージです
type s3ClientInitMsg struct {
	client *aws.S3Client
}

// bucketsMsg はバケットリストのメッセージです
type bucketsMsg struct {
	buckets []string
}

// objectsMsg はオブジェクトリストのメッセージです
type objectsMsg struct {
	objects []string
}

// errorMsg はエラーメッセージです
type errorMsg struct {
	err error
}

// downloadedMsg はダウンロード完了メッセージです
type downloadedMsg struct {
	bucket    string
	key       string
	outputDir string
}
