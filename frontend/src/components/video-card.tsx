import { Link } from 'react-router-dom'
import type { VideoItem } from '@/lib/types'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { VideoStatusBadge } from './video-status-badge'

type Props = {
  video: VideoItem
}

export function VideoCard({ video }: Props) {
  return (
    <Card className="border-border bg-card">
      <CardHeader>
        <CardTitle className="text-base">{video.title}</CardTitle>
      </CardHeader>
      <CardContent className="flex items-center justify-between gap-4">
        <VideoStatusBadge status={video.status} />
        {video.errorMessage ? <p className="text-sm text-destructive">{video.errorMessage}</p> : null}
      </CardContent>
      <CardFooter>
        <Button asChild variant="secondary">
          <Link to={`/videos/${video.id}`}>Open</Link>
        </Button>
      </CardFooter>
    </Card>
  )
}
