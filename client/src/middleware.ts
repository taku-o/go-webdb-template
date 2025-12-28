import type { NextRequest } from "next/server"
import { auth0 } from "./lib/auth0"

export async function middleware(request: NextRequest) {
  return await auth0.middleware(request)
}

export const config = {
  matcher: [
    // Match all paths except static files, Next.js internals, and custom auth/profile route
    "/((?!_next/static|_next/image|favicon.ico|auth/profile).*)",
  ],
}
