import { NextResponse } from "next/server"
import { auth } from "@/auth"

/**
 * カスタムプロフィールエンドポイント
 * 未ログイン時も200を返し、ブラウザコンソールの401エラーを防ぐ
 */
export async function GET() {
  try {
    const session = await auth()
    if (!session?.user) {
      // 未ログイン時は200でnullを返す（401ではなく）
      return NextResponse.json(null)
    }
    return NextResponse.json(session.user)
  } catch {
    // エラー時も200でnullを返す
    return NextResponse.json(null)
  }
}
