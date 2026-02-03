# SKILLS 見直し実装タスク一覧

## 概要
SKILLS 見直しの実装を段階的に進めるためのタスク一覧。要件定義書と設計書に基づき、実施順序に沿ってタスクを分解した。

## 実装フェーズ

### Phase 1: go-test-generator の修正

#### - [ ] タスク 1.1: go-test-generator の APP_ENV=test 対応
**目的**: テスト実行時は APP_ENV=test が必須であることを SKILL に反映する

**作業内容**:
- go-test-generator の SKILL.md を開く
- テスト実行コマンドに記載されている環境変数を APP_ENV=develop から APP_ENV=test に変更する
- 「テスト時は必ず APP_ENV=test を指定すること。指定しないと認証エラー（401）が発生する」に相当する注意を追加する
- .kiro/steering/tech.md の「テスト実行ルール（必須）」への言及を追加する

**受け入れ基準**:
- テスト実行コマンドの環境変数が APP_ENV=test である
- 認証エラー（401）に関する注意が含まれる
- _Requirements: 6.2_

---

### Phase 2: api-endpoint-creator の修正

#### - [ ] タスク 2.1: api-endpoint-creator の 4層・Usecase・パス・型定義の反映
**目的**: Handler → Usecase → Service → Repository の4層と現行実装（Huma/Echo、inputs/outputs、taku-o パス）に合わせる

**作業内容**:
- api-endpoint-creator の SKILL.md を開く
- description（frontmatter）で「Handler/Service/Repositoryの3層」を「Handler → Usecase → Service → Repository の4層」に変更する
- アーキテクチャ図・説明を4層に変更し、Handler は Usecase を保持し Service を直接持たない旨を明記する
- ディレクトリ構成・参照ファイルで types.go を inputs.go, outputs.go に、ハンドラー例を dm_user_handler.go 等に合わせる
- インポートパスを github.com/taku-o/go-webdb-template に統一する
- Handler の構造体・登録関数の例で Usecase を保持し Usecase を呼ぶ形に変更する
- 認証エラー時のレスポンス例に huma.Error403Forbidden を追加する
- Service パターンに「Handler から直接呼ばない」「Usecase 経由で呼ぶ」注記を追加する
- 設計書 5.1 の参照ファイル（dm_user_handler.go, dm_user_usecase.go 等）を参照して記載する

**受け入れ基準**:
- description に4層が含まれる
- アーキテクチャで Handler が Usecase を保持することが明記されている
- 参照ファイルに inputs.go, outputs.go および実在するハンドラーが記載されている
- インポートパスが taku-o である
- コード例で Usecase を保持・呼び出す形である
- 認証エラー例に huma.Error403Forbidden が含まれる
- _Requirements: 6.1_

---

### Phase 3: repository-generator の修正

#### - [ ] タスク 3.1: repository-generator の dm_*・UUID・GORM 対応
**目的**: 参照ファイルを dm_* に、接続・テーブル名取得を UUID ベース API に、ID を UUID、構成を GORM 主に合わせる

**作業内容**:
- repository-generator の SKILL.md を開く
- 参照ファイルを dm_user_repository.go, dm_post_repository.go, dm_news_repository.go に変更する
- 接続取得の例を GetShardingConnectionByUUID に変更する
- テーブル名取得の例を GetTableNameFromUUID（戻り値 (string, error)）に変更する
- TableSelector の初期化を db.NewTableSelector(db.DBShardingTableCount, db.DBShardingTablesPerDB) に変更する
- エンティティ ID 生成を idgen.GenerateUUIDv7() とし、型を UUID 文字列（string）として扱う
- 本プロジェクトでは GORM を標準とし、標準SQL版は削除または他プロジェクト向けと明記する
- CRUD 例を GORM の実際の呼び出し（conn.DB.WithContext(ctx).Table(tableName).Create(entity) 等）に合わせる
- 設計書 5.1 の参照ファイル（dm_user_repository.go, dm_post_repository.go 等）を参照する

**受け入れ基準**:
- 参照ファイルに dm_user_repository.go, dm_post_repository.go, dm_news_repository.go が記載されている
- 接続取得で GetShardingConnectionByUUID が使用されている
- テーブル名取得で GetTableNameFromUUID が使用されている
- TableSelector で db.DBShardingTableCount, db.DBShardingTablesPerDB が使用されている
- ID 生成で idgen.GenerateUUIDv7() が使用されている
- GORM を標準とすることが明記されている
- _Requirements: 6.3_

---

### Phase 4: sharding-pattern の修正

#### - [ ] タスク 4.1: sharding-pattern の UUID・ByUUID/GetTableNameFromUUID 対応
**目的**: シャードキー・接続取得・テーブル名取得を UUID ベースの API に合わせる

