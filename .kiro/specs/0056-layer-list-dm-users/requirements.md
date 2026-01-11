# server/cmd/list-dm-usersのレイヤー構造修正の要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Feature名**: 0056-layer-list-dm-users
- **作成日**: 2026-01-27
- **関連Issue**: https://github.com/taku-o/go-webdb-template/issues/116

### 1.2 目的
`server/cmd/list-dm-users`の実装を、APIサーバーと同じレイヤー構造（usecase -> service -> repository -> model）に変更する。これにより、CLIコマンドとAPIサーバーで一貫したアーキテクチャを実現し、コードの保守性と再利用性を向上させる。

### 1.3 スコープ
- CLIコマンドのバリデーションと入出力制御をCLI層に集約
- `server/internal/usecase/cli`ディレクトリの作成とCLI用usecaseの実装
- 既存のservice、repository、model層の活用
- main.goの簡素化（エントリーポイントと入出力制御のみ）
- 関連ドキュメントの更新（アーキテクチャ、プロジェクト構造、CLIツールドキュメント）

**本実装の範囲外**:
- 既存のservice、repository、model層の変更（既存実装をそのまま使用）
- 既存のAPIサーバーのusecase層への影響（APIサーバーは変更しない）
- 他のCLIコマンドへの影響（本実装は`list-dm-users`のみを対象）

## 2. 背景・現状分析

### 2.1 現在の状況
- **CLI実装**: `server/cmd/list-dm-users/main.go`に直接実装されている
- **レイヤー構造**: main.go → service → repository → model という構造
- **バリデーション**: main.go内の`validateLimit()`関数で実装
- **入出力制御**: main.go内の`printDmUsersTSV()`関数で実装
- **ビジネスロジック**: service層（`DmUserService.ListDmUsers()`）を直接呼び出し
- **usecase層**: CLI用のusecase層が存在しない

### 2.2 課題点
1. **レイヤー構造の不一致**: APIサーバーはusecase層を使用しているが、CLIはservice層を直接呼び出している
2. **責務の混在**: main.goにバリデーション、入出力制御、ビジネスロジックの呼び出しが混在している
3. **再利用性の低さ**: CLI用のビジネスロジックがmain.goに直接実装されており、他のCLIコマンドから再利用できない
4. **テストの困難さ**: main.goにロジックが集中しているため、単体テストが困難

### 2.3 本実装による改善点
1. **アーキテクチャの一貫性**: APIサーバーとCLIで同じレイヤー構造を採用
2. **責務の明確化**: CLI層（main.go）はエントリーポイントと入出力制御のみを担当
3. **再利用性の向上**: CLI用のusecase層を作成することで、他のCLIコマンドからも再利用可能
4. **テスト容易性**: usecase層を独立してテストできるようになる

## 3. 機能要件

### 3.1 CLI用usecase層の作成

#### 3.1.1 ディレクトリ構造
- **目的**: CLI用のusecase層を独立したディレクトリに配置
- **実装内容**:
  - `server/internal/usecase/cli`ディレクトリを作成
  - `server/internal/usecase/cli/list_dm_users_usecase.go`を作成
  - `server/internal/usecase/cli/list_dm_users_usecase_test.go`を作成（テストファイル）

#### 3.1.2 ListDmUsersUsecaseの実装
- **目的**: CLI用のdm_user一覧取得usecaseを実装
- **実装内容**:
  - `ListDmUsersUsecase`構造体を定義
  - `DmUserServiceInterface`を依存として注入
  - `ListDmUsers(ctx context.Context, limit, offset int) ([]*model.DmUser, error)`メソッドを実装
  - service層の`ListDmUsers()`を呼び出して結果を返す
- **インターフェース**: 既存の`DmUserServiceInterface`を使用（`internal/usecase/dm_user_usecase.go`で定義済み）

### 3.2 main.goの簡素化

#### 3.2.1 エントリーポイントの責務
- **目的**: main.goをエントリーポイントと入出力制御のみに限定
- **実装内容**:
  - コマンドライン引数の解析（`flag`パッケージを使用）
  - 引数のバリデーション（`validateLimit()`関数）
  - 設定ファイルの読み込み（`config.Load()`）
  - GroupManagerの初期化（`db.NewGroupManager()`）
  - usecase層の初期化と実行
  - 結果の出力（`printDmUsersTSV()`関数）

