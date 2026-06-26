export interface Comic {
  id: number
  name: string
  url: string
  cover: string
  sort: number
  created_at: string
  updated_at: string
}

export interface CreateComicRequest {
  name: string
  url: string
  cover?: string
  sort?: number
}

export interface UpdateComicRequest {
  name?: string
  url?: string
  cover?: string
  sort?: number
}

export interface ComicListResponse {
  list: Comic[]
  total: number
}
