# SKILLS 見直し計画書

作成日: 2025-02-03  
目的: `.claude/skills/` の現状調査に基づく修正・追加計画

---

## 1. 調査結果サマリ

### 1.1 プロジェクトの現状（tech.md / structure.md ベース）

| 項目 | 内容 |
|------|------|
| API層 | Huma v2 + Echo (humaecho アダプター)。ルーターは `internal/api/router/router.go` で Echo + Huma を設定 |
| レイヤー | Handler → **Usecase** → Service → Repository → DB（4層。Handler は Usecase を保持） |
| モジュールパス | `github.com/taku-o/go-webdb-template` |
| Huma型定義 | `internal/api/huma/` に `inputs.go` / `outputs.go`（types.go は存在しない） |
| Repository | `dm_user_repository.go`, `dm_post_repository.go`, `dm_news_repository.go` のみ（user_repository / post_repository は存在しない） |
| シャーディング | UUID ベース。`GetShardingConnectionByUUID`, `GetTableNameFromUUID`。ID は UUIDv7。GORM 使用 |
| DB定数 | `db.DBShardingTableCount`, `db.DBShardingTablesPerDB` |
| テスト実行 | **APP_ENV=test** 必須。指定しないと認証エラー（401）が発生（tech.md 明記） |
| マイグレーション | `db/migrations/` に master, sharding_1〜4 の他、master-mysql, sharding_*_mysql, view_master 等あり。スキーマは `db/schema/` の .hcl |

### 1.2 既存 SKILL 一覧と乖離の有無

| SKILL | 主な乖離・問題 | 対応方針 |
|-------|----------------|----------|
| **api-endpoint-creator** | ① 3層記載（実際は Handler→Usecase→Service→Repository）② インポートが `example` ③ Huma型が types.go ④ Handler が Service を直接保持している記載 | 修正 |
| **go-test-generator** | テスト実行コマンドが `APP_ENV=develop` になっている（**APP_ENV=test** が必須） | 修正 |
| **migration-helper** | スキーマパスが .sql の例。実際は .hcl も使用。master-mysql 等の存在は未記載 | 軽微修正 |
| **repository-generator** | ① user_repository / post_repository 参照（実際は dm_*）② GetShardingConnectionByID / GetTableName（実際は ByUUID / GetTableNameFromUUID）③ int64 ID（実際は UUID）④ 標準SQL版とGORM版の両方（実装は GORM のみ） | 修正 |
| **sharding-pattern** | ① GetShardingConnectionByID / GetTableName（実際は ByUUID / GetTableNameFromUUID）② user_id int64 / id int64 ベースの例（実際は UUID 文字列） | 修正 |
| **effective-go** | 汎用。プロジェクト固有のパスやスタックに依存していない | 修正不要 |
| **frontend-design** | 汎用。プロジェクト固有に依存していない | 修正不要 |
| **skill-creator** | 汎用。スキル作成手順のため変更不要 | 修正不要 |

---

## 2. 修正計画（既存 SKILL）

### 2.1 api-endpoint-creator

**修正内容:**

1. **description（frontmatter）**  
   - 「Handler/Service/Repositoryの3層」→「Handler → Usecase → Service → Repository の4層」に変更。

2. **アーキテクチャ**  
   - 図を「Handler → Usecase → Service → Repository」に変更。  
   - Handler は **Usecase** を保持し、Service を直接持たない旨を明記。

3. **使用フレームワーク**  
   - 「Echo」「Huma v2」「humaecho」は現状のまま（実装と一致）。  
   - 必要なら「Echo はルーティング・ミドルウェア、Huma は API 定義・OpenAPI」と補足。

4. **ディレクトリ構成**  
   - `internal/api/huma/types.go` → `internal/api/huma/inputs.go`, `outputs.go` に変更。  
   - 参照ファイルを実在するハンドラー（例: `dm_user_handler.go`）に合わせる。

5. **インポートパス**  
   - `github.com/example/go-webdb-template` → `github.com/taku-o/go-webdb-template` に統一。

6. **Handler パターン**  
   - 構造体は `entityUsecase *usecaseapi.EntityUsecase` を保持。  
   - 登録関数内で `h.entityUsecase.CreateEntity(...)` を呼ぶ形に変更。  
   - 認証エラー時は `huma.Error403Forbidden` を例に追加（既存ハンドラーに合わせる）。

