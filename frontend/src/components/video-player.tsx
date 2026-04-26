import { useEffect, useRef, useState } from 'react'
import { createShakaPlayer } from '@/lib/shaka-player'

type Props = {
  title: string
  manifestUrl: string
}

export function VideoPlayer({ title, manifestUrl }: Props) {
  const videoRef = useRef<HTMLVideoElement>(null)
  const [error, setError] = useState<string | null>(null)
  const [status, setStatus] = useState<'initializing' | 'ready' | 'failed'>('initializing')

  useEffect(() => {
    const video = videoRef.current
    if (!video) return

    let active = true
    let cleanup: (() => Promise<void>) | undefined

    setError(null)
    setStatus('initializing')

    void createShakaPlayer({
      videoElement: video,
      manifestUrl,
      onError: (message) => {
        if (!active) return
        setError(message)
        setStatus('failed')
      },
    }).then((result) => {
      if (!result.ok) return
      if (!active) {
        void result.cleanup().catch(() => {})
        return
      }
      cleanup = result.cleanup
      setStatus('ready')
    })

    return () => {
      active = false
      if (cleanup) {
        void cleanup().catch(() => {})
      }
    }
  }, [manifestUrl])

  return (
    <div className="flex flex-col gap-3">
      <video ref={videoRef} aria-label={`${title} playback`} controls={status === 'ready'} className="w-full rounded-lg bg-black" />
      <p role="status" aria-live="polite" className="text-sm text-muted-foreground">
        {status === 'initializing' ? 'Loading player...' : status === 'ready' ? 'Playback ready.' : error ? 'Playback failed.' : ''}
      </p>
      {error ? <p className="text-sm text-destructive">{error}</p> : null}
    </div>
  )
}
