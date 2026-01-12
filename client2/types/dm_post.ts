export interface DmPost {
  id: string
  user_id: string
  title: string
  content: string
  created_at: string
  updated_at: string
}

export interface CreateDmPostRequest {
  user_id: string
  title: string
  content: string
}

export interface UpdateDmPostRequest {
  title?: string
  content?: string
}

export interface DmUserPost {
  post_id: string
  post_title: string
  post_content: string
  user_id: string
  user_name: string
  user_email: string
  created_at: string
}
