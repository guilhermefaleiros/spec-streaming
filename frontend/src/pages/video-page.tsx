import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { getVideo } from '../lib/api'
import { VideoItem } from '../lib/types'
import { VideoPlayer } from '../components/video-player'

export function VideoPage() {
  const { id } = useParams<{ id: string }>()
  const [video, setVideo] = useState<VideoItem | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!id) return
    getVideo(id)
      .then(setVideo)
      .catch((err) => setError(err instanceof Error ? err.message : 'Failed to load video'))
  }, [id])

  if (error) {
    return (
      <div style={{ maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
        <Link to="/" style={{ color: '#2196f3', textDecoration: 'none' }}>← Back to list</Link>
        <div style={{ marginTop: '20px', padding: '20px', backgroundColor: '#ffebee', color: '#c62828', borderRadius: '8px' }}>
          Error: {error}
        </div>
      </div>
    )
  }

  if (!video) {
    return (
      <div style={{ maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
        <Link to="/" style={{ color: '#2196f3', textDecoration: 'none' }}>← Back to list</Link>
        <p style={{ marginTop: '20px' }}>Loading...</p>
      </div>
    )
  }

  const manifestUrl = `${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/videos/${video.id}/stream/manifest.mpd`

  const statusColors: Record<string, string> = {
    uploaded: '#ff9800',
    processing: '#2196f3',
    ready: '#4caf50',
    failed: '#f44336'
  }

  return (
    <div style={{ maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
      <Link to="/" style={{ color: '#2196f3', textDecoration: 'none' }}>← Back to list</Link>
      
      <div style={{ marginTop: '20px', padding: '20px', border: '1px solid #ddd', borderRadius: '8px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
          <h1 style={{ margin: 0 }}>{video.title}</h1>
          <span style={{ 
            padding: '4px 12px', 
            borderRadius: '12px',
            backgroundColor: statusColors[video.status] || '#999',
            color: 'white',
            fontSize: '14px',
            fontWeight: 'bold',
            textTransform: 'uppercase'
          }}>
            {video.status}
          </span>
        </div>

        {video.status === 'ready' ? (
          <VideoPlayer manifestUrl={manifestUrl} />
        ) : (
          <div style={{ 
            padding: '40px', 
            textAlign: 'center', 
            backgroundColor: '#f5f5f5',
            borderRadius: '8px'
          }}>
            <p>Video is not ready for playback.</p>
            <p style={{ color: '#666' }}>Current status: {video.status}</p>
          </div>
        )}

        {video.errorMessage && (
          <div style={{ marginTop: '20px', padding: '10px', backgroundColor: '#ffebee', color: '#c62828', borderRadius: '4px' }}>
            Error: {video.errorMessage}
          </div>
        )}
      </div>
    </div>
  )
}
