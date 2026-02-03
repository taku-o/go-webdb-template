# SKILLS 見直し設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、`.claude/skills/` 配下の既存 SKILL の修正内容と新規 SKILL の構成を詳細に定義する。プロジェクトの現状（tech.md / structure.md / 実装）と SKILL の記載を一致させ、useEffect 検討とテスト認証エラー時の案内を SKILL で担保する。

### 1.2 設計の範囲
- 既存 SKILL 5件（api-endpoint-creator, go-test-generator, repository-generator, sharding-pattern, migration-helper）の修正内容の詳細
- 新規 SKILL 2件（test-auth-env, react-use-effect-guard）の構成・記載内容
- 実施順序と修正時の参照先

### 1.3 設計方針
- **実装準拠**: 各 SKILL のコード例・パス・API 名は server/internal/ 等の実装に合わせる
- **ステアリング整合**: .kiro/steering/tech.md, structure.md と矛盾しない記載とする
- **description の明確化**: トリガー判定に使う description に使用場面とキーワードを適切に含める
- **変更最小限**: effective-go, frontend-design, skill-creator は変更しない

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

#### 2.1.1 変更前の構造
```
.claude/skills/
├── api-endpoint-creator/
│   └── SKILL.md
├── effective-go/
│   └── SKILL.md
├── frontend-design/
│   └── SKILL.md
├── go-test-generator/
│   └── SKILL.md
├── migration-helper/
│   └── SKILL.md
├── repository-generator/
│   └── SKILL.md
├── sharding-pattern/
│   └── SKILL.md
└── skill-creator/
    └── SKILL.md
```

#### 2.1.2 変更後の構造
```
.claude/skills/
├── api-endpoint-creator/
│   └── SKILL.md          # 修正
├── effective-go/
│   └── SKILL.md          # 変更なし
├── frontend-design/
│   └── SKILL.md          # 変更なし
├── go-test-generator/
│   └── SKILL.md          # 修正
├── migration-helper/
│   └── SKILL.md          # 修正
├── repository-generator/
│   └── SKILL.md          # 修正
├── sharding-pattern/
│   └── SKILL.md          # 修正
├── skill-creator/
│   └── SKILL.md          # 変更なし
├── test-auth-env/        # 新規
│   └── SKILL.md
└── react-use-effect-guard/  # 新規
    └── SKILL.md
```

### 2.2 ファイル構成

#### 2.2.1 SKILL の共通形式
- **SKILL.md**: 必須。YAML frontmatter（name, description）と Markdown 本文で構成
- **frontmatter**: name は SKILL 識別子、description はトリガー判定用（使用場面・内容を明確に記載）
- **本文**: 手順・コード例・参照ファイル・注意事項

#### 2.2.2 修正対象 SKILL 一覧
| SKILL | 操作 | 主な変更内容 |
|-------|------|----------------|
| api-endpoint-creator | 修正 | 4層・Usecase・inputs/outputs・taku-o パス・403 例 |
| go-test-generator | 修正 | APP_ENV=test 必須・認証エラー注意 |
| repository-generator | 修正 | dm_* 参照・ByUUID・GetTableNameFromUUID・UUID・GORM 主 |
| sharding-pattern | 修正 | UUID 主・GetTableNameFromUUID・GetShardingConnectionByUUID・定数 |
| migration-helper | 修正 | .hcl / master-mysql 等の補足 |

#### 2.2.3 新規 SKILL 一覧
| SKILL | 目的 |
|-------|------|
| test-auth-env | テストで認証エラーが出たときに APP_ENV=test 未指定の可能性を指摘し、対処法を案内 |
| react-use-effect-guard | クライアントで useEffect を使おうとしたときに、本当に必要か検討を促す |

### 2.3 実施フロー

```
実施順序:
1. go-test-generator 修正
2. api-endpoint-creator 修正
3. repository-generator 修正
4. sharding-pattern 修正
5. migration-helper 修正
6. test-auth-env 新規作成
7. react-use-effect-guard 新規作成
```

### 2.4 ステアリングとの整合

