import { describe, expect, it, vi } from 'vitest'
import { fireEvent, render, screen } from '@testing-library/react'
import { UploadForm } from './upload-form'

describe('UploadForm', () => {
  it('submits title and file', async () => {
    const onSubmit = vi.fn().mockResolvedValue(undefined)
    render(<UploadForm onSubmit={onSubmit} />)

    fireEvent.change(screen.getByLabelText(/title/i), {
      target: { value: 'Trailer' },
    })

    const file = new File(['video'], 'trailer.mp4', { type: 'video/mp4' })
    fireEvent.change(screen.getByLabelText(/file/i), {
      target: { files: [file] },
    })

    fireEvent.click(screen.getByRole('button', { name: /upload/i }))

    expect(onSubmit).toHaveBeenCalledWith({ title: 'Trailer', file })
  })
})
