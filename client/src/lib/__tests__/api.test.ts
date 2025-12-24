import { apiClient } from '../api'
import { User, CreateUserRequest } from '@/types/user'
import { Post, CreatePostRequest } from '@/types/post'

// Mock fetch
global.fetch = jest.fn()

describe('apiClient', () => {
  beforeEach(() => {
    ;(fetch as jest.Mock).mockClear()
  })

  describe('User API', () => {
    it('creates a user', async () => {
      const mockUser: User = {
        id: 1,
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
        'http://localhost:8080/api/users',
        expect.objectContaining({
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(createRequest),
        })
      )
    })

    it('gets all users', async () => {
      const mockUsers: User[] = [
        {
          id: 1,
          name: 'User 1',
          email: 'user1@example.com',
          created_at: '2024-01-15T10:00:00Z',
          updated_at: '2024-01-15T10:00:00Z',
        },
        {
          id: 2,
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
        'http://localhost:8080/api/users?limit=20&offset=0',
        expect.objectContaining({
          headers: { 'Content-Type': 'application/json' },
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

      await apiClient.deleteUser(1)

      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/users/1',
        expect.objectContaining({
          method: 'DELETE',
        })
      )
    })
  })

  describe('Post API', () => {
    it('creates a post', async () => {
      const mockPost: Post = {
        id: 1,
        user_id: 1,
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
        user_id: 1,
        title: 'Test Post',
        content: 'Test content',
      }

      const result = await apiClient.createPost(createRequest)

      expect(result).toEqual(mockPost)
      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/posts',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(createRequest),
        })
      )
    })

    it('gets all posts', async () => {
      const mockPosts: Post[] = [
        {
          id: 1,
          user_id: 1,
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

      await apiClient.deletePost(1, 1)

      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/posts/1?user_id=1',
        expect.objectContaining({
          method: 'DELETE',
        })
      )
    })
  })
})