7. **Service パターン**  
   - 「Handler から直接呼ばない」「Usecase 経由で呼ぶ」ことを注記。  
   - 必要なら Usecase の薄い例を追加（既存 dm_user_usecase を参照）。

### 2.2 go-test-generator

**修正内容:**

1. **テスト実行コマンド**  
   - `APP_ENV=develop` → **`APP_ENV=test`** に変更。  
   - コメントで「テスト時は必ず APP_ENV=test を指定すること。指定しないと認証エラー（401）が発生する」を追加。

2. **参照**  
   - tech.md の「テスト実行ルール（必須）」への言及を追加してもよい。

### 2.3 migration-helper

**修正内容（軽微）:**

1. スキーマが `.hcl` の場合の `atlas migrate diff` の `--to` の例を、プロジェクトの `db/schema/` 構成に合わせて追記（必要なら）。  
2. `db/migrations/` に master-mysql, view_master 等があることを「補足」程度で記載（主な対象は master / sharding_1〜4 のまま）。

### 2.4 repository-generator

**修正内容:**

1. **参照ファイル**  
   - `user_repository.go` / `post_repository.go` → `dm_user_repository.go`, `dm_post_repository.go`, `dm_news_repository.go` に変更。  
   - モデルは `model.DmUser`, `model.CreateDmUserRequest` 等に合わせる。

2. **構成の記載**  
   - 「標準SQL版とGORM版」→ 本プロジェクトでは **GORM 版を標準** とし、標準SQL版の記述は削除または「他プロジェクト向け」に変更。

3. **API の一致**  
   - `GetShardingConnectionByID(id int64, ...)` → `GetShardingConnectionByUUID(uuid string, tableBaseName string)`。  
   - `GetTableName("entities", id)` → `GetTableNameFromUUID("dm_users", uuid)`（戻り値 `(string, error)`）。  
   - `db.NewTableSelector(32, 8)` → `db.NewTableSelector(db.DBShardingTableCount, db.DBShardingTablesPerDB)`。

4. **ID**  
   - 作成時の ID は `idgen.GenerateUUIDv7()` を使用。  
   - 型は `int64` ではなく UUID 文字列（`string`）。

5. **CRUD 例**  
   - `conn.DB.WithContext(ctx).Table(tableName).Create(entity)` 等、GORM の実際の呼び出しに合わせる。  
   - List で全シャードを扱う場合は、既存 `dm_post_repository` 等の実装を参照して記載。

### 2.5 sharding-pattern

**修正内容:**

1. **シャードキー**  
   - 「user_id または エンティティの id」に加え、**UUID（string）** が主であることを明記。

2. **テーブル名取得**  
   - `tableSelector.GetTableName("users", userID)` → `tableSelector.GetTableNameFromUUID("dm_users", uuid)`（戻り値は `(string, error)`）。

3. **接続取得**  
   - `groupManager.GetShardingConnectionByID(userID, "users")` → `groupManager.GetShardingConnectionByUUID(uuid, "dm_users")`。

4. **参照ファイル**  
   - `server/internal/db/sharding.go`（GetTableNameFromUUID, ValidateTableName）, `server/internal/db/group_manager.go`。  
   - Repository は `dm_user_repository.go`, `dm_post_repository.go`。

5. **定数**  
   - テーブル数は `db.DBShardingTableCount`, `db.DBShardingTablesPerDB` を使用。

6. **クロステーブルクエリ**  
   - 実装が GORM と UUID ベースであることを反映し、必要なら「GetAllShardingConnections と GetTableNameFromUUID の組み合わせ」等を簡潔に記載。

---

## 3. 新規 SKILL 追加計画

### 3.1 useEffect 検討を促す SKILL（仮名: react-use-effect-guard）

**目的:**  
クライアント（React/Next.js）で `useEffect` を使おうとしたときに、本当に必要かどうか検討を促す。

**根拠:**  
- CLAUDE.local.md に「どうしても必要な時以外は useEffect を使用してはならない」とある。  
- React 18 / Next.js App Router では、データ取得は Server Component や Server Actions で行う方が推奨される。

**内容案:**

