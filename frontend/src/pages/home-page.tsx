import { useEffect, useState } from 'react'
import { listVideos, uploadVideo } from '../lib/api'
import { VideoItem } from '../lib/types'
import { UploadForm } from '../components/upload-form'
import { VideoList } from '../components/video-list'

export function HomePage() {
  const [videos, setVideos] = useState<VideoItem[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let active = true

    async function refresh() {
      try {
        setError(null)
        const next = await listVideos()
        if (active) {
          setVideos(next)
          setIsLoading(false)
        }
      } catch (err) {
        if (active) {
          setError(err instanceof Error ? err.message : 'Failed to load videos')
          setIsLoading(false)
        }
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
    try {
      setError(null)
      await uploadVideo(input)
      const next = await listVideos()
      setVideos(next)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Upload failed')
    }
  }

  return (
    <div style={{ maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
      <h1 style={{ marginBottom: '30px' }}>Spec Streaming</h1>
      
      <section style={{ marginBottom: '40px', padding: '20px', border: '1px solid #ddd', borderRadius: '8px' }}>
        <h2 style={{ marginTop: 0, marginBottom: '20px' }}>Upload Video</h2>
        <UploadForm onSubmit={handleUpload} />
        {error && (
          <div style={{ marginTop: '15px', padding: '10px', backgroundColor: '#ffebee', color: '#c62828', borderRadius: '4px' }}>
            Error: {error}
          </div>
        )}
      </section>

      <section>
        <h2 style={{ marginBottom: '20px' }}>Videos</h2>
        {isLoading ? (
          <p>Loading videos...</p>
        ) : (
          <VideoList videos={videos} />
        )}
      </section>
    </div>
  )
}
