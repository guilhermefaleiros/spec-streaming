import { describe, expect, it } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { VideoCard } from './video-card'

describe('VideoCard', () => {
  it('renders title, status badge, and player link', () => {
    render(
      <MemoryRouter>
        <VideoCard video={{ id: '1', title: 'Trailer', status: 'ready' }} />
      </MemoryRouter>,
    )

    expect(screen.getByText('Trailer')).toBeInTheDocument()
    expect(screen.getByText(/ready/i)).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /open/i })).toHaveAttribute('href', '/videos/1')
  })
})
