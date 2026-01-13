# GitHub Pages対応の設計書

## Overview

### 目的
プロジェクトのGitHub Pagesページを構築し、プロジェクトの概要・目的、セットアップ手順、各種URL情報を多言語対応（日本語・英語）で公開する。シンプルでモダンなデザインを実現し、プロジェクトの理解と導入を容易にする。

### ユーザー
- **外部開発者**: プロジェクトの概要とセットアップ手順を理解し、プロジェクトを導入する
- **プロジェクトメンバー**: プロジェクトの情報を外部に公開する

### 影響
現在のシステム状態を以下のように変更する：
- `docs/`ディレクトリにGitHub Pages用のMarkdownファイルを追加
- Jekyll設定ファイル（`_config.yml`）を追加（必要に応じて）
- `README.md`にGitHub Pagesへのリンクを追加
- 既存のdocsディレクトリのファイルは変更しない

### Goals
- GitHub Pages用のページ構成の実現
- 多言語対応（日本語・英語）の実現
- シンプルでモダンなデザインの実現
- 既存ドキュメントとの整合性の確保

### Non-Goals
- GitHub Pagesの自動デプロイ設定（GitHub Actions等）
- 既存のdocsディレクトリの変更
- プロジェクトの機能追加・変更
- 特殊なJekyllプラグインの使用

## Architecture

### ディレクトリ構造

#### 追加されるファイル構造
```
docs/
├── _config.yml           # Jekyll設定ファイル（必要に応じて）
├── index.md              # トップページ（言語選択）
├── ja/                   # 日本語版ページ
│   ├── index.md         # 日本語版トップページ
│   ├── about.md         # プロジェクト概要・目的
│   └── setup.md         # セットアップ手順
└── en/                   # 英語版ページ
    ├── index.md         # 英語版トップページ
    ├── about.md         # プロジェクト概要・目的
    └── setup.md         # セットアップ手順
```

#### 既存ファイルとの関係
- 既存のdocsディレクトリ内のファイル（Admin.md、API.md等）は変更しない
- `README.md`にGitHub Pagesへのリンクを追加（`https://taku-o.github.io/go-webdb-template/`）
- 新しいファイルのみを追加する
- 既存ファイルへのリンクは必要に応じて追加する

### ページ構成設計

#### 1. トップページ（docs/index.md）
**目的**: 言語選択ページとして機能し、ユーザーを適切な言語版に誘導する

**構成要素**:
- プロジェクト名の表示
- 言語選択リンク
  - 日本語版へのリンク: `/ja/`
  - 英語版へのリンク: `/en/`
- シンプルで分かりやすいデザイン

**実装方法**:
- Markdownファイルとして実装
- Jekyllのfront matterを使用してレイアウトを指定
- 言語選択ボタンまたはリンクを配置

#### 2. 各言語版のトップページ（docs/ja/index.md、docs/en/index.md）
**目的**: 各言語版のエントリーポイントとして機能し、主要なページへのナビゲーションを提供する

**構成要素**:
- プロジェクト名と簡単な説明
- ナビゲーションメニュー
  - About（プロジェクト概要・目的）へのリンク
  - Setup（セットアップ手順）へのリンク
- 言語切り替えリンク
  - 日本語版: `/en/`へのリンク
  - 英語版: `/ja/`へのリンク
- トップページへの戻りリンク

**実装方法**:
- Markdownファイルとして実装
- Jekyllのfront matterを使用
- ナビゲーションメニューをMarkdownのリンクで実装

#### 3. Aboutページ（docs/ja/about.md、docs/en/about.md）
**目的**: プロジェクトの概要・目的、技術スタック、主要な機能を紹介する

**構成要素**:
- プロジェクトの概要・目的
  - README.mdの「プロジェクト概要」セクションの内容を参考にする
- 技術スタックの概要
  - サーバー: Go言語、レイヤードアーキテクチャ、Database Sharding対応
  - クライアント: Next.js 14 (App Router)、TypeScript
  - データベース: PostgreSQL/MySQL
  - テスト: Go testing、Jest、Playwright
