# docker-composeファイル整理の設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、docker-compose.admin、api、clientの各環境用ファイルを整理するための詳細設計を定義する。production環境とstaging環境用のファイルを削除し、develop環境用のファイルをリネームしてシンプルな構成にする。また、ドキュメントなどで参照している箇所を修正する。

### 1.2 設計の範囲
- production環境用のdocker-composeファイル（6ファイル）の削除
- staging環境用のdocker-composeファイル（6ファイル）の削除
- develop環境用のdocker-composeファイル（3ファイル）のリネーム（`.develop`を削除）
- ドキュメントなどで参照している箇所の修正

### 1.3 設計方針
- **段階的な作業**: 削除 → リネーム → 参照箇所修正の順序で作業を進める
- **網羅的な検索**: すべてのファイルで参照箇所を検索
- **整合性の確保**: すべての参照箇所を新しいファイル名に更新
- **既存機能の維持**: docker-composeファイルの内容は変更しない

## 2. 作業手順設計

### 2.1 全体フロー

```
┌─────────────────────────────────────────────────────────────┐
│              docker-composeファイル整理フロー                   │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │ 1. 削除対象ファイルの確認          │
        │    - 6ファイルの存在確認          │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │ 2. production/stagingファイル削除│
        │    - 6ファイルを削除              │
        │    - 削除後の確認                 │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │ 3. リネーム対象ファイルの確認      │
        │    - 3ファイルの存在確認          │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │ 4. developファイルのリネーム      │
        │    - 3ファイルをリネーム          │
        │    - リネーム後の確認             │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │ 5. 参照箇所の検索                 │
        │    - 全ファイルで検索            │
        │    - 参照箇所のリストアップ        │
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │ 6. 参照箇所の修正                 │
        │    - develop環境用の修正          │
        │    - staging/production環境の削除│
        └─────────────────────────────────┘
                          │
                          ▼
        ┌─────────────────────────────────┐
        │ 7. 動作確認                       │
        │    - docker-composeコマンド確認    │
        │    - ドキュメントの整合性確認      │
        └─────────────────────────────────┘
```

### 2.2 削除対象ファイル

#### 2.2.1 削除するファイル一覧
以下の6ファイルを削除する：
- `docker-compose.admin.production.yml`
- `docker-compose.admin.staging.yml`
- `docker-compose.api.production.yml`
- `docker-compose.api.staging.yml`
- `docker-compose.client.production.yml`
- `docker-compose.client.staging.yml`

#### 2.2.2 削除の確認方法
```bash
# 削除前の確認
ls -la docker-compose.admin.production.yml
ls -la docker-compose.admin.staging.yml
ls -la docker-compose.api.production.yml
ls -la docker-compose.api.staging.yml
ls -la docker-compose.client.production.yml
ls -la docker-compose.client.staging.yml

# 削除実行
rm docker-compose.admin.production.yml
rm docker-compose.admin.staging.yml
rm docker-compose.api.production.yml
rm docker-compose.api.staging.yml
rm docker-compose.client.production.yml
rm docker-compose.client.staging.yml

# 削除後の確認
ls -la docker-compose.admin.production.yml  # ファイルが存在しないことを確認
ls -la docker-compose.admin.staging.yml     # ファイルが存在しないことを確認
ls -la docker-compose.api.production.yml    # ファイルが存在しないことを確認
ls -la docker-compose.api.staging.yml       # ファイルが存在しないことを確認
ls -la docker-compose.client.production.yml # ファイルが存在しないことを確認
ls -la docker-compose.client.staging.yml    # ファイルが存在しないことを確認
```

### 2.3 リネーム対象ファイル

#### 2.3.1 リネームするファイル一覧
以下の3ファイルをリネームする：
- `docker-compose.admin.develop.yml` → `docker-compose.admin.yml`
- `docker-compose.api.develop.yml` → `docker-compose.api.yml`
- `docker-compose.client.develop.yml` → `docker-compose.client.yml`

#### 2.3.2 リネームの確認方法
```bash
# リネーム前の確認
ls -la docker-compose.admin.develop.yml
ls -la docker-compose.api.develop.yml
ls -la docker-compose.client.develop.yml

# リネーム実行
mv docker-compose.admin.develop.yml docker-compose.admin.yml
mv docker-compose.api.develop.yml docker-compose.api.yml
mv docker-compose.client.develop.yml docker-compose.client.yml

# リネーム後の確認
ls -la docker-compose.admin.yml      # 新しいファイルが存在することを確認
ls -la docker-compose.api.yml        # 新しいファイルが存在することを確認
ls -la docker-compose.client.yml     # 新しいファイルが存在することを確認
ls -la docker-compose.admin.develop.yml # 元のファイルが存在しないことを確認
ls -la docker-compose.api.develop.yml   # 元のファイルが存在しないことを確認
ls -la docker-compose.client.develop.yml # 元のファイルが存在しないことを確認
```

