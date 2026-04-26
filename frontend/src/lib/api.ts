import { VideoItem } from './types'

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080'

type VideoPayload = Partial<VideoItem> & {
  ID?: string
  Title?: string
  Status?: string
  ErrorMessage?: string
  OriginalFilename?: string
  SourceStorageKey?: string
  ManifestKey?: string
}

function normalizeVideo(payload: VideoPayload): VideoItem {
  return {
    id: payload.id ?? payload.ID ?? '',
    title: payload.title ?? payload.Title ?? '',
    status: (payload.status ?? payload.Status ?? 'uploaded') as VideoItem['status'],
    errorMessage: payload.errorMessage ?? payload.ErrorMessage ?? '',
  }
}

export async function uploadVideo(input: { title: string; file: File }): Promise<VideoItem> {
  const formData = new FormData()
  formData.append('title', input.title)
  formData.append('file', input.file)

  const res = await fetch(`${API_BASE}/videos`, {
    method: 'POST',
    body: formData,
  })
  if (!res.ok) {
    throw new Error('Upload failed')
  }
  return normalizeVideo(await res.json())
}

export async function listVideos(): Promise<VideoItem[]> {
  const res = await fetch(`${API_BASE}/videos`)
  if (!res.ok) {
    throw new Error('Failed to list videos')
  }
  const data = await res.json()
  return Array.isArray(data) ? data.map((item) => normalizeVideo(item)) : []
}

export async function getVideo(id: string): Promise<VideoItem> {
  const res = await fetch(`${API_BASE}/videos/${id}`)
  if (!res.ok) {
    throw new Error('Failed to get video')
  }
  return normalizeVideo(await res.json())
}