- 主要な機能の紹介
  - README.mdの「特徴」セクションの内容を参考にする
  - Sharding対応、GORM対応、GoAdmin管理画面、レイヤー分離など

**実装方法**:
- Markdownファイルとして実装
- README.mdの内容を参考にし、適切に要約・整理
- コードブロックやリストを使用して見やすく表示

#### 4. Setupページ（docs/ja/setup.md、docs/en/setup.md）
**目的**: クライアントサーバーを動作させるまでの詳細なセットアップ手順を提供する

**構成要素**:
1. **前提条件**
   - Go 1.21+
   - Node.js 18+
   - Docker（PostgreSQLコンテナ用）
   - Atlas CLI
   - Redis（オプション）

2. **初期セットアップ（docs/Initial-Setup.mdの内容）**
   - Docker、Cursor、Goのインストール
   - Homebrewのインストール
   - GitHub CLIのインストール
   - Atlasのインストール
   - Node.js（nvm）のインストール
   - Claude Codeのインストール
   - uvのインストール
   - Serenaの設定

3. **依存関係のインストール**
   - サーバー側: `cd server && go mod download`
   - クライアント側: `cd client && npm install --legacy-peer-deps`

4. **データベースのセットアップ**
   - PostgreSQLの起動: `./scripts/start-postgres.sh start`
   - マイグレーションの適用: `./scripts/migrate.sh all`
   - 接続情報の記載

5. **Redisの起動（オプション）**
   - Redisの起動: `./scripts/start-redis.sh start`
   - Redis Insightの起動（オプション）: `./scripts/start-redis-insight.sh start`

6. **Auth0アカウントの設定**
   - Auth0ダッシュボードでの設定手順
   - コールバックURLの設定: `http://localhost:3000/api/auth/callback/auth0`
   - ログアウトURLの設定: `http://localhost:3000`
   - Web Originsの設定: `http://localhost:3000`

7. **クライアント環境変数の設定**
   - `client/.env.local`の作成
   - `AUTH_SECRET`の生成: `npm run cli:generate-secret`
   - 環境変数の設定内容:
     ```
     AUTH_SECRET=<生成した秘密鍵>
     AUTH_URL=http://localhost:3000
     AUTH0_ISSUER=https://your-tenant.auth0.com
     AUTH0_CLIENT_ID=your-client-id
     AUTH0_CLIENT_SECRET=your-client-secret
     AUTH0_AUDIENCE=https://your-api-audience
     NEXT_PUBLIC_API_KEY=your-api-key
     NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
     APP_ENV=test
     ```

8. **サーバーの起動**
   - APIサーバーの起動: `cd server && APP_ENV=develop go run cmd/server/main.go`
   - Adminサーバーの起動: `cd server && APP_ENV=develop go run cmd/admin/main.go`
   - クライアントサーバーの起動: `cd client && npm run dev`

9. **各種URL情報**
   - クライアントURL: http://localhost:3000
   - APIサーバーdoc URL: http://localhost:8080/docs
   - AdminサーバーURL: http://localhost:8081/admin
     - ID: admin
     - Password: admin123

**実装方法**:
- Markdownファイルとして実装
- README.mdとdocs/Initial-Setup.mdの内容を統合・整理
- コードブロックを使用してコマンドを表示
- 表を使用してURL情報を整理

### Jekyll設定設計

#### _config.ymlの設計
**目的**: Jekyllの動作を制御し、テーマやプラグインを設定する

**設定内容**:
```yaml
# サイト設定
title: "Go WebDB Template"
description: "Go + Next.js + Database Sharding対応のサンプルプロジェクト"
url: "https://taku-o.github.io"
baseurl: "/go-webdb-template"

# テーマ設定
theme: minima

# プラグイン設定
plugins:
  - jekyll-feed
  - jekyll-sitemap

# Markdown設定
markdown: kramdown
kramdown:
  syntax_highlighter: rouge

# 除外ファイル
exclude:
  - README.md
  - Gemfile
  - Gemfile.lock
  - node_modules
  - vendor
```

**実装方針**:
- GitHub Pagesの標準機能のみを使用
- minimaテーマを使用（デフォルトテーマ）
- 特殊なプラグインは使用しない

