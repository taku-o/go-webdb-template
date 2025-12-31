import { User, CreateUserRequest, UpdateUserRequest } from '@/types/user'
import { Post, CreatePostRequest, UpdatePostRequest, UserPost } from '@/types/post'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'

class ApiClient {
  private baseURL: string
  private apiKey: string | null

  constructor(baseURL: string) {
    this.baseURL = baseURL
    this.apiKey = process.env.NEXT_PUBLIC_API_KEY || null

    // APIキーが設定されていない場合、エラーを投げる
    if (!this.apiKey) {
      throw new Error('NEXT_PUBLIC_API_KEY is not set')
    }
  }

  private async request<T>(
    endpoint: string,
    options?: RequestInit,
    jwt?: string
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`

    // JWTの取得（引数で渡された場合はそれを使用、なければapiKeyを使用）
    const token = jwt || this.apiKey

    // Authorizationヘッダーを追加
    const headers = {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
      ...options?.headers,
    }

    const response = await fetch(url, {
      ...options,
      headers,
    })

    if (!response.ok) {
      // エラーレスポンスの処理
      if (response.status === 401 || response.status === 403) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.message || response.statusText)
      }
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
    return this.request<User[]>(`/api/dm-users?limit=${limit}&offset=${offset}`)
  }

  async getUser(id: string): Promise<User> {
    return this.request<User>(`/api/dm-users/${id}`)
  }

  async createUser(data: CreateUserRequest): Promise<User> {
    return this.request<User>('/api/dm-users', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updateUser(id: string, data: UpdateUserRequest): Promise<User> {
    return this.request<User>(`/api/dm-users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async deleteUser(id: string): Promise<void> {
    return this.request<void>(`/api/dm-users/${id}`, {
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
    return this.request<Post[]>(`/api/dm-posts?${params.toString()}`)
  }

  async getPost(id: string, userId: string): Promise<Post> {
    return this.request<Post>(`/api/dm-posts/${id}?user_id=${userId}`)
  }

  async createPost(data: CreatePostRequest): Promise<Post> {
    return this.request<Post>('/api/dm-posts', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updatePost(id: string, userId: string, data: UpdatePostRequest): Promise<Post> {
    return this.request<Post>(`/api/dm-posts/${id}?user_id=${userId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async deletePost(id: string, userId: string): Promise<void> {
    return this.request<void>(`/api/dm-posts/${id}?user_id=${userId}`, {
      method: 'DELETE',
    })
  }

  // User-Post JOIN API
  async getUserPosts(limit = 20, offset = 0): Promise<UserPost[]> {
    return this.request<UserPost[]>(`/api/dm-user-posts?limit=${limit}&offset=${offset}`)
  }

  // Today API (private - requires Auth0 JWT)
  async getToday(jwt: string): Promise<{ date: string }> {
    return this.request<{ date: string }>('/api/today', undefined, jwt)
  }

  // CSV Download API
  async downloadUsersCSV(): Promise<void> {
    const url = `${this.baseURL}/api/export/dm-users/csv`
    const token = this.apiKey

    const response = await fetch(url, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })

    if (!response.ok) {
      if (response.status === 401 || response.status === 403) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.message || response.statusText)
      }
      const error = await response.text()
      throw new Error(error || response.statusText)
    }

    // Blobとして取得
    const blob = await response.blob()

    // Content-Dispositionヘッダーからファイル名を取得
    const contentDisposition = response.headers.get('Content-Disposition')
    let filename = 'dm-users.csv'
    if (contentDisposition) {
      const filenameMatch = contentDisposition.match(/filename="?(.+?)"?$/i)
      if (filenameMatch) {
        filename = filenameMatch[1]
      }
    }

    // Blob URLを生成
    const blobUrl = URL.createObjectURL(blob)

    // <a>要素を作成してダウンロード
    const link = document.createElement('a')
    link.href = blobUrl
    link.download = filename
    document.body.appendChild(link)
    link.click()

    // クリーンアップ
    document.body.removeChild(link)
    URL.revokeObjectURL(blobUrl)
  }
}

export const apiClient = new ApiClient(API_BASE_URL)
