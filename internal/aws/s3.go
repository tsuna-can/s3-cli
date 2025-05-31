package aws

import (
	"context"
	"fmt"
	"os"

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

// 将来的には環境変数から取得するためのエンドポイントURL
// 現時点ではハードコードで定義
const defaultEndpointURL = "http://localhost:4566"

// NewS3Client creates a new S3 client using AWS configuration from ~/.aws/config
func NewS3Client(profile string) (*S3Client, error) {
	var loadOptions []func(*config.LoadOptions) error

	// プロファイルが指定されている場合は使用
	if profile != "" {
		loadOptions = append(loadOptions, config.WithSharedConfigProfile(profile))
	}

	// カスタムエンドポイントリゾルバーを設定
	// 将来的にはここを環境変数から取得するように変更予定
	loadOptions = append(loadOptions, config.WithEndpointResolverWithOptions(
		aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:               defaultEndpointURL,
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
		endpointURL: defaultEndpointURL,
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
