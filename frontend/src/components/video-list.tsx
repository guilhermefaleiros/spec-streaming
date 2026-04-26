import { Link } from 'react-router-dom'
import { VideoItem } from '../lib/types'

type Props = {
  videos: VideoItem[]
}

export function VideoList({ videos }: Props) {
  if (!videos || videos.length === 0) {
    return <p style={{ color: '#666', fontStyle: 'italic' }}>No videos yet. Upload one!</p>
  }

  return (
    <ul style={{ listStyle: 'none', padding: 0, margin: 0 }}>
      {videos.map((video) => (
        <li 
          key={video.id} 
          style={{ 
            marginBottom: '10px', 
            padding: '15px', 
            border: '1px solid #ddd',
            borderRadius: '8px',
            backgroundColor: '#fafafa'
          }}
        >
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Link 
              to={`/videos/${video.id}`}
              style={{ 
                fontSize: '18px', 
                fontWeight: 'bold',
                textDecoration: 'none',
                color: '#2196f3'
              }}
            >
              {video.title}
            </Link>
            <StatusBadge status={video.status} />
          </div>
          {video.errorMessage && (
            <p style={{ marginTop: '8px', color: '#c62828', fontSize: '14px' }}>
              Error: {video.errorMessage}
            </p>
          )}
        </li>
      ))}
    </ul>
  )
}

function StatusBadge({ status }: { status: string }) {
  const colors: Record<string, string> = {
    uploaded: '#ff9800',
    processing: '#2196f3',
    ready: '#4caf50',
    failed: '#f44336'
  }

  return (
    <span 
      style={{ 
        padding: '4px 12px', 
        borderRadius: '12px',
        backgroundColor: colors[status] || '#999',
        color: 'white',
        fontSize: '14px',
        fontWeight: 'bold',
        textTransform: 'uppercase'
      }}
    >
      {status}
    </span>
  )
}
