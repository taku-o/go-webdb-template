export interface Post {
  id: string
  user_id: string
  title: string
  content: string
  created_at: string
  updated_at: string
}

export interface CreatePostRequest {
  user_id: string
  title: string
  content: string
}

export interface UpdatePostRequest {
  title?: string
  content?: string
}

export interface UserPost {
  post_id: string
  post_title: string
  post_content: string
  user_id: string
  user_name: string
  user_email: string
  created_at: string
}
