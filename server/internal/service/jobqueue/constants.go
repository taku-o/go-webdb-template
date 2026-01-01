package jobqueue

// ジョブタイプ定数
const (
	// JobTypeDelayPrint は遅延出力ジョブのタイプ
	// 参考コードとして利用するため、将来の実装に影響しない名前を使用
	JobTypeDelayPrint = "demo:delay_print"
)

// デフォルトの遅延時間（3分 = 180秒）
const DefaultDelaySeconds = 180

// デフォルトの最大リトライ回数
const DefaultMaxRetry = 10
