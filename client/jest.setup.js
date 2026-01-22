// Learn more: https://github.com/testing-library/jest-dom
import '@testing-library/jest-dom'

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

// Mock NextAuth
jest.mock('next-auth/react', () => ({
  useSession: jest.fn(() => ({
    data: null,
    status: 'unauthenticated',
  })),
  signIn: jest.fn(),
  signOut: jest.fn(),
}))

jest.mock('next-auth', () => {
  const mockAuth = jest.fn(() => null)
  mockAuth.signIn = jest.fn()
  mockAuth.signOut = jest.fn()
  return {
    default: jest.fn(() => ({
      handlers: {
        GET: {},
        POST: {},
      },
      auth: mockAuth,
      signIn: mockAuth.signIn,
      signOut: mockAuth.signOut,
    })),
    auth: mockAuth,
    signIn: mockAuth.signIn,
    signOut: mockAuth.signOut,
  }
})

// Mock auth.ts module (must be before any imports that use it)
jest.mock('@/auth', () => {
  const mockAuth = jest.fn(() => Promise.resolve(null))
  return {
    handlers: {
      GET: {},
      POST: {},
    },
    auth: mockAuth,
    signIn: jest.fn(),
    signOut: jest.fn(),
  }
})

// Mock Uppy (ESM module, not compatible with Jest)
jest.mock('@uppy/core', () => {
  const createMockUppyInstance = () => {
    const mockInstance = {
      use: jest.fn(() => mockInstance),
      on: jest.fn(() => mockInstance),
      off: jest.fn(() => mockInstance),
      upload: jest.fn(() => Promise.resolve()),
      addFile: jest.fn(),
      removeFile: jest.fn(),
      getFiles: jest.fn(() => []),
      close: jest.fn(),
      destroy: jest.fn(),
    }
    return mockInstance
  }
  return {
    __esModule: true,
    default: jest.fn(() => createMockUppyInstance()),
  }
})

jest.mock('@uppy/tus', () => ({
  __esModule: true,
  default: jest.fn(),
}))

jest.mock('@uppy/react/dashboard', () => ({
  __esModule: true,
  default: jest.fn(() => null),
}))

// Mock environment variables
process.env.NEXT_PUBLIC_API_BASE_URL = 'http://localhost:8080'
process.env.NEXT_PUBLIC_API_KEY = 'test-api-key'
process.env.AUTH_SECRET = 'test-auth-secret'
process.env.AUTH_URL = 'http://localhost:3000'
process.env.APP_ENV = 'test'
