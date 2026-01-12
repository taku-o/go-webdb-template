import NextAuth from "next-auth"
import Auth0 from "next-auth/providers/auth0"

export const {
  handlers,
  auth,
  signIn,
  signOut,
} = NextAuth({
  // 認証設定
  providers: [
    Auth0({
      clientId: process.env.AUTH0_CLIENT_ID,
      clientSecret: process.env.AUTH0_CLIENT_SECRET,
      issuer: process.env.AUTH0_ISSUER,
      authorization: {
        params: {
          audience: process.env.AUTH0_AUDIENCE,
        },
      },
    }),
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
