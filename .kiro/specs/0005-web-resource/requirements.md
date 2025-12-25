# Webリソース（CSS・画像）配置要件定義書

## 1. 概要

### 1.1 プロジェクト情報
- **プロジェクト名**: go-webdb-template
- **Issue番号**: #7
- **Issueタイトル**: CSS、imageファイルの置き場所を決める
- **Feature名**: 0005-web-resource
- **作成日**: 2025-01-27

### 1.2 目的
CSSファイルと画像ファイルの配置場所を決定し、クライアント側（Next.js）から参照できるように実装する。このプロジェクトは参考用のサンプルプロジェクトであるため、ファイルの中身は適当で良い。

### 1.3 スコープ
- CSSファイルの配置場所の決定と実装
- 画像ファイルの配置場所の決定と実装
- クライアント側（Next.js）からの参照方法の実装
- サンプルCSSファイルと画像ファイルの作成
- 既存コードへの統合

## 2. 背景・現状分析

### 2.1 現在の実装
- **クライアント側**: Next.js 14 (App Router) + TypeScript
- **サーバー側**: Go言語、gorilla/muxを使用したAPIサーバー
- **静的ファイル配信**: 未実装
- **CSSファイル**: 未配置
- **画像ファイル**: 未配置

### 2.2 課題点
1. **静的ファイルの配置場所が未決定**: CSSファイルや画像ファイルをどこに配置すべきか明確でない
2. **参照方法が未実装**: クライアント側から静的ファイルを参照する方法が実装されていない
3. **サンプルファイルの不足**: 参考用プロジェクトとして、実際に動作するサンプルファイルがない

### 2.3 本実装による改善点
1. **明確な配置場所**: Next.jsの標準的な配置場所（`client/public/`）に静的ファイルを配置
2. **簡単な参照方法**: Next.jsの標準的な参照方法で静的ファイルにアクセス可能
3. **実用的なサンプル**: 実際に動作するサンプルファイルを提供し、参考用プロジェクトとしての価値を向上

## 3. 機能要件

### 3.1 CSSファイルの配置と参照

#### 3.1.1 配置場所
- **ディレクトリ**: `client/public/css/`
- **ファイル例**: `client/public/css/style.css`
- **命名規則**: ケバブケース（例: `main-style.css`, `component-style.css`）

#### 3.1.2 参照方法
- Next.jsの`<link>`タグまたは`next/head`を使用して参照
- パス: `/css/style.css`（`public/`ディレクトリがルートとして扱われる）
- 例: `<link rel="stylesheet" href="/css/style.css" />`

#### 3.1.3 サンプルCSSファイル
- 最小限のスタイル定義を含むサンプルファイルを作成
- 内容は適当で良い（参考用プロジェクトのため）
- 例: 基本的なリセットCSS、レイアウト用のスタイルなど

### 3.2 画像ファイルの配置と参照

#### 3.2.1 配置場所
- **ディレクトリ**: `client/public/images/`
- **ファイル例**: `client/public/images/logo.png`, `client/public/images/icon.jpg`
- **命名規則**: ケバブケース（例: `user-avatar.png`, `header-logo.svg`）

#### 3.2.2 参照方法
- Next.jsの`<img>`タグまたは`next/image`コンポーネントを使用して参照
- パス: `/images/logo.png`（`public/`ディレクトリがルートとして扱われる）
- 例: `<img src="/images/logo.png" alt="Logo" />` または `<Image src="/images/logo.png" alt="Logo" />`

#### 3.2.3 サンプル画像ファイル
- 最小限の画像ファイルを作成（1x1ピクセルの透明PNGなど、または適当な画像）
- 内容は適当で良い（参考用プロジェクトのため）
- 複数の形式（PNG、JPG、SVG）のサンプルを用意

### 3.3 既存コードへの統合

#### 3.3.1 レイアウトファイルへの統合
- `client/src/app/layout.tsx`にCSSファイルの参照を追加
- グローバルスタイルとして適用

#### 3.3.2 ページコンポーネントへの統合
- 既存のページコンポーネント（`client/src/app/page.tsx`など）に画像参照のサンプルを追加
- 実際に動作することを確認できるようにする

## 4. 非機能要件

### 4.1 パフォーマンス
- 静的ファイルの配信はNext.jsの標準機能を使用（最適化済み）
- 画像ファイルは必要に応じて`next/image`コンポーネントを使用して最適化

### 4.2 保守性
- 明確なディレクトリ構造を維持
- 命名規則に従ったファイル名を使用
- 必要に応じてサブディレクトリで整理（例: `css/components/`, `images/icons/`）

### 4.3 互換性
- Next.js 14の標準的な機能のみを使用
- 既存のコードベースとの互換性を維持

