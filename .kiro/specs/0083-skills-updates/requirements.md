# SKILLS 見直し要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: なし
- **Issueタイトル**: —
- **Feature名**: 0083-skills-updates
- **作成日**: 2026-02-03
- **参照**: 同ディレクトリの `SKILLS_REVIEW_PLAN.md`（本要件の元となった計画書）

### 1.2 目的
`.claude/skills/` 配下の SKILL をプロジェクトの現状（tech.md / structure.md / 実装）に合わせて修正し、不足している SKILL を追加する。あわせて、クライアントの useEffect 使用検討を促す SKILL と、テスト認証エラー時に APP_ENV=test を指摘する SKILL を新規追加する。

### 1.3 スコープ
- 既存 SKILL の修正: api-endpoint-creator, go-test-generator, repository-generator, sharding-pattern, migration-helper（軽微）
- 新規 SKILL の追加: test-auth-env, react-use-effect-guard
- 対象ディレクトリ: `.claude/skills/`
- 参照するプロジェクト定義: `.kiro/steering/tech.md`, `.kiro/steering/structure.md`, 実装コード

**本実装の範囲外**:
- effective-go, frontend-design, skill-creator の内容変更（汎用のため修正不要と判断）
- .kiro/steering/ の tech.md / structure.md の変更
- アプリケーション本体の機能変更

## 2. 背景・現状分析

### 2.1 現在の実装

#### 2.1.1 プロジェクトの現状（SKILL が従うべき基準）
- **API層**: Huma v2 + Echo (humaecho)。Handler → Usecase → Service → Repository（4層）
- **モジュールパス**: `github.com/taku-o/go-webdb-template`
- **Huma型定義**: `internal/api/huma/inputs.go`, `outputs.go`（types.go は存在しない）
- **Repository**: dm_user_repository, dm_post_repository, dm_news_repository。GORM 使用。UUID（UUIDv7）ベース
- **シャーディング**: GetShardingConnectionByUUID, GetTableNameFromUUID。db.DBShardingTableCount, db.DBShardingTablesPerDB
- **テスト実行**: `APP_ENV=test` 必須。指定しないと認証エラー（401）
- **マイグレーション**: db/migrations/ に master, sharding_1〜4 の他、master-mysql, view_master 等。スキーマは db/schema/ の .hcl

#### 2.1.2 既存 SKILL の構成
- `.claude/skills/` 配下に api-endpoint-creator, effective-go, frontend-design, go-test-generator, migration-helper, repository-generator, sharding-pattern, skill-creator が存在
- 各 SKILL は SKILL.md の YAML frontmatter（name, description）と本文で構成

### 2.2 課題点
1. **api-endpoint-creator**: 3層記載・example パス・types.go 参照・Handler が Service を直接保持する記載になっており、現行の4層・Usecase・inputs/outputs・taku-o パスと乖離している
2. **go-test-generator**: テスト実行コマンドが `APP_ENV=develop` と記載されており、tech.md で必須とされる `APP_ENV=test` と一致していない
3. **repository-generator**: user_repository/post_repository 参照・GetShardingConnectionByID・int64 ID・標準SQL/GORM 両方の記載で、実装の dm_*・ByUUID・GetTableNameFromUUID・UUID・GORM 主と乖離している
4. **sharding-pattern**: GetShardingConnectionByID/GetTableName・int64 ベースの記載で、実装の ByUUID/GetTableNameFromUUID・UUID ベースと乖離している
5. **migration-helper**: スキーマが .sql の例が中心で、.hcl や master-mysql 等の補足がない
6. **useEffect 検討の不足**: CLAUDE.local.md で「どうしても必要な時以外は useEffect を使用してはならない」とあるが、useEffect 使用前に検討を促す SKILL が存在しない
7. **テスト認証エラー時の案内不足**: テストで認証エラーが発生した際に APP_ENV=test 未指定の可能性を指摘する SKILL が存在しない

### 2.3 本実装による改善点
1. **SKILL と実装の一致**: 既存 SKILL の記載がプロジェクトの現状と一致し、AI が正しいパターンでコードを生成・修正できる
2. **テスト実行の誤り防止**: go-test-generator で APP_ENV=test 必須が明示され、認証エラーを防ぐ
3. **useEffect 使用の抑制**: react-use-effect-guard により、useEffect 使用前に Server Component / イベントハンドラでの代替を検討させる
4. **認証エラー時の対処案内**: test-auth-env により、テストで認証エラーが出た際に APP_ENV=test の確認を促す

