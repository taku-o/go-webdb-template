// 投稿の型
export interface DmFeedPost {
  id: string
  userId: string
  userName: string
  userHandle: string
  content: string
  createdAt: string // ISO 8601形式
  likeCount: number
  liked: boolean
}

// 返信の型
export interface DmFeedReply {
  id: string
  postId: string
  userId: string
  userName: string
  userHandle: string
  content: string
  createdAt: string // ISO 8601形式
}

// 新規投稿リクエストの型
export interface CreateDmFeedPostRequest {
  content: string
}

// 返信リクエストの型
export interface CreateDmFeedReplyRequest {
  content: string
}
