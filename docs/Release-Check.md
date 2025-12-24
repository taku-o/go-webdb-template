# リリース前確認項目

本番環境へのリリース前に実施すべき確認項目をまとめたドキュメントです。

---

## 1. テストの実行

### 1.1 サーバーサイドテスト

#### ユニットテスト（必須）

```bash
cd server

# 全ユニットテストを実行
go test ./internal/... -v

# カバレッジ付きで実行
go test ./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out

# 特定のパッケージのみ実行
go test ./internal/db/... -v          # Shardingロジック
go test ./internal/repository/... -v  # Repository層
go test ./internal/service/... -v     # Service層
```

**期待結果**: 全テストがPASSすること

#### 統合テスト（必須）

```bash
cd server

# 統合テストを実行
go test ./test/integration/... -v

# より詳細なログを出力
go test ./test/integration/... -v -count=1
```

**確認項目**:
- ✅ User CRUD Flow: ユーザーの作成・取得・更新・削除
- ✅ Post CRUD Flow: 投稿の作成・取得・更新・削除
- ✅ Cross-Shard Operations: クロスシャードクエリの動作
- ✅ User-Post JOIN: ユーザーと投稿のJOIN操作

#### E2Eテスト（推奨）

```bash
cd server

# E2Eテストを実行
go test ./test/e2e/... -v
```

**確認項目**:
- ✅ API全エンドポイントの動作
- ✅ HTTP status codeの確認
- ✅ エラーハンドリングの確認

---

### 1.2 クライアントサイドテスト

#### ユニットテスト（必須）

```bash
cd client

# 全ユニットテストを実行
npm test

# カバレッジ付きで実行
npm run test:coverage

# ウォッチモード（開発時）
npm run test:watch

# 統合テストを除外して実行
npm test -- --testPathIgnorePatterns="integration"
```

**期待結果**:
- API Client Tests: 7/7 passed
- Home Page Tests: 4/4 passed

#### E2Eテスト（推奨）

```bash
cd client

# サーバーとクライアントを起動してから実行
npm run e2e

# UIモードで実行（デバッグ用）
npm run e2e:ui

# 特定のブラウザのみ実行
npx playwright test --project=chromium
```

**確認項目**:
- ✅ User management flow
- ✅ Post management flow
- ✅ Cross-shard JOIN display

---

## 2. ビルド確認

### 2.1 サーバービルド

```bash
cd server

# ビルドの実行
go build -o bin/server cmd/server/main.go

# バイナリの確認
ls -lh bin/server

# 動作確認
./bin/server
```

**確認項目**:
- ✅ ビルドエラーがないこと
- ✅ バイナリが正常に起動すること
- ✅ 全Shardへの接続が確立されること

### 2.2 クライアントビルド

```bash
cd client

# 本番ビルドの実行
npm run build

# ビルド結果の確認
ls -lh .next

# ビルド後の動作確認（推奨）
npm start
```

**確認項目**:
- ✅ ビルドエラーがないこと
- ✅ TypeScriptの型エラーがないこと
- ✅ 最適化されたバンドルが生成されていること

---

## 3. 環境設定の確認

### 3.1 サーバー環境変数

**本番環境用設定ファイル**: `server/config/production.yaml`

```bash
# 環境変数の確認
cat server/config/production.yaml

# パスワードや機密情報が環境変数から読み込まれることを確認
echo $DB_SHARD1_PASSWORD
echo $DB_SHARD2_PASSWORD
```

**確認項目**:
- ✅ `production.yaml.example`をコピーして`production.yaml`を作成済み
- ✅ データベース接続情報が正しく設定されていること
- ✅ 本番DBのホスト名、ポート、認証情報が正しいこと
- ✅ 機密情報が環境変数から読み込まれること
- ✅ `production.yaml`が`.gitignore`に含まれていること

### 3.2 クライアント環境変数

```bash
# 環境変数の確認
echo $NEXT_PUBLIC_API_BASE_URL
```

**確認項目**:
- ✅ APIのベースURLが本番環境のURLに設定されていること
- ✅ 本番用の環境変数ファイルが用意されていること

---

## 4. データベース確認

### 4.1 マイグレーション

```bash
# 各Shardでマイグレーションを実行
# Shard 1
sqlite3 /path/to/production/shard1.db < db/migrations/shard1/001_init.sql

# Shard 2
sqlite3 /path/to/production/shard2.db < db/migrations/shard2/001_init.sql

# PostgreSQL/MySQLの場合
# psql -h db-shard1.example.com -U app_user -d app_shard1 -f db/migrations/shard1/001_init.sql
```

**確認項目**:
- ✅ 全Shardで同一のスキーマが作成されていること
- ✅ テーブル、インデックス、制約が正しく作成されていること

### 4.2 接続テスト

```bash
cd server

# 本番環境で接続テスト
APP_ENV=production go run cmd/server/main.go
```

**確認項目**:
- ✅ "Successfully connected to all database shards" のログが出力されること
- ✅ 各Shardへのping/接続が成功すること

---

## 5. 動作確認（手動テスト）

### 5.1 サーバー起動

```bash
cd server

# 本番環境で起動
APP_ENV=production ./bin/server

# または
APP_ENV=production go run cmd/server/main.go
```

**確認項目**:
- ✅ サーバーが正常に起動すること
- ✅ ポート8080でリッスンしていること
- ✅ エラーログが出力されていないこと

### 5.2 クライアント起動

```bash
cd client

# 本番ビルドで起動
NEXT_PUBLIC_API_BASE_URL=<本番APIのURL> npm start
```

**確認項目**:
- ✅ クライアントが正常に起動すること
- ✅ ポート3000でアクセス可能なこと