- **tech.md**: テスト実行は APP_ENV=test 必須、レイヤーは Handler → Usecase → Service → Repository、モジュールパスは go.mod に準拠。各 SKILL の記載はこれに合わせる。
- **structure.md**: ディレクトリ構成・命名規則に合わせ、参照パスは server/internal/ を前提とする。
- **変更しない**: effective-go, frontend-design, skill-creator は汎用のため触れない。

## 3. コンポーネント設計

### 3.1 api-endpoint-creator の修正設計

#### 3.1.1 frontmatter
- **description**: 「Handler/Service/Repositoryの3層」を「Handler → Usecase → Service → Repository の4層」に変更。Huma API、Echo、エンドポイント登録パターンは維持。

#### 3.1.2 アーキテクチャ図・説明
- 3層の図を4層に変更: Handler → Usecase → Service → Repository
- 「Handler は Usecase を保持し、Service を直接持たない」旨を明記

#### 3.1.3 ディレクトリ構成・参照ファイル
- `internal/api/huma/types.go` を `internal/api/huma/inputs.go`, `outputs.go` に変更
- ハンドラー例を `dm_user_handler.go` 等の実在ファイルに合わせる

#### 3.1.4 コード例
- インポート: `github.com/example/go-webdb-template` → `github.com/taku-o/go-webdb-template`
- Handler 構造体: `entityService` ではなく `entityUsecase *usecaseapi.EntityUsecase`
- 登録関数内: `h.entityService.CreateEntity` ではなく `h.entityUsecase.CreateEntity`
- 認証エラー時: `huma.Error403Forbidden` を例に追加
- Service パターン: 「Handler から直接呼ばない」「Usecase 経由で呼ぶ」注記を追加

### 3.2 go-test-generator の修正設計

#### 3.2.1 テスト実行コマンド
- 全ての `APP_ENV=develop` を `APP_ENV=test` に変更
- コメントまたは注意として「テスト時は必ず APP_ENV=test を指定すること。指定しないと認証エラー（401）が発生する」を追加
- 「.kiro/steering/tech.md のテスト実行ルール（必須）を参照」と記載

### 3.3 repository-generator の修正設計

#### 3.3.1 参照ファイル
- user_repository.go, post_repository.go → dm_user_repository.go, dm_post_repository.go, dm_news_repository.go
- モデル例を model.DmUser, model.CreateDmUserRequest 等に合わせる

#### 3.3.2 API ・定数
- GetShardingConnectionByID(id int64, ...) → GetShardingConnectionByUUID(uuid string, tableBaseName string)
- GetTableName("entities", id) → GetTableNameFromUUID("dm_users", uuid)。戻り値 (string, error) を扱う
- db.NewTableSelector(32, 8) → db.NewTableSelector(db.DBShardingTableCount, db.DBShardingTablesPerDB)

#### 3.3.3 ID ・CRUD 方針
- 作成時 ID: idgen.GenerateUUIDv7() を使用。型は string（UUID）
- 本プロジェクトでは GORM 版を標準とし、標準SQL版の記述は削除または「他プロジェクト向け」に変更
- CRUD 例は conn.DB.WithContext(ctx).Table(tableName).Create(entity) 等の GORM 呼び出しに合わせる

### 3.4 sharding-pattern の修正設計

#### 3.4.1 シャードキー・API
- シャードキー: 「user_id または エンティティの id」に加え、「UUID（string）が主」と明記
- テーブル名: tableSelector.GetTableName("users", userID) → tableSelector.GetTableNameFromUUID("dm_users", uuid)。戻り値 (string, error)
- 接続取得: groupManager.GetShardingConnectionByID(userID, "users") → groupManager.GetShardingConnectionByUUID(uuid, "dm_users")

#### 3.4.2 参照ファイル・定数
- 参照: server/internal/db/sharding.go（GetTableNameFromUUID, ValidateTableName）, group_manager.go, dm_user_repository.go, dm_post_repository.go
- 定数: db.DBShardingTableCount, db.DBShardingTablesPerDB

### 3.5 migration-helper の修正設計

