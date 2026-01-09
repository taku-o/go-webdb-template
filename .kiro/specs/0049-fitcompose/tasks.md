# docker-composeファイル整理の実装タスク一覧

## 概要
docker-compose.admin、api、clientの各環境用ファイルを整理するためのタスク一覧。要件定義書と設計書に基づき、実装可能な粒度でタスクを分解。

## 実装フェーズ

### Phase 1: ファイル削除

#### タスク 1.1: 削除対象ファイルの確認
**目的**: 削除対象となるdocker-composeファイルが存在することを確認する。

**作業内容**:
- 以下の6ファイルが存在することを確認
  - `docker-compose.admin.production.yml`
  - `docker-compose.admin.staging.yml`
  - `docker-compose.api.production.yml`
  - `docker-compose.api.staging.yml`
  - `docker-compose.client.production.yml`
  - `docker-compose.client.staging.yml`

**確認コマンド**:
```bash
ls -la docker-compose.admin.production.yml
ls -la docker-compose.admin.staging.yml
ls -la docker-compose.api.production.yml
ls -la docker-compose.api.staging.yml
ls -la docker-compose.client.production.yml
ls -la docker-compose.client.staging.yml
```

**受け入れ基準**:
- すべての削除対象ファイルが存在することを確認できた
- ファイルの存在確認が記録されている

- _Requirements: 3.1.1, 6.1_
- _Design: 2.2.1, 2.2.2_

---

#### タスク 1.2: production/staging環境用ファイルの削除
**目的**: production環境とstaging環境用のdocker-composeファイルを削除する。

**作業内容**:
- 以下の6ファイルを削除
  - `docker-compose.admin.production.yml`
  - `docker-compose.admin.staging.yml`
  - `docker-compose.api.production.yml`
  - `docker-compose.api.staging.yml`
  - `docker-compose.client.production.yml`
  - `docker-compose.client.staging.yml`

**削除コマンド**:
```bash
rm docker-compose.admin.production.yml
rm docker-compose.admin.staging.yml
rm docker-compose.api.production.yml
rm docker-compose.api.staging.yml
rm docker-compose.client.production.yml
rm docker-compose.client.staging.yml
```

**受け入れ基準**:
- すべての削除対象ファイルが削除されている
- 削除後にファイルが存在しないことを確認できた

- _Requirements: 3.1.1, 3.1.2, 6.1_
- _Design: 2.2.1, 2.2.2_

---

#### タスク 1.3: 削除後の確認
**目的**: 削除対象ファイルが正しく削除されたことを確認する。

**作業内容**:
- 削除対象ファイルが存在しないことを確認
- 削除が正常に完了したことを記録

**確認コマンド**:
```bash
ls -la docker-compose.admin.production.yml  # ファイルが存在しないことを確認
ls -la docker-compose.admin.staging.yml      # ファイルが存在しないことを確認
ls -la docker-compose.api.production.yml    # ファイルが存在しないことを確認
ls -la docker-compose.api.staging.yml       # ファイルが存在しないことを確認
ls -la docker-compose.client.production.yml # ファイルが存在しないことを確認
ls -la docker-compose.client.staging.yml    # ファイルが存在しないことを確認
```

**受け入れ基準**:
- すべての削除対象ファイルが存在しないことを確認できた
- 削除が正常に完了したことが記録されている

- _Requirements: 3.1.2, 6.1_
- _Design: 2.2.2_

---

### Phase 2: ファイルリネーム

#### タスク 2.1: リネーム対象ファイルの確認
**目的**: リネーム対象となるdocker-composeファイルが存在することを確認する。

**作業内容**:
- 以下の3ファイルが存在することを確認
  - `docker-compose.admin.develop.yml`
  - `docker-compose.api.develop.yml`
  - `docker-compose.client.develop.yml`

**確認コマンド**:
```bash
ls -la docker-compose.admin.develop.yml
ls -la docker-compose.api.develop.yml
ls -la docker-compose.client.develop.yml
```

**受け入れ基準**:
- すべてのリネーム対象ファイルが存在することを確認できた
- ファイルの存在確認が記録されている

- _Requirements: 3.2.1, 6.2_
- _Design: 2.3.1, 2.3.2_

---

#### タスク 2.2: develop環境用ファイルのリネーム
**目的**: develop環境用のdocker-composeファイルをリネームする。

**作業内容**:
- 以下の3ファイルをリネーム
  - `docker-compose.admin.develop.yml` → `docker-compose.admin.yml`
  - `docker-compose.api.develop.yml` → `docker-compose.api.yml`
  - `docker-compose.client.develop.yml` → `docker-compose.client.yml`

