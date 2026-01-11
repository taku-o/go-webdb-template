# APIサーバーのusecaseのソースコードの位置を変更するの要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Feature名**: 0060-layer-api
- **作成日**: 2026-01-27
- **関連Issue**: https://github.com/taku-o/go-webdb-template/issues/124

### 1.2 目的
APIサーバーのusecaseのソースコードの位置を一段下げ、`server/internal/usecase/api`ディレクトリに移動する。これにより、APIサーバー用のusecaseと、Adminアプリ用のusecase（`server/internal/usecase/admin`）、CLI用のusecase（`server/internal/usecase/cli`）を明確に分離し、ディレクトリ構造の一貫性を向上させる。

### 1.3 スコープ
- APIサーバー用のusecaseファイルを`server/internal/usecase/`から`server/internal/usecase/api/`に移動
- 移動対象ファイルのimport文を修正
- 移動対象ファイルを参照している全てのファイルのimport文を修正
- 関連ドキュメントの更新（プロジェクト構造、ファイル組織）

**本実装の範囲外**:
- Adminアプリ用のusecase（`server/internal/usecase/admin`）への影響（変更しない）
- CLI用のusecase（`server/internal/usecase/cli`）への影響（変更しない）
- 既存のビジネスロジックの変更（既存のロジックを維持）

## 2. 背景・現状分析

### 2.1 現在の状況

#### 2.1.1 usecase層のディレクトリ構造
- **APIサーバー用のusecase**: `server/internal/usecase/`に直接配置されている
  - `dm_user_usecase.go`
  - `dm_user_usecase_test.go`
  - `dm_post_usecase.go`
  - `dm_post_usecase_test.go`
  - `dm_jobqueue_usecase.go`
  - `dm_jobqueue_usecase_test.go`
  - `email_usecase.go`
  - `email_usecase_test.go`
  - `today_usecase.go`
  - `today_usecase_test.go`
- **Adminアプリ用のusecase**: `server/internal/usecase/admin/`に配置されている
- **CLI用のusecase**: `server/internal/usecase/cli/`に配置されている

#### 2.1.2 現在のディレクトリ構造の問題点
- **一貫性の欠如**: Adminアプリ用とCLI用のusecaseはサブディレクトリに配置されているが、APIサーバー用のusecaseはルートディレクトリに直接配置されている
- **識別の困難さ**: APIサーバー用のusecaseと、他の用途のusecaseを区別しにくい
- **拡張性の低さ**: 将来的に他の用途のusecaseを追加する際に、ディレクトリ構造が複雑になる可能性がある

### 2.2 現在の処理内容

#### 2.2.1 移動対象のusecase
1. **DmUserUsecase**: ユーザー管理のビジネスロジックを担当
2. **DmPostUsecase**: 投稿管理のビジネスロジックを担当
3. **DmJobqueueUsecase**: ジョブキュー管理のビジネスロジックを担当
4. **EmailUsecase**: メール送信のビジネスロジックを担当
5. **TodayUsecase**: 日付取得のビジネスロジックを担当

#### 2.2.2 移動対象ファイル
- `server/internal/usecase/dm_user_usecase.go` → `server/internal/usecase/api/dm_user_usecase.go`
- `server/internal/usecase/dm_user_usecase_test.go` → `server/internal/usecase/api/dm_user_usecase_test.go`
- `server/internal/usecase/dm_post_usecase.go` → `server/internal/usecase/api/dm_post_usecase.go`
- `server/internal/usecase/dm_post_usecase_test.go` → `server/internal/usecase/api/dm_post_usecase_test.go`
- `server/internal/usecase/dm_jobqueue_usecase.go` → `server/internal/usecase/api/dm_jobqueue_usecase.go`
- `server/internal/usecase/dm_jobqueue_usecase_test.go` → `server/internal/usecase/api/dm_jobqueue_usecase_test.go`
- `server/internal/usecase/email_usecase.go` → `server/internal/usecase/api/email_usecase.go`
- `server/internal/usecase/email_usecase_test.go` → `server/internal/usecase/api/email_usecase_test.go`
- `server/internal/usecase/today_usecase.go` → `server/internal/usecase/api/today_usecase.go`
- `server/internal/usecase/today_usecase_test.go` → `server/internal/usecase/api/today_usecase_test.go`

