package aws

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Client provides an interface to AWS S3 operations
type S3Client struct {
	client      *s3.Client
	region      string
	profile     string
	endpointURL string
}

// NewS3Client creates a new S3 client using AWS configuration from ~/.aws/config
func NewS3Client(profile string, endpointURL string) (*S3Client, error) {
	var loadOptions []func(*config.LoadOptions) error

	// プロファイルが指定されている場合は使用
	if profile != "" {
		loadOptions = append(loadOptions, config.WithSharedConfigProfile(profile))
	}

	// カスタムエンドポイントリゾルバーを設定
	loadOptions = append(loadOptions, config.WithEndpointResolverWithOptions(
		aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:               endpointURL,
				HostnameImmutable: true,
			}, nil
		}),
	))

	// 設定を読み込む
	cfg, err := config.LoadDefaultConfig(context.Background(), loadOptions...)
	if err != nil {
		return nil, fmt.Errorf("AWS設定の読み込みに失敗しました: %w", err)
	}

	// 設定から使用されているリージョンを取得
	region := cfg.Region
	if region == "" {
		// LocalStackなどのローカル環境ではリージョン指定が必要
		region = "us-east-1" // デフォルトリージョン
	}

	// 使用しているプロファイルを特定
	usedProfile := profile
	if usedProfile == "" {
		usedProfile = os.Getenv("AWS_PROFILE")
		if usedProfile == "" {
			usedProfile = "default"
		}
	}

	client := s3.NewFromConfig(cfg)
	return &S3Client{
		client:      client,
		region:      region,
		profile:     usedProfile,
		endpointURL: endpointURL,
	}, nil
}

// GetRegion returns the region being used by the client
func (c *S3Client) GetRegion() string {
	return c.region
}

// GetProfile returns the profile being used by the client
func (c *S3Client) GetProfile() string {
	return c.profile
}

// GetEndpointURL returns the endpoint URL being used by the client
func (c *S3Client) GetEndpointURL() string {
	return c.endpointURL
}

// ListBuckets returns a list of all S3 buckets
func (c *S3Client) ListBuckets(ctx context.Context) ([]string, error) {
	result, err := c.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	var bucketNames []string
	for _, bucket := range result.Buckets {
		bucketNames = append(bucketNames, *bucket.Name)
	}

	return bucketNames, nil
}

// ListObjects returns a list of objects in the specified bucket
func (c *S3Client) ListObjects(ctx context.Context, bucketName string) ([]string, error) {
	result, err := c.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	})
	if err != nil {
		return nil, err
	}

	var objectKeys []string
	for _, object := range result.Contents {
		objectKeys = append(objectKeys, *object.Key)
	}

	return objectKeys, nil
}

// DownloadObject は指定したバケット・キーのオブジェクトをローカルにダウンロードします
func (c *S3Client) DownloadObject(ctx context.Context, bucketName, key, outputDir string) error {
	outputPath := filepath.Join(outputDir, key)

	// 同名ファイルが既に存在するかチェック
	if _, err := os.Stat(outputPath); err == nil {
		return fmt.Errorf("ファイルが既に存在します: %s", outputPath)
	}

	// ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	resp, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(outFile, resp.Body)
	return err
}
