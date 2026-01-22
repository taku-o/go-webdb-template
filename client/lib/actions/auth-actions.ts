"use server"

import { signIn, signOut } from '@/auth'
import { redirect } from "next/navigation"

export async function signInAction() {
  await signIn('auth0')
}

export async function signOutAction() {
  // 環境変数の取得
  const auth0Issuer = process.env.AUTH0_ISSUER
  const auth0ClientId = process.env.AUTH0_CLIENT_ID
  const appBaseUrl = process.env.NEXT_PUBLIC_APP_BASE_URL

  // 環境変数の検証
  if (!auth0Issuer) {
    throw new Error('AUTH0_ISSUER is not set')
  }
  if (!auth0ClientId) {
    throw new Error('AUTH0_CLIENT_ID is not set')
  }
  if (!appBaseUrl) {
    throw new Error('NEXT_PUBLIC_APP_BASE_URL is not set')
  }

  // AUTH0_ISSUERから/v2/logoutを追加してAuth0ログアウトURLを構築
  const auth0LogoutUrl = `${auth0Issuer}/v2/logout`

  // Auth0ログアウトURLにreturnToパラメータを追加
  const returnToUrl = `${appBaseUrl}`
  const logoutUrl = `${auth0LogoutUrl}?client_id=${auth0ClientId}&returnTo=${encodeURIComponent(returnToUrl)}`

  // next-authのsignOutにパラメータを渡して
  await signOut({
    redirect: false,
    redirectTo: logoutUrl
  })
  // その後、リダイレクト
  redirect(logoutUrl);
}
