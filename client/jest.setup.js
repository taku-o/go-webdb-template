// Learn more: https://github.com/testing-library/jest-dom
import '@testing-library/jest-dom'

// Polyfill for TextEncoder/TextDecoder (required for MSW in Node environment)
import { TextEncoder, TextDecoder } from 'util'
global.TextEncoder = TextEncoder
global.TextDecoder = TextDecoder

// Mock Next.js router
jest.mock('next/navigation', () => ({
  useRouter() {
    return {
      push: jest.fn(),
      replace: jest.fn(),
      prefetch: jest.fn(),
      back: jest.fn(),
    }
  },
  usePathname() {
    return '/'
  },
  useSearchParams() {
    return new URLSearchParams()
  },
}))

// Mock environment variables
process.env.NEXT_PUBLIC_API_BASE_URL = 'http://localhost:8080'