## 3. 機能要件

### 3.1 既存 SKILL の修正

#### 3.1.1 api-endpoint-creator の修正
- **目的**: Handler → Usecase → Service → Repository の4層と現行実装（Huma/Echo、inputs/outputs、taku-o パス）に合わせる
- **description（frontmatter）**: 「Handler/Service/Repositoryの3層」を「Handler → Usecase → Service → Repository の4層」に変更
- **アーキテクチャ**: 図・説明を「Handler → Usecase → Service → Repository」に変更。Handler は Usecase を保持し、Service を直接持たない旨を明記
- **ディレクトリ・参照ファイル**: `internal/api/huma/types.go` ではなく `inputs.go`, `outputs.go` および実在するハンドラー（例: dm_user_handler.go）を参照
- **インポートパス**: `github.com/example/go-webdb-template` を `github.com/taku-o/go-webdb-template` に統一
- **Handler パターン**: 構造体は entityUsecase（Usecase）を保持し、登録関数内で Usecase を呼ぶ形の例に変更
- **認証エラー時**: レスポンス例に `huma.Error403Forbidden` を追加

#### 3.1.2 go-test-generator の修正
- **目的**: テスト実行時は APP_ENV=test が必須であることを SKILL に反映する
- **テスト実行コマンド**: 記載されている環境変数を `APP_ENV=develop` から `APP_ENV=test` に変更
- **注意記載**: 「テスト時は必ず APP_ENV=test を指定すること。指定しないと認証エラー（401）が発生する」に相当する注意を追加。tech.md の「テスト実行ルール（必須）」への言及を追加してもよい

#### 3.1.3 repository-generator の修正
- **目的**: 参照ファイルを dm_* に、接続・テーブル名取得を UUID ベース API に、ID を UUID、構成を GORM 主に合わせる
- **参照ファイル**: user_repository.go / post_repository.go ではなく dm_user_repository.go, dm_post_repository.go, dm_news_repository.go を参照
- **接続取得**: GetShardingConnectionByID(id int64) ではなく GetShardingConnectionByUUID(uuid string, tableBaseName string) の例に変更
- **テーブル名取得**: GetTableName ではなく GetTableNameFromUUID(baseName, uuid)（戻り値 (string, error)）の例に変更
- **TableSelector**: db.NewTableSelector(32, 8) ではなく db.NewTableSelector(db.DBShardingTableCount, db.DBShardingTablesPerDB) に変更
- **ID 生成**: 作成時の ID は idgen.GenerateUUIDv7() を使用し、型は UUID 文字列（string）
- **CRUD 方針**: 本プロジェクトでは GORM を標準とし、標準SQL版の記述は削除または他プロジェクト向けと明記

#### 3.1.4 sharding-pattern の修正
- **目的**: シャードキー・接続取得・テーブル名取得を UUID ベースの API に合わせる
- **シャードキー**: user_id またはエンティティ id に加え、UUID（string）が主であることを明記
- **テーブル名取得**: GetTableName ではなく GetTableNameFromUUID(baseName, uuid)（戻り値 (string, error)）の例に変更
- **接続取得**: GetShardingConnectionByID ではなく GetShardingConnectionByUUID(uuid, tableBaseName) の例に変更
- **参照ファイル**: server/internal/db/sharding.go（GetTableNameFromUUID, ValidateTableName）, group_manager.go および dm_user_repository.go, dm_post_repository.go を参照
- **定数**: db.DBShardingTableCount, db.DBShardingTablesPerDB を使用

#### 3.1.5 migration-helper の修正
- **目的**: スキーマが .hcl の場合や master-mysql 等の存在を補足する
- **スキーマ**: db/schema/ の .hcl 構成に合わせた atlas migrate diff の --to の例を追記する
- **マイグレーションディレクトリ**: db/migrations/ に master-mysql, view_master 等があることを補足する

### 3.2 新規 SKILL の追加

#### 3.2.1 test-auth-env の追加
- **目的**: テスト実行で認証エラー（401 等）が発生したときに、APP_ENV=test を指定していない可能性を指摘し、対処法を案内する
- **配置**: `.claude/skills/test-auth-env/SKILL.md` を新規作成
- **description（frontmatter）**: Go テストや E2E/統合テストで認証エラーが発生したときに使用すること、および APP_ENV=test 未指定の可能性を指摘することを含む
- **本文**: テスト実行コマンドに APP_ENV=test が付いているか確認する手順、コマンド例（例: `APP_ENV=test go test ./...` または `cd server && APP_ENV=test go test ./...`）、.kiro/steering/tech.md の「テスト実行ルール（必須）」の確認を含む
- **方針**: 認証エラーが1件でも出た場合は「今回の修正とは関係ない」と判断せず原因を調査する旨を記載