- **name:** `react-use-effect-guard` など  
- **description:**  
  - クライアントの React/Next.js コードで `useEffect` を追加・使用しようとするときに発動。  
  - 「本当に useEffect が必要か」を検討するよう促す。  
  - データ取得は Server Component / fetch in Server / Server Actions でできないか、イベントハンドラで十分でないかを確認する。

- **本文:**  
  - useEffect を使う前に確認すること:  
    - データ取得 → Server Component や Server Actions で可能か。  
    - イベントに紐づく処理 → イベントハンドラで十分か。  
    - 外部システムとの同期 → 本当にマウント/更新時に毎回必要か。  
  - どうしても必要な場合（例: ブラウザ API の購読、フォーカス制御、クライアント専用の初回実行）のみ useEffect を使用する。

**配置:**  
`.claude/skills/react-use-effect-guard/SKILL.md`（新規作成）

---

### 3.2 テスト認証エラー時に APP_ENV=test を指摘する SKILL（仮名: test-auth-env）

**目的:**  
テスト実行で認証エラー（401 等）が発生したときに、`APP_ENV=test` を指定していない可能性を指摘する。

**根拠:**  
- tech.md「テスト実行ルール（必須）」: `APP_ENV=test go test ./...` を指定すること。指定しないと認証エラー（401）が発生する。  
- CLAUDE.local.md: 認証エラーは「確認した」とは言わない。確認できなかった場合はそのタスクは未完了。  
- 同「頻発」: テストで認証エラーが起きたら APP_ENV=test を指定していない可能性を .kiro/steering/tech.md で確認すること。

**内容案:**

- **name:** `test-auth-env` など  
- **description:**  
  - Go テストや E2E/統合テストの実行で認証エラー（401 Unauthorized 等）が発生したときに使用。  
  - `APP_ENV=test` を指定していない可能性を指摘し、対処法（コマンド例・tech.md の確認）を案内する。

- **本文:**  
  - 認証エラーがテストで出た場合:  
    1. テスト実行コマンドに **APP_ENV=test** が付いているか確認する。  
    2. 例: `APP_ENV=test go test ./...`（server ディレクトリで実行する場合は `cd server && APP_ENV=test go test ./...`）。  
    3. `.kiro/steering/tech.md` の「テスト実行ルール（必須）」を確認する。  
  - 認証エラーが 1 件でも出た場合は「今回の修正とは関係ない」と判断せず、原因を調査する。

**配置:**  
`.claude/skills/test-auth-env/SKILL.md`（新規作成）

---

## 4. 実施順序の提案

1. **既存 SKILL の修正（優先度順）**  
   - go-test-generator（修正が少なく、テスト実行の誤りを防ぐため最優先）  
   - api-endpoint-creator（レイヤー・パス・型定義の乖離が大きい）  
   - repository-generator（UUID/GORM/参照ファイルの乖離が大きい）  
   - sharding-pattern（UUID/API 名の乖離）  
   - migration-helper（軽微な追記）

2. **新規 SKILL の追加**  
   - test-auth-env（認証エラーと APP_ENV=test の指摘）  
   - react-use-effect-guard（useEffect の検討を促す）

3. **見直し後の確認**  
   - 各 SKILL の description でトリガーが期待どおりか確認。  
   - 必要なら skill-creator の手順に従いパッケージングや検証を実行。

---

## 5. まとめ

| 種別 | 対象 | 主な対応 |
|------|------|----------|
| 修正 | api-endpoint-creator | 4層・Usecase・パス・Huma型・403 |
| 修正 | go-test-generator | APP_ENV=test 必須に変更 |
| 修正 | repository-generator | dm_*・UUID・GetTableNameFromUUID・GORM 主 |
| 修正 | sharding-pattern | UUID・GetShardingConnectionByUUID・GetTableNameFromUUID |
| 軽微 | migration-helper | .hcl / 他マイグレーションの補足（任意） |
| 新規 | test-auth-env | 認証エラー時の APP_ENV=test 指摘 |
| 新規 | react-use-effect-guard | useEffect 使用前の検討を促す |

この計画に沿って修正・追加を進めると、SKILLS がプロジェクトの現状と一致し、useEffect とテスト認証エラーに関するルールも SKILL で担保できます。
