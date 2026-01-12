import NextAuth from "next-auth"

export const {
  handlers,
  auth,
  signIn,
  signOut,
} = NextAuth({
  // 認証設定
  providers: [
    // 必要最小限のプロバイダーを設定
    // 現時点ではプロバイダーは設定しない（環境変数設定後に追加）
  ],
})
