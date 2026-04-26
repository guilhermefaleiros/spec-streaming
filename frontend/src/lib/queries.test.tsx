import { describe, expect, it } from 'vitest'
import { videoKeys } from './queries'

describe('videoKeys', () => {
  it('builds stable list and detail keys', () => {
    expect(videoKeys.all).toEqual(['videos'])
    expect(videoKeys.detail('abc')).toEqual(['videos', 'abc'])
  })
})