#### 3.2.2 バリデーション関数
- **目的**: コマンドライン引数のバリデーションをmain.goに保持
- **実装内容**:
  - `validateLimit(limit int) (int, error, bool)`関数を維持
  - 最小値チェック（1以上）
  - 最大値チェック（100以下、超過時は警告を出して100に制限）
  - エラーと警告の返却

#### 3.2.3 出力関数
- **目的**: TSV形式での出力をmain.goに保持
- **実装内容**:
  - `printDmUsersTSV(dmUsers []*model.DmUser)`関数を維持
  - ヘッダー行の出力（`ID\tName\tEmail\tCreatedAt\tUpdatedAt`）
  - 各ユーザー情報のTSV形式での出力
  - 日時はRFC3339形式で出力

### 3.3 依存関係の注入

#### 3.3.1 usecase層の初期化
- **目的**: usecase層を適切に初期化してmain.goで使用
- **実装内容**:
  - Repository層の初期化（`repository.NewDmUserRepository(groupManager)`）
  - Service層の初期化（`service.NewDmUserService(dmUserRepo)`）
  - Usecase層の初期化（`cli.NewListDmUsersUsecase(dmUserService)`）
  - usecase層の`ListDmUsers()`メソッドを呼び出し

### 3.4 ドキュメントの更新

#### 3.4.1 アーキテクチャドキュメントの更新
- **目的**: CLIコマンドのレイヤー構造をアーキテクチャドキュメントに反映
- **修正対象**: `docs/Architecture.md`
- **実装内容**:
  - CLIコマンドのレイヤー構造を追加（usecase層を含む）
  - CLIコマンドのアーキテクチャ図を更新
  - CLI用usecase層の説明を追加

#### 3.4.2 プロジェクト構造ドキュメントの更新
- **目的**: 新規作成する`server/internal/usecase/cli`ディレクトリをプロジェクト構造に反映
- **修正対象**: `docs/Project-Structure.md`
- **実装内容**:
  - `server/internal/usecase/cli`ディレクトリを追加
  - `server/internal/usecase/cli/list_dm_users_usecase.go`を追加

#### 3.4.3 CLIツールドキュメントの更新
- **目的**: CLIコマンドのアーキテクチャ図を更新してusecase層を含める
- **修正対象**: `docs/Command-Line-Tool.md`
- **実装内容**:
  - CLIコマンドのアーキテクチャ図を更新（usecase層を追加）
  - レイヤー構造の説明を更新（main.go → usecase → service → repository → model）

#### 3.4.4 ファイル組織ドキュメントの更新
- **目的**: 新規作成する`server/internal/usecase/cli`ディレクトリをファイル組織に反映
- **修正対象**: `.kiro/steering/structure.md`
- **実装内容**:
  - `server/internal/usecase/cli`ディレクトリを追加
  - CLI用usecase層の説明を追加

## 4. 非機能要件

### 4.1 パフォーマンス
- **既存機能の維持**: 既存のservice層とrepository層のパフォーマンスを維持
- **オーバーヘッド**: usecase層の追加によるパフォーマンスオーバーヘッドは無視できるレベル

### 4.2 信頼性
- **エラーハンドリング**: 既存のエラーハンドリングを維持
- **後方互換性**: 既存のCLIコマンドの動作を維持（出力形式、エラーメッセージなど）

### 4.3 保守性
- **コードの可読性**: レイヤー構造が明確で、各レイヤーの責務が明確
- **一貫性**: APIサーバーと同じレイヤー構造を採用することで、コードベース全体の一貫性を向上
- **テスト容易性**: usecase層を独立してテストできる

### 4.4 互換性
- **既存機能**: 既存のservice、repository、model層に影響を与えない
- **CLIコマンドの動作**: 既存のCLIコマンドの動作（引数、出力形式、エラーメッセージ）を維持

## 5. 制約事項

### 5.1 技術的制約
- **既存のservice層**: 既存の`DmUserService`とそのインターフェース（`DmUserServiceInterface`）を使用
- **既存のrepository層**: 既存の`DmUserRepository`を使用
- **既存のmodel層**: 既存の`model.DmUser`を使用

