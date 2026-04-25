import { useState } from 'react'

type Props = {
  onSubmit: (input: { title: string; file: File }) => Promise<void>
}

export function UploadForm({ onSubmit }: Props) {
  const [title, setTitle] = useState('')
  const [file, setFile] = useState<File | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!file) return
    await onSubmit({ title, file })
    setTitle('')
    setFile(null)
  }

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label htmlFor="title">Title</label>
        <input
          id="title"
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />
      </div>
      <div>
        <label htmlFor="file">File</label>
        <input
          id="file"
          type="file"
          accept="video/mp4"
          onChange={(e) => setFile(e.target.files?.[0] ?? null)}
        />
      </div>
      <button type="submit">Upload</button>
    </form>
  )
}
