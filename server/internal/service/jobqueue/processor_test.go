package jobqueue

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
)

func TestProcessDelayPrintJob_WithMessage(t *testing.T) {
	// 標準出力をキャプチャ
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// テスト用ペイロードを作成
	payload := DelayPrintPayload{
		Message: "Test message",
	}
	payloadBytes, err := json.Marshal(payload)
	assert.NoError(t, err)

	// タスクを作成
	task := asynq.NewTask(JobTypeDelayPrint, payloadBytes)

	// ジョブを処理
	err = ProcessDelayPrintJob(context.Background(), task)
	assert.NoError(t, err)

	// 標準出力をリストア
	w.Close()
	os.Stdout = oldStdout

	// 出力を確認
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	assert.Contains(t, output, "Test message")
}

func TestProcessDelayPrintJob_EmptyPayload(t *testing.T) {
	// 標準出力をキャプチャ
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 空のペイロードでタスクを作成
	task := asynq.NewTask(JobTypeDelayPrint, nil)

	// ジョブを処理
	err := ProcessDelayPrintJob(context.Background(), task)
	assert.NoError(t, err)

	// 標準出力をリストア
	w.Close()
	os.Stdout = oldStdout

	// 出力を確認
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// デフォルトメッセージが出力されること
	assert.Contains(t, output, "Job executed successfully")
}

func TestProcessDelayPrintJob_EmptyMessageField(t *testing.T) {
	// 標準出力をキャプチャ
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 空のメッセージフィールドを持つペイロードを作成
	payload := DelayPrintPayload{
		Message: "",
	}
	payloadBytes, err := json.Marshal(payload)
	assert.NoError(t, err)

	task := asynq.NewTask(JobTypeDelayPrint, payloadBytes)

	// ジョブを処理
	err = ProcessDelayPrintJob(context.Background(), task)
	assert.NoError(t, err)

	// 標準出力をリストア
	w.Close()
	os.Stdout = oldStdout

	// 出力を確認
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// デフォルトメッセージが出力されること
	assert.Contains(t, output, "Job executed successfully")
}

func TestProcessDelayPrintJob_InvalidJSON(t *testing.T) {
	// 不正なJSONでタスクを作成
	task := asynq.NewTask(JobTypeDelayPrint, []byte("invalid json"))

	// ジョブを処理
	err := ProcessDelayPrintJob(context.Background(), task)

	// エラーが発生すること
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "failed to unmarshal payload"))
}

func TestProcessDelayPrintJob_OutputFormat(t *testing.T) {
	// 標準出力をキャプチャ
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	payload := DelayPrintPayload{
		Message: "Format test",
	}
	payloadBytes, err := json.Marshal(payload)
	assert.NoError(t, err)

	task := asynq.NewTask(JobTypeDelayPrint, payloadBytes)

	// ジョブを処理
	err = ProcessDelayPrintJob(context.Background(), task)
	assert.NoError(t, err)

	// 標準出力をリストア
	w.Close()
	os.Stdout = oldStdout

	// 出力を確認
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// タイムスタンプ形式が含まれていることを確認
	// 出力形式: [YYYY-MM-DD HH:MM:SS] Message
	assert.Contains(t, output, "[")
	assert.Contains(t, output, "]")
	assert.Contains(t, output, "Format test")
}