### 2.3 課題点
1. **ディレクトリ構造の不一致**: Adminアプリ用とCLI用のusecaseはサブディレクトリに配置されているが、APIサーバー用のusecaseはルートディレクトリに直接配置されている
2. **識別の困難さ**: APIサーバー用のusecaseと、他の用途のusecaseを区別しにくい
3. **拡張性の低さ**: 将来的に他の用途のusecaseを追加する際に、ディレクトリ構造が複雑になる可能性がある

### 2.4 本実装による改善点
1. **ディレクトリ構造の一貫性**: APIサーバー用、Adminアプリ用、CLI用のusecaseを全てサブディレクトリに配置することで、一貫性を向上
2. **識別の容易さ**: ディレクトリ名から用途を明確に識別できる
3. **拡張性の向上**: 将来的に他の用途のusecaseを追加する際に、ディレクトリ構造が明確になる

## 3. 機能要件

### 3.1 API用usecaseディレクトリの作成

#### 3.1.1 ディレクトリ構造
- **目的**: APIサーバー用のusecaseディレクトリを新規作成
- **実装内容**:
  - `server/internal/usecase/api`ディレクトリを新規作成

### 3.2 usecaseファイルの移動

#### 3.2.1 移動対象ファイル
- **目的**: APIサーバー用のusecaseファイルを`server/internal/usecase/`から`server/internal/usecase/api/`に移動
- **実装内容**:
  - `server/internal/usecase/dm_user_usecase.go`を`server/internal/usecase/api/dm_user_usecase.go`に移動
  - `server/internal/usecase/dm_user_usecase_test.go`を`server/internal/usecase/api/dm_user_usecase_test.go`に移動
  - `server/internal/usecase/dm_post_usecase.go`を`server/internal/usecase/api/dm_post_usecase.go`に移動
  - `server/internal/usecase/dm_post_usecase_test.go`を`server/internal/usecase/api/dm_post_usecase_test.go`に移動
  - `server/internal/usecase/dm_jobqueue_usecase.go`を`server/internal/usecase/api/dm_jobqueue_usecase.go`に移動
  - `server/internal/usecase/dm_jobqueue_usecase_test.go`を`server/internal/usecase/api/dm_jobqueue_usecase_test.go`に移動
  - `server/internal/usecase/email_usecase.go`を`server/internal/usecase/api/email_usecase.go`に移動
  - `server/internal/usecase/email_usecase_test.go`を`server/internal/usecase/api/email_usecase_test.go`に移動
  - `server/internal/usecase/today_usecase.go`を`server/internal/usecase/api/today_usecase.go`に移動
  - `server/internal/usecase/today_usecase_test.go`を`server/internal/usecase/api/today_usecase_test.go`に移動

#### 3.2.2 パッケージ名の変更
- **目的**: パッケージ名を`package usecase`から`package api`に変更（admin、cliと実装を統一するため）
- **実装内容**:
  - 移動後のファイルのパッケージ名を`package usecase`から`package api`に変更
  - 全てのファイルでパッケージ名を変更

### 3.3 import文の修正