#### 3.2.2 react-use-effect-guard の追加
- **目的**: クライアント（React/Next.js）で useEffect を使おうとしたときに、本当に必要かどうか検討を促す
- **配置**: `.claude/skills/react-use-effect-guard/SKILL.md` を新規作成
- **description（frontmatter）**: クライアントの React/Next.js コードで useEffect を追加・使用しようとするときに発動すること、「本当に useEffect が必要か」を検討するよう促すことを含む
- **本文**: データ取得は Server Component や Server Actions で可能か、イベントに紐づく処理はイベントハンドラで十分でないか、外部システムとの同期は本当にマウント/更新時に毎回必要か、に触れる
- **許容場合**: どうしても必要な場合（例: ブラウザ API の購読、フォーカス制御、クライアント専用の初回実行）のみ useEffect を使用する旨を記載

## 4. 非機能要件

### 4.1 実施順序
1. go-test-generator の修正（修正量が少なく、テスト実行の誤り防止のため最優先）
2. api-endpoint-creator の修正
3. repository-generator の修正
4. sharding-pattern の修正
5. migration-helper の修正
6. test-auth-env の新規追加
7. react-use-effect-guard の新規追加

### 4.2 品質
- 修正後の SKILL の description は、トリガーが期待どおりに働くよう、内容とキーワードが一致していること
- 既存の effective-go, frontend-design, skill-creator のファイルは変更しないこと

### 4.3 ドキュメント整合性
- 各 SKILL の記載は .kiro/steering/tech.md および .kiro/steering/structure.md と矛盾しないこと
- コード例・パス・API 名は実装（server/internal/ 等）と一致すること

## 5. 制約事項

