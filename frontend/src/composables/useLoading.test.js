import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { ref } from 'vue'

// Simple composable test
describe('useLoading', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('should return loading state functions', async () => {
    // Import the composable
    const { useLoading } = await import('../src/composables/useLoading')

    const { loading, withLoading } = useLoading()

    expect(loading.value).toBe(false)
  })

  it('should set loading to true during async operation', async () => {
    const { useLoading } = await import('../src/composables/useLoading')

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
})