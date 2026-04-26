import { useEffect, useRef } from 'react'

type Props = {
  manifestUrl: string
}

export function VideoPlayer({ manifestUrl }: Props) {
  const videoRef = useRef<HTMLVideoElement>(null)

  useEffect(() => {
    const video = videoRef.current
    if (!video) return

    let player: any = null

    const initPlayer = async () => {
      try {
        const dashjs = await import('dashjs')
        const mediaPlayer = (dashjs as any).MediaPlayer || (dashjs as any).default?.MediaPlayer
        
        if (!mediaPlayer) {
          console.error('dashjs MediaPlayer not available')
          return
        }

        player = mediaPlayer().create()
        player.initialize(video, manifestUrl, true)
        player.setAutoPlay(true)
      } catch (err) {
        console.error('Failed to load dashjs:', err)
      }
    }

    initPlayer()

    return () => {
      if (player) {
        player.reset()
      }
    }
  }, [manifestUrl])

  return (
    <div style={{ marginTop: '20px' }}>
      <video
        ref={videoRef}
        controls
        style={{
          width: '100%',
          maxWidth: '800px',
          borderRadius: '8px',
          backgroundColor: '#000'
        }}
      />
    </div>
  )
}