### 5.2 実装上の制約
- **ディレクトリ構造**: `server/internal/usecase/cli`ディレクトリを作成
- **命名規則**: 既存のusecase層の命名規則に従う（`ListDmUsersUsecase`）
- **インターフェース**: 既存の`DmUserServiceInterface`を使用（新規インターフェースは作成しない）

### 5.3 動作環境
- **ローカル環境**: ローカル環境でCLIコマンドが正常に動作することを確認
- **CI環境**: CI環境でもCLIコマンドが正常に動作することを確認（該当する場合）

## 6. 受け入れ基準

### 6.1 CLI用usecase層の作成
- [ ] `server/internal/usecase/cli`ディレクトリが作成されている
- [ ] `server/internal/usecase/cli/list_dm_users_usecase.go`が作成されている
- [ ] `ListDmUsersUsecase`構造体が定義されている
- [ ] `DmUserServiceInterface`を依存として注入している
- [ ] `ListDmUsers(ctx context.Context, limit, offset int) ([]*model.DmUser, error)`メソッドが実装されている
- [ ] service層の`ListDmUsers()`を呼び出して結果を返している

### 6.2 main.goの簡素化
- [ ] main.goがエントリーポイントと入出力制御のみを担当している
- [ ] usecase層を初期化して使用している
- [ ] `validateLimit()`関数が維持されている
- [ ] `printDmUsersTSV()`関数が維持されている
- [ ] 既存のバリデーションロジックが維持されている
- [ ] 既存の出力形式（TSV形式）が維持されている

### 6.3 依存関係の注入
- [ ] Repository層が適切に初期化されている
- [ ] Service層が適切に初期化されている
- [ ] Usecase層が適切に初期化されている
- [ ] usecase層の`ListDmUsers()`メソッドが呼び出されている

### 6.4 動作確認
- [ ] ローカル環境でCLIコマンドが正常に動作する
- [ ] 既存の引数（`-limit`）が正常に動作する
- [ ] 既存の出力形式（TSV形式）が維持されている
- [ ] 既存のエラーメッセージが維持されている
- [ ] 既存のテストが全て通過する
- [ ] CI環境でCLIコマンドが正常に動作する（該当する場合）

### 6.5 テスト
- [ ] `server/internal/usecase/cli/list_dm_users_usecase_test.go`が作成されている
- [ ] usecase層の単体テストが実装されている
- [ ] 既存のテスト（`server/cmd/list-dm-users/main_test.go`）が全て通過する

### 6.6 ドキュメントの更新
- [ ] `docs/Architecture.md`にCLIコマンドのレイヤー構造が追加されている
- [ ] `docs/Architecture.md`のCLIコマンドのアーキテクチャ図が更新されている
- [ ] `docs/Project-Structure.md`に`server/internal/usecase/cli`ディレクトリが追加されている
- [ ] `docs/Command-Line-Tool.md`のアーキテクチャ図が更新されている（usecase層を含む）
- [ ] `docs/Command-Line-Tool.md`のレイヤー構造の説明が更新されている
- [ ] `.kiro/steering/structure.md`に`server/internal/usecase/cli`ディレクトリが追加されている

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 新規作成が必要なファイル
- `server/internal/usecase/cli/list_dm_users_usecase.go`: CLI用usecase層の実装
- `server/internal/usecase/cli/list_dm_users_usecase_test.go`: CLI用usecase層のテスト

#### 修正が必要なファイル
- `server/cmd/list-dm-users/main.go`: usecase層を使用するように修正
- `docs/Architecture.md`: CLIコマンドのレイヤー構造を追加
- `docs/Project-Structure.md`: `server/internal/usecase/cli`ディレクトリを追加
- `docs/Command-Line-Tool.md`: CLIコマンドのアーキテクチャ図を更新
- `.kiro/steering/structure.md`: `server/internal/usecase/cli`ディレクトリを追加

#### 確認が必要なファイル
- 既存のservice層（`server/internal/service/dm_user_service.go`）: 変更不要だが、正常に動作することを確認
- 既存のrepository層（`server/internal/repository/dm_user_repository.go`）: 変更不要だが、正常に動作することを確認
- 既存のテスト（`server/cmd/list-dm-users/main_test.go`）: 全て通過することを確認

