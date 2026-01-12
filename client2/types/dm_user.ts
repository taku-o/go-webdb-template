export interface DmUser {
  id: string
  name: string
  email: string
  created_at: string
  updated_at: string
}

export interface CreateDmUserRequest {
  name: string
  email: string
}

export interface UpdateDmUserRequest {
  name?: string
  email?: string
}
