import { User } from '@auth0/nextjs-auth0'

type Auth0User = User

export async function getAuthToken(auth0user: Auth0User | undefined): Promise<string> {
  if (auth0user) {
    const response = await fetch('/auth/token')
    if (!response.ok) {
      throw new Error('Failed to get access token')
    }
    const data = await response.json()
    return data.accessToken
  } else {
    const apiKey = process.env.NEXT_PUBLIC_API_KEY
    if (!apiKey) {
      throw new Error('NEXT_PUBLIC_API_KEY is not set')
    }
    return apiKey
  }
}