### 7.2 既存機能への影響
- **既存のservice層**: 影響なし（既存実装をそのまま使用）
- **既存のrepository層**: 影響なし（既存実装をそのまま使用）
- **既存のmodel層**: 影響なし（既存実装をそのまま使用）
- **既存のAPIサーバー**: 影響なし（APIサーバーは変更しない）
- **既存のCLIコマンドの動作**: 影響なし（動作は維持される）

## 8. 実装上の注意事項

### 8.1 usecase層の実装
- **インターフェースの使用**: 既存の`DmUserServiceInterface`を使用（`internal/usecase/dm_user_usecase.go`で定義済み）
- **依存関係の注入**: コンストラクタで`DmUserServiceInterface`を注入
- **エラーハンドリング**: service層から返されたエラーをそのまま返す（エラーのラップは不要）

### 8.2 main.goの修正
- **usecase層の初期化**: Repository → Service → Usecaseの順で初期化
- **既存のバリデーション**: `validateLimit()`関数はmain.goに保持
- **既存の出力**: `printDmUsersTSV()`関数はmain.goに保持
- **エラーハンドリング**: 既存のエラーハンドリングを維持

### 8.3 テストの実装
- **usecase層のテスト**: `DmUserServiceInterface`のモックを使用してテスト
- **既存のテスト**: `main_test.go`の既存テストが全て通過することを確認

### 8.4 ドキュメントの更新
- **アーキテクチャドキュメント**: CLIコマンドのレイヤー構造を明確に記載
- **プロジェクト構造ドキュメント**: 新規作成するディレクトリを反映
- **CLIツールドキュメント**: アーキテクチャ図を更新してusecase層を含める
- **ファイル組織ドキュメント**: 新規作成するディレクトリを反映
- **一貫性**: 全てのドキュメントで同じレイヤー構造を記載

## 9. 参考情報

### 9.1 関連ドキュメント
- `docs/Architecture.md`: アーキテクチャドキュメント
- `docs/Project-Structure.md`: プロジェクト構造ドキュメント
- `.kiro/steering/structure.md`: ファイル組織とコードパターン
- `.kiro/steering/tech.md`: 技術スタックとアーキテクチャ

### 9.2 既存実装の参考
- `server/internal/usecase/dm_user_usecase.go`: 既存のusecase層の実装パターン
- `server/cmd/list-dm-users/main.go`: 既存のCLI実装
- `server/internal/service/dm_user_service.go`: 既存のservice層の実装
- `server/internal/repository/dm_user_repository.go`: 既存のrepository層の実装

### 9.3 技術スタック
- **言語**: Go
- **アーキテクチャ**: レイヤードアーキテクチャ（usecase -> service -> repository -> model）
- **テスト**: `testing`（標準ライブラリ）、`github.com/stretchr/testify`（アサーション、モック）

### 9.4 レイヤー構造の比較

| 項目 | 現在（修正前） | 修正後 |
|------|---------------|--------|
| CLI層 | main.go（バリデーション、入出力、ビジネスロジック呼び出し） | main.go（エントリーポイント、バリデーション、入出力） |
| Usecase層 | なし | `usecase/cli/ListDmUsersUsecase` |
| Service層 | `service.DmUserService` | `service.DmUserService`（変更なし） |
| Repository層 | `repository.DmUserRepository` | `repository.DmUserRepository`（変更なし） |
| Model層 | `model.DmUser` | `model.DmUser`（変更なし） |

### 9.5 APIサーバーとの比較

| 項目 | APIサーバー | CLI（修正後） |
|------|------------|--------------|
| エントリーポイント | `server/cmd/server/main.go` | `server/cmd/list-dm-users/main.go` |
| バリデーション | API Layer（Handler） | CLI層（main.go） |
| Usecase層 | `usecase.DmUserUsecase` | `usecase/cli.ListDmUsersUsecase` |
| Service層 | `service.DmUserService` | `service.DmUserService` |
| Repository層 | `repository.DmUserRepository` | `repository.DmUserRepository` |
| Model層 | `model.DmUser` | `model.DmUser` |
