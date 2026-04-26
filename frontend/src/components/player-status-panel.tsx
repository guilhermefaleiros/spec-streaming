import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

type Props = {
  title: string
  status: string
}

export function PlayerStatusPanel({ title, status }: Props) {
  return (
    <Card className="border-border bg-card">
      <CardHeader>
        <CardTitle className="text-base">{title}</CardTitle>
      </CardHeader>
      <CardContent>
        <p role="status" aria-live="polite" className="text-sm text-muted-foreground">Current status: {status}</p>
      </CardContent>
    </Card>
  )
}
