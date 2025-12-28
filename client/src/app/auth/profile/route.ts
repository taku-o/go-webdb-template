import { NextResponse } from "next/server"
import { auth0 } from "@/lib/auth0"

/**
 * カスタムプロフィールエンドポイント
 * 未ログイン時も200を返し、ブラウザコンソールの401エラーを防ぐ
 */
export async function GET() {
  try {
    const session = await auth0.getSession()
    if (!session) {
      // 未ログイン時は200でnullを返す（401ではなく）
      return NextResponse.json(null)
    }
    return NextResponse.json(session.user)
  } catch {
    // エラー時も200でnullを返す
    return NextResponse.json(null)
  }
}
