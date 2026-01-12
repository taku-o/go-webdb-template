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
    // 既存のAuth0設定を参考に、NextAuth (Auth.js) v5のプロバイダーを設定
    // 環境変数設定後に追加（例: Auth0プロバイダー）
  ],
  callbacks: {
    async session({ session, token }) {
      // セッションにアクセストークンを含める
      if (token.accessToken) {
        session.accessToken = token.accessToken as string
      }
      return session
    },
    async jwt({ token, account }) {
      // アカウント情報からアクセストークンを取得
      if (account?.access_token) {
        token.accessToken = account.access_token
      }
      return token
    },
  },
})