### 5.3 基本機能の動作確認

#### ユーザー管理
1. http://localhost:3000/users にアクセス
2. 新規ユーザーを作成
3. ユーザー一覧に表示されることを確認
4. ユーザーを削除
5. 一覧から消えることを確認

#### 投稿管理
1. http://localhost:3000/posts にアクセス
2. ユーザーを選択して新規投稿を作成
3. 投稿一覧に表示されることを確認
4. 投稿を削除
5. 一覧から消えることを確認

#### クロスシャードクエリ
1. http://localhost:3000/user-posts にアクセス
2. ユーザーと投稿がJOINされて表示されることを確認
3. 複数Shardからデータが取得されていることを確認

---

## 6. パフォーマンステスト（推奨）

### 6.1 負荷テスト

```bash
# Apache Benchを使用した簡易負荷テスト
ab -n 1000 -c 10 http://localhost:8080/api/users

# より詳細な負荷テスト（別途ツールが必要）
# k6, JMeter, Gatling等を使用
```

**確認項目**:
- ✅ レスポンスタイムが許容範囲内であること
- ✅ エラー率が低いこと（<1%）
- ✅ メモリリークがないこと

---

## 7. セキュリティチェック

### 7.1 依存関係の脆弱性スキャン

```bash
# サーバー（Go）
cd server
go list -json -m all | nancy sleuth

# クライアント（npm）
cd client
npm audit
npm audit fix  # 自動修正可能な場合
```

**確認項目**:
- ✅ 重大な脆弱性がないこと
- ✅ 依存パッケージが最新の安定版であること

### 7.2 設定ファイルの確認

```bash
# 機密情報がGitにコミットされていないことを確認
git log --all --full-history -- "*production.yaml"
git log --all --full-history -- "*.env"

# .gitignoreの確認
cat .gitignore | grep -E "(production.yaml|.env|*.db)"
```

**確認項目**:
- ✅ `production.yaml`, `.env` 等の機密ファイルがコミットされていないこと
- ✅ データベースファイルがコミットされていないこと
- ✅ `.gitignore`が適切に設定されていること

---

## 8. デプロイ準備

### 8.1 成果物の確認

```bash
# サーバー
ls -lh server/bin/server

# クライアント
ls -lh client/.next/

# 設定ファイル
ls -lh server/config/production.yaml
```

**確認項目**:
- ✅ サーバーバイナリが生成されていること
- ✅ クライアントビルドが生成されていること
- ✅ 本番用設定ファイルが用意されていること

### 8.2 ドキュメントの確認

```bash
# 必要なドキュメントが揃っていることを確認
ls -lh docs/
```

**確認項目**:
- ✅ API.md: APIドキュメント
- ✅ Architecture.md: アーキテクチャドキュメント
- ✅ Sharding.md: Sharding戦略ドキュメント
- ✅ Testing.md: テスト戦略ドキュメント
- ✅ Project-Structure.md: プロジェクト構造
- ✅ Release-Check.md: 本ドキュメント

---

## 9. リリースチェックリスト

リリース前に以下の項目を確認してください：

### テスト
- [ ] サーバーユニットテスト: PASS
- [ ] サーバー統合テスト: PASS
- [ ] クライアントユニットテスト: PASS
- [ ] E2Eテスト: PASS（推奨）

### ビルド
- [ ] サーバービルド: 成功
- [ ] クライアントビルド: 成功
- [ ] TypeScript型チェック: エラーなし

### 環境設定
- [ ] 本番用設定ファイル作成済み
- [ ] 環境変数設定済み
- [ ] 機密情報がGitにコミットされていないことを確認

### データベース
- [ ] 全Shardへの接続確認
- [ ] マイグレーション実行済み
- [ ] スキーマ整合性確認

### 動作確認
- [ ] サーバー起動確認
- [ ] クライアント起動確認
- [ ] 基本機能の手動テスト完了

### セキュリティ
- [ ] 依存関係の脆弱性スキャン実施
- [ ] セキュリティ設定確認

### ドキュメント
- [ ] APIドキュメント更新済み
- [ ] README.md更新済み
- [ ] リリースノート作成済み

---

## 10. トラブルシューティング

### テストが失敗する場合

```bash
# キャッシュをクリアして再実行
go clean -testcache
go test ./... -v

# クライアント
rm -rf client/node_modules client/.next
npm install
npm test
```

### ビルドが失敗する場合

```bash
# 依存関係を再インストール
cd server && go mod tidy
cd client && rm -rf node_modules && npm install
```

### データベース接続エラー

```bash
# 接続情報を確認
cat server/config/production.yaml

# 手動で接続テスト
psql -h <host> -U <user> -d <database>

# ログを確認
tail -f /var/log/app/server.log
```

---

## 11. ロールバック手順

リリース後に問題が発生した場合の緊急対応手順：

1. **即座にサービスを停止**
   ```bash
   pkill -f server
   ```

2. **前バージョンに切り戻し**
   ```bash
   # バイナリを前バージョンに戻す
   cp bin/server.backup bin/server

   # 設定を前バージョンに戻す
   git checkout HEAD~1 server/config/production.yaml
   ```

3. **サービスを再起動**
   ```bash
   ./bin/server
   ```

4. **問題の調査**
   - ログファイルの確認
   - エラーメッセージの収集
   - データベースの状態確認

---

## まとめ

このチェックリストに従ってリリース前確認を実施することで、本番環境への安全なデプロイが可能になります。

**重要**:
- 全てのテストが成功していることを確認してからリリースしてください
- 本番環境での初回デプロイ時は、特に慎重に確認を行ってください
- 問題が発生した場合に備えて、ロールバック手順を事前に確認しておいてください