### 2.4 参照箇所の検索

#### 2.4.1 検索パターン
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

#### 2.4.2 検索コマンド
```bash
# 全ファイルで検索
grep -r "docker-compose\.\(admin\|api\|client\)\.\(develop\|staging\|production\)" .

# 特定のファイルで検索
grep "docker-compose\.\(admin\|api\|client\)\.\(develop\|staging\|production\)" docs/Docker.md
grep "docker-compose\.\(admin\|api\|client\)\.\(develop\|staging\|production\)" README.md
```

#### 2.4.3 検索対象ファイル
以下のファイルで参照箇所を検索する：
- `docs/Docker.md`（多数の参照箇所あり）
- `README.md`（参照箇所がある場合）
- その他のドキュメントファイル（参照箇所がある場合）
- スクリプトファイル（参照箇所がある場合）

## 3. 参照箇所の修正設計

### 3.1 docs/Docker.mdの修正

#### 3.1.1 ファイル一覧テーブルの修正
**修正前**:
```markdown
| ファイル | サービス | 環境 |
|---------|---------|------|
| `docker-compose.api.develop.yml` | APIサーバー | develop |
| `docker-compose.api.staging.yml` | APIサーバー | staging |
| `docker-compose.api.production.yml` | APIサーバー | production |
| `docker-compose.admin.develop.yml` | Adminサーバー | develop |
| `docker-compose.admin.staging.yml` | Adminサーバー | staging |
| `docker-compose.admin.production.yml` | Adminサーバー | production |
| `docker-compose.client.develop.yml` | クライアント | develop |
| `docker-compose.client.staging.yml` | クライアント | staging |
| `docker-compose.client.production.yml` | クライアント | production |
```

**修正後**:
```markdown
| ファイル | サービス | 環境 |
|---------|---------|------|
| `docker-compose.api.yml` | APIサーバー | develop |
| `docker-compose.admin.yml` | Adminサーバー | develop |
| `docker-compose.client.yml` | クライアント | develop |
```

#### 3.1.2 開発環境（develop）セクションの修正
**修正前**:
```bash
# ビルド
docker-compose -f docker-compose.api.develop.yml build
docker-compose -f docker-compose.admin.develop.yml build
docker-compose -f docker-compose.client.develop.yml build

# 起動
docker-compose -f docker-compose.api.develop.yml up -d
docker-compose -f docker-compose.admin.develop.yml up -d
docker-compose -f docker-compose.client.develop.yml up -d

# 停止
docker-compose -f docker-compose.api.develop.yml down
docker-compose -f docker-compose.admin.develop.yml down
docker-compose -f docker-compose.client.develop.yml down

# ログ確認
docker-compose -f docker-compose.api.develop.yml logs -f
docker-compose -f docker-compose.admin.develop.yml logs -f
docker-compose -f docker-compose.client.develop.yml logs -f
```

**修正後**:
```bash
# ビルド
docker-compose -f docker-compose.api.yml build
docker-compose -f docker-compose.admin.yml build
docker-compose -f docker-compose.client.yml build

# 起動
docker-compose -f docker-compose.api.yml up -d
docker-compose -f docker-compose.admin.yml up -d
docker-compose -f docker-compose.client.yml up -d

# 停止
docker-compose -f docker-compose.api.yml down
docker-compose -f docker-compose.admin.yml down
docker-compose -f docker-compose.client.yml down

# ログ確認
docker-compose -f docker-compose.api.yml logs -f
docker-compose -f docker-compose.admin.yml logs -f
docker-compose -f docker-compose.client.yml logs -f
```

#### 3.1.3 ステージング環境（staging）セクションの削除
**削除対象**: ステージング環境（staging）セクション全体を削除する。

```markdown
### ステージング環境（staging）

```bash
# ビルド
docker-compose -f docker-compose.api.staging.yml build
docker-compose -f docker-compose.admin.staging.yml build
docker-compose -f docker-compose.client.staging.yml build

# 起動
docker-compose -f docker-compose.api.staging.yml up -d
docker-compose -f docker-compose.admin.staging.yml up -d
docker-compose -f docker-compose.client.staging.yml up -d

# 停止
docker-compose -f docker-compose.api.staging.yml down
docker-compose -f docker-compose.admin.staging.yml down
docker-compose -f docker-compose.client.staging.yml down
```
```

