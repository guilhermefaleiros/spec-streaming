export type VideoStatus = 'uploaded' | 'processing' | 'ready' | 'failed'

export interface VideoItem {
  id: string
  title: string
  status: VideoStatus
  errorMessage?: string
}
