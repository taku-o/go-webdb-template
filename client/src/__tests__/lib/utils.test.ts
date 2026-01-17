import {
  cn,
  formatRelativeTime,
  nFormatter,
  capitalize,
  truncate,
} from '@/lib/utils'

describe('cn', () => {
  it('merges class names', () => {
    const result = cn('class1', 'class2')
    expect(result).toBe('class1 class2')
  })

  it('handles conditional classes', () => {
    const result = cn('base', false && 'conditional', 'another')
    expect(result).toBe('base another')
  })

  it('merges tailwind classes correctly', () => {
    const result = cn('px-2 py-1', 'px-4')
    expect(result).toBe('py-1 px-4')
  })
})

describe('formatRelativeTime', () => {
  beforeEach(() => {
    // Mock Date.now to return a fixed time
    jest.useFakeTimers()
    jest.setSystemTime(new Date('2024-01-15T12:00:00Z'))
  })

  afterEach(() => {
    jest.useRealTimers()
  })

  it('returns "たった今" for less than 1 minute ago', () => {
    const result = formatRelativeTime('2024-01-15T11:59:30Z')
    expect(result).toBe('たった今')
  })

  it('returns minutes for less than 1 hour ago', () => {
    const result = formatRelativeTime('2024-01-15T11:30:00Z')
    expect(result).toBe('30分前')
  })

  it('returns hours for less than 24 hours ago', () => {
    const result = formatRelativeTime('2024-01-15T06:00:00Z')
    expect(result).toBe('6時間前')
  })

  it('returns days for less than 7 days ago', () => {
    const result = formatRelativeTime('2024-01-12T12:00:00Z')
    expect(result).toBe('3日前')
  })

  it('returns month and day for same year, more than 7 days ago', () => {
    const result = formatRelativeTime('2024-01-01T12:00:00Z')
    expect(result).toBe('1月1日')
  })

  it('returns full date for different year', () => {
    const result = formatRelativeTime('2023-06-15T12:00:00Z')
    expect(result).toBe('2023年6月15日')
  })
})

describe('nFormatter', () => {
  it('returns "0" for zero', () => {
    expect(nFormatter(0)).toBe('0')
  })

  it('returns number as is for small numbers', () => {
    expect(nFormatter(500)).toBe('500')
  })

  it('formats thousands with K', () => {
    expect(nFormatter(1000)).toBe('1K')
    expect(nFormatter(1500)).toBe('1.5K')
  })

  it('formats millions with M', () => {
    expect(nFormatter(1000000)).toBe('1M')
  })

  it('formats billions with G', () => {
    expect(nFormatter(1000000000)).toBe('1G')
  })

  it('respects digits parameter', () => {
    expect(nFormatter(1234, 2)).toBe('1.23K')
  })
})

describe('capitalize', () => {
  it('capitalizes first letter', () => {
    expect(capitalize('hello')).toBe('Hello')
  })

  it('returns empty string for empty input', () => {
    expect(capitalize('')).toBe('')
  })

  it('handles single character', () => {
    expect(capitalize('a')).toBe('A')
  })

  it('preserves rest of string', () => {
    expect(capitalize('hELLO')).toBe('HELLO')
  })

  it('handles non-string input gracefully', () => {
    // @ts-ignore - testing edge case
    expect(capitalize(null)).toBe(null)
    // @ts-ignore - testing edge case
    expect(capitalize(undefined)).toBe(undefined)
  })
})

describe('truncate', () => {
  it('does not truncate short strings', () => {
    expect(truncate('hello', 10)).toBe('hello')
  })

  it('truncates long strings with ellipsis', () => {
    expect(truncate('hello world', 5)).toBe('hello...')
  })

  it('handles exact length', () => {
    expect(truncate('hello', 5)).toBe('hello')
  })

  it('handles empty string', () => {
    expect(truncate('', 5)).toBe('')
  })

  it('handles null/undefined gracefully', () => {
    // @ts-ignore - testing edge case
    expect(truncate(null, 5)).toBe(null)
    // @ts-ignore - testing edge case
    expect(truncate(undefined, 5)).toBe(undefined)
  })
})
