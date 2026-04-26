import { describe, expect, it } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { VideoList } from './video-list'

describe('VideoList', () => {
  it('renders ready and processing videos', () => {
    render(
      <MemoryRouter>
        <VideoList
          videos={[
            { id: '1', title: 'Trailer', status: 'ready' },
            { id: '2', title: 'Interview', status: 'processing' },
          ]}
        />
      </MemoryRouter>,
    )

    expect(screen.getByText('Trailer')).toBeInTheDocument()
    expect(screen.getByText('Interview')).toBeInTheDocument()
  })
})