### デザイン設計

#### 全体的なデザイン方針
- **シンプルで読みやすい**: 余計な装飾を避け、コンテンツに集中
- **モダンな見た目**: 適切な余白とタイポグラフィを使用
- **レスポンシブデザイン**: モバイル・タブレット・デスクトップで適切に表示

#### ナビゲーションデザイン
- **トップページ**: 言語選択を中央に配置
- **各言語版トップページ**: ナビゲーションメニューを上部またはサイドに配置
- **言語切り替え**: 各ページの上部または下部に配置
- **トップページへの戻り**: 各ページの上部に配置

#### タイポグラフィ
- 見出しは適切な階層で使用（h1, h2, h3）
- コードブロックは適切にシンタックスハイライト
- リストは適切にインデント

#### カラースキーム
- Jekyllのminimaテーマのデフォルトカラーを使用
- 必要に応じてカスタムCSSを追加（最小限）

### 多言語対応設計

#### 言語切り替えの実装
- **トップページ（docs/index.md）**: 言語選択リンクを提供
- **各言語版ページ**: 言語切り替えリンクを各ページに配置
  - 日本語版ページ: `/en/`へのリンク
  - 英語版ページ: `/ja/`へのリンク

#### コンテンツの整合性
- 日本語版と英語版の内容は一致させる
- 技術用語は適切に翻訳または英語表記を併記
- 既存のREADME.mdやdocs/配下のドキュメントと整合性を保つ

### コンテンツ設計

#### Aboutページのコンテンツ
**日本語版**:
- プロジェクトの概要・目的
- 技術スタックの概要
- 主要な機能の紹介（README.mdの「特徴」セクションを参考）

**英語版**:
- Project Overview and Purpose
- Technology Stack Overview
- Key Features (referring to the "Features" section of README.md)

#### Setupページのコンテンツ
**日本語版**:
- セットアップ手順を詳細に記載
- README.mdとdocs/Initial-Setup.mdの内容を統合

**英語版**:
- Detailed setup instructions
- Integrate content from README.md and docs/Initial-Setup.md

## Implementation Strategy

### 実装フェーズ

#### Phase 1: ディレクトリ構造の作成
1. `docs/ja/`ディレクトリの作成
2. `docs/en/`ディレクトリの作成

#### Phase 2: Jekyll設定ファイルの作成
1. `docs/_config.yml`の作成（必要に応じて）
2. テーマとプラグインの設定

#### Phase 3: トップページの作成
1. `docs/index.md`の作成
2. 言語選択リンクの実装

#### Phase 4: 日本語版ページの作成
1. `docs/ja/index.md`の作成
2. `docs/ja/about.md`の作成
3. `docs/ja/setup.md`の作成

#### Phase 5: 英語版ページの作成
1. `docs/en/index.md`の作成
2. `docs/en/about.md`の作成
3. `docs/en/setup.md`の作成

#### Phase 6: デザインの調整
1. ナビゲーションの確認・調整
2. レイアウトの確認・調整
3. レスポンシブデザインの確認

#### Phase 7: コンテンツの確認・修正
1. 既存ドキュメントとの整合性確認
2. 多言語対応の整合性確認
3. 誤字脱字の確認

#### Phase 8: README.mdへのGitHub Pagesリンク追加
1. `README.md`にGitHub Pagesへのリンクを追加
   - リンクURL: `https://taku-o.github.io/go-webdb-template/`
   - プロジェクト概要セクションまたは適切な場所に追加
   - リンクテキストは「GitHub Pages」または「Documentation」など

### 実装の優先順位
1. **高優先度**: ディレクトリ構造の作成、トップページ、日本語版ページ
2. **中優先度**: 英語版ページ、Jekyll設定
3. **低優先度**: デザインの微調整

## Technical Details

### Jekyll設定

#### テーマの選択
- **minima**: GitHub Pagesのデフォルトテーマ
- シンプルで読みやすい
- カスタマイズが容易