#### 3.3.1 API Handler層のimport文修正
- **目的**: API Handler層のimport文を修正
- **実装内容**:
  - `server/internal/api/handler/dm_user_handler.go`のimport文と型参照を修正（`internal/usecase` → `internal/usecase/api`、`usecase.DmUserUsecase` → `api.DmUserUsecase`）
  - `server/internal/api/handler/dm_user_handler_test.go`のimport文と型参照を修正（該当する場合）
  - `server/internal/api/handler/dm_post_handler.go`のimport文と型参照を修正（`internal/usecase` → `internal/usecase/api`、`usecase.DmPostUsecase` → `api.DmPostUsecase`）
  - `server/internal/api/handler/dm_post_handler_test.go`のimport文と型参照を修正（該当する場合）
  - `server/internal/api/handler/dm_jobqueue_handler.go`のimport文と型参照を修正（`internal/usecase` → `internal/usecase/api`、`usecase.DmJobqueueUsecase` → `api.DmJobqueueUsecase`）
  - `server/internal/api/handler/dm_jobqueue_handler_test.go`のimport文と型参照を修正（`internal/usecase` → `internal/usecase/api`、`usecase.DmJobqueueUsecase` → `api.DmJobqueueUsecase`）
  - `server/internal/api/handler/email_handler.go`のimport文と型参照を修正（`internal/usecase` → `internal/usecase/api`、`usecase.EmailUsecase` → `api.EmailUsecase`）
  - `server/internal/api/handler/email_handler_test.go`のimport文と型参照を修正（`internal/usecase` → `internal/usecase/api`、`usecase.EmailUsecase` → `api.EmailUsecase`）
  - `server/internal/api/handler/today_handler.go`のimport文と型参照を修正（`internal/usecase` → `internal/usecase/api`、`usecase.TodayUsecase` → `api.TodayUsecase`）
  - `server/internal/api/handler/today_handler_test.go`のimport文と型参照を修正（`internal/usecase` → `internal/usecase/api`、`usecase.TodayUsecase` → `api.TodayUsecase`）

#### 3.3.2 main.goのimport文修正
- **目的**: main.goのimport文と型参照を修正
- **実装内容**:
  - `server/cmd/server/main.go`のimport文と型参照を修正（`internal/usecase` → `internal/usecase/api`、`usecase.DmJobqueueUsecase` → `api.DmJobqueueUsecase`）

#### 3.3.3 テストユーティリティのimport文修正
- **目的**: テストユーティリティのimport文とインターフェース参照を修正
- **実装内容**:
  - `server/test/testutil/db.go`のimport文とインターフェース参照を修正（`internal/usecase` → `internal/usecase/api`、`usecase.EmailServiceInterface` → `api.EmailServiceInterface`、`usecase.TemplateServiceInterface` → `api.TemplateServiceInterface`）

#### 3.3.4 Admin/CLI用usecaseのimport文修正
- **目的**: Admin/CLI用usecaseが参照しているインターフェースのimport文とインターフェース参照を修正
- **実装内容**:
  - `server/internal/usecase/admin/dm_user_register_usecase.go`のimport文とインターフェース参照を修正（`internal/usecase` → `internal/usecase/api`、`usecase.DmUserServiceInterface` → `api.DmUserServiceInterface`）
  - `server/internal/usecase/cli/list_dm_users_usecase.go`のimport文とインターフェース参照を修正（`internal/usecase` → `internal/usecase/api`、`usecase.DmUserServiceInterface` → `api.DmUserServiceInterface`）

### 3.4 ドキュメントの更新

#### 3.4.1 プロジェクト構造ドキュメントの更新
- **目的**: 新規作成するディレクトリをプロジェクト構造に反映
- **修正対象**: `docs/Project-Structure.md`
- **実装内容**:
  - `server/internal/usecase/api`ディレクトリを追加
  - 移動したファイルのパスを更新

#### 3.4.2 ファイル組織ドキュメントの更新
- **目的**: 新規作成するディレクトリをファイル組織に反映
- **修正対象**: `.kiro/steering/structure.md`
- **実装内容**:
  - `server/internal/usecase/api`ディレクトリを追加
  - 移動したファイルのパスを更新

## 4. 非機能要件

### 4.1 パフォーマンス
- **既存機能の維持**: 既存のAPIサーバーのパフォーマンスを維持
- **オーバーヘッド**: ディレクトリ移動によるパフォーマンスオーバーヘッドは無視できるレベル