## 5. 制約事項

### 5.1 技術的制約
- Next.js 14 (App Router)の標準機能のみを使用
- サーバー側（Go）での静的ファイル配信は対象外（クライアント側のみ）

### 5.2 プロジェクト制約
- 参考用プロジェクトのため、ファイルの中身は適当で良い
- 実用的な機能よりも、配置場所と参照方法の例示を重視

### 5.3 ディレクトリ構造
- `client/public/`ディレクトリ配下に配置（Next.jsの標準）
- 既存のディレクトリ構造を維持

## 6. 受け入れ基準

### 6.1 機能要件
- [ ] `client/public/css/`ディレクトリが作成されている
- [ ] `client/public/images/`ディレクトリが作成されている
- [ ] サンプルCSSファイル（`client/public/css/style.css`）が作成されている
- [ ] サンプル画像ファイル（`client/public/images/`配下）が作成されている
- [ ] `client/src/app/layout.tsx`にCSSファイルの参照が追加されている
- [ ] 既存のページコンポーネントに画像参照のサンプルが追加されている
- [ ] クライアント側からCSSファイルが正常に読み込まれる
- [ ] クライアント側から画像ファイルが正常に表示される

### 6.2 非機能要件
- [ ] 既存のコードベースとの互換性が維持されている
- [ ] 既存のテストが正常に動作する
- [ ] ディレクトリ構造が明確で保守しやすい

### 6.3 ドキュメント
- [ ] README.mdに静的ファイルの配置場所と参照方法が記載されている（必要に応じて）

## 7. 影響範囲

### 7.1 新規追加が必要なファイル

#### ディレクトリ構造
- `client/public/css/`: CSSファイル用ディレクトリ
- `client/public/images/`: 画像ファイル用ディレクトリ

#### サンプルファイル
- `client/public/css/style.css`: サンプルCSSファイル
- `client/public/images/logo.png`: サンプル画像ファイル（PNG形式）
- `client/public/images/icon.jpg`: サンプル画像ファイル（JPG形式、オプション）
- `client/public/images/icon.svg`: サンプル画像ファイル（SVG形式、オプション）

### 7.2 変更が必要なファイル

#### レイアウトファイル
- `client/src/app/layout.tsx`: CSSファイルの参照を追加

#### ページコンポーネント
- `client/src/app/page.tsx`: 画像参照のサンプルを追加（オプション）

### 7.3 削除されるファイル
- なし

### 7.4 ドキュメント更新
- `README.md`: 静的ファイルの配置場所と参照方法を追加（必要に応じて）

## 8. 実装上の注意事項

### 8.1 Next.jsの静的ファイル配信
- Next.jsは`public/`ディレクトリ配下のファイルを自動的に配信する
- `public/`ディレクトリがルート（`/`）として扱われるため、パスは`/css/style.css`のように記述
- `public/`という文字列をパスに含めない

### 8.2 CSSファイルの参照
- `layout.tsx`で`<link>`タグを使用する場合、`next/head`は不要（App Routerでは直接記述可能）
- グローバルスタイルとして適用する場合は`layout.tsx`に記述

### 8.3 画像ファイルの参照
- `<img>`タグを使用する場合は通常のHTMLと同様
- `next/image`コンポーネントを使用する場合は最適化機能が有効
- 参考用プロジェクトのため、どちらの方法でも良い

### 8.4 サンプルファイルの内容
- ファイルの中身は適当で良い（参考用プロジェクトのため）
- 最小限の内容で動作確認ができる程度のもの
- 実際のプロジェクトで使用する際の参考になるように、コメントなどで説明を追加

### 8.5 ディレクトリ構造の拡張性
- 将来的にファイルが増えることを考慮して、サブディレクトリで整理可能な構造にする
- 例: `css/components/`, `css/pages/`, `images/icons/`, `images/avatars/`など

## 9. 参考情報

### 9.1 Next.js公式ドキュメント
- Next.js Static File Serving: https://nextjs.org/docs/app/building-your-application/optimizing/static-assets
- Next.js Image Optimization: https://nextjs.org/docs/app/building-your-application/optimizing/images

### 9.2 関連Issue
- GitHub Issue #7: CSS、imageファイルの置き場所を決める

### 9.3 既存ドキュメント
- `docs/Architecture.md`: システムアーキテクチャ
- `docs/Project-Structure.md`: プロジェクト構造計画
- `.kiro/steering/structure.md`: ファイル組織とコードパターン
- `.kiro/steering/tech.md`: 技術スタックとアーキテクチャ

### 9.4 既存実装
- `client/src/app/layout.tsx`: ルートレイアウト
- `client/src/app/page.tsx`: トップページ
- `client/next.config.js`: Next.js設定ファイル

