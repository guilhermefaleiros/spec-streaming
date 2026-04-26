import { Card, CardContent } from '@/components/ui/card'

type Props = {
  title: string
  description: string
}

export function EmptyState({ title, description }: Props) {
  return (
    <Card className="border-border bg-card">
      <CardContent className="p-6">
        <p className="text-base font-semibold">{title}</p>
        <p className="mt-2 text-sm text-muted-foreground">{description}</p>
      </CardContent>
    </Card>
  )
}