### 5.1 既存システムとの関係
- **effective-go, frontend-design, skill-creator**: 汎用 SKILL のため内容を変更しない
- **.kiro/steering/**: tech.md / structure.md は本実装では変更しない
- **アプリケーション本体**: サーバー・クライアントの機能変更は行わない

### 5.2 SKILL の形式
- 各 SKILL は SKILL.md の YAML frontmatter（name, description）と Markdown 本文で構成
- description はトリガー判定に使われるため、使用場面と内容を明確に記載する

### 5.3 参照する実装
- 修正内容は既存の server/internal/api/, server/internal/repository/, server/internal/db/ 等の実装に合わせる
- モジュールパスは go.mod の module に合わせる（github.com/taku-o/go-webdb-template）

## 6. 受け入れ基準

### 6.1 api-endpoint-creator の修正
- [ ] description に「Handler → Usecase → Service → Repository の4層」が含まれる
- [ ] アーキテクチャ説明で Handler は Usecase を保持し Service を直接持たない旨が明記されている
- [ ] 参照ファイルに inputs.go, outputs.go および実在するハンドラー（例: dm_user_handler.go）が記載されている
- [ ] インポートパスが github.com/taku-o/go-webdb-template である
- [ ] Handler の構造体・登録関数の例で Usecase を保持し Usecase を呼ぶ形である
- [ ] 認証エラー時のレスポンス例に huma.Error403Forbidden が含まれる

### 6.2 go-test-generator の修正
- [ ] テスト実行コマンドの環境変数が APP_ENV=test である（APP_ENV=develop ではない）
- [ ] テスト時は必ず APP_ENV=test を指定すること、指定しないと認証エラー（401）が発生する旨の注意が含まれる

### 6.3 repository-generator の修正
- [ ] 参照ファイルに dm_user_repository.go, dm_post_repository.go, dm_news_repository.go が記載されている
- [ ] 接続取得の例で GetShardingConnectionByUUID が使用されている
- [ ] テーブル名取得の例で GetTableNameFromUUID が使用され、戻り値が (string, error) として扱われている
- [ ] TableSelector の初期化で db.NewTableSelector(db.DBShardingTableCount, db.DBShardingTablesPerDB) が使用されている
- [ ] エンティティ ID 生成で idgen.GenerateUUIDv7() が使用され、型が UUID 文字列（string）として扱われている
- [ ] 本プロジェクトでは GORM を標準とし、標準SQL版は削除または他プロジェクト向けと明記されている

### 6.4 sharding-pattern の修正
- [ ] シャードキーで UUID（string）が主であることが明記されている
- [ ] テーブル名取得の例で GetTableNameFromUUID が使用され、戻り値が (string, error) として扱われている
- [ ] 接続取得の例で GetShardingConnectionByUUID が使用されている
- [ ] 参照ファイルに sharding.go, group_manager.go, dm_user_repository.go, dm_post_repository.go が記載されている
- [ ] 定数として db.DBShardingTableCount, db.DBShardingTablesPerDB が使用されている

### 6.5 migration-helper の修正
- [ ] .hcl スキーマおよび master-mysql 等の補足が追加されている

### 6.6 test-auth-env の新規追加
- [ ] .claude/skills/test-auth-env/SKILL.md が存在する
- [ ] description に認証エラー発生時の使用、APP_ENV=test 未指定の可能性の指摘が含まれる
- [ ] 本文に APP_ENV=test の確認手順、コマンド例、.kiro/steering/tech.md の確認が含まれる
- [ ] 認証エラーが1件でも出た場合は原因を調査する旨が含まれる

### 6.7 react-use-effect-guard の新規追加
- [ ] .claude/skills/react-use-effect-guard/SKILL.md が存在する
- [ ] description に useEffect 追加・使用時に発動、「本当に useEffect が必要か」の検討を促すことが含まれる
- [ ] 本文にデータ取得・イベントハンドラ・外部同期の確認事項が含まれる
- [ ] どうしても必要な場合のみ useEffect を使用する旨が含まれる

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 既存 SKILL（修正）
- `.claude/skills/api-endpoint-creator/SKILL.md`
- `.claude/skills/go-test-generator/SKILL.md`
- `.claude/skills/repository-generator/SKILL.md`
- `.claude/skills/sharding-pattern/SKILL.md`
- `.claude/skills/migration-helper/SKILL.md`

### 7.2 新規追加が必要なディレクトリ・ファイル

#### ディレクトリ・ファイル
- `.claude/skills/test-auth-env/`: test-auth-env SKILL 用
  - `.claude/skills/test-auth-env/SKILL.md`
- `.claude/skills/react-use-effect-guard/`: react-use-effect-guard SKILL 用
  - `.claude/skills/react-use-effect-guard/SKILL.md`

### 7.3 変更しないファイル
- `.claude/skills/effective-go/SKILL.md`
- `.claude/skills/frontend-design/SKILL.md`
- `.claude/skills/skill-creator/SKILL.md`
- `.kiro/steering/tech.md`
- `.kiro/steering/structure.md`

## 8. 実装上の注意事項

### 8.1 修正時の参照
- Handler の実例は `server/internal/api/handler/dm_user_handler.go` 等を参照する
- Repository の実例は `server/internal/repository/dm_user_repository.go`, `dm_post_repository.go` 等を参照する
- シャーディングは `server/internal/db/sharding.go`, `group_manager.go` を参照する
- モジュールパスは `server/go.mod` の module に合わせる

### 8.2 SKILL の description（frontmatter）
- description は AI が SKILL をいつ使うかの判定に使うため、使用場面と内容を明確に記載する
- キーワード（例: useEffect, 認証エラー, APP_ENV=test）を適切に含める

### 8.3 実施順序
- go-test-generator → api-endpoint-creator → repository-generator → sharding-pattern → migration-helper → test-auth-env → react-use-effect-guard の順で実施する

## 9. 参考情報

### 9.1 関連ドキュメント
- 計画書: `.kiro/specs/0083-skills-updates/SKILLS_REVIEW_PLAN.md`
- ステアリング: `.kiro/steering/tech.md`, `.kiro/steering/structure.md`
- 開発ルール: `CLAUDE.local.md`（useEffect 禁止方針、認証エラー時の対応）

### 9.2 既存 SKILL の参照
- `api-endpoint-creator`: Handler/Usecase の実例は dm_user_handler, dm_user_usecase を参照
- `repository-generator`: dm_user_repository, dm_post_repository の実装を参照
- `sharding-pattern`: db/sharding.go, db/group_manager.go を参照

### 9.3 技術スタック（本実装で触れる範囲）
- **SKILL 形式**: YAML frontmatter + Markdown（SKILL.md）
- **プロジェクト**: Go (server), Next.js/React (client), Huma v2 + Echo, GORM, UUIDv7 シャーディング
