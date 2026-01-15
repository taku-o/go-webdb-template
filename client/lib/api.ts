import { DmUser, CreateDmUserRequest, UpdateDmUserRequest } from '@/types/dm_user'
import { DmPost, CreateDmPostRequest, UpdateDmPostRequest, DmUserPost } from '@/types/dm_post'
import { DmFeedPost, DmFeedReply } from '@/types/dm_feed'
import { RegisterJobRequest, RegisterJobResponse } from '@/types/jobqueue'
import { getAuthToken } from './auth'

// フィード用メモリ内データストア
interface FeedDataStore {
  dmFeedPosts: Map<string, DmFeedPost>
  replies: Map<string, DmFeedReply[]>
  likes: Map<string, Set<string>>
}

// データストアのインスタンス（ページリロード時に初期化される）
const feedDataStore: FeedDataStore = {
  dmFeedPosts: new Map(),
  replies: new Map(),
  likes: new Map(),
}

// ランダムデータ生成用の定数
const RANDOM_USER_NAMES = [
  '田中太郎', '佐藤花子', '鈴木一郎', '高橋美咲', '伊藤健太',
  '渡辺さくら', '山本翔太', '中村優子', '小林大輔', '加藤真理',
  '吉田拓也', '山田美穂', '斎藤隼人', '松本理恵', '井上誠',
]

const RANDOM_USER_HANDLES = [
  'tanaka_taro', 'sato_hanako', 'suzuki_ichiro', 'takahashi_misaki', 'ito_kenta',
  'watanabe_sakura', 'yamamoto_shota', 'nakamura_yuko', 'kobayashi_daisuke', 'kato_mari',
  'yoshida_takuya', 'yamada_miho', 'saito_hayato', 'matsumoto_rie', 'inoue_makoto',
]

const RANDOM_POST_CONTENTS = [
  '今日はとても良い天気ですね！散歩日和です。',
  '新しいプロジェクトを始めました。楽しみです！',
  'コーヒーを飲みながらコーディング中。最高の時間。',
  '週末は映画を見に行く予定です。何を見ようかな。',
  '美味しいラーメン屋を発見しました！また行きたい。',
  'プログラミングの勉強は終わりがないですね。でも楽しい！',
  '今日のランチは手作り弁当。節約と健康のために。',
  '新しい本を読み始めました。とても面白いです。',
  '運動不足を解消するためにジムに通い始めました。',
  '友達と久しぶりに会えて嬉しかった！',
  'TypeScriptの型システムは奥が深いですね。',
  'リモートワークも慣れてきました。自分のペースで働けるのが良い。',
  '週末のハイキング計画中。天気が良いといいな。',
  '新しい技術を学ぶのはワクワクしますね！',
  '今日は早起きできた！朝活最高！',
]

const RANDOM_REPLY_CONTENTS = [
  'とても共感します！',
  '私も同じ経験があります。',
  '素晴らしいですね！',
  '参考になりました。ありがとうございます。',
  'いいですね！羨ましいです。',
  '頑張ってください！応援してます。',
  '面白そう！詳しく教えてください。',
  'それは良い考えですね。',
  '私もやってみたいです。',
  '楽しそうですね！',
]

// ランダムな要素を配列から取得
function getRandomElement<T>(array: T[]): T {
  return array[Math.floor(Math.random() * array.length)]
}

// ランダムなIDを生成
function generateRandomId(): string {
  return `${Date.now()}-${Math.random().toString(36).substring(2, 11)}`
}

// ランダムな過去の日時を生成（分単位で過去に遡る）
function generateRandomPastDate(minutesAgo: number): string {
  const now = new Date()
  const pastTime = now.getTime() - minutesAgo * 60 * 1000
  return new Date(pastTime).toISOString()
}