### 4.2 信頼性
- **エラーハンドリング**: 既存のエラーハンドリングを維持
- **後方互換性**: 既存のAPIサーバーの動作を維持（出力形式、エラーメッセージなど）
- **データ整合性**: 既存のデータ処理ロジックを維持

### 4.3 保守性
- **コードの可読性**: ディレクトリ構造が明確で、各用途のusecaseを識別しやすい
- **一貫性**: Adminアプリ用、CLI用のusecaseと同じディレクトリ構造を採用することで、コードベース全体の一貫性を向上
- **テスト容易性**: 既存のテストが正常に動作することを確認

### 4.4 互換性
- **既存機能**: 既存のAPIサーバーに影響を与えない（動作は維持）
- **Admin/CLI用usecase**: 既存のAdmin/CLI用usecaseに影響を与えない（インターフェースの参照は維持）

## 5. 制約事項

### 5.1 技術的制約
- **既存のビジネスロジック**: 既存のビジネスロジックを変更しない
- **パッケージ名**: パッケージ名を`package usecase`から`package api`に変更（admin、cliと実装を統一するため）
- **インターフェース**: 既存のインターフェース定義を維持（`DmUserServiceInterface`、`DmPostServiceInterface`など）し、パッケージ名を`api`に変更

### 5.2 実装上の制約
- **ディレクトリ構造**: 既存のディレクトリ構造に従う（`server/internal/usecase/api`）
- **命名規則**: 既存の命名規則に従う（ファイル名、構造体名など）
- **import文**: 全ての参照箇所のimport文を修正する必要がある

### 5.3 動作環境
- **ローカル環境**: ローカル環境でAPIサーバーが正常に動作することを確認
- **CI環境**: CI環境でもAPIサーバーが正常に動作することを確認（該当する場合）
- **データベース**: 既存のデータベース接続が正常に動作することを前提

## 6. 受け入れ基準

### 6.1 API用usecaseディレクトリの作成
- [ ] `server/internal/usecase/api`ディレクトリが作成されている

### 6.2 usecaseファイルの移動
- [ ] `server/internal/usecase/dm_user_usecase.go`が`server/internal/usecase/api/dm_user_usecase.go`に移動されている
- [ ] `server/internal/usecase/dm_user_usecase_test.go`が`server/internal/usecase/api/dm_user_usecase_test.go`に移動されている
- [ ] `server/internal/usecase/dm_post_usecase.go`が`server/internal/usecase/api/dm_post_usecase.go`に移動されている
- [ ] `server/internal/usecase/dm_post_usecase_test.go`が`server/internal/usecase/api/dm_post_usecase_test.go`に移動されている
- [ ] `server/internal/usecase/dm_jobqueue_usecase.go`が`server/internal/usecase/api/dm_jobqueue_usecase.go`に移動されている
- [ ] `server/internal/usecase/dm_jobqueue_usecase_test.go`が`server/internal/usecase/api/dm_jobqueue_usecase_test.go`に移動されている
- [ ] `server/internal/usecase/email_usecase.go`が`server/internal/usecase/api/email_usecase.go`に移動されている
- [ ] `server/internal/usecase/email_usecase_test.go`が`server/internal/usecase/api/email_usecase_test.go`に移動されている
- [ ] `server/internal/usecase/today_usecase.go`が`server/internal/usecase/api/today_usecase.go`に移動されている
- [ ] `server/internal/usecase/today_usecase_test.go`が`server/internal/usecase/api/today_usecase_test.go`に移動されている
- [ ] 移動後のファイルのパッケージ名が`package api`に変更されている

