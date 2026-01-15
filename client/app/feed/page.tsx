import { redirect } from 'next/navigation'

// 固定のダミーuserId（ダミー実装のため）
// 将来的には、認証ユーザーのIDを使用する想定
const DUMMY_USER_ID = 'dummy-user-001'

export default function FeedRedirectPage() {
  redirect(`/feed/${DUMMY_USER_ID}`)
}
