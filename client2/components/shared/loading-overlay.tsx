"use client"

import { LoadingSpinner } from "./loading-spinner"

interface LoadingOverlayProps {
  message?: string
  "aria-label"?: string
}

export function LoadingOverlay({ message = "Loading...", "aria-label": ariaLabel }: LoadingOverlayProps) {
  return (
    <div 
      className="flex flex-col items-center justify-center gap-4 p-8"
      role="status"
      aria-live="polite"
      aria-label={ariaLabel}
    >
      <LoadingSpinner size="lg" />
      <p className="text-sm text-muted-foreground">{message}</p>
    </div>
  )
}
