import type { VideoItem } from '@/lib/types'
import { VideoCard } from './video-card'

type Props = {
  videos: VideoItem[]
}

export function VideoGrid({ videos }: Props) {
  return <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">{videos.map((video) => <VideoCard key={video.id} video={video} />)}</div>
}
