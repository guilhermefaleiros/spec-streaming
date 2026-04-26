import { useParams } from 'react-router-dom'
import { Link } from 'react-router-dom'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { PlayerStatusPanel } from '@/components/player-status-panel'
import { VideoPlayer } from '@/components/video-player'
import { useVideoQuery } from '@/lib/queries'

export function VideoPage() {
  const { id } = useParams<{ id: string }>()
  const query = useVideoQuery(id ?? '')

  if (query.isLoading) {
    return (
      <Card className="border-border bg-card">
        <CardContent className="flex flex-col gap-3 p-6">
          <p className="text-sm text-muted-foreground">Loading...</p>
          <Button asChild variant="secondary" className="w-fit">
            <Link to="/catalog">Back to catalog</Link>
          </Button>
        </CardContent>
      </Card>
    )
  }

  if (query.isError || !query.data) {
    return (
      <Card className="border-border bg-card">
        <CardContent className="flex flex-col gap-3 p-6">
          <p className="text-sm text-destructive">Failed to load video.</p>
          <Button asChild variant="secondary" className="w-fit">
            <Link to="/catalog">Back to catalog</Link>
          </Button>
        </CardContent>
      </Card>
    )
  }

  const video = query.data
  const manifestUrl = `${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/videos/${video.id}/stream/manifest.mpd`
  const failure = video.status === 'failed'

  return (
    <Card className="border-border bg-card">
      <CardHeader>
        <CardTitle>{video.title}</CardTitle>
      </CardHeader>
      <CardContent className="flex flex-col gap-4">
        <div className="flex flex-col gap-3">
          <PlayerStatusPanel title="Playback" status={video.status} />
          {video.errorMessage ? <p className="text-sm text-muted-foreground">{video.errorMessage}</p> : null}
        </div>
        {failure ? (
          <div className="flex flex-col gap-3 rounded-lg border border-border bg-muted p-4">
            <p className="text-sm text-destructive">Playback failed.</p>
            {video.errorMessage ? <p className="text-sm text-muted-foreground">{video.errorMessage}</p> : null}
            <Button asChild variant="secondary" className="w-fit">
              <Link to="/catalog">Back to catalog</Link>
            </Button>
          </div>
        ) : video.status === 'ready' ? (
          <VideoPlayer title={video.title} manifestUrl={manifestUrl} />
        ) : (
          <div className="flex flex-col gap-3 rounded-lg border border-border bg-muted p-4">
            <p className="text-sm text-muted-foreground">
              Video is not ready for playback.
            </p>
            <Button asChild variant="secondary" className="w-fit">
              <Link to="/catalog">Back to catalog</Link>
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
