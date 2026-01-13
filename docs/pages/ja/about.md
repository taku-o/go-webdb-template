---
layout: default
title: プロジェクト概要
lang: ja
---

# プロジェクト概要

Go WebDB Templateは、Go + Next.js + Database Sharding対応のサンプルプロジェクトです。大規模プロジェクト向けの構成を採用しています。

---

## 目的

本プロジェクトは以下の目的で作成されました：

- **スケーラビリティの実証**: データベースシャーディングによる水平スケーリングの実装例を提供
- **ベストプラクティスの提示**: レイヤードアーキテクチャ、テスト戦略、環境別設定管理の実装例
- **学習リソース**: 大規模アプリケーション開発の参考となる実装パターン

---

## 技術スタック

### サーバー

| 項目 | 技術 |
|------|------|
| 言語 | Go 1.21+ |
| アーキテクチャ | レイヤードアーキテクチャ |
| データベース | PostgreSQL/MySQL |
| ORM | GORM v1.25.12 |
| HTTPルーター | Echo v4 |
| API仕様 | Huma (OpenAPI自動生成) |

### クライアント

| 項目 | 技術 |
|------|------|
| フレームワーク | Next.js 14 (App Router) |
| 言語 | TypeScript 5+ |
| UIコンポーネント | shadcn/ui |
| 認証 | NextAuth (Auth.js) v5 |
| スタイリング | Tailwind CSS |

### テスト

| 項目 | 技術 |
|------|------|
| サーバー | Go testing, testify |
| クライアント単体 | Jest, React Testing Library |
| E2E | Playwright |
| APIモック | MSW |

---

## 主要な機能

- **Sharding対応**: テーブルベースシャーディング（32分割）で複数DBへデータ分散
- **GORM対応**: Writer/Reader分離をサポート
- **GoAdmin管理画面**: Webベースの管理画面でデータ管理
- **レイヤー分離**: API層、Usecase層、Service層、Repository層、DB層で責務を明確化
- **環境別設定**: develop/staging/production環境で設定切り替え
- **型安全**: TypeScriptによる型定義
- **テスト**: ユニット/統合/E2Eテスト対応
- **レートリミット**: IPアドレス単位でのAPI呼び出し制限
- **ジョブキュー**: Redis + Asynqを使用したバックグラウンドジョブ処理
- **メール送信**: 標準出力、Mailpit、AWS SES対応のメール送信機能
- **ファイルアップロード**: TUSプロトコルによる大容量ファイルアップロード
- **ログ機能**: アクセスログ、メール送信ログ、SQLログの出力
- **Docker対応**: APIサーバー、Adminサーバー、クライアントサーバーのDocker化

---

## ナビゲーション

- [ホーム]({{ site.baseurl }}/pages/ja/)
- [セットアップ手順]({{ site.baseurl }}/pages/ja/setup)
- [English]({{ site.baseurl }}/pages/en/about)
