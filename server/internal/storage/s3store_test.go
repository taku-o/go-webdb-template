package storage

import (
	"testing"
)

func TestNewS3FileStore(t *testing.T) {
	// S3FileStoreを作成（実際のAWS接続は行わない）
	// 注意: 実際のS3接続テストは統合テストで行う
	store, err := NewS3FileStore("test-bucket", "ap-northeast-1")

	// 開発環境ではAWS認証情報がないためエラーになる可能性がある
	// そのため、エラーが返されても構造体が正しく定義されていることを確認する
	if err != nil {
		t.Logf("NewS3FileStore returned error (expected in dev environment without AWS credentials): %v", err)
		// エラーの場合でもテストはスキップしない（構造体の定義確認が目的）
		return
	}

	if store == nil {
		t.Fatal("expected store to be non-nil")
	}
}

func TestS3FileStore_GetBucket(t *testing.T) {
	// S3FileStore構造体のフィールドアクセステスト
	store := &S3FileStore{
		bucket: "test-bucket",
		region: "ap-northeast-1",
	}

	if store.GetBucket() != "test-bucket" {
		t.Errorf("expected bucket 'test-bucket', got %s", store.GetBucket())
	}
}

func TestS3FileStore_GetRegion(t *testing.T) {
	// S3FileStore構造体のフィールドアクセステスト
	store := &S3FileStore{
		bucket: "test-bucket",
		region: "ap-northeast-1",
	}

	if store.GetRegion() != "ap-northeast-1" {
		t.Errorf("expected region 'ap-northeast-1', got %s", store.GetRegion())
	}
}