#### 3.1.4 本番環境（production）セクションの削除
**削除対象**: 本番環境（production）セクション全体を削除する。

```markdown
### 本番環境（production）

```bash
# ビルド
docker-compose -f docker-compose.api.production.yml build
docker-compose -f docker-compose.admin.production.yml build
docker-compose -f docker-compose.client.production.yml build

# 起動
docker-compose -f docker-compose.api.production.yml up -d
docker-compose -f docker-compose.admin.production.yml up -d
docker-compose -f docker-compose.client.production.yml up -d

# 停止
docker-compose -f docker-compose.api.production.yml down
docker-compose -f docker-compose.admin.production.yml down
docker-compose -f docker-compose.client.production.yml down
```
```

#### 3.1.5 対応環境テーブルの修正
**修正前**:
```markdown
| 環境 | 用途 | データベース |
|------|------|-------------|
| develop | 開発環境 | PostgreSQL/MySQL |
| staging | ステージング環境 | PostgreSQL/MySQL |
| production | 本番環境 | PostgreSQL/MySQL |
```

**修正後**:
```markdown
| 環境 | 用途 | データベース |
|------|------|-------------|
| develop | 開発環境 | PostgreSQL/MySQL |
```

#### 3.1.6 イメージの再ビルドセクションの修正
**修正前**:
```bash
docker-compose -f docker-compose.api.develop.yml build --no-cache
```

**修正後**:
```bash
docker-compose -f docker-compose.api.yml build --no-cache
```

#### 3.1.7 コンテナレジストリへのプッシュセクションの修正
**修正前**:
```bash
# APIサーバー
docker-compose -f docker-compose.api.production.yml build

# Adminサーバー
docker-compose -f docker-compose.admin.production.yml build

# クライアントサーバー
docker-compose -f docker-compose.client.production.yml build
```

**修正後**:
このセクションは、production環境用のファイルが削除されるため、削除または大幅に修正する必要がある。ただし、本番環境へのデプロイに関する記述は残す必要がある場合は、別の方法（直接イメージをビルドする方法など）を記載する。

### 3.2 README.mdの修正

#### 3.2.1 参照箇所の確認
README.mdに以下の参照箇所がある：
- APIサーバーのビルドと起動コマンド
- Adminサーバーのビルドと起動コマンド
- クライアントサーバーのビルドと起動コマンド

#### 3.2.2 修正内容
**修正前**:
```bash
# APIサーバーのビルドと起動
docker-compose -f docker-compose.api.develop.yml build
docker-compose -f docker-compose.api.develop.yml up -d

# Adminサーバーのビルドと起動
docker-compose -f docker-compose.admin.develop.yml build
docker-compose -f docker-compose.admin.develop.yml up -d

# クライアントサーバーのビルドと起動
docker-compose -f docker-compose.client.develop.yml build
docker-compose -f docker-compose.client.develop.yml up -d
```

**修正後**:
```bash
# APIサーバーのビルドと起動
docker-compose -f docker-compose.api.yml build
docker-compose -f docker-compose.api.yml up -d

# Adminサーバーのビルドと起動
docker-compose -f docker-compose.admin.yml build
docker-compose -f docker-compose.admin.yml up -d

# クライアントサーバーのビルドと起動
docker-compose -f docker-compose.client.yml build
docker-compose -f docker-compose.client.yml up -d
```

### 3.3 その他のファイルの修正

#### 3.3.1 スクリプトファイル
スクリプトファイルに参照箇所がある場合は、同様に修正する。

#### 3.3.2 その他のドキュメントファイル
その他のドキュメントファイルに参照箇所がある場合は、同様に修正する。

## 4. 修正パターン設計

### 4.1 develop環境用ファイルの参照修正

#### 4.1.1 パターン1: ファイル名の直接参照
**修正前**: `docker-compose.admin.develop.yml`
**修正後**: `docker-compose.admin.yml`

**修正前**: `docker-compose.api.develop.yml`
**修正後**: `docker-compose.api.yml`

**修正前**: `docker-compose.client.develop.yml`
**修正後**: `docker-compose.client.yml`

#### 4.1.2 パターン2: docker-composeコマンド内の参照
**修正前**: `docker-compose -f docker-compose.api.develop.yml build`
**修正後**: `docker-compose -f docker-compose.api.yml build`

**修正前**: `docker-compose -f docker-compose.admin.develop.yml up -d`
**修正後**: `docker-compose -f docker-compose.admin.yml up -d`

