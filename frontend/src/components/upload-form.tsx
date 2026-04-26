import { useState, type FormEvent } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

type Props = {
  onSubmit: (input: { title: string; file: File }) => Promise<void>
}

export function UploadForm({ onSubmit }: Props) {
  const [title, setTitle] = useState('')
  const [file, setFile] = useState<File | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (!file) return
    setError('')
    setIsSubmitting(true)
    try {
      await onSubmit({ title, file })
      setTitle('')
      setFile(null)
    } catch {
      setError('Upload failed. Please try again.')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <form className="flex flex-col gap-4" onSubmit={handleSubmit} data-testid="upload-form">
      <div className="flex flex-col gap-2">
        <Label htmlFor="title">Title</Label>
        <Input id="title" value={title} onChange={(e) => setTitle(e.target.value)} placeholder="Enter video title" required />
      </div>
      <div className="flex flex-col gap-2">
        <Label htmlFor="file">File</Label>
        <Input id="file" type="file" accept="video/mp4" onChange={(e) => setFile(e.target.files?.[0] ?? null)} required />
        {file ? <p className="text-sm text-muted-foreground">Selected: {file.name} ({(file.size / 1024 / 1024).toFixed(2)} MB)</p> : null}
      </div>
      {error ? (
        <p className="text-sm text-destructive" role="alert" aria-live="polite">
          {error}
        </p>
      ) : null}
      <Button type="submit" disabled={isSubmitting || !file}>
        {isSubmitting ? 'Uploading...' : 'Upload'}
      </Button>
    </form>
  )
}
