import Link from 'next/link'
import { auth } from '@/auth'
import TodayApiButton from '@/components/TodayApiButton'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { AuthButtons } from '@/components/auth/auth-buttons'

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
    {
      title: 'フィード',
      description: 'Twitter風のフィードUIコンポーネント',
      href: '/dm_feed',
    },
    {
      title: '動画プレイヤー',
      description: '動画プレイヤーコンポーネントのデモ',
      href: '/dm_videoplayer',
    },
  ]

  return (
    <main className="min-h-screen p-4 sm:p-6 md:p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-3xl sm:text-4xl font-bold mb-6 sm:mb-8">Go DB Project Sample</h1>

        {/* プロジェクト説明 */}
        <p className="mb-6 sm:mb-8 text-muted-foreground text-sm sm:text-base">
          Go + Next.js + Sharding対応のサンプルプロジェクトです。
        </p>

        <Separator className="my-6 sm:my-8" />

        {/* データ操作機能 */}
        <section aria-label="機能一覧">
          <h2 className="sr-only">利用可能な機能</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 sm:gap-6 mb-6 sm:mb-8">
            {features.map((feature) => (
              <Link key={feature.href} href={feature.href} aria-label={`${feature.title}: ${feature.description}`}>
                <Card className="h-full hover:shadow-lg transition-shadow focus-within:ring-2 focus-within:ring-ring">
                  <CardHeader>
                    <CardTitle className="text-xl sm:text-2xl">{feature.title}</CardTitle>
                    <CardDescription className="text-sm">{feature.description}</CardDescription>
                  </CardHeader>
                </Card>
              </Link>
            ))}
          </div>
        </section>

        <Separator className="my-6 sm:my-8" />

        {/* 認証状態の表示 */}
        <section aria-label="認証状態">
          <Card className="mb-6 sm:mb-8">
            <CardHeader>
              <CardTitle>認証状態</CardTitle>
            </CardHeader>
            <CardContent>
              <AuthButtons user={session?.user ?? null} />
            </CardContent>
          </Card>
        </section>

        <Separator className="my-6 sm:my-8" />

        {/* Today API (Private Endpoint) */}
        <TodayApiButton />
      </div>
    </main>
  )
}