**作業内容**:
- sharding-pattern の SKILL.md を開く
- シャードキーで UUID（string）が主であることを明記する
- テーブル名取得の例を GetTableNameFromUUID（戻り値 (string, error)）に変更する
- 接続取得の例を GetShardingConnectionByUUID に変更する
- 参照ファイルに sharding.go（GetTableNameFromUUID, ValidateTableName）, group_manager.go, dm_user_repository.go, dm_post_repository.go を記載する
- 定数として db.DBShardingTableCount, db.DBShardingTablesPerDB を使用する
- 設計書 5.1 の参照ファイル（sharding.go, group_manager.go 等）を参照する

**受け入れ基準**:
- シャードキーで UUID（string）が主であることが明記されている
- テーブル名取得で GetTableNameFromUUID が使用されている
- 接続取得で GetShardingConnectionByUUID が使用されている
- 参照ファイルに sharding.go, group_manager.go, dm_*_repository が記載されている
- 定数が使用されている
- _Requirements: 6.4_

---

### Phase 5: migration-helper の修正

#### - [ ] タスク 5.1: migration-helper の .hcl・他マイグレーションの補足
**目的**: スキーマが .hcl の場合や master-mysql 等の存在を補足する

**作業内容**:
- migration-helper の SKILL.md を開く
- db/schema/ の .hcl 構成に合わせた atlas migrate diff の --to の例を追記する
- db/migrations/ に master-mysql, view_master 等があることを補足する

**受け入れ基準**:
- .hcl スキーマおよび master-mysql 等の補足が追加されている
- _Requirements: 6.5_

---

### Phase 6: test-auth-env の新規追加

#### - [ ] タスク 6.1: test-auth-env SKILL の作成
**目的**: テストで認証エラーが発生したときに APP_ENV=test 未指定の可能性を指摘し、対処法を案内する SKILL を追加する

**作業内容**:
- test-auth-env 用のディレクトリを作成する
- SKILL.md を新規作成する
- frontmatter に name: test-auth-env、description に認証エラー発生時の使用と APP_ENV=test 未指定の指摘を含める
- 本文に以下を含める:
  - 認証エラーがテストで出た場合、テスト実行コマンドに APP_ENV=test が付いているか確認する手順
  - コマンド例: APP_ENV=test go test ./...、cd server && APP_ENV=test go test ./...
  - .kiro/steering/tech.md の「テスト実行ルール（必須）」の確認
  - 認証エラーが1件でも出た場合は「今回の修正とは関係ない」と判断せず原因を調査する旨

**受け入れ基準**:
- test-auth-env/SKILL.md が存在する
- description に認証エラー時の使用と APP_ENV=test 未指定の指摘が含まれる
- 本文に確認手順・コマンド例・tech.md の確認が含まれる
- 原因調査の方針が含まれる
- _Requirements: 6.6_

---

### Phase 7: react-use-effect-guard の新規追加

#### - [ ] タスク 7.1: react-use-effect-guard SKILL の作成
**目的**: クライアントで useEffect を使おうとしたときに、本当に必要か検討を促す SKILL を追加する

**作業内容**:
- react-use-effect-guard 用のディレクトリを作成する
- SKILL.md を新規作成する
- frontmatter に name: react-use-effect-guard、description に useEffect 追加・使用時に発動することと「本当に useEffect が必要か」の検討を促すことを含める
- 本文に以下を含める:
  - useEffect を使う前に確認すること: データ取得は Server Component や Server Actions で可能か、イベントに紐づく処理はイベントハンドラで十分か、外部システムとの同期は本当にマウント/更新時に毎回必要か
  - どうしても必要な場合（例: ブラウザ API の購読、フォーカス制御、クライアント専用の初回実行）のみ useEffect を使用する旨

**受け入れ基準**:
- react-use-effect-guard/SKILL.md が存在する
- description に useEffect 使用時に発動することと検討を促すことが含まれる
- 本文にデータ取得・イベントハンドラ・外部同期の確認事項が含まれる
- どうしても必要な場合のみ useEffect を使用する旨が含まれる
- _Requirements: 6.7_

---

## 実装順序の推奨

1. Phase 1: go-test-generator の修正
2. Phase 2: api-endpoint-creator の修正
3. Phase 3: repository-generator の修正
4. Phase 4: sharding-pattern の修正
5. Phase 5: migration-helper の修正
6. Phase 6: test-auth-env の新規追加
7. Phase 7: react-use-effect-guard の新規追加

## 注意事項

- effective-go, frontend-design, skill-creator は変更しないこと
- 各修正タスクでは、設計書「5.1 参照する実装ファイル」に挙げた server/internal/ のファイルを参照してから記載すること
- 新規・修正いずれも、description に使用場面とキーワードを明確に含めること
