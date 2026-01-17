import { renderHook, act } from '@testing-library/react'
import useLocalStorage from '@/lib/hooks/use-local-storage'

describe('useLocalStorage', () => {
  beforeEach(() => {
    // Clear localStorage before each test
    localStorage.clear()
    jest.clearAllMocks()
  })

  it('returns initial value when localStorage is empty', () => {
    const { result } = renderHook(() => useLocalStorage('test-key', 'initial'))

    expect(result.current[0]).toBe('initial')
  })

  it('returns stored value from localStorage', () => {
    localStorage.setItem('test-key', JSON.stringify('stored-value'))

    const { result } = renderHook(() => useLocalStorage('test-key', 'initial'))

    expect(result.current[0]).toBe('stored-value')
  })

  it('updates localStorage when setValue is called', () => {
    const { result } = renderHook(() => useLocalStorage('test-key', 'initial'))

    act(() => {
      result.current[1]('new-value')
    })

    expect(result.current[0]).toBe('new-value')
    expect(localStorage.getItem('test-key')).toBe(JSON.stringify('new-value'))
  })

  it('handles object values', () => {
    const initialValue = { name: 'John', age: 30 }

    const { result } = renderHook(() => useLocalStorage('user', initialValue))

    expect(result.current[0]).toEqual(initialValue)

    const newValue = { name: 'Jane', age: 25 }
    act(() => {
      result.current[1](newValue)
    })

    expect(result.current[0]).toEqual(newValue)
    expect(JSON.parse(localStorage.getItem('user')!)).toEqual(newValue)
  })

  it('handles array values', () => {
    const initialValue = [1, 2, 3]

    const { result } = renderHook(() => useLocalStorage('numbers', initialValue))

    expect(result.current[0]).toEqual(initialValue)

    const newValue = [4, 5, 6]
    act(() => {
      result.current[1](newValue)
    })

    expect(result.current[0]).toEqual(newValue)
  })

  it('handles boolean values', () => {
    const { result } = renderHook(() => useLocalStorage('flag', false))

    expect(result.current[0]).toBe(false)

    act(() => {
      result.current[1](true)
    })

    expect(result.current[0]).toBe(true)
  })

  it('returns initial value when localStorage has invalid JSON', () => {
    localStorage.setItem('test-key', 'invalid-json')

    const { result } = renderHook(() => useLocalStorage('test-key', 'initial'))

    expect(result.current[0]).toBe('initial')
  })
})
