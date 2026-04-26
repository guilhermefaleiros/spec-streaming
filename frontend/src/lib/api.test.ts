import { beforeEach, describe, expect, it, vi } from 'vitest'
import { listVideos } from './api'

describe('listVideos', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
  })

  it('surfaces fetch failures', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('network down')))

    await expect(listVideos()).rejects.toThrow('network down')
  })

  it('normalizes uppercase backend keys to the frontend shape', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: true,
        json: async () => [
          {
            ID: 'video-1',
            Title: 'Trailer',
            Status: 'ready',
            ErrorMessage: '',
          },
        ],
      }),
    )

    await expect(listVideos()).resolves.toEqual([
      {
        id: 'video-1',
        title: 'Trailer',
        status: 'ready',
        errorMessage: '',
      },
    ])
  })
})
