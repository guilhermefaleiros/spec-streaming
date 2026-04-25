import { VideoItem } from './types'

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080'

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
  return res.json()
}

export async function listVideos(): Promise<VideoItem[]> {
  const res = await fetch(`${API_BASE}/videos`)
  if (!res.ok) {
    throw new Error('Failed to list videos')
  }
  return res.json()
}

export async function getVideo(id: string): Promise<VideoItem> {
  const res = await fetch(`${API_BASE}/videos/${id}`)
  if (!res.ok) {
    throw new Error('Failed to get video')
  }
  return res.json()
}
