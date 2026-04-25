import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { getVideo } from '../lib/api'
import { VideoItem } from '../lib/types'
import { VideoPlayer } from '../components/video-player'

export function VideoPage() {
  const { id } = useParams<{ id: string }>()
  const [video, setVideo] = useState<VideoItem | null>(null)

  useEffect(() => {
    if (!id) return
    getVideo(id).then(setVideo).catch(() => setVideo(null))
  }, [id])

  if (!video) {
    return <div>Loading...</div>
  }

  const manifestUrl = `${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/videos/${video.id}/stream/manifest.mpd`

  return (
    <div>
      <h1>{video.title}</h1>
      <VideoPlayer manifestUrl={manifestUrl} />
    </div>
  )
}
