export interface RegisterJobRequest {
  message?: string
  delay_seconds?: number
  max_retry?: number
}

export interface RegisterJobResponse {
  job_id: string
  status: string
}