### 6.3 import文の修正
- [ ] `server/internal/api/handler/dm_user_handler.go`のimport文と型参照が修正されている（`usecase.DmUserUsecase` → `api.DmUserUsecase`）
- [ ] `server/internal/api/handler/dm_user_handler_test.go`のimport文と型参照が修正されている（該当する場合）
- [ ] `server/internal/api/handler/dm_post_handler.go`のimport文と型参照が修正されている（`usecase.DmPostUsecase` → `api.DmPostUsecase`）
- [ ] `server/internal/api/handler/dm_post_handler_test.go`のimport文と型参照が修正されている（該当する場合）
- [ ] `server/internal/api/handler/dm_jobqueue_handler.go`のimport文と型参照が修正されている（`usecase.DmJobqueueUsecase` → `api.DmJobqueueUsecase`）
- [ ] `server/internal/api/handler/dm_jobqueue_handler_test.go`のimport文と型参照が修正されている（`usecase.DmJobqueueUsecase` → `api.DmJobqueueUsecase`）
- [ ] `server/internal/api/handler/email_handler.go`のimport文と型参照が修正されている（`usecase.EmailUsecase` → `api.EmailUsecase`）
- [ ] `server/internal/api/handler/email_handler_test.go`のimport文と型参照が修正されている（`usecase.EmailUsecase` → `api.EmailUsecase`）
- [ ] `server/internal/api/handler/today_handler.go`のimport文と型参照が修正されている（`usecase.TodayUsecase` → `api.TodayUsecase`）
- [ ] `server/internal/api/handler/today_handler_test.go`のimport文と型参照が修正されている（`usecase.TodayUsecase` → `api.TodayUsecase`）
- [ ] `server/cmd/server/main.go`のimport文と型参照が修正されている（`usecase.DmJobqueueUsecase` → `api.DmJobqueueUsecase`）
- [ ] `server/test/testutil/db.go`のimport文とインターフェース参照が修正されている（`usecase.EmailServiceInterface` → `api.EmailServiceInterface`、`usecase.TemplateServiceInterface` → `api.TemplateServiceInterface`）
- [ ] `server/internal/usecase/admin/dm_user_register_usecase.go`のimport文とインターフェース参照が修正されている（`usecase.DmUserServiceInterface` → `api.DmUserServiceInterface`）
- [ ] `server/internal/usecase/cli/list_dm_users_usecase.go`のimport文とインターフェース参照が修正されている（`usecase.DmUserServiceInterface` → `api.DmUserServiceInterface`）

### 6.4 動作確認
- [ ] ローカル環境でAPIサーバーが正常に動作する
- [ ] 既存のAPIエンドポイントが正常に動作する
- [ ] 既存のテストが全て通過する
- [ ] CI環境でAPIサーバーが正常に動作する（該当する場合）

### 6.5 テスト
- [ ] 既存のテストが全て通過する
- [ ] 移動後のファイルのテストが正常に動作する

### 6.6 ドキュメントの更新
- [ ] `docs/Project-Structure.md`に新規作成するディレクトリが追加されている
- [ ] `docs/Project-Structure.md`の移動したファイルのパスが更新されている
- [ ] `.kiro/steering/structure.md`に新規作成するディレクトリが追加されている
- [ ] `.kiro/steering/structure.md`の移動したファイルのパスが更新されている

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 移動が必要なファイル
- `server/internal/usecase/dm_user_usecase.go` → `server/internal/usecase/api/dm_user_usecase.go`
- `server/internal/usecase/dm_user_usecase_test.go` → `server/internal/usecase/api/dm_user_usecase_test.go`
- `server/internal/usecase/dm_post_usecase.go` → `server/internal/usecase/api/dm_post_usecase.go`
- `server/internal/usecase/dm_post_usecase_test.go` → `server/internal/usecase/api/dm_post_usecase_test.go`
- `server/internal/usecase/dm_jobqueue_usecase.go` → `server/internal/usecase/api/dm_jobqueue_usecase.go`
- `server/internal/usecase/dm_jobqueue_usecase_test.go` → `server/internal/usecase/api/dm_jobqueue_usecase_test.go`
- `server/internal/usecase/email_usecase.go` → `server/internal/usecase/api/email_usecase.go`
- `server/internal/usecase/email_usecase_test.go` → `server/internal/usecase/api/email_usecase_test.go`
- `server/internal/usecase/today_usecase.go` → `server/internal/usecase/api/today_usecase.go`
- `server/internal/usecase/today_usecase_test.go` → `server/internal/usecase/api/today_usecase_test.go`