#### プラグイン
- **jekyll-feed**: RSSフィードの生成（オプション）
- **jekyll-sitemap**: サイトマップの生成（オプション）
- 標準的なプラグインのみを使用

#### Markdown処理
- **kramdown**: JekyllのデフォルトMarkdownプロセッサ
- **rouge**: シンタックスハイライト

### ファイル命名規則
- すべて小文字で統一
- ハイフン（-）を使用して単語を区切る
- 拡張子は`.md`（Markdown）

### リンク構造
- 相対パスを使用
- 言語プレフィックス（`/ja/`、`/en/`）を使用
- 既存のdocsディレクトリ内のファイルへのリンクは必要に応じて追加

## Testing Strategy

### テスト項目

#### 機能テスト
- [ ] トップページが正常に表示される
- [ ] 言語選択リンクが正常に動作する
- [ ] 各言語版のトップページが正常に表示される
- [ ] ナビゲーションメニューが正常に動作する
- [ ] 言語切り替えリンクが正常に動作する
- [ ] Aboutページが正常に表示される
- [ ] Setupページが正常に表示される

#### コンテンツテスト
- [ ] 日本語版のコンテンツが適切に記載されている
- [ ] 英語版のコンテンツが適切に記載されている
- [ ] 日本語版と英語版の内容が一致している
- [ ] 既存ドキュメントとの整合性が保たれている

#### デザインテスト
- [ ] シンプルで読みやすいデザインになっている
- [ ] モダンな見た目になっている
- [ ] 適切な余白とタイポグラフィが使用されている
- [ ] コードブロックが適切に表示される
- [ ] レスポンシブデザインが適切に機能する

#### 整合性テスト
- [ ] 既存のREADME.mdやdocs/配下のドキュメントと整合性が保たれている
- [ ] 記載されていない情報が追加されていない

### テスト環境
- **ローカル環境**: Jekyllをローカルで実行して確認
  ```bash
  cd docs
  bundle install
  bundle exec jekyll serve
  ```
- **GitHub Pages**: GitHubにプッシュして実際のGitHub Pagesで確認

## Risk Management

### リスクと対策

#### リスク1: 既存のdocsディレクトリとの競合
**影響**: 既存のドキュメントが正しく表示されなくなる可能性
**対策**: 
- 既存ファイルは変更しない
- 新しいファイルのみを追加
- 既存ファイルへのリンクは必要に応じて追加

#### リスク2: Jekyll設定の不備
**影響**: GitHub Pagesが正常に動作しない可能性
**対策**:
- 標準的な設定のみを使用
- 特殊なプラグインは使用しない
- GitHub Pagesのドキュメントを参考にする

#### リスク3: コンテンツの不整合
**影響**: 既存ドキュメントと内容が一致しない可能性
**対策**:
- 既存のREADME.mdやdocs/配下のドキュメントを参考にする
- 記載されていない情報は追加しない
- レビュー時に整合性を確認する

#### リスク4: 多言語対応の不整合
**影響**: 日本語版と英語版の内容が一致しない可能性
**対策**:
- 日本語版と英語版の内容を一致させる
- レビュー時に両方のバージョンを確認する

## Dependencies

### 外部依存
- **GitHub Pages**: ホスティングサービス
- **Jekyll**: 静的サイトジェネレーター（GitHub Pagesのデフォルト）
- **minimaテーマ**: Jekyllのデフォルトテーマ

### 内部依存
- **既存のdocsディレクトリ**: 参考情報として使用
- **README.md**: コンテンツの参考情報として使用
- **docs/Initial-Setup.md**: セットアップ手順の参考情報として使用

## Future Considerations

### 将来の拡張可能性
- 追加の言語対応（中国語、韓国語など）
- 追加のページ（APIドキュメント、アーキテクチャ説明など）
- カスタムテーマの適用
- 検索機能の追加

### 注意事項
- 本実装では将来の拡張は考慮しない
- 必要に応じて将来追加する

## Conclusion

本設計書では、GitHub Pages対応の実装方針を定義した。シンプルでモダンなデザインを実現し、多言語対応（日本語・英語）を実現する。既存のdocsディレクトリとの整合性を保ちながら、新しいページを追加する。