// ランダムな投稿データを生成
function generateRandomDmFeedPost(minutesAgo: number, currentUserId: string): DmFeedPost {
  const userIndex = Math.floor(Math.random() * RANDOM_USER_NAMES.length)
  const liked = Math.random() > 0.7
  return {
    id: generateRandomId(),
    userId: `user-${userIndex + 1}`,
    userName: RANDOM_USER_NAMES[userIndex],
    userHandle: `@${RANDOM_USER_HANDLES[userIndex]}`,
    content: getRandomElement(RANDOM_POST_CONTENTS),
    createdAt: generateRandomPastDate(minutesAgo),
    likeCount: Math.floor(Math.random() * 100),
    liked: liked,
  }
}

// ランダムな返信データを生成
function generateRandomDmFeedReply(postId: string, minutesAgo: number): DmFeedReply {
  const userIndex = Math.floor(Math.random() * RANDOM_USER_NAMES.length)
  return {
    id: generateRandomId(),
    postId: postId,
    userId: `user-${userIndex + 1}`,
    userName: RANDOM_USER_NAMES[userIndex],
    userHandle: `@${RANDOM_USER_HANDLES[userIndex]}`,
    content: getRandomElement(RANDOM_REPLY_CONTENTS),
    createdAt: generateRandomPastDate(minutesAgo),
  }
}

// 初期データの生成（投稿30件、各投稿に0〜5件の返信）
function initializeFeedData(currentUserId: string): void {
  if (feedDataStore.dmFeedPosts.size > 0) {
    return // 既に初期化済み
  }

  // 30件の投稿を生成
  for (let i = 0; i < 30; i++) {
    const minutesAgo = i * 15 + Math.floor(Math.random() * 10) // 15分間隔 + ランダム
    const post = generateRandomDmFeedPost(minutesAgo, currentUserId)
    feedDataStore.dmFeedPosts.set(post.id, post)

    // 各投稿に0〜5件の返信を生成
    const replyCount = Math.floor(Math.random() * 6)
    const replies: DmFeedReply[] = []
    for (let j = 0; j < replyCount; j++) {
      const replyMinutesAgo = minutesAgo - (j + 1) * 5 // 投稿より後の時刻
      const reply = generateRandomDmFeedReply(post.id, Math.max(0, replyMinutesAgo))
      replies.push(reply)
    }
    if (replies.length > 0) {
      feedDataStore.replies.set(post.id, replies)
    }

    // ランダムにいいねを設定
    if (post.liked) {
      const likeSet = new Set<string>()
      likeSet.add(currentUserId)
      feedDataStore.likes.set(post.id, likeSet)
    }
  }
}