#### 修正が必要なファイル
- `server/internal/api/handler/dm_user_handler.go`: import文の修正
- `server/internal/api/handler/dm_user_handler_test.go`: import文の修正（該当する場合）
- `server/internal/api/handler/dm_post_handler.go`: import文の修正
- `server/internal/api/handler/dm_post_handler_test.go`: import文の修正（該当する場合）
- `server/internal/api/handler/dm_jobqueue_handler.go`: import文の修正
- `server/internal/api/handler/dm_jobqueue_handler_test.go`: import文の修正
- `server/internal/api/handler/email_handler.go`: import文の修正
- `server/internal/api/handler/email_handler_test.go`: import文の修正
- `server/internal/api/handler/today_handler.go`: import文の修正
- `server/internal/api/handler/today_handler_test.go`: import文の修正
- `server/cmd/server/main.go`: import文の修正
- `server/test/testutil/db.go`: import文の修正
- `server/internal/usecase/admin/dm_user_register_usecase.go`: import文の修正
- `server/internal/usecase/cli/list_dm_users_usecase.go`: import文の修正
- `docs/Project-Structure.md`: ディレクトリ構造の更新
- `.kiro/steering/structure.md`: ディレクトリ構造の更新

### 7.2 既存機能への影響
- **既存のAPIサーバー**: 動作は維持されるが、import文の修正が必要
- **既存のAdmin/CLI用usecase**: インターフェースの参照は維持されるが、import文の修正が必要
- **既存のビジネスロジック**: 影響なし（ロジックは維持される）

## 8. 実装上の注意事項

### 8.1 usecaseファイルの移動
- **パッケージ名の変更**: 移動後のファイルのパッケージ名を`package usecase`から`package api`に変更
- **ファイル内容の変更**: パッケージ名を変更するため、ファイル内容も修正が必要

### 8.2 import文の修正
- **全ての参照箇所の修正**: 移動したusecaseを参照している全てのファイルのimport文と型参照を修正する必要がある
- **import文の形式**: `internal/usecase` → `internal/usecase/api`に変更
- **型参照の変更**: 型の参照は`usecase.DmUserUsecase`から`api.DmUserUsecase`に変更
- **インターフェース参照の変更**: インターフェースの参照も`usecase.DmUserServiceInterface`から`api.DmUserServiceInterface`に変更
- **エイリアスの確認**: import文にエイリアスが使用されている場合は、エイリアスも適切に修正

### 8.3 テストの確認
- **既存テストの動作確認**: 移動後、既存のテストが全て正常に動作することを確認
- **テストファイルの移動**: テストファイルも同時に移動する必要がある

### 8.4 ドキュメントの更新
- **プロジェクト構造ドキュメント**: 新規作成するディレクトリと移動したファイルのパスを反映
- **ファイル組織ドキュメント**: 新規作成するディレクトリと移動したファイルのパスを反映
- **一貫性**: 全てのドキュメントで同じディレクトリ構造を記載

## 9. 参考情報

### 9.1 関連ドキュメント
- `docs/Project-Structure.md`: プロジェクト構造ドキュメント
- `.kiro/steering/structure.md`: ファイル組織とコードパターン
- `.kiro/steering/tech.md`: 技術スタックとアーキテクチャ

### 9.2 既存実装の参考
- `server/internal/usecase/admin/`: Adminアプリ用のusecase層の実装パターン
- `server/internal/usecase/cli/`: CLI用のusecase層の実装パターン
- `server/internal/api/handler/`: API Handler層の実装パターン

### 9.3 技術スタック
- **言語**: Go
- **アーキテクチャ**: レイヤードアーキテクチャ（handler -> usecase -> service -> repository -> db）
- **テスト**: `testing`（標準ライブラリ）、`github.com/stretchr/testify`（アサーション、モック）

