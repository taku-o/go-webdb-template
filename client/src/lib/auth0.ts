import { Auth0Client } from "@auth0/nextjs-auth0/server"

export const auth0 = new Auth0Client({
  authorizationParameters: {
    audience: process.env.AUTH0_AUDIENCE,
  },
})

/**
 * Server ComponentsでJWTを取得する
 * 注意: 本実装ではJWT取得機能を提供するが、API呼び出しでの使用は次のissueで対応
 */
export async function getAccessToken(): Promise<{ accessToken: string | undefined }> {
  try {
    const tokenResponse = await auth0.getAccessToken()
    return { accessToken: tokenResponse?.token }
  } catch (error) {
    console.error('Failed to get access token:', error)
    return { accessToken: undefined }
  }
}

/**
 * Server Componentsでセッション情報を取得する
 */
export async function getSession() {
  return auth0.getSession()
}