#### 3.5.1 補足内容
- スキーマが .hcl の場合の atlas migrate diff の --to の例を、db/schema/ 構成に合わせて追記する
- db/migrations/ に master-mysql, view_master 等があることを補足する

### 3.6 test-auth-env の新規設計

#### 3.6.1 frontmatter
- **name**: test-auth-env
- **description**: Go テストや E2E/統合テストの実行で認証エラー（401 Unauthorized 等）が発生したときに使用。APP_ENV=test を指定していない可能性を指摘し、対処法（コマンド例・tech.md の確認）を案内する。

#### 3.6.2 本文構成
1. **認証エラーがテストで出た場合の手順**
   - テスト実行コマンドに APP_ENV=test が付いているか確認する
   - コマンド例: `APP_ENV=test go test ./...`、`cd server && APP_ENV=test go test ./...`
   - .kiro/steering/tech.md の「テスト実行ルール（必須）」を確認する
2. **方針**
   - 認証エラーが1件でも出た場合は「今回の修正とは関係ない」と判断せず、原因を調査する

### 3.7 react-use-effect-guard の新規設計

#### 3.7.1 frontmatter
- **name**: react-use-effect-guard
- **description**: クライアントの React/Next.js コードで useEffect を追加・使用しようとするときに発動。「本当に useEffect が必要か」を検討するよう促す。データ取得は Server Component / Server Actions でできないか、イベントハンドラで十分でないかを確認する。

#### 3.7.2 本文構成
1. **useEffect を使う前に確認すること**
   - データ取得 → Server Component や Server Actions で可能か
   - イベントに紐づく処理 → イベントハンドラで十分か
   - 外部システムとの同期 → 本当にマウント/更新時に毎回必要か
2. **使用を許容する場合**
   - どうしても必要な場合（例: ブラウザ API の購読、フォーカス制御、クライアント専用の初回実行）のみ useEffect を使用する

## 4. エラーハンドリング・注意事項

### 4.1 修正時の参照ミス
- 実装を参照せずに修正すると、パス・API 名が実装とずれる。必ず server/internal/ の該当ファイルを確認してから記載する。

### 4.2 description のトリガー
- description が曖昧だと、意図しない場面で SKILL が発動したり、発動しなかったりする。使用場面とキーワードを明確に含める。

### 4.3 既存 SKILL の変更禁止
- effective-go, frontend-design, skill-creator は変更しない。誤って編集しないよう、修正対象リストから除外する。

## 5. 実装上の注意事項

### 5.1 参照する実装ファイル
- Handler: `server/internal/api/handler/dm_user_handler.go`
- Usecase: `server/internal/usecase/api/dm_user_usecase.go`
- Huma 型: `server/internal/api/huma/inputs.go`, `outputs.go`
- Repository: `server/internal/repository/dm_user_repository.go`, `dm_post_repository.go`
- DB: `server/internal/db/sharding.go`, `group_manager.go`
- モジュール: `server/go.mod` の module 行

### 5.2 新規 SKILL の配置
- test-auth-env: `.claude/skills/test-auth-env/SKILL.md` を新規作成
- react-use-effect-guard: `.claude/skills/react-use-effect-guard/SKILL.md` を新規作成
- いずれも SKILL.md のみでよい（references/ 等は本実装では不要）

### 5.3 実施順序の遵守
- go-test-generator → api-endpoint-creator → repository-generator → sharding-pattern → migration-helper → test-auth-env → react-use-effect-guard の順で実施する

## 6. 参考情報

### 6.1 関連ドキュメント
- 要件定義書: `.kiro/specs/0083-skills-updates/requirements.md`
- 計画書: `.kiro/specs/0083-skills-updates/SKILLS_REVIEW_PLAN.md`
- ステアリング: `.kiro/steering/tech.md`, `.kiro/steering/structure.md`
- 開発ルール: `CLAUDE.local.md`

### 6.2 技術スタック（参照先）
- API: Huma v2, Echo, humaecho
- レイヤー: Handler → Usecase → Service → Repository
- Repository: GORM, UUIDv7, GetShardingConnectionByUUID, GetTableNameFromUUID
- テスト: APP_ENV=test 必須
