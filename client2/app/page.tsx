import Link from 'next/link'
import { auth, signIn, signOut } from '@/auth'
import TodayApiButton from '@/components/TodayApiButton'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'

export default async function Home() {
  const session = await auth()

  const features = [
    {
      title: 'ユーザー管理',
      description: 'ユーザーの一覧・作成・編集・削除',
      href: '/dm-users',
    },
    {
      title: '投稿管理',
      description: '投稿の一覧・作成・編集・削除',
      href: '/dm-posts',
    },
    {
      title: 'ユーザーと投稿',
      description: 'ユーザーと投稿をJOINして表示（クロスシャードクエリ）',
      href: '/dm-user-posts',
    },
    {
      title: '動画アップロード',
      description: '動画ファイルのアップロード（TUSプロトコル）',
      href: '/dm_movie/upload',
    },
    {
      title: 'メール送信',
      description: 'ウェルカムメールの送信',
      href: '/dm_email/send',
    },
    {
      title: 'ジョブキュー',
      description: '遅延ジョブの登録（参考実装）',
      href: '/dm-jobqueue',
    },
  ]

  return (
    <main className="min-h-screen p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold mb-8">Go DB Project Sample</h1>

        {/* プロジェクト説明 */}
        <p className="mb-8 text-gray-600">
          Go + Next.js + Sharding対応のサンプルプロジェクトです。
        </p>

        <Separator className="my-8" />

        {/* データ操作機能 */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
          {features.map((feature) => (
            <Link key={feature.href} href={feature.href}>
              <Card className="h-full hover:shadow-lg transition-shadow">
                <CardHeader>
                  <CardTitle className="text-2xl">{feature.title}</CardTitle>
                  <CardDescription>{feature.description}</CardDescription>
                </CardHeader>
              </Card>
            </Link>
          ))}
        </div>

        <Separator className="my-8" />

        {/* 認証状態の表示 */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle>認証状態</CardTitle>
          </CardHeader>
          <CardContent>
            {session?.user ? (
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-semibold">ログイン中: {session.user.name}</p>
                  {session.user.email && (
                    <p className="text-sm text-gray-600">{session.user.email}</p>
                  )}
                </div>
                <form action={async () => {
                  "use server"
                  await signOut()
                }}>
                  <Button type="submit" variant="destructive">ログアウト</Button>
                </form>
              </div>
            ) : (
              <div className="flex items-center justify-between">
                <p className="text-gray-600">ログインしていません</p>
                <form action={async () => {
                  "use server"
                  await signIn()
                }}>
                  <Button type="submit">ログイン</Button>
                </form>
              </div>
            )}
          </CardContent>
        </Card>

        <Separator className="my-8" />

        {/* Today API (Private Endpoint) */}
        <TodayApiButton />
      </div>
    </main>
  )
}
