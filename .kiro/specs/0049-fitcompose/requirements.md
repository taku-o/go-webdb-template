# docker-composeファイル整理の要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #101
- **Issueタイトル**: docker-compose.admin、api、clientの整理
- **Feature名**: 0049-fitcompose
- **作成日**: 2026-01-10

### 1.2 目的
docker-compose.admin、api、clientの各環境用ファイルを整理し、production環境とstaging環境用のファイルを削除し、develop環境用のファイルをリネームしてシンプルな構成にする。

### 1.3 スコープ
- production環境用のdocker-composeファイルの削除
- staging環境用のdocker-composeファイルの削除
- develop環境用のdocker-composeファイルのリネーム（`.develop`を削除）
- ドキュメントなどで参照している箇所の修正

**本実装の範囲外**:
- 他のdocker-composeファイル（postgres、cloudbeaver、redis等）の変更
- docker-composeファイルの内容の変更
- 設定ファイルやコード内の変更（参照箇所の修正のみ）

## 2. 背景・現状分析

### 2.1 現在の状況
- **docker-composeファイル**: 以下の9ファイルが存在する
  - `docker-compose.admin.develop.yml`
  - `docker-compose.admin.staging.yml`
  - `docker-compose.admin.production.yml`
  - `docker-compose.api.develop.yml`
  - `docker-compose.api.staging.yml`
  - `docker-compose.api.production.yml`
  - `docker-compose.client.develop.yml`
  - `docker-compose.client.staging.yml`
  - `docker-compose.client.production.yml`
- **参照箇所**: `docs/Docker.md`などで多数の参照がある

### 2.2 課題点
1. **ファイル数の多さ**: 9ファイルが存在し、管理が煩雑
2. **不要なファイル**: production環境とstaging環境用のファイルが不要
3. **命名規則**: `.develop`という接尾辞が冗長
4. **参照箇所の不整合**: ファイル名変更後、参照箇所が古いままになる可能性

### 2.3 本実装による改善点
1. **ファイル数の削減**: 9ファイルから3ファイルに削減
2. **命名規則の簡素化**: `.develop`を削除してシンプルな名前に統一
3. **管理の簡素化**: 必要なファイルのみを残すことで管理が容易になる
4. **参照箇所の整合性**: すべての参照箇所を新しいファイル名に更新

## 3. 機能要件

### 3.1 docker-composeファイルの削除

#### 3.1.1 削除対象ファイル
以下の6ファイルを削除する：
- `docker-compose.admin.production.yml`
- `docker-compose.admin.staging.yml`
- `docker-compose.api.production.yml`
- `docker-compose.api.staging.yml`
- `docker-compose.client.production.yml`
- `docker-compose.client.staging.yml`

#### 3.1.2 削除の確認
- 削除前にファイルが存在することを確認
- 削除後にファイルが存在しないことを確認

### 3.2 docker-composeファイルのリネーム

#### 3.2.1 リネーム対象ファイル
以下の3ファイルをリネームする：
- `docker-compose.admin.develop.yml` → `docker-compose.admin.yml`
- `docker-compose.api.develop.yml` → `docker-compose.api.yml`
- `docker-compose.client.develop.yml` → `docker-compose.client.yml`

#### 3.2.2 リネームの確認
- リネーム前に元のファイルが存在することを確認
- リネーム後に新しいファイルが存在することを確認
- リネーム後に元のファイルが存在しないことを確認

### 3.3 参照箇所の修正

#### 3.3.1 参照箇所の検索
以下のパターンで参照箇所を検索する：
- `docker-compose.admin.develop.yml`
- `docker-compose.admin.staging.yml`
- `docker-compose.admin.production.yml`
- `docker-compose.api.develop.yml`
- `docker-compose.api.staging.yml`
- `docker-compose.api.production.yml`
- `docker-compose.client.develop.yml`
- `docker-compose.client.staging.yml`
- `docker-compose.client.production.yml`

#### 3.3.2 参照箇所の修正内容
- **develop環境用ファイルの参照**: `.develop`を削除
  - `docker-compose.admin.develop.yml` → `docker-compose.admin.yml`
  - `docker-compose.api.develop.yml` → `docker-compose.api.yml`
  - `docker-compose.client.develop.yml` → `docker-compose.client.yml`
- **staging環境用ファイルの参照**: 該当箇所を削除または修正
- **production環境用ファイルの参照**: 該当箇所を削除または修正

#### 3.3.3 修正対象の確認
- `docs/Docker.md`内の参照箇所を修正
- その他のドキュメントファイル内の参照箇所を修正
- README.mdやその他のMarkdownファイル内の参照箇所を修正
- スクリプトファイル内の参照箇所を修正（該当する場合）

## 4. 非機能要件

### 4.1 ファイル管理
- **一貫性**: すべてのdocker-composeファイルで命名規則が統一されていること
- **明確性**: ファイル名から用途が明確であること

### 4.2 ドキュメントの整合性
- **正確性**: ドキュメントの内容が実際のファイル構成と一致していること
- **完全性**: すべての参照箇所が修正されていること

### 4.3 保守性
- **更新容易性**: 将来の変更に対応しやすい構造であること
- **検索性**: 参照箇所が残っていないことを確認しやすいこと

## 5. 制約事項

