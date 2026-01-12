"use client"

import { LoadingSpinner } from "./loading-spinner"

interface LoadingOverlayProps {
  message?: string
}

export function LoadingOverlay({ message = "Loading..." }: LoadingOverlayProps) {
  return (
    <div className="flex flex-col items-center justify-center gap-4 p-8">
      <LoadingSpinner size="lg" />
      <p className="text-sm text-gray-600">{message}</p>
    </div>
  )
}
