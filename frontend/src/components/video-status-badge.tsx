import { Badge } from '@/components/ui/badge'
import { formatStatus } from '@/lib/format'
import type { VideoStatus } from '@/lib/types'

type Props = {
  status: VideoStatus
}

export function VideoStatusBadge({ status }: Props) {
  return <Badge variant={status === 'ready' ? 'default' : 'secondary'}>{formatStatus(status)}</Badge>
}
