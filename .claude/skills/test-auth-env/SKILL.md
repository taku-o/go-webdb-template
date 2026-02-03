---
name: test-auth-env
description: Go テストや E2E/統合テストの実行で認証エラー（401 Unauthorized 等）が発生したときに使用。APP_ENV=test を指定していない可能性を指摘し、対処法（コマンド例・tech.md の確認）を案内する。
allowed-tools: Read, Glob
---

# テスト認証エラー時の APP_ENV=test 確認

テスト実行中に認証エラー（401 Unauthorized 等）が発生した場合、**テスト実行コマンドに `APP_ENV=test` が付いているか**をまず確認する。

## 認証エラーがテストで出た場合の手順

1. **テスト実行コマンドを確認する**
   - 実行したコマンドに `APP_ENV=test` が含まれているか確認する。
   - 指定していないと認証ミドルウェアがテスト用の認証をスキップせず、401 が返る。

2. **正しいコマンド例**
   ```bash
   # サーバー配下の全テスト
   cd server && APP_ENV=test go test ./...

   # ルートから全パッケージ
   APP_ENV=test go test ./...

   # 特定パッケージのみ
   cd server && APP_ENV=test go test ./internal/repository/...
   cd server && APP_ENV=test go test -v ./...
   ```

3. **ステアリングの確認**
   - `.kiro/steering/tech.md` の「**テスト実行ルール（必須）**」を確認する。
   - 本プロジェクトではテスト実行時に `APP_ENV=test` の指定が必須である旨が記載されている。

## 方針

- 認証エラーが **1件でも** 出た場合は、「今回の修正とは関係ない」と判断せず、**原因を調査する**。
- まずは上記のとおり `APP_ENV=test` が付いているか確認し、付いていなければ付けて再実行する。
- 付いているにもかかわらず 401 が出る場合は、認証設定・テスト用トークン・ミドルウェアの挙動などを確認する。
