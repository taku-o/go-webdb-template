# Webリソース（CSS・画像）配置設計書

## 1. 概要

### 1.1 設計の目的
要件定義書に基づき、Next.jsアプリケーションにおけるCSSファイルと画像ファイルの配置場所と参照方法を設計する。参考用プロジェクトとして、実際に動作するサンプルファイルを提供し、静的ファイルの配置と参照のベストプラクティスを示す。

### 1.2 設計の範囲
- 静的ファイルのディレクトリ構造設計
- CSSファイルの配置と参照方法の設計
- 画像ファイルの配置と参照方法の設計
- 既存コードへの統合設計
- サンプルファイルの作成方針

## 2. アーキテクチャ設計

### 2.1 ディレクトリ構造

Next.jsの標準的な`public/`ディレクトリを使用して静的ファイルを配置する。

```
client/
├── public/                    # 静的ファイルのルートディレクトリ
│   ├── css/                   # CSSファイル用ディレクトリ
│   │   └── style.css          # サンプルCSSファイル
│   └── images/                # 画像ファイル用ディレクトリ
│       ├── logo.svg           # サンプルSVG画像
│       ├── logo.png           # サンプルPNG画像
│       └── icon.jpg           # サンプルJPG画像
└── src/
    └── app/
        ├── layout.tsx         # ルートレイアウト（CSS参照を追加）
        └── page.tsx           # トップページ（画像参照のサンプルを追加）
```

### 2.2 Next.jsの静的ファイル配信メカニズム

Next.jsは`public/`ディレクトリ配下のファイルを自動的に配信する。

```
ブラウザリクエスト: /css/style.css
         ↓
Next.js: public/css/style.css を配信
         ↓
ブラウザ: CSSファイルを読み込み
```

**重要なポイント**:
- `public/`ディレクトリがルート（`/`）として扱われる
- パスに`public/`を含めない（`/css/style.css`が正しい）
- 開発環境と本番環境の両方で同じ動作

### 2.3 ファイル参照の設計

#### CSSファイルの参照
- **場所**: `client/src/app/layout.tsx`
- **方法**: `<link>`タグを直接記述（App Routerでは`next/head`は不要）
- **パス**: `/css/style.css`

#### 画像ファイルの参照
- **場所**: `client/src/app/page.tsx`
- **方法**: `<img>`タグまたは`next/image`コンポーネント
- **パス**: `/images/logo.png`など

## 3. コンポーネント設計

### 3.1 レイアウトコンポーネント（layout.tsx）

| 項目 | 詳細 |
|------|------|
| ファイル | `client/src/app/layout.tsx` |
| 変更内容 | CSSファイルの参照を追加 |
| 要件 | 3.3.1 |

**変更内容**:
- `<head>`タグ内に`<link rel="stylesheet" href="/css/style.css" />`を追加
- グローバルスタイルとして適用

**実装例**:
```tsx
<html lang="en">
  <head>
    <link rel="stylesheet" href="/css/style.css" />
  </head>
  <body>{children}</body>
</html>
```

### 3.2 ページコンポーネント（page.tsx）

| 項目 | 詳細 |
|------|------|
| ファイル | `client/src/app/page.tsx` |
| 変更内容 | 画像参照のサンプルを追加 |
| 要件 | 3.3.2 |

**変更内容**:
- サンプル画像（SVG、PNG、JPG）を表示するセクションを追加
- 各画像形式の参照例を示す

**実装例**:
```tsx
<div>
  <img src="/images/logo.svg" alt="Logo SVG" />
  <img src="/images/logo.png" alt="Logo PNG" />
  <img src="/images/icon.jpg" alt="Icon JPG" />
</div>
```

## 4. データモデル

この機能はデータモデルを使用しないため、該当なし。

## 5. エラーハンドリング

### 5.1 ファイルが見つからない場合

Next.jsは自動的に404エラーを返す。開発環境ではコンソールに警告が表示される。

### 5.2 画像読み込みエラー

`<img>`タグの`onError`属性を使用してエラーハンドリングを実装可能（今回はサンプルのため実装しない）。

## 6. テスト戦略

### 6.1 ユニットテスト
- 不要（静的ファイルの配置のみ）

### 6.2 統合テスト
- 不要（静的ファイルの配置のみ）

### 6.3 E2Eテスト
- オプション: Playwrightで画像が正しく表示されることを確認
- 今回は実装しない（参考用プロジェクトのため）

## 7. 実装上の注意事項

### 7.1 Next.js App Routerでの`<head>`タグ

Next.js 14のApp Routerでは、`layout.tsx`で`<head>`タグを直接使用できる。ただし、`metadata`オブジェクトを使用する方法もあるが、CSSファイルの参照には`<link>`タグの直接記述が適切。

### 7.2 画像最適化

本番環境では`next/image`コンポーネントを使用することで画像最適化が可能。参考用プロジェクトのため、今回は`<img>`タグを使用する。

### 7.3 ファイル命名規則

- CSSファイル: ケバブケース（例: `main-style.css`）
- 画像ファイル: ケバブケース（例: `user-avatar.png`）
- 将来的な拡張を考慮してサブディレクトリで整理可能

## 8. 参考情報

### 8.1 Next.js公式ドキュメント
- Static File Serving: https://nextjs.org/docs/app/building-your-application/optimizing/static-assets
- Image Optimization: https://nextjs.org/docs/app/building-your-application/optimizing/images

### 8.2 既存実装
- `client/src/app/layout.tsx`: ルートレイアウト
- `client/src/app/page.tsx`: トップページ