### 9.4 ディレクトリ構造の比較

| 項目 | 現在（修正前） | 修正後 |
|------|---------------|--------|
| APIサーバー用usecase | `usecase/dm_user_usecase.go` | `usecase/api/dm_user_usecase.go` |
| Adminアプリ用usecase | `usecase/admin/dm_user_register_usecase.go` | `usecase/admin/dm_user_register_usecase.go`（変更なし） |
| CLI用usecase | `usecase/cli/list_dm_users_usecase.go` | `usecase/cli/list_dm_users_usecase.go`（変更なし） |
| パッケージ名 | `package usecase` | `package api` |
| importパス | `internal/usecase` | `internal/usecase/api` |
| 型参照 | `usecase.DmUserUsecase` | `api.DmUserUsecase` |

### 9.5 移動対象ファイル一覧

| 現在のパス | 移動後のパス |
|-----------|-------------|
| `server/internal/usecase/dm_user_usecase.go` | `server/internal/usecase/api/dm_user_usecase.go` |
| `server/internal/usecase/dm_user_usecase_test.go` | `server/internal/usecase/api/dm_user_usecase_test.go` |
| `server/internal/usecase/dm_post_usecase.go` | `server/internal/usecase/api/dm_post_usecase.go` |
| `server/internal/usecase/dm_post_usecase_test.go` | `server/internal/usecase/api/dm_post_usecase_test.go` |
| `server/internal/usecase/dm_jobqueue_usecase.go` | `server/internal/usecase/api/dm_jobqueue_usecase.go` |
| `server/internal/usecase/dm_jobqueue_usecase_test.go` | `server/internal/usecase/api/dm_jobqueue_usecase_test.go` |
| `server/internal/usecase/email_usecase.go` | `server/internal/usecase/api/email_usecase.go` |
| `server/internal/usecase/email_usecase_test.go` | `server/internal/usecase/api/email_usecase_test.go` |
| `server/internal/usecase/today_usecase.go` | `server/internal/usecase/api/today_usecase.go` |
| `server/internal/usecase/today_usecase_test.go` | `server/internal/usecase/api/today_usecase_test.go` |

### 9.6 import文修正対象ファイル一覧

| ファイルパス | 修正内容 |
|------------|---------|
| `server/internal/api/handler/dm_user_handler.go` | `internal/usecase` → `internal/usecase/api` |
| `server/internal/api/handler/dm_user_handler_test.go` | `internal/usecase` → `internal/usecase/api`（該当する場合） |
| `server/internal/api/handler/dm_post_handler.go` | `internal/usecase` → `internal/usecase/api` |
| `server/internal/api/handler/dm_post_handler_test.go` | `internal/usecase` → `internal/usecase/api`（該当する場合） |
| `server/internal/api/handler/dm_jobqueue_handler.go` | `internal/usecase` → `internal/usecase/api` |
| `server/internal/api/handler/dm_jobqueue_handler_test.go` | `internal/usecase` → `internal/usecase/api` |
| `server/internal/api/handler/email_handler.go` | `internal/usecase` → `internal/usecase/api` |
| `server/internal/api/handler/email_handler_test.go` | `internal/usecase` → `internal/usecase/api` |
| `server/internal/api/handler/today_handler.go` | `internal/usecase` → `internal/usecase/api` |
| `server/internal/api/handler/today_handler_test.go` | `internal/usecase` → `internal/usecase/api` |
| `server/cmd/server/main.go` | `internal/usecase` → `internal/usecase/api` |
| `server/test/testutil/db.go` | `internal/usecase` → `internal/usecase/api` |
| `server/internal/usecase/admin/dm_user_register_usecase.go` | `internal/usecase` → `internal/usecase/api` |
| `server/internal/usecase/cli/list_dm_users_usecase.go` | `internal/usecase` → `internal/usecase/api` |
