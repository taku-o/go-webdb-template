import { NextResponse } from "next/server"
import { auth } from "@/auth"

export async function GET() {
  try {
    const session = await auth()
    if (!session?.accessToken) {
      return NextResponse.json(
        { error: "No access token available" },
        { status: 401 }
      )
    }
    return NextResponse.json({ accessToken: session.accessToken })
  } catch (error) {
    console.error("Failed to get access token:", error)
    return NextResponse.json(
      { error: "Failed to get access token" },
      { status: 500 }
    )
  }
}
