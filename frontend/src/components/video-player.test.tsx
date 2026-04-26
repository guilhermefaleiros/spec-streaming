import { describe, expect, it, vi } from 'vitest'
import { render, waitFor } from '@testing-library/react'
import { VideoPlayer } from './video-player'

const mocks = vi.hoisted(() => {
  const cleanup = vi.fn().mockResolvedValue(undefined)
  const createShakaPlayer = vi.fn().mockImplementation(
    () => new Promise(() => undefined),
  )
  return { cleanup, createShakaPlayer }
})

vi.mock('@/lib/shaka-player', () => ({
  createShakaPlayer: mocks.createShakaPlayer,
}))

describe('VideoPlayer', () => {
  it('renders a video element', () => {
    const { container } = render(<VideoPlayer title="Trailer" manifestUrl="/manifest.mpd" />)

    expect(container.querySelector('video')).toBeInTheDocument()
  })

  it('shows a loading message and accessible name', async () => {
    const { getByText, container } = render(<VideoPlayer title="Trailer" manifestUrl="/manifest.mpd" />)

    await waitFor(() => {
      expect(getByText(/loading player/i)).toBeInTheDocument()
    })
    expect(container.querySelector('video')).toHaveAttribute('aria-label', 'Trailer playback')
  })
})
