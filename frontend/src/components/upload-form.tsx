import { useState } from 'react'

type Props = {
  onSubmit: (input: { title: string; file: File }) => Promise<void>
}

export function UploadForm({ onSubmit }: Props) {
  const [title, setTitle] = useState('')
  const [file, setFile] = useState<File | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!file) return
    setIsSubmitting(true)
    try {
      await onSubmit({ title, file })
      setTitle('')
      setFile(null)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} data-testid="upload-form">
      <div style={{ marginBottom: '15px' }}>
        <label 
          htmlFor="title" 
          style={{ display: 'block', marginBottom: '5px', fontWeight: 'bold' }}
        >
          Title
        </label>
        <input
          id="title"
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="Enter video title"
          style={{ 
            width: '100%', 
            padding: '10px', 
            fontSize: '16px',
            border: '1px solid #ccc',
            borderRadius: '4px',
            boxSizing: 'border-box'
          }}
          required
        />
      </div>
      <div style={{ marginBottom: '15px' }}>
        <label 
          htmlFor="file" 
          style={{ display: 'block', marginBottom: '5px', fontWeight: 'bold' }}
        >
          File
        </label>
        <input
          id="file"
          type="file"
          accept="video/mp4"
          onChange={(e) => setFile(e.target.files?.[0] ?? null)}
          style={{ 
            width: '100%', 
            padding: '10px', 
            fontSize: '16px',
            border: '1px solid #ccc',
            borderRadius: '4px',
            boxSizing: 'border-box'
          }}
          required
        />
        {file && (
          <p style={{ marginTop: '5px', fontSize: '14px', color: '#666' }}>
            Selected: {file.name} ({(file.size / 1024 / 1024).toFixed(2)} MB)
          </p>
        )}
      </div>
      <button 
        type="submit" 
        disabled={isSubmitting || !file}
        style={{ 
          padding: '12px 24px', 
          fontSize: '16px',
          backgroundColor: isSubmitting || !file ? '#ccc' : '#2196f3',
          color: 'white',
          border: 'none',
          borderRadius: '4px',
          cursor: isSubmitting || !file ? 'not-allowed' : 'pointer'
        }}
      >
        {isSubmitting ? 'Uploading...' : 'Upload'}
      </button>
    </form>
  )
}
