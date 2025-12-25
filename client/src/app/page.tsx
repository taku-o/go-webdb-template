import Link from 'next/link'

export default function Home() {
  return (
    <main className="min-h-screen p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold mb-8">Go DB Project Sample</h1>

        <p className="mb-8 text-gray-600">
          Go + Next.js + Sharding対応のサンプルプロジェクトです。
        </p>

        {/* サンプル画像ファイルの参照例 */}
        {/* Next.jsのpublicディレクトリ配下の画像ファイルは /images/logo.png のように参照できます */}
        <div className="mb-8 p-4 border rounded-lg bg-gray-50">
          <h3 className="font-semibold mb-4">静的ファイルの参照例</h3>
          <div className="flex gap-4 items-center">
            <div>
              <p className="text-sm text-gray-600 mb-2">SVG画像:</p>
              <img src="/images/logo.svg" alt="Logo SVG" className="w-20 h-20" />
            </div>
            <div>
              <p className="text-sm text-gray-600 mb-2">PNG画像:</p>
              <img src="/images/logo.png" alt="Logo PNG" className="w-20 h-20 border" />
            </div>
            <div>
              <p className="text-sm text-gray-600 mb-2">JPG画像:</p>
              <img src="/images/icon.jpg" alt="Icon JPG" className="w-20 h-20 border" />
            </div>
          </div>
          <p className="text-xs text-gray-500 mt-4">
            これらの画像は client/public/images/ ディレクトリに配置されています。
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <Link
            href="/users"
            className="p-6 border rounded-lg hover:border-blue-500 hover:shadow-lg transition-all"
          >
            <h2 className="text-2xl font-semibold mb-2">ユーザー管理</h2>
            <p className="text-gray-600">ユーザーの一覧・作成・編集・削除</p>
          </Link>

          <Link
            href="/posts"
            className="p-6 border rounded-lg hover:border-blue-500 hover:shadow-lg transition-all"
          >
            <h2 className="text-2xl font-semibold mb-2">投稿管理</h2>
            <p className="text-gray-600">投稿の一覧・作成・編集・削除</p>
          </Link>

          <Link
            href="/user-posts"
            className="p-6 border rounded-lg hover:border-blue-500 hover:shadow-lg transition-all"
          >
            <h2 className="text-2xl font-semibold mb-2">ユーザーと投稿</h2>
            <p className="text-gray-600">ユーザーと投稿をJOINして表示（クロスシャードクエリ）</p>
          </Link>

          <div className="p-6 border rounded-lg bg-gray-50">
            <h2 className="text-2xl font-semibold mb-2">技術スタック</h2>
            <ul className="text-gray-600 space-y-1">
              <li>• Go (Sharding対応)</li>
              <li>• Next.js 14 (App Router)</li>
              <li>• TypeScript</li>
              <li>• SQLite (開発環境)</li>
            </ul>
          </div>
        </div>

        <div className="mt-8 p-4 border rounded-lg bg-blue-50">
          <h3 className="font-semibold mb-2">開発サーバーの起動方法</h3>
          <div className="space-y-2 text-sm">
            <div>
              <span className="font-mono bg-gray-100 px-2 py-1 rounded">cd server && go run cmd/server/main.go</span>
              <span className="ml-2 text-gray-600">- APIサーバー起動 (Port 8080)</span>
            </div>
            <div>
              <span className="font-mono bg-gray-100 px-2 py-1 rounded">cd client && npm run dev</span>
              <span className="ml-2 text-gray-600">- フロントエンド起動 (Port 3000)</span>
            </div>
          </div>
        </div>
      </div>
    </main>
  )
}
