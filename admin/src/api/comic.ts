import request from "@/utils/request"
import type { Comic, ComicListResponse, CreateComicRequest, UpdateComicRequest } from "@/types/comic"

const BASE = "/admin/comics"

export const getComics = (params?: { page?: number; page_size?: number }) =>
  request.get<ComicListResponse>(BASE, { params }) as unknown as Promise<ComicListResponse>

export const getComicDetail = (id: number) =>
  request.get<Comic>(BASE + "/" + id) as unknown as Promise<Comic>

export const createComic = (data: CreateComicRequest) =>
  request.post<Comic>(BASE, data) as unknown as Promise<Comic>

export const updateComic = (id: number, data: UpdateComicRequest) =>
  request.put<Comic>(BASE + "/" + id, data) as unknown as Promise<Comic>

export const deleteComic = (id: number) =>
  request.delete(BASE + "/" + id)
