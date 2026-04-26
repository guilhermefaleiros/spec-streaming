import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'

type Props = {
  title: string
  description: string
}

export function ErrorState({ title, description }: Props) {
  return (
    <Alert>
      <AlertTitle>{title}</AlertTitle>
      <AlertDescription>{description}</AlertDescription>
    </Alert>
  )
}
