import { User, CreateUserRequest, UpdateUserRequest } from '@/types/user'
import { Post, CreatePostRequest, UpdatePostRequest, UserPost } from '@/types/post'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'

class ApiClient {
  private baseURL: string

  constructor(baseURL: string) {
    this.baseURL = baseURL
  }

  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    })

    if (!response.ok) {
      const error = await response.text()
      throw new Error(error || response.statusText)
    }

    if (response.status === 204) {
      return {} as T
    }

    return response.json()
  }

  // User API
  async getUsers(limit = 20, offset = 0): Promise<User[]> {
    return this.request<User[]>(`/api/users?limit=${limit}&offset=${offset}`)
  }

  async getUser(id: string): Promise<User> {
    return this.request<User>(`/api/users/${id}`)
  }

  async createUser(data: CreateUserRequest): Promise<User> {
    return this.request<User>('/api/users', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updateUser(id: string, data: UpdateUserRequest): Promise<User> {
    return this.request<User>(`/api/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async deleteUser(id: string): Promise<void> {
    return this.request<void>(`/api/users/${id}`, {
      method: 'DELETE',
    })
  }

  // Post API
  async getPosts(limit = 20, offset = 0, userId?: string): Promise<Post[]> {
    const params = new URLSearchParams({
      limit: limit.toString(),
      offset: offset.toString(),
    })
    if (userId) {
      params.append('user_id', userId)
    }
    return this.request<Post[]>(`/api/posts?${params.toString()}`)
  }

  async getPost(id: string, userId: string): Promise<Post> {
    return this.request<Post>(`/api/posts/${id}?user_id=${userId}`)
  }

  async createPost(data: CreatePostRequest): Promise<Post> {
    return this.request<Post>('/api/posts', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updatePost(id: string, userId: string, data: UpdatePostRequest): Promise<Post> {
    return this.request<Post>(`/api/posts/${id}?user_id=${userId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async deletePost(id: string, userId: string): Promise<void> {
    return this.request<void>(`/api/posts/${id}?user_id=${userId}`, {
      method: 'DELETE',
    })
  }

  // User-Post JOIN API
  async getUserPosts(limit = 20, offset = 0): Promise<UserPost[]> {
    return this.request<UserPost[]>(`/api/user-posts?limit=${limit}&offset=${offset}`)
  }
}

export const apiClient = new ApiClient(API_BASE_URL)