// @ts-ignore - Uppyの型定義が正しく解決されない場合があるため
import Uppy from '@uppy/core'
// @ts-ignore
import Tus from '@uppy/tus'

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
    options?: RequestInit
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`

    // 認証トークンを取得
    const token = await getAuthToken()

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
  async getDmUsers(limit = 20, offset = 0): Promise<DmUser[]> {
    return this.request<DmUser[]>(`/api/dm-users?limit=${limit}&offset=${offset}`)
  }

  async getDmUser(id: string): Promise<DmUser> {
    return this.request<DmUser>(`/api/dm-users/${id}`)
  }

  async createDmUser(data: CreateDmUserRequest): Promise<DmUser> {
    return this.request<DmUser>('/api/dm-users', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updateDmUser(id: string, data: UpdateDmUserRequest): Promise<DmUser> {
    return this.request<DmUser>(`/api/dm-users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async deleteDmUser(id: string): Promise<void> {
    return this.request<void>(`/api/dm-users/${id}`, {
      method: 'DELETE',
    })
  }

  // DmPost API
  async getDmPosts(limit = 20, offset = 0, userId?: string): Promise<DmPost[]> {
    const params = new URLSearchParams({
      limit: limit.toString(),
      offset: offset.toString(),
    })
    if (userId) {
      params.append('user_id', userId)
    }
    return this.request<DmPost[]>(`/api/dm-posts?${params.toString()}`)
  }

  async getDmPost(id: string, userId: string): Promise<DmPost> {
    return this.request<DmPost>(`/api/dm-posts/${id}?user_id=${userId}`)
  }

  async createDmPost(data: CreateDmPostRequest): Promise<DmPost> {
    return this.request<DmPost>('/api/dm-posts', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  async updateDmPost(id: string, userId: string, data: UpdateDmPostRequest): Promise<DmPost> {
    return this.request<DmPost>(`/api/dm-posts/${id}?user_id=${userId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  async deleteDmPost(id: string, userId: string): Promise<void> {
    return this.request<void>(`/api/dm-posts/${id}?user_id=${userId}`, {
      method: 'DELETE',
    })
  }

  // DmUser-DmPost JOIN API
  async getDmUserPosts(limit = 20, offset = 0): Promise<DmUserPost[]> {
    return this.request<DmUserPost[]>(`/api/dm-user-posts?limit=${limit}&offset=${offset}`)
  }

  // Today API (private - requires JWT)
  async getToday(): Promise<{ date: string }> {
    return this.request<{ date: string }>('/api/today')
  }

  // Email API
  async sendEmail(
    to: string[],
    template: string,
    data: Record<string, unknown>
  ): Promise<{ success: boolean; message: string }> {
    return this.request<{ success: boolean; message: string }>('/api/email/send', {
      method: 'POST',
      body: JSON.stringify({ to, template, data }),
    })
  }

  // JobQueue API (Demo)
  async registerJob(data?: RegisterJobRequest): Promise<RegisterJobResponse> {
    return this.request<RegisterJobResponse>('/api/dm-jobqueue/register', {
      method: 'POST',
      body: JSON.stringify(data || {}),
    })
  }

  // CSV Download API
  async downloadDmUsersCSV(): Promise<void> {
    const url = `${this.baseURL}/api/export/dm-users/csv`
    const token = await getAuthToken()

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

  // DmFeed API (ダミーデータ、メモリ管理)
  async getDmFeedPosts(
    userId: string,
    limit: number,
    fromPostId: string
  ): Promise<DmFeedPost[]> {
    // 初期データの生成
    initializeFeedData(userId)

    // 全投稿を時系列順（新しいものが先頭）でソート
    const allPosts = Array.from(feedDataStore.dmFeedPosts.values()).sort(
      (a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
    )

    // fromPostIdが空文字列の場合は最新の投稿を取得
    if (!fromPostId) {
      return allPosts.slice(0, limit)
    }

    // fromPostIdが指定された場合はその投稿より古い投稿を取得
    const fromIndex = allPosts.findIndex((post) => post.id === fromPostId)
    if (fromIndex === -1) {
      return allPosts.slice(0, limit)
    }

    return allPosts.slice(fromIndex + 1, fromIndex + 1 + limit)
  }

  async getDmFeedPostById(userId: string, postId: string): Promise<DmFeedPost | null> {
    // 初期データの生成
    initializeFeedData(userId)

    // 指定されたIDの投稿を取得
    return feedDataStore.dmFeedPosts.get(postId) || null
  }

  async createDmFeedPost(userId: string, content: string): Promise<DmFeedPost> {
    // 初期データの生成
    initializeFeedData(userId)

    // 新規投稿を作成
    const newPost: DmFeedPost = {
      id: generateRandomId(),
      userId: userId,
      userName: 'あなた',
      userHandle: '@you',
      content: content,
      createdAt: new Date().toISOString(),
      likeCount: 0,
      liked: false,
    }

    // メモリに追加
    feedDataStore.dmFeedPosts.set(newPost.id, newPost)

    return newPost
  }

  async getDmFeedReplies(
    userId: string,
    postId: string,
    limit: number,
    fromReplyId: string
  ): Promise<DmFeedReply[]> {
    // 初期データの生成
    initializeFeedData(userId)

    // 指定された投稿の返信を取得
    const replies = feedDataStore.replies.get(postId) || []

    // 返信を時系列順（古いものが先頭、新しいものが末尾）でソート
    const sortedReplies = [...replies].sort(
      (a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
    )

    // fromReplyIdが空文字列の場合は最新の返信を取得（下方向スクロール用）
    if (!fromReplyId) {
      return sortedReplies.slice(-limit)
    }

    // fromReplyIdが指定された場合はその返信より新しい返信を取得（上方向スクロール用）
    const fromIndex = sortedReplies.findIndex((reply) => reply.id === fromReplyId)
    if (fromIndex === -1) {
      return sortedReplies.slice(-limit)
    }

    // 指定された返信より古い返信を取得
    return sortedReplies.slice(Math.max(0, fromIndex - limit), fromIndex)
  }

  async replyToDmPost(
    userId: string,
    postId: string,
    content: string
  ): Promise<DmFeedReply> {
    // 初期データの生成
    initializeFeedData(userId)

    // 新規返信を作成
    const newReply: DmFeedReply = {
      id: generateRandomId(),
      postId: postId,
      userId: userId,
      userName: 'あなた',
      userHandle: '@you',
      content: content,
      createdAt: new Date().toISOString(),
    }

    // メモリに追加
    const existingReplies = feedDataStore.replies.get(postId) || []
    existingReplies.push(newReply)
    feedDataStore.replies.set(postId, existingReplies)

    return newReply
  }

  async toggleLikeDmPost(
    userId: string,
    postId: string
  ): Promise<{ liked: boolean; likeCount: number }> {
    // 初期データの生成
    initializeFeedData(userId)

    // 投稿を取得
    const post = feedDataStore.dmFeedPosts.get(postId)
    if (!post) {
      throw new Error('Post not found')
    }

    // いいね状態を取得または初期化
    let likeSet = feedDataStore.likes.get(postId)
    if (!likeSet) {
      likeSet = new Set<string>()
      feedDataStore.likes.set(postId, likeSet)
    }

    // いいねのトグル
    const wasLiked = likeSet.has(userId)
    if (wasLiked) {
      likeSet.delete(userId)
      post.likeCount = Math.max(0, post.likeCount - 1)
      post.liked = false
    } else {
      likeSet.add(userId)
      post.likeCount = post.likeCount + 1
      post.liked = true
    }

    // 投稿を更新
    feedDataStore.dmFeedPosts.set(postId, post)

    return {
      liked: post.liked,
      likeCount: post.likeCount,
    }
  }
}

export const apiClient = new ApiClient(API_BASE_URL)

// 動画アップロード用のUppyインスタンスを作成する関数
export interface MovieUploadCallbacks {
  onUploadProgress?: (percent: number) => void
  onUploadSuccess?: () => void
  onUploadError?: (error: string) => void
  onUploadStart?: () => void
}

export function createMovieUploader(callbacks: MovieUploadCallbacks = {}): Uppy {
  const {
    onUploadProgress,
    onUploadSuccess,
    onUploadError,
    onUploadStart,
  } = callbacks

  const uppyInstance = new Uppy({
    id: 'dm_movie_uploader',
    autoProceed: false,
    restrictions: {
      maxFileSize: 2147483648, // 2GB
      allowedFileTypes: ['.mp4'],
      maxNumberOfFiles: 1,
    },
  })
    .use(Tus, {
      endpoint: `${API_BASE_URL}/api/upload/dm_movie`,
      chunkSize: 5 * 1024 * 1024, // 5MB
      retryDelays: [0, 1000, 3000, 5000],
      onBeforeRequest: async (req: any) => {
        const token = await getAuthToken()
        req.setHeader('Authorization', `Bearer ${token}`)
      },
    })
    .on('upload-progress', (file: any, progress: any) => {
      if (progress.bytesTotal && onUploadProgress) {
        const percent = Math.round((progress.bytesUploaded / progress.bytesTotal) * 100)
        onUploadProgress(percent)
      }
    })
    .on('upload-success', (file: any, response: any) => {
      if (onUploadSuccess) {
        onUploadSuccess()
      }
    })
    .on('upload-error', (file: any, error: any, response: any) => {
      if (onUploadError) {
        onUploadError(error?.message || 'アップロードに失敗しました')
      }
    })
    .on('upload', () => {
      if (onUploadStart) {
        onUploadStart()
      }
    })

  return uppyInstance
}