**リネームコマンド**:
```bash
git mv docker-compose.admin.develop.yml docker-compose.admin.yml
git mv docker-compose.api.develop.yml docker-compose.api.yml
git mv docker-compose.client.develop.yml docker-compose.client.yml
```

**注意事項**:
- `git mv`コマンドを使用することで、Git履歴を保持できる

**受け入れ基準**:
- すべてのリネーム対象ファイルがリネームされている
- リネーム後のファイルが存在することを確認できた
- リネーム前のファイルが存在しないことを確認できた

- _Requirements: 3.2.1, 3.2.2, 6.2_
- _Design: 2.3.1, 2.3.2_

---

#### タスク 2.3: リネーム後の確認
**目的**: リネーム対象ファイルが正しくリネームされたことを確認する。

**作業内容**:
- リネーム後のファイルが存在することを確認
- リネーム前のファイルが存在しないことを確認
- リネーム後もファイル内容が変更されていないことを確認

**確認コマンド**:
```bash
# リネーム後のファイルが存在することを確認
ls -la docker-compose.admin.yml
ls -la docker-compose.api.yml
ls -la docker-compose.client.yml

# リネーム前のファイルが存在しないことを確認
ls -la docker-compose.admin.develop.yml # ファイルが存在しないことを確認
ls -la docker-compose.api.develop.yml    # ファイルが存在しないことを確認
ls -la docker-compose.client.develop.yml # ファイルが存在しないことを確認
```

**受け入れ基準**:
- リネーム後のファイルが存在することを確認できた
- リネーム前のファイルが存在しないことを確認できた
- リネーム後もファイル内容が変更されていないことを確認できた

- _Requirements: 3.2.2, 6.2_
- _Design: 2.3.2_

---

### Phase 3: 参照箇所の検索

#### タスク 3.1: 参照箇所の網羅的検索
**目的**: すべてのファイルからdocker-composeファイルの参照箇所を検索する。

**作業内容**:
- 以下のパターンで参照箇所を検索
  - `docker-compose.admin.develop.yml`
  - `docker-compose.admin.staging.yml`
  - `docker-compose.admin.production.yml`
  - `docker-compose.api.develop.yml`
  - `docker-compose.api.staging.yml`
  - `docker-compose.api.production.yml`
  - `docker-compose.client.develop.yml`
  - `docker-compose.client.staging.yml`
  - `docker-compose.client.production.yml`

**検索コマンド**:
```bash
# 全ファイルで検索
grep -r "docker-compose\.\(admin\|api\|client\)\.\(develop\|staging\|production\)" .

# 特定のファイルで検索
grep "docker-compose\.\(admin\|api\|client\)\.\(develop\|staging\|production\)" docs/Docker.md
grep "docker-compose\.\(admin\|api\|client\)\.\(develop\|staging\|production\)" README.md
```

**検索対象ファイル**:
- `docs/Docker.md`（多数の参照箇所あり）
- `README.md`（参照箇所がある場合）
- その他のドキュメントファイル（参照箇所がある場合）
- スクリプトファイル（参照箇所がある場合）

**受け入れ基準**:
- すべてのファイルで参照箇所が検索されている
- 検索結果がファイルごとに記録されている
- 参照箇所のリストが作成されている

- _Requirements: 3.3.1, 6.3_
- _Design: 2.4_

---

#### タスク 3.2: 検索結果の確認と分類
**目的**: 検索結果を確認し、修正対象の参照箇所を分類する。

**作業内容**:
- 検索結果を確認し、参照箇所を特定
- 参照箇所の種類を分類:
  - develop環境用ファイルの参照（`.develop`を削除して修正）
  - staging環境用ファイルの参照（削除または修正）
  - production環境用ファイルの参照（削除または修正）
- 各参照箇所の文脈を確認
- 修正対象の参照箇所をリストアップ

**受け入れ基準**:
- すべての参照箇所が特定されている
- 参照箇所の種類が分類されている
- 各参照箇所の文脈が確認されている
- 修正対象の参照箇所がリストアップされている

- _Requirements: 3.3.2, 6.3_
- _Design: 4.1, 4.2_

---

### Phase 4: 参照箇所の修正

#### タスク 4.1: docs/Docker.mdの修正 - ファイル一覧テーブル
**目的**: `docs/Docker.md`のファイル一覧テーブルを修正する。

