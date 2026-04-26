import { describe, expect, it } from 'vitest'
import { shouldPollVideo, shouldPollVideos } from './polling'

describe('polling helpers', () => {
  it('keeps polling a processing video detail page', () => {
    expect(shouldPollVideo({ id: '1', title: 'Trailer', status: 'processing' })).toBe(3000)
  })

  it('stops polling a ready video detail page', () => {
    expect(shouldPollVideo({ id: '1', title: 'Trailer', status: 'ready' })).toBe(false)
  })

  it('keeps polling while any video is processing', () => {
    expect(
      shouldPollVideos([
        { id: '1', title: 'Trailer', status: 'ready' },
        { id: '2', title: 'Interview', status: 'processing' },
      ]),
    ).toBe(3000)
  })

  it('stops polling when every video is ready', () => {
    expect(shouldPollVideos([{ id: '1', title: 'Trailer', status: 'ready' }])).toBe(false)
  })
})
