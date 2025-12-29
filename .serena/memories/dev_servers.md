# 開発環境サーバー構成

開発環境では3つのサーバーを起動する必要があります。

## サーバー一覧

| サーバー | ポート | ディレクトリ | 起動コマンド |
|---------|-------|-------------|-------------|
| API サーバー | 8080 | `server/cmd/server` | `APP_ENV=develop go run ./cmd/server/main.go` |
| クライアント | 3000 | `client/` | `npm run dev` |
| Admin | 8081 | `server/cmd/admin` | `APP_ENV=develop go run ./cmd/admin/main.go` |

## 注意事項

- 「サーバーを起動して」と言われた場合、上記3つ全てを起動すること
- クライアントはNext.jsアプリケーション（port 3000）
- API サーバーとAdminサーバーはGoアプリケーション