**修正前**: `docker-compose -f docker-compose.client.develop.yml down`
**修正後**: `docker-compose -f docker-compose.client.yml down`

### 4.2 staging/production環境用ファイルの参照削除

#### 4.2.1 パターン1: セクション全体の削除
staging環境やproduction環境に関するセクション全体を削除する。

#### 4.2.2 パターン2: テーブル行の削除
テーブルからstaging環境やproduction環境の行を削除する。

#### 4.2.3 パターン3: コマンド例の削除
staging環境やproduction環境のコマンド例を削除する。

## 5. 整合性確認設計

### 5.1 ファイル削除の確認
- [ ] `docker-compose.admin.production.yml`が削除されている
- [ ] `docker-compose.admin.staging.yml`が削除されている
- [ ] `docker-compose.api.production.yml`が削除されている
- [ ] `docker-compose.api.staging.yml`が削除されている
- [ ] `docker-compose.client.production.yml`が削除されている
- [ ] `docker-compose.client.staging.yml`が削除されている

### 5.2 ファイルリネームの確認
- [ ] `docker-compose.admin.develop.yml`が`docker-compose.admin.yml`にリネームされている
- [ ] `docker-compose.api.develop.yml`が`docker-compose.api.yml`にリネームされている
- [ ] `docker-compose.client.develop.yml`が`docker-compose.client.yml`にリネームされている
- [ ] リネーム後のファイルが存在する
- [ ] リネーム前のファイルが存在しない

### 5.3 参照箇所修正の確認
- [ ] `docs/Docker.md`内のすべての参照箇所が修正されている
- [ ] `README.md`内の参照箇所が修正されている（該当する場合）
- [ ] その他のドキュメントファイル内の参照箇所が修正されている（該当する場合）
- [ ] スクリプトファイル内の参照箇所が修正されている（該当する場合）
- [ ] staging環境用ファイルの参照が削除または修正されている
- [ ] production環境用ファイルの参照が削除または修正されている

### 5.4 動作確認
- [ ] リネーム後のファイルでdocker-composeコマンドが正常に動作する
  ```bash
  docker-compose -f docker-compose.api.yml config
  docker-compose -f docker-compose.admin.yml config
  docker-compose -f docker-compose.client.yml config
  ```
- [ ] ドキュメントを読んで、新しいファイル名が正しく記載されている
- [ ] ドキュメントに矛盾や不整合がない

## 6. 実装上の注意事項

### 6.1 ファイル削除の注意事項
- **削除前の確認**: 削除対象ファイルが実際に存在することを確認
- **削除後の確認**: 削除後にファイルが存在しないことを確認
- **Git履歴**: ファイルはGitで管理されているため、履歴は残る

### 6.2 ファイルリネームの注意事項
- **リネーム前の確認**: リネーム対象ファイルが実際に存在することを確認
- **リネーム後の確認**: リネーム後に新しいファイルが存在することを確認
- **内容の確認**: リネーム後もファイル内容が変更されていないことを確認
- **Gitでの追跡**: `git mv`コマンドを使用することで、Git履歴を保持できる

### 6.3 参照箇所修正の注意事項
- **網羅的な検索**: すべてのファイルで参照箇所を検索する
- **正確な置換**: ファイル名を正確に置換する
- **文脈の確認**: 参照箇所の文脈を確認し、適切に修正する
- **staging/production環境の記述**: staging環境やproduction環境に関する記述を削除または修正する際は、文脈を確認して適切に対応する
- **セクション全体の削除**: セクション全体を削除する場合は、前後の文脈が自然になるように調整する

### 6.4 動作確認の注意事項
- **docker-composeコマンドの確認**: リネーム後のファイルでdocker-composeコマンドが正常に動作することを確認
- **ドキュメントの確認**: ドキュメントを読んで、新しいファイル名が正しく記載されていることを確認
- **整合性の確認**: ドキュメントに矛盾や不整合がないことを確認

## 7. 参考情報

### 7.1 関連Issue
- GitHub Issue #101: docker-compose.admin、api、clientの整理

### 7.2 既存ファイル
- `docker-compose.admin.develop.yml`: Adminサーバー用のdocker-composeファイル（リネーム対象）
- `docker-compose.api.develop.yml`: APIサーバー用のdocker-composeファイル（リネーム対象）
- `docker-compose.client.develop.yml`: クライアント用のdocker-composeファイル（リネーム対象）
- `docs/Docker.md`: Docker関連のドキュメント（多数の参照箇所あり）

### 7.3 技術スタック
- **コンテナ管理**: Docker Compose
- **ドキュメント形式**: Markdown