### 5.1 対象外のファイル
以下のファイルは変更対象外とする：
- `docker-compose.postgres.yml`
- `docker-compose.cloudbeaver.yml`
- `docker-compose.apache-superset.yml`
- `docker-compose.redis-cluster.yml`
- `docker-compose.redis.yml`
- `docker-compose.redis-insight.yml`
- `docker-compose.mailpit.yml`
- `docker-compose.metabase.yml`
- その他のdocker-composeファイル

### 5.2 ファイル内容の変更
- docker-composeファイルの内容（サービス定義など）は変更しない
- ファイル名の変更のみを行う

### 5.3 既存の機能への影響
- 既存の機能が動作することを確認する
- ファイル名変更による影響がないことを確認する

## 6. 受け入れ基準

### 6.1 docker-composeファイルの削除
- [ ] `docker-compose.admin.production.yml`が削除されている
- [ ] `docker-compose.admin.staging.yml`が削除されている
- [ ] `docker-compose.api.production.yml`が削除されている
- [ ] `docker-compose.api.staging.yml`が削除されている
- [ ] `docker-compose.client.production.yml`が削除されている
- [ ] `docker-compose.client.staging.yml`が削除されている

### 6.2 docker-composeファイルのリネーム
- [ ] `docker-compose.admin.develop.yml`が`docker-compose.admin.yml`にリネームされている
- [ ] `docker-compose.api.develop.yml`が`docker-compose.api.yml`にリネームされている
- [ ] `docker-compose.client.develop.yml`が`docker-compose.client.yml`にリネームされている
- [ ] リネーム後のファイルが存在する
- [ ] リネーム前のファイルが存在しない

### 6.3 参照箇所の修正
- [ ] `docs/Docker.md`内のすべての参照箇所が修正されている
- [ ] その他のドキュメントファイル内の参照箇所が修正されている
- [ ] README.mdやその他のMarkdownファイル内の参照箇所が修正されている
- [ ] スクリプトファイル内の参照箇所が修正されている（該当する場合）
- [ ] staging環境用ファイルの参照が削除または修正されている
- [ ] production環境用ファイルの参照が削除または修正されている

### 6.4 動作確認
- [ ] リネーム後のファイルでdocker-composeコマンドが正常に動作する
- [ ] ドキュメントを読んで、新しいファイル名が正しく記載されている
- [ ] ドキュメントに矛盾や不整合がない

## 7. 影響範囲

### 7.1 変更が必要なファイル

#### 削除するファイル
- `docker-compose.admin.production.yml`
- `docker-compose.admin.staging.yml`
- `docker-compose.api.production.yml`
- `docker-compose.api.staging.yml`
- `docker-compose.client.production.yml`
- `docker-compose.client.staging.yml`

#### リネームするファイル
- `docker-compose.admin.develop.yml` → `docker-compose.admin.yml`
- `docker-compose.api.develop.yml` → `docker-compose.api.yml`
- `docker-compose.client.develop.yml` → `docker-compose.client.yml`

#### 修正が必要なファイル
- `docs/Docker.md`（多数の参照箇所あり）
- その他のドキュメントファイル（参照箇所がある場合）
- README.md（参照箇所がある場合）
- スクリプトファイル（参照箇所がある場合）

### 7.2 既存ファイルの扱い
- 各参照箇所を確認し、新しいファイル名に更新
- staging環境とproduction環境に関する記述を削除または修正

## 8. 実装上の注意事項

### 8.1 ファイル削除の注意事項
- **削除前の確認**: 削除対象ファイルが実際に存在することを確認
- **削除後の確認**: 削除後にファイルが存在しないことを確認
- **バックアップ**: 必要に応じて削除前にバックアップを取る（履歴はGitで管理されているため必須ではない）

### 8.2 ファイルリネームの注意事項
- **リネーム前の確認**: リネーム対象ファイルが実際に存在することを確認
- **リネーム後の確認**: リネーム後に新しいファイルが存在することを確認
- **内容の確認**: リネーム後もファイル内容が変更されていないことを確認

### 8.3 参照箇所修正の注意事項
- **網羅的な検索**: すべてのファイルで参照箇所を検索する
- **正確な置換**: ファイル名を正確に置換する
- **文脈の確認**: 参照箇所の文脈を確認し、適切に修正する
- **staging/production環境の記述**: staging環境やproduction環境に関する記述を削除または修正する際は、文脈を確認して適切に対応する

### 8.4 動作確認の注意事項
- **docker-composeコマンドの確認**: リネーム後のファイルでdocker-composeコマンドが正常に動作することを確認
- **ドキュメントの確認**: ドキュメントを読んで、新しいファイル名が正しく記載されていることを確認
- **整合性の確認**: ドキュメントに矛盾や不整合がないことを確認

## 9. 参考情報

### 9.1 関連Issue
- GitHub Issue #101: docker-compose.admin、api、clientの整理

### 9.2 既存ファイル
- `docker-compose.admin.develop.yml`: Adminサーバー用のdocker-composeファイル
- `docker-compose.api.develop.yml`: APIサーバー用のdocker-composeファイル
- `docker-compose.client.develop.yml`: クライアント用のdocker-composeファイル
- `docs/Docker.md`: Docker関連のドキュメント

### 9.3 技術スタック
- **コンテナ管理**: Docker Compose
- **ドキュメント形式**: Markdown
