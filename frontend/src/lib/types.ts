export type VideoStatus = 'uploaded' | 'processing' | 'ready' | 'failed'

export interface VideoItem {
  id: string
  title: string
  status: VideoStatus
  originalFilename?: string
  sourceStorageKey?: string
  manifestKey?: string
  errorMessage?: string
}
