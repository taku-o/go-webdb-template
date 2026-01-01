'use client'

import Link from 'next/link'
import { useUser } from '@auth0/nextjs-auth0'
import TodayApiButton from '@/components/TodayApiButton'

export default function Home() {
  const { user, error, isLoading } = useUser()

  // 未認証(Unauthorized)は正常な未ログイン状態として扱う
  const isUnauthorized = error?.message === 'Unauthorized'
  const hasRealError = error && !isUnauthorized

  return (
    <main className="min-h-screen p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold mb-8">Go DB Project Sample</h1>

        {/* プロジェクト説明 */}
        <p className="mb-8 text-gray-600">
          Go + Next.js + Sharding対応のサンプルプロジェクトです。
        </p>

        <hr className="my-8 border-gray-200" />

        {/* データ操作機能 */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
          <Link
            href="/dm-users"
            className="p-6 border rounded-lg hover:border-blue-500 hover:shadow-lg transition-all"
          >
            <h2 className="text-2xl font-semibold mb-2">ユーザー管理</h2>
            <p className="text-gray-600">ユーザーの一覧・作成・編集・削除</p>
          </Link>

          <Link
            href="/dm-posts"
            className="p-6 border rounded-lg hover:border-blue-500 hover:shadow-lg transition-all"
          >
            <h2 className="text-2xl font-semibold mb-2">投稿管理</h2>
            <p className="text-gray-600">投稿の一覧・作成・編集・削除</p>
          </Link>

          <Link
            href="/dm-user-posts"
            className="p-6 border rounded-lg hover:border-blue-500 hover:shadow-lg transition-all"
          >
            <h2 className="text-2xl font-semibold mb-2">ユーザーと投稿</h2>
            <p className="text-gray-600">ユーザーと投稿をJOINして表示（クロスシャードクエリ）</p>
          </Link>

          <Link
            href="/dm_movie/upload"
            className="p-6 border rounded-lg hover:border-blue-500 hover:shadow-lg transition-all"
          >
            <h2 className="text-2xl font-semibold mb-2">動画アップロード</h2>
            <p className="text-gray-600">動画ファイルのアップロード（TUSプロトコル）</p>
          </Link>
        </div>

        <hr className="my-8 border-gray-200" />

        {/* 認証状態の表示 */}
        <div className="mb-8 p-4 border rounded-lg bg-gray-50">
          {isLoading ? (
            <p className="text-gray-600">Loading...</p>
          ) : hasRealError ? (
            <div className="text-red-600">
              <p>認証エラーが発生しました: {error.message}</p>
              <a href="/auth/login" className="text-blue-600 underline">
                再度ログイン
              </a>
            </div>
          ) : user ? (
            <div className="flex items-center justify-between">
              <div>
                <p className="font-semibold">ログイン中: {user.name}</p>
                {user.email && (
                  <p className="text-sm text-gray-600">{user.email}</p>
                )}
              </div>
              <a
                href="/auth/logout"
                className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600 transition-colors"
              >
                ログアウト
              </a>
            </div>
          ) : (
            <div className="flex items-center justify-between">
              <p className="text-gray-600">ログインしていません</p>
              <a
                href="/auth/login"
                className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors"
              >
                ログイン
              </a>
            </div>
          )}
        </div>

        <hr className="my-8 border-gray-200" />

        {/* Today API (Private Endpoint) */}
        <TodayApiButton />

        <hr className="my-8 border-gray-200" />

        {/* サンプル画像ファイルの参照例 */}
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
        </div>
      </div>
    </main>
  )
}
