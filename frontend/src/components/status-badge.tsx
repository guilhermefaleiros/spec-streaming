import { VideoStatus } from '../lib/types'

type Props = {
  status: VideoStatus
}

export function StatusBadge({ status }: Props) {
  return <span data-testid="status-badge">{status}</span>
}
