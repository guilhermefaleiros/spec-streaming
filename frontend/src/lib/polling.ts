import type { VideoItem } from './types'

export function shouldPollVideo(video?: VideoItem | null) {
  return video?.status === 'processing' ? 3000 : false
}

export function shouldPollVideos(videos: VideoItem[]) {
  return videos.some((video) => video.status === 'processing') ? 3000 : false
}
