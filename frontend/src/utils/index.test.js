import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import {
  debounce,
  throttle,
  formatDate,
  formatFileSize,
  formatNumber,
  deepClone,
  validateEmail,
  validateUrl,
  downloadFile
} from '../utils'

describe('debounce', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('should delay function execution', () => {
    const fn = vi.fn()
    const debounced = debounce(fn, 100)
    debounced()
    expect(fn).not.toHaveBeenCalled()
    vi.advanceTimersByTime(100)
    expect(fn).toHaveBeenCalledTimes(1)
  })

  it('should cancel previous call if invoked again within wait time', () => {
    const fn = vi.fn()
    const debounced = debounce(fn, 100)
    debounced()
    debounced()
    debounced()
    vi.advanceTimersByTime(100)
    expect(fn).toHaveBeenCalledTimes(1)
  })

  it('should pass arguments to the debounced function', () => {
    const fn = vi.fn()
    const debounced = debounce(fn, 50)
    debounced('hello', 42)
    vi.advanceTimersByTime(50)
    expect(fn).toHaveBeenCalledWith('hello', 42)
  })

  it('should use default wait time of 300ms', () => {
    const fn = vi.fn()
    const debounced = debounce(fn)
    debounced()
    vi.advanceTimersByTime(299)
    expect(fn).not.toHaveBeenCalled()
    vi.advanceTimersByTime(1)
    expect(fn).toHaveBeenCalledTimes(1)
  })
})

describe('throttle', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('should execute function immediately on first call', () => {
    const fn = vi.fn()
    const throttled = throttle(fn, 100)
    throttled()
    expect(fn).toHaveBeenCalledTimes(1)
  })

  it('should ignore calls within the limit period', () => {
    const fn = vi.fn()
    const throttled = throttle(fn, 100)
    throttled()
    throttled()
    throttled()
    expect(fn).toHaveBeenCalledTimes(1)
  })

  it('should allow execution after the limit period', () => {
    const fn = vi.fn()
    const throttled = throttle(fn, 100)
    throttled()
    vi.advanceTimersByTime(100)
    throttled()
    expect(fn).toHaveBeenCalledTimes(2)
  })

  it('should use default limit of 300ms', () => {
    const fn = vi.fn()
    const throttled = throttle(fn)
    throttled()
    throttled()
    expect(fn).toHaveBeenCalledTimes(1)
    vi.advanceTimersByTime(299)
    throttled()
    expect(fn).toHaveBeenCalledTimes(1)
    vi.advanceTimersByTime(1)
    throttled()
    expect(fn).toHaveBeenCalledTimes(2)
  })
})

describe('formatDate', () => {
  it('should return "--" for falsy input', () => {
    expect(formatDate(null)).toBe('--')
    expect(formatDate(undefined)).toBe('--')
    expect(formatDate('')).toBe('--')
  })

  it('should return "--" for invalid date', () => {
    expect(formatDate('not-a-date')).toBe('--')
  })

  it('should format as date (YYYY-MM-DD)', () => {
    expect(formatDate('2024-03-15T10:30:00Z', 'date')).toBe('2024-03-15')
  })

  it('should format as time (HH:mm:ss)', () => {
    // Use a fixed offset date to avoid timezone issues
    const date = new Date(2024, 2, 15, 10, 30, 45) // local time
    expect(formatDate(date, 'time')).toBe('10:30:45')
  })

  it('should format as datetime (YYYY-MM-DD HH:mm:ss)', () => {
    // Note: timezone dependent, so we just check the date part
    const result = formatDate('2024-03-15T10:30:45Z', 'datetime')
    expect(result).toMatch(/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$/)
  })

  it('should return relative time', () => {
    const now = new Date()
    const result = formatDate(now, 'relative')
    expect(result).toBe('刚刚')
  })

  it('should use toLocaleString for unknown format', () => {
    const result = formatDate('2024-03-15T10:30:00Z', 'unknown')
    // Just verify it returns a non-empty string
    expect(typeof result).toBe('string')
    expect(result.length).toBeGreaterThan(0)
  })
})

describe('formatFileSize', () => {
  it('should return "0 B" for zero or falsy input', () => {
    expect(formatFileSize(0)).toBe('0 B')
    expect(formatFileSize(null)).toBe('0 B')
    expect(formatFileSize(undefined)).toBe('0 B')
  })

  it('should format bytes', () => {
    expect(formatFileSize(500)).toBe('500 B')
  })

  it('should format KB', () => {
    expect(formatFileSize(1024)).toBe('1 KB')
    expect(formatFileSize(1536)).toBe('1.5 KB')
  })

  it('should format MB', () => {
    expect(formatFileSize(1048576)).toBe('1 MB')
  })

  it('should format GB', () => {
    expect(formatFileSize(1073741824)).toBe('1 GB')
  })

  it('should format TB', () => {
    expect(formatFileSize(1099511627776)).toBe('1 TB')
  })
})

