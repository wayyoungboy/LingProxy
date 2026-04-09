import { describe, it, expect, beforeEach } from 'vitest'
import { ref } from 'vue'
import { useLoading } from './useLoading'

describe('useLoading', () => {
  it('should initialize with false by default', () => {
    const { loading } = useLoading()
    expect(loading.value).toBe(false)
  })

  it('should initialize with custom value', () => {
    const { loading } = useLoading(true)
    expect(loading.value).toBe(true)
  })

  it('should set loading via setLoading', () => {
    const { loading, setLoading } = useLoading()
    setLoading(true)
    expect(loading.value).toBe(true)
    setLoading(false)
    expect(loading.value).toBe(false)
  })

  it('should start and stop loading', () => {
    const { loading, startLoading, stopLoading } = useLoading()
    expect(loading.value).toBe(false)
    startLoading()
    expect(loading.value).toBe(true)
    stopLoading()
    expect(loading.value).toBe(false)
  })

  it('should set loading to true during async operation', async () => {
    const { loading, withLoading } = useLoading()

    let resolved = false
    const promise = withLoading(async () => {
      await new Promise(resolve => setTimeout(resolve, 10))
      resolved = true
      return 'result'
    })

    expect(loading.value).toBe(true)
    const result = await promise
    expect(loading.value).toBe(false)
    expect(result).toBe('result')
    expect(resolved).toBe(true)
  })

  it('should reset loading even if async operation throws', async () => {
    const { loading, withLoading } = useLoading()

    await expect(
      withLoading(async () => {
        throw new Error('test error')
      })
    ).rejects.toThrow('test error')

    expect(loading.value).toBe(false)
  })
})
