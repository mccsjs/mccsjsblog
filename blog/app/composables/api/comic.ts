import { $fetch } from 'ofetch'
import type { Comic } from '@@/types/comic'

interface ApiResponse {
  code: number
  data: Comic[]
}

const getBaseURL = () => useRuntimeConfig().public.apiUrl as string

export const getComics = async (): Promise<Comic[]> => {
  try {
    const res = await $fetch<ApiResponse>('/comics', {
      baseURL: getBaseURL(),
    })
    return res.data || []
  } catch {
    return []
  }
}