**作業内容**:
- ファイル一覧テーブルからstaging環境とproduction環境の行を削除
- develop環境用ファイルの参照を`.develop`を削除して修正

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

**受け入れ基準**:
- ファイル一覧テーブルが正しく修正されている
- staging環境とproduction環境の行が削除されている
- develop環境用ファイルの参照が`.develop`を削除して修正されている

- _Requirements: 3.3.2, 6.3_
- _Design: 3.1.1_

---

#### タスク 4.2: docs/Docker.mdの修正 - 開発環境（develop）セクション
**目的**: `docs/Docker.md`の開発環境（develop）セクションのコマンド例を修正する。

**作業内容**:
- 開発環境（develop）セクションのすべてのコマンド例で、`.develop`を削除して修正

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

**受け入れ基準**:
- 開発環境（develop）セクションのすべてのコマンド例が修正されている
- `.develop`が削除されている

- _Requirements: 3.3.2, 6.3_
- _Design: 3.1.2_

---

#### タスク 4.3: docs/Docker.mdの修正 - ステージング環境（staging）セクションの削除
**目的**: `docs/Docker.md`のステージング環境（staging）セクションを削除する。

**作業内容**:
- ステージング環境（staging）セクション全体を削除
- 前後の文脈が自然になるように調整

**削除対象**:
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

**受け入れ基準**:
- ステージング環境（staging）セクションが削除されている
- 前後の文脈が自然になっている

- _Requirements: 3.3.2, 6.3_
- _Design: 3.1.3_

---

#### タスク 4.4: docs/Docker.mdの修正 - 本番環境（production）セクションの削除
**目的**: `docs/Docker.md`の本番環境（production）セクションを削除する。

**作業内容**:
- 本番環境（production）セクション全体を削除
- 前後の文脈が自然になるように調整

**削除対象**:
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

**受け入れ基準**:
- 本番環境（production）セクションが削除されている
- 前後の文脈が自然になっている

- _Requirements: 3.3.2, 6.3_
- _Design: 3.1.4_

---

#### タスク 4.5: docs/Docker.mdの修正 - 対応環境テーブル
**目的**: `docs/Docker.md`の対応環境テーブルを修正する。

**作業内容**:
- 対応環境テーブルからstaging環境とproduction環境の行を削除

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

**受け入れ基準**:
- 対応環境テーブルが正しく修正されている
- staging環境とproduction環境の行が削除されている

- _Requirements: 3.3.2, 6.3_
- _Design: 3.1.5_

---

#### タスク 4.6: docs/Docker.mdの修正 - イメージの再ビルドセクション
**目的**: `docs/Docker.md`のイメージの再ビルドセクションを修正する。

**作業内容**:
- イメージの再ビルドセクションのコマンド例で、`.develop`を削除して修正

**修正前**:
```bash
docker-compose -f docker-compose.api.develop.yml build --no-cache
```

**修正後**:
```bash
docker-compose -f docker-compose.api.yml build --no-cache
```

**受け入れ基準**:
- イメージの再ビルドセクションのコマンド例が修正されている
- `.develop`が削除されている

- _Requirements: 3.3.2, 6.3_
- _Design: 3.1.6_

---

#### タスク 4.7: docs/Docker.mdの修正 - コンテナレジストリへのプッシュセクション
**目的**: `docs/Docker.md`のコンテナレジストリへのプッシュセクションを修正する。

**作業内容**:
- コンテナレジストリへのプッシュセクションのコマンド例で、production環境用ファイルの参照を削除または修正
- 本番環境へのデプロイに関する記述は残す必要がある場合は、別の方法（直接イメージをビルドする方法など）を記載

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

**受け入れ基準**:
- コンテナレジストリへのプッシュセクションが修正されている
- production環境用ファイルの参照が削除または修正されている
- 本番環境へのデプロイに関する記述が適切に修正されている（該当する場合）

- _Requirements: 3.3.2, 6.3_
- _Design: 3.1.7_

---

#### タスク 4.8: README.mdの修正
**目的**: `README.md`の参照箇所を修正する。

**作業内容**:
- `README.md`内のすべての参照箇所で、`.develop`を削除して修正

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

**受け入れ基準**:
- `README.md`内のすべての参照箇所が修正されている
- `.develop`が削除されている

- _Requirements: 3.3.2, 3.3.3, 6.3_
- _Design: 3.2_

---

#### タスク 4.9: その他のファイルの修正
**目的**: その他のファイル（スクリプトファイル、その他のドキュメントファイル）の参照箇所を修正する。

