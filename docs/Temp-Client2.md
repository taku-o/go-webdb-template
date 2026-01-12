# Client2 App

新clientアプリケーション（precedentテンプレートベース）

## セットアップ

1. 依存関係のインストール
   ```bash
   cd client2
   npm install
   ```

2. 環境変数の設定
   `.env.local`を作成して以下の環境変数を設定：
   ```
   AUTH_SECRET=your-secret-key-here
   AUTH_URL=http://localhost:3000
   ```
   
   **注意**: `AUTH_SECRET`は適切な秘密鍵を生成する必要があります。開発環境では`openssl rand -base64 32`などで生成できます。

3. 開発サーバーの起動
   ```bash
   npm run dev
   ```
   
   開発サーバーは`http://localhost:3000`で起動します。

## 技術スタック

- Next.js 14+ (App Router)
- TypeScript 5+
- shadcn/ui
- NextAuth (Auth.js)
- Tailwind CSS

## 利用可能なスクリプト

- `npm run dev` - 開発サーバーを起動（ポート3000）
- `npm run build` - プロダクションビルドを実行
- `npm run start` - プロダクションビルドを起動（ポート3000）
- `npm run lint` - ESLintを実行
- `npm run type-check` - TypeScript型チェックを実行
- `npm run format` - Prettierでフォーマットを確認
- `npm run format:write` - Prettierでフォーマットを適用

## 注意事項

このドキュメントは一時的なもので、`client`から`client2`への移行が完了したら、この内容をREADMEに移植する想定です。
