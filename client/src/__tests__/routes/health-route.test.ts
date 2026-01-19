import { GET } from '@/app/health/route'

describe('GET /health', () => {
  it('should return 200 OK with "OK" body', async () => {
    const response = await GET()

    expect(response.status).toBe(200)
    expect(await response.text()).toBe('OK')
    expect(response.headers.get('Content-Type')).toBe('text/plain')
  })

  it('should not require authentication', async () => {
    const response = await GET()

    // 認証エラーが発生しないことを確認
    expect(response.status).toBe(200)
  })
})
