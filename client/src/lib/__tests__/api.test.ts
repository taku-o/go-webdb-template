import { User, CreateUserRequest } from '@/types/user'
import { Post, CreatePostRequest } from '@/types/post'
import { apiClient } from '../api'

// テスト用APIキー（jest.setup.jsで設定済み）
const TEST_API_KEY = 'test-api-key'

// Mock fetch
global.fetch = jest.fn()

describe('apiClient', () => {
  beforeEach(() => {
    ;(fetch as jest.Mock).mockClear()
  })

  describe('Authorization header', () => {
    it('includes Authorization header in all requests', async () => {
      ;(fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => [],
      })

      await apiClient.getUsers()

      expect(fetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.objectContaining({
            Authorization: `Bearer ${TEST_API_KEY}`,
          }),
        })
      )
    })
  })

  describe('Error handling for 401/403', () => {
    it('handles 401 Unauthorized error with message from response', async () => {
      ;(fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        status: 401,
        statusText: 'Unauthorized',
        json: async () => ({ code: 401, message: 'Invalid API key' }),
      })

      await expect(apiClient.getUsers()).rejects.toThrow('Invalid API key')
    })

    it('handles 403 Forbidden error with message from response', async () => {
      ;(fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        status: 403,
        statusText: 'Forbidden',
        json: async () => ({ code: 403, message: 'Insufficient scope' }),
      })

      await expect(apiClient.getUsers()).rejects.toThrow('Insufficient scope')
    })

    it('uses statusText when json parsing fails for 401', async () => {
      ;(fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        status: 401,
        statusText: 'Unauthorized',
        json: async () => {
          throw new Error('JSON parse error')
        },
      })

      await expect(apiClient.getUsers()).rejects.toThrow('Unauthorized')
    })
  })

  describe('User API', () => {
    it('creates a user', async () => {
      const mockUser: User = {
        id: '1',
        name: 'Test User',
        email: 'test@example.com',
        created_at: '2024-01-15T10:00:00Z',
        updated_at: '2024-01-15T10:00:00Z',
      }

      ;(fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockUser,
      })

      const createRequest: CreateUserRequest = {
        name: 'Test User',
        email: 'test@example.com',
      }

      const result = await apiClient.createUser(createRequest)

      expect(result).toEqual(mockUser)
      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/dm-users',
        expect.objectContaining({
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${TEST_API_KEY}`,
          },
          body: JSON.stringify(createRequest),
        })
      )
    })

    it('gets all users', async () => {
      const mockUsers: User[] = [
        {
          id: '1',
          name: 'User 1',
          email: 'user1@example.com',
          created_at: '2024-01-15T10:00:00Z',
          updated_at: '2024-01-15T10:00:00Z',
        },
        {
          id: '2',
          name: 'User 2',
          email: 'user2@example.com',
          created_at: '2024-01-15T11:00:00Z',
          updated_at: '2024-01-15T11:00:00Z',
        },
      ]

      ;(fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockUsers,
      })

      const result = await apiClient.getUsers()

      expect(result).toEqual(mockUsers)
      // getUsers has default parameters limit=20 and offset=0
      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/dm-users?limit=20&offset=0',
        expect.objectContaining({
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${TEST_API_KEY}`,
          },
        })
      )
    })

    it('handles API errors', async () => {
      ;(fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error',
        text: async () => 'Internal server error',
      })

      await expect(apiClient.getUsers()).rejects.toThrow('Internal server error')
    })

    it('deletes a user', async () => {
      ;(fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        status: 204,
      })

      await apiClient.deleteUser('1')

      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/dm-users/1',
        expect.objectContaining({
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${TEST_API_KEY}`,
          },
        })
      )
    })
  })

  describe('Post API', () => {
    it('creates a post', async () => {
      const mockPost: Post = {
        id: '1',
        user_id: '1',
        title: 'Test Post',
        content: 'Test content',
        created_at: '2024-01-15T12:00:00Z',
        updated_at: '2024-01-15T12:00:00Z',
      }

      ;(fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockPost,
      })

      const createRequest: CreatePostRequest = {
        user_id: '1',
        title: 'Test Post',
        content: 'Test content',
      }

      const result = await apiClient.createPost(createRequest)

      expect(result).toEqual(mockPost)
      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/dm-posts',
        expect.objectContaining({
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${TEST_API_KEY}`,
          },
          body: JSON.stringify(createRequest),
        })
      )
    })

    it('gets all posts', async () => {
      const mockPosts: Post[] = [
        {
          id: '1',
          user_id: '1',
          title: 'Post 1',
          content: 'Content 1',
          created_at: '2024-01-15T12:00:00Z',
          updated_at: '2024-01-15T12:00:00Z',
        },
      ]

      ;(fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockPosts,
      })

      const result = await apiClient.getPosts()

      expect(result).toEqual(mockPosts)
    })

    it('deletes a post', async () => {
      ;(fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        status: 204,
      })

      await apiClient.deletePost('1', '1')

      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/dm-posts/1?user_id=1',
        expect.objectContaining({
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${TEST_API_KEY}`,
          },
        })
      )
    })
  })

  describe('CSV Download API', () => {
    // DOMのモックを設定
    let mockCreateObjectURL: jest.Mock
    let mockRevokeObjectURL: jest.Mock
    let mockAppendChild: jest.Mock
    let mockRemoveChild: jest.Mock
    let mockClick: jest.Mock
    let mockLink: HTMLAnchorElement

    beforeEach(() => {
      // URL.createObjectURL/revokeObjectURLのモック
      mockCreateObjectURL = jest.fn(() => 'blob:http://localhost/mock-blob-url')
      mockRevokeObjectURL = jest.fn()
      global.URL.createObjectURL = mockCreateObjectURL
      global.URL.revokeObjectURL = mockRevokeObjectURL

      // document.createElementのモック
      mockClick = jest.fn()
      mockLink = {
        href: '',
        download: '',
        click: mockClick,
      } as unknown as HTMLAnchorElement

      jest.spyOn(document, 'createElement').mockReturnValue(mockLink)
      mockAppendChild = jest.spyOn(document.body, 'appendChild').mockImplementation(() => mockLink)
      mockRemoveChild = jest.spyOn(document.body, 'removeChild').mockImplementation(() => mockLink)
    })

    afterEach(() => {
      jest.restoreAllMocks()
    })

    it('downloads CSV file successfully', async () => {
      const mockBlob = new Blob(['ID,Name,Email\n1,Test,test@example.com'], { type: 'text/csv' })

      ;(fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        headers: new Headers({
          'Content-Disposition': 'attachment; filename="dm-users.csv"',
        }),
        blob: async () => mockBlob,
      })

      await apiClient.downloadUsersCSV()

      // fetchが正しく呼ばれたか確認
      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/export/dm-users/csv',
        expect.objectContaining({
          method: 'GET',
          headers: expect.objectContaining({
            Authorization: `Bearer ${TEST_API_KEY}`,
          }),
        })
      )

      // Blob URLが生成されたか確認
      expect(mockCreateObjectURL).toHaveBeenCalledWith(mockBlob)

      // リンクが正しく設定されたか確認
      expect(mockLink.href).toBe('blob:http://localhost/mock-blob-url')
      expect(mockLink.download).toBe('dm-users.csv')

      // クリックが実行されたか確認
      expect(mockClick).toHaveBeenCalled()

      // クリーンアップが行われたか確認
      expect(mockRemoveChild).toHaveBeenCalledWith(mockLink)
      expect(mockRevokeObjectURL).toHaveBeenCalledWith('blob:http://localhost/mock-blob-url')
    })

    it('uses default filename when Content-Disposition is not present', async () => {
      const mockBlob = new Blob(['ID,Name,Email\n1,Test,test@example.com'], { type: 'text/csv' })

      ;(fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        headers: new Headers({}),
        blob: async () => mockBlob,
      })

      await apiClient.downloadUsersCSV()

      // デフォルトのファイル名が使用されているか確認
      expect(mockLink.download).toBe('dm-users.csv')
    })

    it('handles 401 error correctly', async () => {
      ;(fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        status: 401,
        statusText: 'Unauthorized',
        json: async () => ({ message: 'Invalid token' }),
      })

      await expect(apiClient.downloadUsersCSV()).rejects.toThrow('Invalid token')
    })

    it('handles 403 error correctly', async () => {
      ;(fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        status: 403,
        statusText: 'Forbidden',
        json: async () => ({ message: 'Access denied' }),
      })

      await expect(apiClient.downloadUsersCSV()).rejects.toThrow('Access denied')
    })

    it('handles 500 error correctly', async () => {
      ;(fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error',
        text: async () => 'Server error occurred',
      })

      await expect(apiClient.downloadUsersCSV()).rejects.toThrow('Server error occurred')
    })
  })

  describe('Email API', () => {
    it('sends email successfully', async () => {
      const mockResponse = { success: true, message: 'メールを送信しました' }

      ;(fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      })

      const result = await apiClient.sendEmail(
        ['test@example.com'],
        'welcome',
        { Name: 'テスト太郎', Email: 'test@example.com' }
      )

      expect(result).toEqual(mockResponse)
      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/email/send',
        expect.objectContaining({
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${TEST_API_KEY}`,
          },
          body: JSON.stringify({
            to: ['test@example.com'],
            template: 'welcome',
            data: { Name: 'テスト太郎', Email: 'test@example.com' },
          }),
        })
      )
    })

    it('sends email to multiple recipients', async () => {
      const mockResponse = { success: true, message: 'メールを送信しました' }

      ;(fetch as jest.Mock).mockResolvedValueOnce({
        ok: true,
        json: async () => mockResponse,
      })

      const result = await apiClient.sendEmail(
        ['user1@example.com', 'user2@example.com'],
        'welcome',
        { Name: 'ユーザーの皆様', Email: 'users@example.com' }
      )

      expect(result).toEqual(mockResponse)
      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/email/send',
        expect.objectContaining({
          body: JSON.stringify({
            to: ['user1@example.com', 'user2@example.com'],
            template: 'welcome',
            data: { Name: 'ユーザーの皆様', Email: 'users@example.com' },
          }),
        })
      )
    })

    it('handles 400 Bad Request error', async () => {
      ;(fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        status: 400,
        statusText: 'Bad Request',
        text: async () => 'invalid email address',
      })

      await expect(
        apiClient.sendEmail(['invalid'], 'welcome', { Name: 'Test', Email: 'invalid' })
      ).rejects.toThrow('invalid email address')
    })

    it('handles 401 Unauthorized error', async () => {
      ;(fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        status: 401,
        statusText: 'Unauthorized',
        json: async () => ({ message: 'Invalid API key' }),
      })

      await expect(
        apiClient.sendEmail(['test@example.com'], 'welcome', { Name: 'Test', Email: 'test@example.com' })
      ).rejects.toThrow('Invalid API key')
    })

    it('handles 500 Internal Server Error', async () => {
      ;(fetch as jest.Mock).mockResolvedValueOnce({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error',
        text: async () => 'Failed to send email',
      })

      await expect(
        apiClient.sendEmail(['test@example.com'], 'welcome', { Name: 'Test', Email: 'test@example.com' })
      ).rejects.toThrow('Failed to send email')
    })
  })
})
