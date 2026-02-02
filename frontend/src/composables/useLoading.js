/**
 * Loading状态管理 Composable
 */
import { ref } from 'vue'

export function useLoading(initialValue = false) {
  const loading = ref(initialValue)

  const setLoading = (value) => {
    loading.value = value
  }

  const startLoading = () => {
    loading.value = true
  }

  const stopLoading = () => {
    loading.value = false
  }

  const withLoading = async (asyncFn) => {
    try {
      loading.value = true
      return await asyncFn()
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    setLoading,
    startLoading,
    stopLoading,
    withLoading
  }
}
