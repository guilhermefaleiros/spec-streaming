import { Link } from 'react-router-dom'
import { VideoItem } from '../lib/types'
import { StatusBadge } from './status-badge'

type Props = {
  videos: VideoItem[]
}

export function VideoList({ videos }: Props) {
  return (
    <ul>
      {videos.map((video) => (
        <li key={video.id}>
          <Link to={`/videos/${video.id}`}>{video.title}</Link>
          <StatusBadge status={video.status} />
        </li>
      ))}
    </ul>
  )
}
