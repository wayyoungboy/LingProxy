/**
 * 分页管理 Composable
 */
import { ref, reactive, computed } from 'vue'
import { PAGINATION } from '../utils/constants'

export function usePagination(initialPageSize = PAGINATION.DEFAULT_PAGE_SIZE) {
  const pagination = reactive({
    currentPage: 1,
    pageSize: initialPageSize,
    total: 0
  })

  const resetPagination = () => {
    pagination.currentPage = 1
    pagination.total = 0
  }

  const setTotal = total => {
    pagination.total = total
  }

  const setPageSize = size => {
    pagination.pageSize = size
    pagination.currentPage = 1
  }

  const setCurrentPage = page => {
    pagination.currentPage = page
  }

  const totalPages = computed(() => {
    return Math.ceil(pagination.total / pagination.pageSize)
  })

  const hasNextPage = computed(() => {
    return pagination.currentPage < totalPages.value
  })

  const hasPrevPage = computed(() => {
    return pagination.currentPage > 1
  })

  return {
    pagination,
    resetPagination,
    setTotal,
    setPageSize,
    setCurrentPage,
    totalPages,
    hasNextPage,
    hasPrevPage
  }
}
