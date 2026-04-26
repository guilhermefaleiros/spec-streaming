import { describe, expect, it, vi } from 'vitest'
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
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

    // Wait for state update
    await waitFor(() => {
      expect(screen.getByText(/Selected: trailer.mp4/)).toBeInTheDocument()
    })

    // Submit the form directly (bypass HTML5 validation)
    fireEvent.submit(screen.getByTestId('upload-form'))

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({ title: 'Trailer', file })
    })
  })
})