describe('formatNumber', () => {
  it('should return "0" for null or undefined', () => {
    expect(formatNumber(null)).toBe('0')
    expect(formatNumber(undefined)).toBe('0')
  })

  it('should format small numbers without commas', () => {
    expect(formatNumber(100)).toBe('100')
    expect(formatNumber(999)).toBe('999')
  })

  it('should format large numbers with commas', () => {
    expect(formatNumber(1000)).toBe('1,000')
    expect(formatNumber(1000000)).toBe('1,000,000')
    expect(formatNumber(1234567)).toBe('1,234,567')
  })
})

describe('deepClone', () => {
  it('should return primitive values as-is', () => {
    expect(deepClone(42)).toBe(42)
    expect(deepClone('hello')).toBe('hello')
    expect(deepClone(true)).toBe(true)
    expect(deepClone(null)).toBe(null)
    expect(deepClone(undefined)).toBe(undefined)
  })

  it('should clone arrays deeply', () => {
    const arr = [1, [2, 3], { a: 4 }]
    const cloned = deepClone(arr)
    expect(cloned).toEqual(arr)
    expect(cloned).not.toBe(arr)
    expect(cloned[1]).not.toBe(arr[1])
    expect(cloned[2]).not.toBe(arr[2])
  })

  it('should clone objects deeply', () => {
    const obj = { a: 1, b: { c: 2 }, d: [3, 4] }
    const cloned = deepClone(obj)
    expect(cloned).toEqual(obj)
    expect(cloned).not.toBe(obj)
    expect(cloned.b).not.toBe(obj.b)
    expect(cloned.d).not.toBe(obj.d)
  })

  it('should clone Date objects', () => {
    const date = new Date('2024-03-15')
    const cloned = deepClone(date)
    expect(cloned).toEqual(date)
    expect(cloned).not.toBe(date)
  })
})

describe('validateEmail', () => {
  it('should validate correct email addresses', () => {
    expect(validateEmail('test@example.com')).toBe(true)
    expect(validateEmail('user.name@domain.org')).toBe(true)
    expect(validateEmail('a@b.co')).toBe(true)
  })

  it('should reject invalid email addresses', () => {
    expect(validateEmail('invalid')).toBe(false)
    expect(validateEmail('invalid@')).toBe(false)
    expect(validateEmail('@domain.com')).toBe(false)
    expect(validateEmail('test@.com')).toBe(false)
    expect(validateEmail('')).toBe(false)
  })
})

describe('validateUrl', () => {
  it('should validate correct URLs', () => {
    expect(validateUrl('https://example.com')).toBe(true)
    expect(validateUrl('http://localhost:3000')).toBe(true)
    expect(validateUrl('https://api.example.com/path?query=1')).toBe(true)
  })

  it('should reject invalid URLs', () => {
    expect(validateUrl('not-a-url')).toBe(false)
    expect(validateUrl('')).toBe(false)
    expect(validateUrl('   ')).toBe(false)
  })
})

describe('downloadFile', () => {
  it('should create a link element and trigger download', () => {
    // Mock URL APIs for jsdom
    const mockCreateObjectURL = vi.fn(() => 'blob:mock-url')
    const mockRevokeObjectURL = vi.fn()
    window.URL.createObjectURL = mockCreateObjectURL
    window.URL.revokeObjectURL = mockRevokeObjectURL

    const blob = new Blob(['test content'], { type: 'text/plain' })

    const linkMock = {
      href: '',
      download: '',
      click: vi.fn(),
      style: {}
    }
    const createElementSpy = vi.spyOn(document, 'createElement').mockReturnValue(linkMock)
    const appendChildSpy = vi.spyOn(document.body, 'appendChild').mockImplementation(() => {})
    const removeChildSpy = vi.spyOn(document.body, 'removeChild').mockImplementation(() => {})

    downloadFile(blob, 'test.txt')

    expect(mockCreateObjectURL).toHaveBeenCalledWith(blob)
    expect(linkMock.download).toBe('test.txt')
    expect(linkMock.click).toHaveBeenCalled()
    expect(appendChildSpy).toHaveBeenCalled()
    expect(removeChildSpy).toHaveBeenCalled()
    expect(mockRevokeObjectURL).toHaveBeenCalledWith('blob:mock-url')

    createElementSpy.mockRestore()
    appendChildSpy.mockRestore()
    removeChildSpy.mockRestore()
  })
})
