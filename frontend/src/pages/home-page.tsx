import { useEffect, useState } from 'react'
import { listVideos } from '../lib/api'
import { VideoItem } from '../lib/types'
import { UploadForm } from '../components/upload-form'
import { VideoList } from '../components/video-list'

export function HomePage() {
  const [videos, setVideos] = useState<VideoItem[]>([])

  useEffect(() => {
    let active = true

    async function refresh() {
      try {
        const next = await listVideos()
        if (active) setVideos(next)
      } catch {
        // ignore polling errors
      }
    }

    refresh()
    const id = window.setInterval(refresh, 3000)
    return () => {
      active = false
      window.clearInterval(id)
    }
  }, [])

  const handleUpload = async (input: { title: string; file: File }) => {
    // Placeholder: actual upload would call uploadVideo
    // For now just refresh list
    const next = await listVideos()
    setVideos(next)
  }

  return (
    <div>
      <h1>Spec Streaming</h1>
      <UploadForm onSubmit={handleUpload} />
      <VideoList videos={videos} />
    </div>
  )
}
