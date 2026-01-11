import { User } from '@auth0/nextjs-auth0/types'
import { DmUser, CreateDmUserRequest, UpdateDmUserRequest } from '@/types/dm_user'
import { DmPost, CreateDmPostRequest, UpdateDmPostRequest, DmUserPost } from '@/types/dm_post'
import { RegisterJobRequest, RegisterJobResponse } from '@/types/jobqueue'
import { getAuthToken } from './auth'

type Auth0User = User

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
    auth0user?: Auth0User | undefined
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`

    // 認証トークンを取得
    const token = await getAuthToken(auth0user)

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

  // DmUser API
  async getDmUsers(limit = 20, offset = 0, auth0user?: Auth0User | undefined): Promise<DmUser[]> {
    return this.request<DmUser[]>(`/api/dm-users?limit=${limit}&offset=${offset}`, undefined, auth0user)
  }

  async getDmUser(id: string, auth0user?: Auth0User | undefined): Promise<DmUser> {
    return this.request<DmUser>(`/api/dm-users/${id}`, undefined, auth0user)
  }

  async createDmUser(data: CreateDmUserRequest, auth0user?: Auth0User | undefined): Promise<DmUser> {
    return this.request<DmUser>('/api/dm-users', {
      method: 'POST',
      body: JSON.stringify(data),
    }, auth0user)
  }

  async updateDmUser(id: string, data: UpdateDmUserRequest, auth0user?: Auth0User | undefined): Promise<DmUser> {
    return this.request<DmUser>(`/api/dm-users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }, auth0user)
  }

  async deleteDmUser(id: string, auth0user?: Auth0User | undefined): Promise<void> {
    return this.request<void>(`/api/dm-users/${id}`, {
      method: 'DELETE',
    }, auth0user)
  }

  // DmPost API
  async getDmPosts(limit = 20, offset = 0, userId?: string, auth0user?: Auth0User | undefined): Promise<DmPost[]> {
    const params = new URLSearchParams({
      limit: limit.toString(),
      offset: offset.toString(),
    })
    if (userId) {
      params.append('user_id', userId)
    }
    return this.request<DmPost[]>(`/api/dm-posts?${params.toString()}`, undefined, auth0user)
  }

  async getDmPost(id: string, userId: string, auth0user?: Auth0User | undefined): Promise<DmPost> {
    return this.request<DmPost>(`/api/dm-posts/${id}?user_id=${userId}`, undefined, auth0user)
  }

  async createDmPost(data: CreateDmPostRequest, auth0user?: Auth0User | undefined): Promise<DmPost> {
    return this.request<DmPost>('/api/dm-posts', {
      method: 'POST',
      body: JSON.stringify(data),
    }, auth0user)
  }

  async updateDmPost(id: string, userId: string, data: UpdateDmPostRequest, auth0user?: Auth0User | undefined): Promise<DmPost> {
    return this.request<DmPost>(`/api/dm-posts/${id}?user_id=${userId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    }, auth0user)
  }

  async deleteDmPost(id: string, userId: string, auth0user?: Auth0User | undefined): Promise<void> {
    return this.request<void>(`/api/dm-posts/${id}?user_id=${userId}`, {
      method: 'DELETE',
    }, auth0user)
  }

  // DmUser-DmPost JOIN API
  async getDmUserPosts(limit = 20, offset = 0, auth0user?: Auth0User | undefined): Promise<DmUserPost[]> {
    return this.request<DmUserPost[]>(`/api/dm-user-posts?limit=${limit}&offset=${offset}`, undefined, auth0user)
  }

  // Today API (private - requires Auth0 JWT)
  async getToday(auth0user?: Auth0User | undefined): Promise<{ date: string }> {
    return this.request<{ date: string }>('/api/today', undefined, auth0user)
  }

  // Email API
  async sendEmail(
    to: string[],
    template: string,
    data: Record<string, unknown>,
    auth0user?: Auth0User | undefined
  ): Promise<{ success: boolean; message: string }> {
    return this.request<{ success: boolean; message: string }>('/api/email/send', {
      method: 'POST',
      body: JSON.stringify({ to, template, data }),
    }, auth0user)
  }

  // JobQueue API (Demo)
  async registerJob(data?: RegisterJobRequest, auth0user?: Auth0User | undefined): Promise<RegisterJobResponse> {
    return this.request<RegisterJobResponse>('/api/dm-jobqueue/register', {
      method: 'POST',
      body: JSON.stringify(data || {}),
    }, auth0user)
  }

  // CSV Download API
  async downloadDmUsersCSV(auth0user?: Auth0User | undefined): Promise<void> {
    const url = `${this.baseURL}/api/export/dm-users/csv`
    const token = await getAuthToken(auth0user)

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
