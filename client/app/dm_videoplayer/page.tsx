import Link from 'next/link'
import { VideoPlayer } from '@/components/video-player/video-player'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { ArrowLeft } from 'lucide-react'

export default function VideoPlayerDemoPage() {
  return (
    <main className="min-h-screen p-4 sm:p-6 md:p-8">
      <div className="max-w-4xl mx-auto">
        <nav aria-label="パンくずリスト">
          <div className="mb-4 sm:mb-6">
            <Link href="/" className="inline-flex items-center text-primary hover:underline text-sm sm:text-base" aria-label="トップページに戻る">
              <ArrowLeft className="mr-2 h-4 w-4" aria-hidden="true" />
              トップページに戻る
            </Link>
          </div>
        </nav>

        <Card>
          <CardHeader>
            <CardTitle>動画プレイヤー</CardTitle>
            <CardDescription>
              動画プレイヤーコンポーネントのデモページです。
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="w-full max-w-3xl mx-auto">
              <VideoPlayer
                videoUrl="/demo-videos/mini-movie-m.mp4"
                thumbnailUrl="/demo-videos/mini-movie-m.png"
              />
            </div>
          </CardContent>
        </Card>
      </div>
    </main>
  )
}
