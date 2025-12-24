export interface Post {
  id: number
  user_id: number
  title: string
  content: string
  created_at: string
  updated_at: string
}

export interface CreatePostRequest {
  user_id: number
  title: string
  content: string
}

export interface UpdatePostRequest {
  title?: string
  content?: string
}

export interface UserPost {
  post_id: number
  post_title: string
  post_content: string
  user_id: number
  user_name: string
  user_email: string
  created_at: string
}
