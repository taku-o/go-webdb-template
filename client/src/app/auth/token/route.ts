import { NextResponse } from "next/server"
import { auth0 } from "@/lib/auth0"

export async function GET() {
  try {
    const tokenResponse = await auth0.getAccessToken()
    if (!tokenResponse?.token) {
      return NextResponse.json(
        { error: "No access token available" },
        { status: 401 }
      )
    }
    return NextResponse.json({ accessToken: tokenResponse.token })
  } catch (error) {
    console.error("Failed to get access token:", error)
    return NextResponse.json(
      { error: "Failed to get access token" },
      { status: 500 }
    )
  }
}