**作業内容**:
- スクリプトファイル内の参照箇所を修正（該当する場合）
- その他のドキュメントファイル内の参照箇所を修正（該当する場合）
- develop環境用ファイルの参照は`.develop`を削除して修正
- staging/production環境用ファイルの参照は削除または修正

**受け入れ基準**:
- すべての参照箇所が修正されている
- develop環境用ファイルの参照が`.develop`を削除して修正されている
- staging/production環境用ファイルの参照が削除または修正されている

- _Requirements: 3.3.2, 3.3.3, 6.3_
- _Design: 3.3_

---

### Phase 5: 動作確認

#### タスク 5.1: docker-composeコマンドの動作確認
**目的**: リネーム後のファイルでdocker-composeコマンドが正常に動作することを確認する。

**作業内容**:
- リネーム後のファイルでdocker-composeコマンドを実行
- コマンドが正常に動作することを確認

**確認コマンド**:
```bash
# 設定ファイルの確認
docker-compose -f docker-compose.api.yml config
docker-compose -f docker-compose.admin.yml config
docker-compose -f docker-compose.client.yml config
```

**受け入れ基準**:
- リネーム後のファイルでdocker-composeコマンドが正常に動作する
- エラーが発生しない

- _Requirements: 6.4_
- _Design: 5.4_

---

#### タスク 5.2: ドキュメントの整合性確認
**目的**: ドキュメントを読んで、新しいファイル名が正しく記載されていることを確認する。

**作業内容**:
- `docs/Docker.md`を読んで、新しいファイル名が正しく記載されていることを確認
- `README.md`を読んで、新しいファイル名が正しく記載されていることを確認
- その他のドキュメントファイルを読んで、新しいファイル名が正しく記載されていることを確認（該当する場合）
- ドキュメントに矛盾や不整合がないことを確認

**受け入れ基準**:
- すべてのドキュメントで、新しいファイル名が正しく記載されている
- ドキュメントに矛盾や不整合がない

- _Requirements: 6.4_
- _Design: 5.4_

---

#### タスク 5.3: 参照箇所の最終確認
**目的**: すべての参照箇所が修正されていることを最終確認する。

**作業内容**:
- 再度参照箇所を検索して、修正漏れがないことを確認
- staging環境用ファイルの参照が削除または修正されていることを確認
- production環境用ファイルの参照が削除または修正されていることを確認

**確認コマンド**:
```bash
# 修正漏れがないことを確認
grep -r "docker-compose\.\(admin\|api\|client\)\.\(develop\|staging\|production\)" .
```

**受け入れ基準**:
- 修正漏れがないことを確認できた
- staging環境用ファイルの参照が削除または修正されている
- production環境用ファイルの参照が削除または修正されている

- _Requirements: 6.3_
- _Design: 5.3_

---

## 受け入れ基準の確認

### 要件定義書の受け入れ基準

#### 6.1 docker-composeファイルの削除
- [ ] `docker-compose.admin.production.yml`が削除されている
- [ ] `docker-compose.admin.staging.yml`が削除されている
- [ ] `docker-compose.api.production.yml`が削除されている
- [ ] `docker-compose.api.staging.yml`が削除されている
- [ ] `docker-compose.client.production.yml`が削除されている
- [ ] `docker-compose.client.staging.yml`が削除されている

#### 6.2 docker-composeファイルのリネーム
- [ ] `docker-compose.admin.develop.yml`が`docker-compose.admin.yml`にリネームされている
- [ ] `docker-compose.api.develop.yml`が`docker-compose.api.yml`にリネームされている
- [ ] `docker-compose.client.develop.yml`が`docker-compose.client.yml`にリネームされている
- [ ] リネーム後のファイルが存在する
- [ ] リネーム前のファイルが存在しない

#### 6.3 参照箇所の修正
- [ ] `docs/Docker.md`内のすべての参照箇所が修正されている
- [ ] その他のドキュメントファイル内の参照箇所が修正されている
- [ ] README.mdやその他のMarkdownファイル内の参照箇所が修正されている
- [ ] スクリプトファイル内の参照箇所が修正されている（該当する場合）
- [ ] staging環境用ファイルの参照が削除または修正されている
- [ ] production環境用ファイルの参照が削除または修正されている

#### 6.4 動作確認
- [ ] リネーム後のファイルでdocker-composeコマンドが正常に動作する
- [ ] ドキュメントを読んで、新しいファイル名が正しく記載されている
- [ ] ドキュメントに矛盾や不整合がない
