import { auth, signIn, signOut } from "@/auth"

// サーバー側での認証状態取得
export async function getServerSession() {
  return await auth()
}

// トークン取得関数
// サーバー側とクライアント側の両方で動作する
export async function getAuthToken(): Promise<string> {
  // サーバー側での実行
  if (typeof window === 'undefined') {
    const session = await auth()
    if (session?.accessToken) {
      return session.accessToken
    }
    
    // 認証なしの場合はAPIキーを使用
    const apiKey = process.env.NEXT_PUBLIC_API_KEY
    if (!apiKey) {
      throw new Error('NEXT_PUBLIC_API_KEY is not set')
    }
    return apiKey
  }
  
  // クライアント側での実行
  const response = await fetch('/api/auth/token')
  if (!response.ok) {
    // 認証なしの場合はAPIキーを使用
    const apiKey = process.env.NEXT_PUBLIC_API_KEY
    if (!apiKey) {
      throw new Error('NEXT_PUBLIC_API_KEY is not set')
    }
    return apiKey
  }
  const data = await response.json()
  return data.accessToken
}

// クライアント側での認証状態取得用フック（ダミー実装で処理を差し込む場所を用意）
export function useAuth() {
  // TODO: クライアント側での認証状態取得を実装
  // 現時点ではダミー実装で処理を差し込む場所を用意
  return {
    user: null,
    isLoading: false,
    signIn: async () => {},
    signOut: async () => {},
  }
}
