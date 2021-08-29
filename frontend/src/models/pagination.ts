export interface Pagination<T> {
  limit?: number,
  totalRows?: number,
  totalPages?: number,
  page?: number,
  rows?: T[]
}
