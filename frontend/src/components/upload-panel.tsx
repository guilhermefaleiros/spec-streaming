import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { UploadForm } from './upload-form'

type Props = {
  onSubmit: (input: { title: string; file: File }) => Promise<void>
}

export function UploadPanel({ onSubmit }: Props) {
  return (
    <Card className="border-border bg-card">
      <CardHeader>
        <CardTitle>Upload</CardTitle>
        <CardDescription>Send MP4 videos and watch their status update.</CardDescription>
      </CardHeader>
      <CardContent>
        <UploadForm onSubmit={onSubmit} />
      </CardContent>
    </Card>
  )
}
