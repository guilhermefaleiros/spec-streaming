import { VideoItem } from '../lib/types'
import { EmptyState } from './empty-state'
import { VideoGrid } from './video-grid'

type Props = {
  videos: VideoItem[]
}

export function VideoList({ videos }: Props) {
  if (!videos || videos.length === 0) {
    return <EmptyState title="No videos yet" description="Upload one to see it appear here." />
  }

  return <VideoGrid videos={videos} />
}
