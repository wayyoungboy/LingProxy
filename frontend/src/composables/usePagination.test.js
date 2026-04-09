import { describe, it, expect, beforeEach } from 'vitest'
import { usePagination } from './usePagination'
import { PAGINATION } from '../utils/constants'

describe('usePagination', () => {
  it('should initialize with default page size', () => {
    const { pagination } = usePagination()
    expect(pagination.currentPage).toBe(1)
    expect(pagination.pageSize).toBe(PAGINATION.DEFAULT_PAGE_SIZE)
    expect(pagination.total).toBe(0)
  })

  it('should initialize with custom page size', () => {
    const { pagination } = usePagination(25)
    expect(pagination.pageSize).toBe(25)
  })

  it('should reset pagination to initial state', () => {
    const { pagination, resetPagination, setTotal, setCurrentPage } = usePagination()
    setTotal(100)
    setCurrentPage(5)
    expect(pagination.currentPage).toBe(5)
    expect(pagination.total).toBe(100)

    resetPagination()
    expect(pagination.currentPage).toBe(1)
    expect(pagination.total).toBe(0)
    // pageSize should remain unchanged
    expect(pagination.pageSize).toBe(PAGINATION.DEFAULT_PAGE_SIZE)
  })

  it('should set total correctly', () => {
    const { pagination, setTotal } = usePagination()
    setTotal(50)
    expect(pagination.total).toBe(50)
  })

  it('should set page size and reset to page 1', () => {
    const { pagination, setPageSize, setCurrentPage } = usePagination()
    setCurrentPage(3)
    setPageSize(50)
    expect(pagination.pageSize).toBe(50)
    expect(pagination.currentPage).toBe(1)
  })

  it('should set current page correctly', () => {
    const { pagination, setCurrentPage } = usePagination()
    setCurrentPage(5)
    expect(pagination.currentPage).toBe(5)
  })

  it('should compute total pages correctly', () => {
    const { pagination, setTotal, setPageSize, totalPages } = usePagination(10)
    setTotal(95)
    expect(totalPages.value).toBe(10) // ceil(95/10)

    setTotal(100)
    expect(totalPages.value).toBe(10)

    setTotal(0)
    expect(totalPages.value).toBe(0)

    setTotal(1)
    expect(totalPages.value).toBe(1)
  })

  it('should compute hasNextPage correctly', () => {
    const { pagination, setTotal, setPageSize, hasNextPage } = usePagination(10)
    setTotal(50)
    expect(hasNextPage.value).toBe(true) // page 1 of 5

    pagination.currentPage = 5
    expect(hasNextPage.value).toBe(false) // last page

    pagination.currentPage = 6
    expect(hasNextPage.value).toBe(false) // beyond last page
  })

  it('should compute hasPrevPage correctly', () => {
    const { pagination, hasPrevPage } = usePagination()
    expect(hasPrevPage.value).toBe(false) // page 1

    pagination.currentPage = 2
    expect(hasPrevPage.value).toBe(true)

    pagination.currentPage = 1
    expect(hasPrevPage.value).toBe(false)
  })
})
