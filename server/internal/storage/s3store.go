package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/s3store"
)

// S3FileStore はAWS S3ストア
type S3FileStore struct {
	bucket   string
	region   string
	s3Client *s3.Client
	store    s3store.S3Store
}

// NewS3FileStore は新しいS3FileStoreを作成する
// AWS認証情報は環境変数またはAWS設定から取得する
func NewS3FileStore(bucket, region string) (*S3FileStore, error) {
	// AWS設定を読み込み
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	// S3クライアントを作成
	s3Client := s3.NewFromConfig(cfg)

	// tusdのS3Storeを作成
	store := s3store.New(bucket, s3Client)

	return &S3FileStore{
		bucket:   bucket,
		region:   region,
		s3Client: s3Client,
		store:    store,
	}, nil
}

// GetStore はtusdのDataStoreを返す
func (s *S3FileStore) GetStore() handler.DataStore {
	return &s.store
}

// GetBucket はS3バケット名を返す
func (s *S3FileStore) GetBucket() string {
	return s.bucket
}

// GetRegion はAWSリージョンを返す
func (s *S3FileStore) GetRegion() string {
	return s.region
}
